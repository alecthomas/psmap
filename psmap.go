package psmap

import (
	"bytes"
	"encoding/binary"
	mmap "github.com/edsrzf/mmap-go"
	"hash/fnv"
	"os"
)

func hash(key []byte) uint32 {
	h := fnv.New32()
	h.Write(key)
	return h.Sum32()
}

const (
	UINT32_SIZE = 4
)

type KeyValue struct {
	Key   []byte
	Value []byte
}

type PersistentStaticMap struct {
	f     *os.File
	m     mmap.MMap
	index map[uint32][]*KeyValue
}

func Open(filename string) (*PersistentStaticMap, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	m, err := mmap.Map(f, 0, mmap.RDONLY)
	if err != nil {
		f.Close()
		return nil, err
	}

	// Build index
	index := make(map[uint32][]*KeyValue)
	i := uint32(0)
	size := uint32(len(m))
	for i < size {
		ks := binary.BigEndian.Uint32(m[i : i+4])
		vs := binary.BigEndian.Uint32(m[i+4 : i+8])
		key := m[i+8 : i+8+ks]
		value := m[i+8+ks : i+8+ks+vs]
		h := hash(key)
		index[h] = append(index[h], &KeyValue{Key: key, Value: value})
		i += 8 + ks + vs
	}

	return &PersistentStaticMap{
		f:     f,
		m:     m,
		index: index,
	}, nil
}

func (p *PersistentStaticMap) Get(key []byte) []byte {
	h := hash(key)
	for _, kv := range p.index[h] {
		if bytes.Compare(key, kv.Key) == 0 {
			return kv.Value
		}
	}
	return nil
}

func (p *PersistentStaticMap) Iterate() chan *KeyValue {
	out := make(chan *KeyValue)
	go func() {
		for _, kvs := range p.index {
			for _, kv := range kvs {
				out <- kv
			}
		}
		close(out)
	}()
	return out
}
