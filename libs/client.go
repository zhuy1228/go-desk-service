package libs

import (
	"sync"

	"github.com/libp2p/go-libp2p/core/host"
)

type Client struct {
	This host.Host
}

var once sync.Once

var clientObj *Client

func InitClient() *Client {
	once.Do(func() {

		clientObj = &Client{}
	})

	return clientObj
}
