package tail

import (
	"fmt"
	"github.com/viant/dsc"
	"github.com/viant/toolbox/data"
	"sync"
	"sync/atomic"
)

type Session struct {
	dsc.Manager
	*sync.WaitGroup
	resourceURL string
	Route       *Route
	table       *dsc.TableDescriptor
	dmlProvider dsc.DmlProvider
	limiter     chan bool
	count       int64
	transferred int64
	err         error
	done        int32
}

func (s *Session) SetError(err error) {
	if err == nil {
		return
	}
	s.err = err
	atomic.StoreInt32(&s.done, 1)
	close(s.limiter)
}

func (s *Session) acquire() {
	s.Add(1)
	s.limiter <- true
}

func (s *Session) release() {
	s.Done()
	<-s.limiter
}

func (s *Session) getDMLProvider(slice *data.CompactedSlice) (dsc.DmlProvider, error) {
	if s.dmlProvider != nil {
		return s.dmlProvider, nil
	}
	table, err := s.getTable(slice)
	if err != nil {
		return nil, err
	}
	s.dmlProvider = dsc.NewMapDmlProvider(table)
	return s.dmlProvider, nil
}

func (s *Session) getTable(slice *data.CompactedSlice) (*dsc.TableDescriptor, error) {
	if s.table != nil {
		return s.table, nil
	}
	table := s.TableDescriptorRegistry().Get(s.Route.Dest.Table)
	if table == nil {
		return nil, fmt.Errorf("target table %v not found", s.Route.Dest.Table)
	}
	if len(table.Columns) == 0 {
		table.Columns = []string{}
		for _, field := range slice.Fields() {
			table.Columns = append(table.Columns, field.Name)
		}
	}
	s.table = table
	return table, nil
}

func NewSession(config *Config, route *Route, resourceURL string) (*Session, error) {
	manager, err := dsc.NewManagerFactory().Create(route.Dest.Config)
	if err != nil {
		return nil, err
	}
	threads := config.Threads
	if threads == 0 {
		threads = 1
	}
	return &Session{
		Manager:     manager,
		Route:       route,
		resourceURL: resourceURL,
		WaitGroup:   &sync.WaitGroup{},
		limiter:     make(chan bool, threads),
	}, nil
}
