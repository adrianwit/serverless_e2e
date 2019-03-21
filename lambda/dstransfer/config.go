package dstransfer

import (
	"errors"
	"github.com/viant/dsc"
	"github.com/viant/toolbox/cred"
)

//Config represetns a service config
type Config struct {
	DbHost             string
	DbName             string
	DbSecretParam      string
	StorageSecretParam string
	DefaultStorageURL  string
}


//Validate checks if config is valid
func (c *Config) Validate() error {
	if c.DbHost == "" {
		return errors.New("dbHost was empty")
	}
	if c.DbName == "" {
		return errors.New("dbName was empty")
	}
	if c.DbSecretParam == "" {
		return errors.New("dbSecretParam was empty")
	}
	if c.StorageSecretParam == "" {
		return errors.New("storageSecretParam was empty")
	}
	if c.DefaultStorageURL == "" {
		return errors.New("defaultStorageURL was empty")
	}
	return nil
}


//DbConfig returns dsc.Config or error
func (c *Config) DbConfig(credConfig *cred.Config) (*dsc.Config, error) {
	dbConifg := &dsc.Config{
		DriverName: "mysql",
		Descriptor: "[username]:[password]@tcp([host]:3306)/[dbname]?parseTime=true",
		Parameters: map[string]interface{}{
			"host":   c.DbHost,
			"dbname": c.DbName,
			"username":credConfig.Username,
			"password": credConfig.Password,
		},
	}
	return dbConifg, dbConifg.Init()
}
