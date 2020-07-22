package download

import (
	"fmt"
	"log"

	"github.com/hanzki/moviebox-server/core"
	"github.com/hekmon/transmissionrpc"
)

type TransmissionClient struct {
	rpc    *transmissionrpc.Client
	config *Config
}

type Config struct {
	//DownloadDir string
	RPCHost     string
	RPCUser     string
	RPCPassword string
}

func NewTransmissionClient(config *Config) *TransmissionClient {
	/*
		if config.DownloadDir == "" {
			config.DownloadDir = "/var/moviebox/downloads"
		}
	*/
	transmissionbt, err := transmissionrpc.New(config.RPCHost, config.RPCUser, config.RPCPassword, nil)
	if err != nil {
		log.Fatal(err)
	}
	/*
		err = os.MkdirAll(config.DownloadDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	*/
	return &TransmissionClient{transmissionbt, config}
}

func (tc *TransmissionClient) StartDownload(searchResult *core.SearchResult) (*core.Download, error) {
	payload := &transmissionrpc.TorrentAddPayload{Filename: &searchResult.Link}
	torrent, err := tc.rpc.TorrentAdd(payload)
	if err != nil {
		return nil, err
	}
	download := &core.Download{
		SearchID: searchResult.ID,
		Status:   core.Downloading,
		Link:     searchResult.Link,
		Location: "",
		Hash:     *torrent.HashString,
	}
	return download, nil
}

func (tc *TransmissionClient) Progress(download *core.Download) (*core.Download, error) {
	torrents, err := tc.rpc.TorrentGetAllForHashes([]string{download.Hash})
	if err != nil {
		return nil, err
	}
	if len(torrents) < 1 {
		return nil, fmt.Errorf("TransmissionClient.Progress: Didn't find torrent for hash %s", download.Hash)
	}
	torrent := torrents[0]
	download.Progress = *torrent.PercentDone
	if *torrent.IsFinished {
		download.Status = core.Complete
	}
	return download, nil
}

func (tc *TransmissionClient) StopDownload(download *core.Download) (*core.Download, error) {
	err := tc.rpc.TorrentStopHashes([]string{download.Hash})
	if err != nil {
		return nil, err
	}
	download.Status = core.Stopped
	return download, nil
}
