package gomap

import (
	"io"
	"unsafe"
)

type KV struct {
	key any
	val any
}

type DataKv []KV

func (d *DataKv) Set(key, val any) {
	if val == nil {
		return
	}

	if b, ok := key.([]byte); ok {
		key = bytesToString(b)
	}
	args := *d
	n := len(args)
	for i := range n {
		kv := &args[i]
		if kv.key == key {
			kv.val = val
			return
		}
	}

	c := cap(args)
	if c > n {
		args = args[:n+1]
		kv := &args[n]
		kv.key = key
		kv.val = val
		*d = args
		return
	}

	kv := KV{}
	kv.key = key
	kv.val = val
	args = append(args, kv)
	*d = args
}

func (d *DataKv) Get(key any) any {
	if b, ok := key.([]byte); ok {
		key = bytesToString(b)
	}
	args := *d
	n := len(args)
	for i := range n {
		kv := &args[i]
		if kv.key == key {
			return kv.val
		}
	}
	return nil
}

func (d *DataKv) Reset() {
	args := *d
	n := len(args)
	for i := range n {
		v := args[i].val
		if vc, ok := v.(io.Closer); ok {
			vc.Close()
		}
		(*d)[i].val = nil
		(*d)[i].key = nil
	}
	*d = (*d)[:0]
}

func (d *DataKv) Remove(key any) {
	if b, ok := key.([]byte); ok {
		key = bytesToString(b)
	}
	args := *d
	n := len(args)
	for i := 0; i < n; i++ {
		kv := &args[i]
		if kv.key == key {
			n--
			args[i], args[n] = args[n], args[i]
			args[n].key = nil
			args[n].val = nil
			args = args[:n]
			*d = args
			return
		}
	}
}

func (d *DataKv) Peek(fn func(k, v any) bool) {
	data := *d
	for i := range data {
		k := data[i].key
		v := data[i].val
		if fn(k, v) {
			return
		}
	}
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
