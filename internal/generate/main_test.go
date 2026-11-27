package main

import (
	"bytes"
	"os"
	"testing"
)

func TestGeneratedUpToDate(t *testing.T) {
	var buf bytes.Buffer

	if err := generate(&buf); err != nil {
		t.Fatal(err)
	}

	want, err := os.ReadFile("../../rpc/gen.go")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(buf.Bytes(), want) {
		t.Error("rpc/gen.go is out of date; run go generate ./...")
	}
}
