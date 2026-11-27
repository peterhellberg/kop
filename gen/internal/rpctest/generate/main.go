// Command generate writes the rpctest package from the defs fixture.
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

	if err := os.WriteFile(os.Args[1], buf.Bytes(), 0o644); err != nil {
		log.Fatal(err)
	}
}
