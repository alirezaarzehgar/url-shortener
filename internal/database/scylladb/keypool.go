package scylladb

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"math/bits"

	"github.com/alirezaarzehgar/url-shortener/internal/database"
	"github.com/gocql/gocql"
)

func allocateRange(session *gocql.Session, id int, keyRange uint64) ([]string, error) {
	var n uint64
	err := session.Query("SELECT counter FROM urlshortener.keypool WHERE key = ?", id).Scan(&n)
	if err != nil {
		return nil, fmt.Errorf("failed to get counter value: %w", err)
	}

	applied, err := session.Query(
		`UPDATE urlshortener.keypool SET counter = ? WHERE key = ? IF counter = ?`, n+keyRange, id, n,
	).ScanCAS(&n)
	if err != nil {
		return nil, fmt.Errorf("failed to update counter: %w", err)
	}

	if !applied {
		return nil, errors.New("do not applied")
	}

	r := []string{}
	for i := n; i < n+keyRange; i++ {
		r = append(r, GenerateID(i))
	}
	return r, nil
}

func GenerateID(counter uint64) string {
	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, counter)
	nbits := (bits.Len64(counter) + 7) >> 3
	nbits = max(nbits, 4)
	return base64.RawURLEncoding.EncodeToString(key[:nbits])
}

type KeyPoolImpl struct {
	ScyllaDBConnection
	clusterID int
	keyRange  uint64
}

func (kp KeyPoolImpl) ReserveRange() (database.KeyRange, error) {
	return allocateRange(kp.session, kp.clusterID, kp.keyRange)
}

func NewKeyPool(conn ScyllaDBConnection, conf Config) database.KeyPool {
	return KeyPoolImpl{
		keyRange:           conf.keyRange,
		clusterID:          conf.clusterID,
		ScyllaDBConnection: conn,
	}
}
