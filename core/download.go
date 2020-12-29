package core

import (
	"fmt"
	"time"
)

type DownloadStatus int

const (
	Pending = iota
	Downloading
	Stopped
	Complete
	Deleted
)

func (s DownloadStatus) String() string {
	return [...]string{"Pending", "Downloading", "Stopped", "Complete", "Deleted"}[s]
}

type DownloadID string

type Download struct {
	ID        DownloadID
	CreatedAt time.Time
	SearchID  SearchID
	Status    DownloadStatus
	Link      string
	Progress  float64
	Location  string
	Hash      string
}

// DownloadClient provides methods for downloading torrents and querying the download progress
type DownloadClient interface {
	StartDownload(searchResult *SearchResult) (*Download, error)
	Progress(download *Download) (*Download, error)
	StopDownload(Download *Download) (*Download, error)
}

// DownloadStorage allows storing and retrieving Download state from storage
type DownloadStorage interface {
	Save(r *Download) (*Download, error)
	Update(r *Download) error
	Load(id DownloadID) (*Download, error)
	Delete(id DownloadID) error
}

// DownloadController orchestrates starting a download
type DownloadController struct {
	Storage DownloadStorage
	Client  DownloadClient
}

// StartNewDownload starts downloading the movie or TV series from the searchResult
func (dc *DownloadController) StartNewDownload(searchResult *SearchResult) (*Download, error) {
	download, err := dc.Client.StartDownload(searchResult)
	if err != nil {
		return nil, fmt.Errorf("DownloadController.StartNewDownload starting download failed, %v", err)
	}
	download, err = dc.Storage.Save(download)
	if err != nil {
		return nil, fmt.Errorf("DownloadController.StartNewDownload failed to save download to storage, %v", err)
	}
	return download, nil
}

// GetProgress returns up to date information on an ongoing download
func (dc *DownloadController) GetProgress(id DownloadID) (*Download, error) {
	download, err := dc.Storage.Load(id)
	if err != nil {
		return nil, fmt.Errorf("DownloadController.GetProgress didn't find download \"%s\", %v", id, err)
	}
	download, err = dc.Client.Progress(download)
	if err != nil {
		return nil, fmt.Errorf("DownloadController.GetProgress failed to query progress for download \"%s\", %v", id, err)
	}
	err = dc.Storage.Update(download)
	if err != nil {
		return nil, fmt.Errorf("DownloadController.GetProgress failed to update download to storage, %v", err)
	}
	return download, nil
}
