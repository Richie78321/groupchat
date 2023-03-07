package replication

import (
	"time"
)

const (
	connectionTimeout    time.Duration = 10 * time.Second
	connectionRetryDelay time.Duration = 10 * time.Second
)

type ReplicationPeer interface {
}

type peer struct {
	id     string
	addr   string
	events chan struct{}
}

func NewPeer(id string, addr string) ReplicationPeer {
	return peer{
		id:     id,
		addr:   addr,
		events: make(chan struct{}),
	}
}
