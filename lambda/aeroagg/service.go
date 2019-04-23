package aeroagg

import (
	"fmt"
	aero "github.com/aerospike/aerospike-client-go"
)

type Service interface {
	Consume(message *Message) error
}

type service struct {
	*aero.Client
	writePolicy *aero.WritePolicy
	Config      *Config
}

func (s *service) Consume(message *Message) error {
	client, err := s.client()
	if err != nil {
		return err
	}
	key, err := aero.NewKey(s.Config.Namespace, s.Config.Dataset, message.Date)
	if err != nil {
		return fmt.Errorf("unable to create key %v", err)
	}
	bin := aero.NewBin(fmt.Sprintf("event_type_%v", message.EventType), 1)
	_, err = client.Operate(s.writePolicy, key, aero.AddOp(bin))
	return err
}

func (s *service) client() (*aero.Client, error) {
	if s.Client != nil {
		return s.Client, nil
	}
	var err error
	s.Client, err = aero.NewClient(s.Config.IP, s.Config.Port)
	return s.Client, err
}

//New create a new service
func New(config *Config) Service {
	return &service{
		Config:      config,
		writePolicy: aero.NewWritePolicy(0, 0),
	}
}
