package storage

import (
	"fmt"
	"time"

	"github.com/beevik/guid"
	"github.com/hanzki/moviebox-server/core"
)

type downloadStorageMock struct {
	data map[core.DownloadID]core.Download
}

// NewMock constructs a new mock storage client
func NewDownloadStorageMock() *downloadStorageMock {
	m := new(downloadStorageMock)
	m.data = make(map[core.DownloadID]core.Download)
	return m
}

func (m *downloadStorageMock) Save(r *core.Download) (*core.Download, error) {
	r.ID = newDownloadID()
	r.CreatedAt = time.Now()
	m.data[r.ID] = *r
	return r, nil
}

func newDownloadID() core.DownloadID {
	return core.DownloadID(guid.NewString())
}

func (m *downloadStorageMock) Update(r *core.Download) error {
	m.data[r.ID] = *r
	return nil
}

func (m *downloadStorageMock) Load(id core.DownloadID) (*core.Download, error) {
	if r, ok := m.data[id]; ok {
		return &r, nil
	}
	return nil, fmt.Errorf("MockStorage.Load: No Download found with ID=%s", id)
}

func (m *downloadStorageMock) Delete(id core.DownloadID) error {
	delete(m.data, id)
	return nil
}
