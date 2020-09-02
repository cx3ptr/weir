package driver

import (
	"context"
	"crypto/tls"

	"github.com/pingcap-incubator/weir/pkg/proxy"
	"github.com/pingcap-incubator/weir/pkg/proxy/backend/client"
)

type Backend interface {
	GetConn(context.Context) (*client.Conn, error)
	PutConn(context.Context, *client.Conn) error
}

type DriverImpl struct {
	backend Backend
}

func NewDriverImpl(backend Backend) *DriverImpl {
	return &DriverImpl{
		backend: backend,
	}
}

func (d *DriverImpl) OpenCtx(connID uint64, capability uint32, collation uint8, dbname string, tlsState *tls.ConnectionState) (proxy.QueryCtx, error) {
	return NewQueryCtxImpl(d.backend), nil
}