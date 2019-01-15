package filemeta

import (
	"bytes"
	"encoding/json"
	"github.com/viant/toolbox/storage"
)

type Service interface {
	UpdateMeta(request *Request) error
}

type service struct {
	storageService storage.Service
}

func (s *service) UpdateMeta(request *Request) error {
	if len(request.ObjectURLs) == 0 {
		return nil
	}
	meta, err := s.loadMeta(request.MetaURL)
	if err != nil {
		return err
	}
	for _, URL := range request.ObjectURLs {
		lineCount, err := countObjectLines(s.storageService, URL)
		if err != nil {
			return err
		}
		meta.Add(URL, lineCount)
	}
	return s.persistMeta(request.MetaURL, meta)
}

func (s *service) persistMeta(URL string, meta *Meta) error {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(meta); err != nil {
		return err
	}
	return s.storageService.Upload(URL, buf)
}

func (s *service) loadMeta(URL string) (*Meta, error) {
	exists, err := s.storageService.Exists(URL)
	if err != nil {
		return nil, err
	}
	if !exists {
		return &Meta{
			Paths:  make(map[string]*FolderInfo),
			Assets: make(map[string]int),
		}, nil
	}
	reader, err := s.storageService.DownloadWithURL(URL)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	meta := &Meta{}
	return meta, json.NewDecoder(reader).Decode(meta)
}

func New(region string) (Service, error) {
	storageService, err := getStorageService(region)
	if err != nil {
		return nil, err
	}
	return &service{
		storageService: storageService,
	}, nil
}
