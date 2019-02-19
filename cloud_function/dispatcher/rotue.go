package dispatcher

import "github.com/viant/toolbox"

//Route represents matching route
type Route struct {
	Filter    toolbox.AnyJSONType `json:"filter"`
	Target    *Target             `json:"target"`
	Fields    map[string]string   `json:"fields"`
	predicate toolbox.Predicate
}
