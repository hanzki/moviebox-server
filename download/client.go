package download

import (
	"github.com/hanzki/moviebox-server/core"
	"github.com/hekmon/transmissionrpc"
)

type TransmissionClient struct {
	rpc *transmissionrpc.Client
}

func NewTransmissionClient() *TransmissionClient {
	transmissionbt, err := transmissionrpc.New("127.0.0.1", "rpcuser", "rpcpass", nil)
	_ = transmissionbt
	_ = err
	return nil
}

func (tc *TransmissionClient) StartDownload(searchResult *core.SearchResult) (*core.Download, error) {
	return nil, nil
}

func (tc *TransmissionClient) Progress(download *core.Download) (*core.Download, error) {
	return nil, nil
}

func (tc *TransmissionClient) StopDownload(download *core.Download) (*core.Download, error) {
	return nil, nil
}
