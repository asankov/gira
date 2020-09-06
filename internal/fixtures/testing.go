package fixtures

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
)

// MarshalBytes unmarshals the payload and returns a []byte with the content,
// or fails the test if unable to
func MarshalBytes(t *testing.T, payload interface{}) []byte {
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Got unexpected error while marshalling payload - %v", err)
	}
	return body
}

// Marshal unmarshalls the payload and returns a *bytes.Buffer with the content,
// or fails the test if unable to
func Marshal(t *testing.T, payload interface{}) *bytes.Buffer {
	return bytes.NewBuffer(MarshalBytes(t, payload))
}

// Decode decodes the reader into into the passed value,
// or fails the test if unable to
func Decode(t *testing.T, r io.Reader, into interface{}) {
	if err := json.NewDecoder(r).Decode(&into); err != nil {
		t.Fatalf("Got unexpected error while decoding response - %v", err)
	}
}
