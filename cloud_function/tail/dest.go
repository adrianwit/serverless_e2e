package tail

import "github.com/viant/dsc"

type Destination struct {
	Table string
	*dsc.Config
}