package download

import (
	"fmt"
	"log"
	"strconv"

	"github.com/davecgh/go-spew/spew"
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
	RPCPort     string
	RPCUser     string
	RPCPassword string
}

func NewTransmissionClient(config *Config) *TransmissionClient {
	/*
		if config.DownloadDir == "" {
			config.DownloadDir = "/var/moviebox/downloads"
		}
	*/
	var ac *transmissionrpc.AdvancedConfig
	if config.RPCPort != "" {
		port, err := strconv.Atoi(config.RPCPort)
		if err != nil {
			log.Fatalf("Failed to parse transmission port %s: %v", config.RPCPort, err)
		}
		ac = &transmissionrpc.AdvancedConfig{Port: uint16(port)}
	}

	transmissionbt, err := transmissionrpc.New(config.RPCHost, config.RPCUser, config.RPCPassword, ac)
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
	if torrent.PercentDone == nil {
		spew.Dump(torrent)
		log.Fatalln("No PercentDone received!")
	}
	download.Progress = *torrent.PercentDone

	if *torrent.IsFinished || *torrent.PercentDone == 1 {
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
