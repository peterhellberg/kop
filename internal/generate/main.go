// Command generate writes rpc/gen.go from the definitions package.
// Run it from the repository root via go generate.
package main

import (
	"bytes"
	"log"
	"os"
	"reflect"

	"github.com/peterhellberg/kop/definitions"
	"github.com/peterhellberg/kop/gen"
)

func main() {
	var buf bytes.Buffer

	if err := generate(&buf); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("rpc/gen.go", buf.Bytes(), 0o644); err != nil {
		log.Fatal(err)
	}
}

func generate(buf *bytes.Buffer) error {
	return gen.RPC(buf, "rpc", reflect.TypeFor[definitions.List]())
}
