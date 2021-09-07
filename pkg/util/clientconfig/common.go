package clientconfig

import (
	"crypto/tls"
)

const (
	NifiDefaultTimeout = int64(5)
)

type Manager interface {
	BuildConfig() (*NifiConfig, error)
	BuildConnect() (ClusterConnect, error)
	IsExternal() bool
}

type ClusterConnect interface {
	//NodeConnection(log logr.Logger, client client.Client) (node nificlient.NifiClient, err error)
	IsInternal() bool
	IsExternal() bool
	ClusterLabelString() string
	IsReady() bool
	Id() string
}

// NifiConfig are the options to creating a new ClusterAdmin client
type NifiConfig struct {
	NodeURITemplate string
	NodesURI        map[int32]NodeUri
	NifiURI         string
	UseSSL          bool
	TLSConfig       *tls.Config

	OperationTimeout   int64
	RootProcessGroupId string
}

type NodeUri struct {
	HostListener string
	RequestHost  string
}
