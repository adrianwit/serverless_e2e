package bq

import (
	"fmt"
	"github.com/viant/toolbox"
	"regexp"
)

type predicate struct {
	Filter
}

func (p *predicate) Apply(source interface{}) bool {
	event, ok := source.(*Event)
	if !ok {
		return false
	}
	if p.Filter.Type != "" {
		if !matches(p.Filter.Type, event.Type()) {
			return false
		}
	}
	if p.Filter.Source != "" {
		if !matches(p.Filter.Source, event.Source()) {
			return false
		}
	}
	if p.Filter.Destination != "" {
		if !matches(p.Filter.Destination, event.Destination()) {
			return false
		}
	}
	if p.Filter.Status != "" {
		if !matches(p.Filter.Status, event.Status.State) {
			return false
		}
	}

	return true
}

func matches(pattern, target string) bool {
	matched, err := regexp.MatchString(pattern, target)
	if err != nil {
		return false
	}
	return matched
}

//NewPredicate creates a new predicate for supplied fileter
func NewPredicate(filter interface{}) (toolbox.Predicate, error) {
	if filter == nil {
		return nil, fmt.Errorf("filter was nil")
	}
	result := &predicate{}
	err := toolbox.DefaultConverter.AssignConverted(&result.Filter, filter)
	return result, err
}
