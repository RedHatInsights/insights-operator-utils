package helpers

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"testing"
)

// MustGobSerialize serializes an object using gob or panics
func MustGobSerialize(t testing.TB, obj interface{}) []byte {
	buf := new(bytes.Buffer)

	err := gob.NewEncoder(buf).Encode(obj)
	FailOnError(t, err)

	res, err := ioutil.ReadAll(buf)
	FailOnError(t, err)

	return res
}
