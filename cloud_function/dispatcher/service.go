package dispatcher

import (
	"fmt"
	"github.com/viant/toolbox"
	"golang.org/x/net/context"
	"strings"
)

type Service interface {
	Handle(ctx context.Context, event interface{}) error
}

type service struct {
	config   *Config
	provider func(filter interface{}) (toolbox.Predicate, error)
}

func (s *service) Handle(ctx context.Context, event interface{}) error {
	for _, route := range s.config.Routes {
		if route.predicate.Apply(event) {
			if err := s.notify(ctx, route, event); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *service) notify(ctx context.Context, route *Route, event interface{}) error {
	switch strings.ToLower(route.Target.Type) {
	case "topic":
		return Publish(ctx, route.Target, event)
	case "function":
		return Call(ctx, route.Target, event)
	default:
		return fmt.Errorf("unsupporter target type:%v", route.Target.Type)
	}
}

//New creates a new dispatcher service
func New(config *Config, provider func(filter interface{}) (toolbox.Predicate, error)) (Service, error) {
	result := &service{
		config:   config,
		provider: provider,
	}
	for _, route := range result.config.Routes {
		filterValue, err := route.Filter.Value()
		if err != nil {
			return nil, err
		}
		if route.predicate, err = provider(filterValue); err != nil {
			return nil, err
		}
	}
	return result, nil
}
