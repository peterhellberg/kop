// Command generate writes gen.go in the current directory from the defs
// fixture. It is run by go generate in package rpctest.
package main

import (
	"bytes"
	"log"
	"os"
	"reflect"

	"github.com/peterhellberg/kop/gen"
	"github.com/peterhellberg/kop/gen/internal/rpctest/defs"
)

func main() {
	var buf bytes.Buffer

	if err := gen.RPC(&buf, "rpctest", reflect.TypeFor[defs.List]()); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("gen.go", buf.Bytes(), 0o644); err != nil {
		log.Fatal(err)
	}
}
