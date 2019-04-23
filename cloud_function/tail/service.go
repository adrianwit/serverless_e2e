package tail

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/viant/dsc"
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/data"
	"github.com/viant/toolbox/storage"
	"context"
	"io"
	"sync/atomic"
)

type Service interface {
	Transfer(ctx context.Context, resourceURL string) error
}

type service struct {
	storageService storage.Service
	Config         *Config
}


func (s *service) writeBatch(session *Session, collection *data.CompactedSlice) {
	defer session.release()
	dmlProvider, err := session.getDMLProvider(collection)
	if err != nil {
		session.SetError(err)
		return
	}
	sqlProvider := func(item interface{}) *dsc.ParametrizedSQL {
		return dmlProvider.Get(dsc.SQLTypeInsert, item)
	}
	connection, err := session.ConnectionProvider().Get()
	if err != nil {
		return
	}
	defer func() {
		_ = connection.Close()
	}()
	count, err := session.PersistData(connection, collection.Ranger(), session.table.Table, session.dmlProvider, sqlProvider)
	atomic.StoreInt64(&session.transferred, int64(count))
	session.SetError(err)
}


func (s *service) UpdateMeta(session *Session) {
	_, name := toolbox.URLSplit(session.resourceURL)
	destURL := toolbox.URLPathJoin(s.Config.ProcessedURL, name+".meta")
	meta := &Meta{
		ResourceURL: session.resourceURL,
		RecordCount: int(atomic.LoadInt64(&session.count)),
		Transferred:int(atomic.LoadInt64(&session.transferred)),
	}
	if err := session.err; err != nil {
		meta.Error = err.Error()
	}
	data, err := json.Marshal(meta)
	if err == nil {
		err = s.storageService.Upload(destURL, bytes.NewReader(data))
	}
	session.SetError(err)
}

func (s *service) transferData(session *Session, reader io.Reader) error {
	var err error
	var record Record
	scanner := bufio.NewScanner(reader)
	slice := data.NewCompactedSlice(true, true)
	for scanner.Scan() {
		err = scanner.Err()
		if err == io.EOF {
			break
		}
		if err = json.Unmarshal(scanner.Bytes(), &record); err != nil {
			session.SetError(err)
			return err
		}
		atomic.AddInt64(&session.count, 1)
		slice.Add(record)
		if slice.Size() > s.Config.BatchSize {
			session.acquire()
			go s.writeBatch(session, slice)
			slice = data.NewCompactedSlice(true, true)
		}
	}
	if slice.Size() > 0 {
		session.acquire()
		s.writeBatch(session, slice)
	}
	return err
}


func (s *service) Transfer(ctx context.Context, resourceURL string) error {
	route, err := s.Config.Routes.Match(resourceURL)
	if err != nil || route == nil {
		return err
	}
	session, err := NewSession(s.Config, route, resourceURL)
	if err != nil {
		return err
	}
	defer s.UpdateMeta(session)
	reader, err := s.storageService.DownloadWithURL(resourceURL)
	if err != nil {
		session.SetError(err)
		return err
	}
	defer func() { _ = reader.Close() }()
	reader, err = useCompressReaderIfNeeded(reader, resourceURL)
	if err == nil {
		err = s.transferData(session, reader)
		session.Wait()
	}
	session.SetError(err)
	return err
}

var srv Service

func New(config *Config) (Service, error) {
	if srv != nil {
		return srv, nil
	}
	storageService, err := storage.NewServiceForURL(config.ProcessedURL, "")
	if err != nil {
		return nil, err
	}
	srv = &service{
		storageService: storageService,
		Config:         config,
	}
	return srv, nil
}
