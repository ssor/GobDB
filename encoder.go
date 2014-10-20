package GobDB

import (
	"bytes"
	"encoding/gob"
)

// FilteredEncoder wraps gob.Encoder in such a way to retain
// all of gobs internal  type definitions, and to return byte
// slices of type-value pairs without any additional data.
//
// Returns bytes of encoded objects, writes only type definitions
// and uniquely-typed objects (i.e. one object of each type).
type FilteredEncoder struct {
	encoder *gob.Encoder
	buffer  *bytes.Buffer
}

// Encode performs gob.Encode twice, compares the two encodings
// to  deduce whether or not the value had been encoded before.
// If it has the gob type definition is returned in position one.
// Regardless,  the actual encoded value (without typedef bytes)
// is returned in position two.
func (f *FilteredEncoder) Encode(e interface{}) ([]byte, []byte, error) {
	// Empty write buffer and ensure that an encoder is present.
	f.ready()

	// Encode value and its size.
	err := f.encoder.Encode(e)
	if err != nil {
		return []byte{}, []byte{}, err
	}
	s1 := f.buffer.Len()

	// Repeat.
	err = f.encoder.Encode(e)
	if err != nil {
		return []byte{}, []byte{}, err
	}
	s2 := f.buffer.Len() - s1

	// Infer whether or not a new type was written. If block sizes
	// are equal, we know that the type being encoded has was
	// encoded beforehand. Otherwise, we know that the parameter to
	// this function call is the first of its type to be encoded
	// through this filter.
	r := make([]byte, s1+s2)
	copy(r, f.buffer.Bytes())
	return r[:s1-s2], r[s1:], nil
}

// Initializes encoder if not yet initialized.
func (f *FilteredEncoder) ready() {
	if f.encoder == nil {
		f.buffer = bytes.NewBuffer([]byte{})
		f.encoder = gob.NewEncoder(f.buffer)
	} else {
		f.buffer.Reset()
	}
}
