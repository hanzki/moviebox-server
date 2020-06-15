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
	link      string
	progress  float64
	location  string
	hash      string
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
