package dstransfer

import (
	"bytes"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/klauspost/compress/gzip"
	"github.com/viant/dsc"
	"github.com/viant/toolbox/storage"
	_ "github.com/viant/toolbox/storage/gs"
)

//Service represents a data transfer service
type Service interface {
	Copy(request *Request) *Response
}

const (
	statusOK    = "ok"
	statusError = "error"
)

type service struct {
	manager     dsc.Manager
	destStorage storage.Service
}

func (s *service) Copy(request *Request) *Response {
	response := &Response{
		Status: statusOK,
	}
	data := new(bytes.Buffer)
	writer := gzip.NewWriter(data)
	record := make(map[string]interface{})
	err := s.manager.ReadAllWithHandler(request.SQL, nil, func(scanner dsc.Scanner) (toContinue bool, err error) {
		err = scanner.Scan(&record)
		if err == nil {
			err = json.NewEncoder(writer).Encode(record)
		}
		response.Count++
		return err == nil, err
	})
	if err == nil {
		if err = writer.Close(); err == nil {
			err = s.destStorage.Upload(request.DestURL, data)
		}
	}
	if err != nil {
		response.Status = statusError
		response.Message = err.Error()
	}
	return response
}

//New creates a new service for supplied config
func New(config *Config) (Service, error) {
	err := config.Validate()
	if err != nil {
		return nil, err
	}
	var result = &service{}
	dbCred, err := getCredConfigFromParam(config.DbSecretParam)
	if err != nil {
		return nil, err
	}
	dbConfig, err := config.DbConfig(dbCred)
	if err != nil {
		return nil, err
	}

	if result.manager, err = dsc.NewManagerFactory().Create(dbConfig); err != nil {
		return nil, err
	}
	storageParam, err := getParameters(config.StorageSecretParam)
	if err != nil {
		return nil, err
	}
	if result.destStorage, err = storage.NewServiceForURL(config.DefaultStorageURL, *storageParam.Value); err != nil {
		return nil, err
	}
	return result, nil
}
