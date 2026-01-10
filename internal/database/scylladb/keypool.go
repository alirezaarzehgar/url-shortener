package scylladb

import "github.com/alirezaarzehgar/url-shortener/internal/database"

type KeyPoolImpl struct {
	ScyllaDBConnection
}

func (kp KeyPoolImpl) ReserveRange() (database.KeyRange, error) {
	panic("not implemented")
}

func NewKeyPool(conn ScyllaDBConnection) database.KeyPool {
	return KeyPoolImpl{
		ScyllaDBConnection: conn,
	}
}
