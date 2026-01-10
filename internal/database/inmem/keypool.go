package inmem

import (
	"encoding/base64"
	"encoding/binary"
	"math/bits"
	"sync"

	"github.com/alirezaarzehgar/url-shortener/internal/database"
)

type KeyPoolImpl struct {
	counter  *uint64
	keyRange uint64
	mutex    *sync.RWMutex
}

func GenerateID(counter uint64) string {
	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, counter)
	nbits := (bits.Len64(counter) + 7) >> 3
	nbits = max(nbits, 4)
	return base64.RawURLEncoding.EncodeToString(key[:nbits])
}

func (kp KeyPoolImpl) ReserveRange() (database.KeyRange, error) {
	kp.mutex.Lock()
	start := *kp.counter
	end := *kp.counter + kp.keyRange
	*kp.counter += kp.keyRange
	kp.mutex.Unlock()

	listOfKeys := []string{}
	for i := start; i < end; i++ {
		listOfKeys = append(listOfKeys, GenerateID(i))
	}
	return listOfKeys, nil
}

func NewKeyPool() database.KeyPool {
	counter := uint64(0)
	return KeyPoolImpl{
		mutex:    &sync.RWMutex{},
		counter:  &counter,
		keyRange: 10,
	}
}
