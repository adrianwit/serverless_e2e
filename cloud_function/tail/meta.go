package tail

import "github.com/viant/toolbox/url"

type Meta struct {
	ResourceURL string
	RecordCount int
	Transferred int
	Error       string
}

func NewMetaFromURL(URL string) (*Meta, error) {
	meta := &Meta{}
	resource := url.NewResource(URL)
	return meta, resource.Decode(meta)
}
