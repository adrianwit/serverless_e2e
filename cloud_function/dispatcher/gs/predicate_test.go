package gs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPredicate_Apply(t *testing.T) {

	var useCases = []struct {
		description string
		filter      interface{}
		event       interface{}
		hasError    bool
		expect      bool
	}{

		{
			description: "prefix match",
			filter: &Filter{
				Prefix: []string{
					"/folder1",
				},
			},
			event: &Event{
				Name: "/folder1/asset1.txt",
			},
			expect: true,
		},
		{
			description: "prefix does not match",
			filter: &Filter{
				Prefix: []string{
					"/folder2",
				},
			},
			event: &Event{
				Name: "/folder1/asset1.txt",
			},
			expect: false,
		},
		{
			description: "suffix match",
			filter: &Filter{
				Suffix: []string{
					".txt",
				},
			},
			event: &Event{
				Name: "/folder1/asset1.txt",
			},
			expect: true,
		},
		{
			description: "suffix does match",
			filter: &Filter{
				Suffix: []string{
					".txt",
				},
			},
			event: &Event{
				Name: "/folder1/asset1.cfg",
			},
			expect: false,
		},
	}

	for _, useCase := range useCases {
		predicate, err := NewPredicate(useCase.filter)
		if useCase.hasError {
			assert.NotNil(t, err, useCase.description)
			continue
		}
		if !assert.Nil(t, err, useCase.description) {
			continue
		}
		actual := predicate.Apply(useCase.event)
		assert.Equal(t, useCase.expect, actual, useCase.description)
	}
}
