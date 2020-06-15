package storage

import (
	"testing"

	"github.com/hanzki/moviebox-server/core"
)

func TestSave(t *testing.T) {
	mockStorage := NewMock()

	r := &core.SearchResult{
		Title: "Testinator II",
	}

	r, _ = mockStorage.Save(r)

	if r.ID == "" {
		t.Errorf("Save didn't set new ID")
	}
	if r.CreatedAt.IsZero() {
		t.Errorf("Save didn't set CreatedAt")
	}
}

func TestLoad(t *testing.T) {
	mockStorage := NewMock()

	r := &core.SearchResult{
		Title: "Testinator II",
	}

	r, _ = mockStorage.Save(r)
	r2, _ := mockStorage.Load(r.ID)

	if r2.Title != "Testinator II" {
		t.Errorf("Load failed. Title = %s, expected = %s", r.Title, "Testinator II")
	}
}

func TestUpdate(t *testing.T) {
	mockStorage := NewMock()

	r := &core.SearchResult{
		Title: "Testinator II",
	}

	r, _ = mockStorage.Save(r)

	r.Title = "Mocky V"

	_ = mockStorage.Update(r)

	r, _ = mockStorage.Load(r.ID)

	if r.Title != "Mocky V" {
		t.Errorf("Update failed. Title = %s, expected = %s", r.Title, "Mocky V")
	}
}
