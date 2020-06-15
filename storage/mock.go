package storage

import (
	"fmt"
	"time"

	"github.com/beevik/guid"
	"github.com/hanzki/moviebox-server/core"
)

type mock struct {
	data map[core.SearchID]core.SearchResult
}

// NewMock constructs a new mock storage client
func NewMock() *mock {
	m := new(mock)
	m.data = make(map[core.SearchID]core.SearchResult)
	return m
}

func (m *mock) Save(r *core.SearchResult) (*core.SearchResult, error) {
	r.ID = newID()
	r.CreatedAt = time.Now()
	m.data[r.ID] = *r
	return r, nil
}

func newID() core.SearchID {
	return core.SearchID(guid.NewString())
}

func (m *mock) Update(r *core.SearchResult) error {
	m.data[r.ID] = *r
	return nil
}

func (m *mock) Load(id core.SearchID) (*core.SearchResult, error) {
	if r, ok := m.data[id]; ok {
		return &r, nil
	}
	return nil, fmt.Errorf("MockStorage.Load: No SearchResult found with ID=%s", id)
}

func (m *mock) Delete(id core.SearchID) error {
	delete(m.data, id)
	return nil
}
