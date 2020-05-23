package fixtures

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
)

// Marshall unmarshalls the payload and returns a *bytes.Buffer with the content,
// or fails the test if unable to
func Marshall(t *testing.T, payload interface{}) *bytes.Buffer {
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Got unexpected error while marshalling payload - %v", err)
	}
	return bytes.NewBuffer(body)
}

// Decode decodes the reader into into the passed value,
// or fails the test if unable to
func Decode(t *testing.T, r io.Reader, into interface{}) {
	if err := json.NewDecoder(r).Decode(&into); err != nil {
		t.Fatalf("Got unexpected error while decoding response - %v", err)
	}
}
