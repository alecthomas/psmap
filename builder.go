package psmap

import (
	"encoding/binary"
	"io"
)

type Builder struct {
	w  io.WriteCloser
	sb []byte
}

func NewBuilder(w io.WriteCloser) *Builder {
	return &Builder{w: w, sb: make([]byte, 4)}
}

func (b *Builder) Add(key []byte, value []byte) error {
	binary.BigEndian.PutUint32(b.sb, uint32(len(key)))
	if _, err := b.w.Write(b.sb); err != nil {
		return err
	}
	binary.BigEndian.PutUint32(b.sb, uint32(len(value)))
	if _, err := b.w.Write(b.sb); err != nil {
		return err
	}
	if _, err := b.w.Write(key); err != nil {
		return err
	}
	if _, err := b.w.Write(value); err != nil {
		return err
	}
	return nil
}

func (b *Builder) AddMap(m map[string][]byte) error {
	for k, v := range m {
		err := b.Add([]byte(k), v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Builder) Close() error {
	return b.w.Close()
}
