package GobDB

import (
	"os"
	"fmt"
	"sync"
	"encoding/gob"
)


// Provides a goroutine-safe method of decoding gobbed objects
// individually. Assumes that gob type definitions are registered
// in the necessary order.
//
// NOTE: Decoder pointers to the same address will be forced to use mutex
// locks when utilized concurrently. On the other hand, if you copy a decoder
// itself, the two copies may be interacted with concurrently and lock-free.
type Decoder struct {
	decoder *gob.Decoder
	reader  hookedReader
	mutex   sync.Mutex
}


// Attempts to decode the given data, placing results into <address>
// if no errors occur. New types are registered implicitly, but they
// must be used in correct order.
func (d *Decoder) Decode(data []byte, address interface{}) error {
	d.ensureInitialized()
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.reader.buffer(data)
	return d.decoder.Decode(address)
}


// Informs internal decoder of the given gob type definition and its
// example object.
func (d *Decoder) Register(data []byte) error {
	d.ensureInitialized()
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.reader.buffer(data)
	return d.decoder.Decode(nil)
}


// Makes sure that a gob decoder has been set up with the correct
// reader.
func (d *Decoder) ensureInitialized() {
	reader := makeAtomicReader()
	if d.decoder == nil {
		d.reader = reader
		d.decoder = gob.NewDecoder(reader)
	}
}


// Helper to access the gob decoders internal reader.
type hookedReader struct {
	data *[]byte
}


// Allocates memory for shared slice.
func makeAtomicReader() hookedReader {
	return hookedReader {data: &[]byte{}}
}


// Implementing the reader interface for gob.decoder...
func (r hookedReader) Read(data []byte) (int, error) {
	if r.data == nil {
		fmt.Println(" - ERROR: called decode before buffer.")
		os.Exit(1)
	}


	n := copy(data, *r.data)
	if n == len(*r.data) {
		*r.data = []byte{}
	} else {
		*r.data = (*r.data)[n:]
	}
	return n, nil
}


// Replaces data in buffer. This should replace data in the gob
// decoders buffer too.
func (r *hookedReader) buffer(data []byte) {
	*r.data = data
}
