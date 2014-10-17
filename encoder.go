package GobDB

import (
	"io"
	"bytes"
	"encoding/gob"
)


// Wraps gob.Encoder in such a way to retain all of gobs internal
// type definitions, and to return byte slices of type-value
// pairs without any additional data.
//
// Returns bytes of encoded objects, writes only type definitions
// and uniquely-typed objects (i.e. one object of each type).
type FilteredEncoder struct {
	essentials []interface{}
	encoder *gob.Encoder
	buffer *bytes.Buffer
	writer io.Writer
}


// Constructs filter with an empty buffer and unused encoder. Only
// data essential to the construction of new encodes are kept. This
// consists of internal gob type definitions and one object of each
// type.
func MakeFilteredEncoder(w io.Writer) FilteredEncoder {
	var r FilteredEncoder
	r.writer = w // where essential data is outputted
	r.buffer = bytes.NewBuffer([]byte{}) // buffer for temporary use
	r.encoder = gob.NewEncoder(r.buffer) // wrapped gob encoder
	return r
}






func (f *FilteredEncoder) Encode(e interface{}) ([]byte, error) {	
	// Empty write buffer.
	f.buffer.Reset()

	// Encode value and its size.
	err := f.encoder.Encode(e)
	if err != nil {
		return []byte{}, err
	}
	s1 := f.buffer.Len()

	// Repeat.
	err = f.encoder.Encode(e)
	if err != nil {
		return []byte{}, err
	}
	s2 := f.buffer.Len() - s1

	// Infer whether or not a new type was written. If block sizes
	// are equal, we know that the type being encoded has was 
	// encoded beforehand. Otherwise, we know that the parameter to
	// this function call is the first of its type to be encoded 
	// through this filter.
	r := make([]byte, s2)
	copy(r, f.buffer.Bytes()[s1:])
	if s1 != s2 {
		f.writer.Write(f.buffer.Bytes())
	}
	return r, nil
}

