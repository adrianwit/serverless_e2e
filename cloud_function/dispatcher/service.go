package dispatcher

import (
	"fmt"
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/data"
	"golang.org/x/net/context"
	"log"
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

func (s *service) transform(event interface{}, route *Route) interface{} {
	if len(route.Fields) == 0 {
		return event
	}
	var eventMap = make(map[string]interface{})
	jsonConverter := toolbox.NewConverter("", "json")
	err := jsonConverter.AssignConverted(&eventMap, event)
	if err != nil {
		log.Printf("unable convert event: %v\n", err)
		return event
	}
	eventDataMap := data.Map(eventMap)
	result := data.NewMap()
	for sourcePath, taragetPath := range route.Fields {
		val, ok := eventDataMap.GetValue(sourcePath)
		if !ok {
			continue
		}
		result.SetValue(taragetPath, val)
	}
	return result
}

func (s *service) notify(ctx context.Context, route *Route, event interface{}) error {
	transformed := s.transform(event, route)
	switch strings.ToLower(route.Target.Type) {
	case "topic":
		return Publish(ctx, route.Target, transformed)
	case "function":
		return Call(ctx, route.Target, transformed)
	case "slack":
		return Post(ctx, route.Target, transformed, fmt.Sprintf("%T", event))
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
