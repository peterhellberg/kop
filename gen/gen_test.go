package gen_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/peterhellberg/kop/gen"
)

type Store interface {
	Put(PutRequest) PutResponse
}

type PutRequest struct {
	Thing  *Thing
	ByName map[string][]Thing `json:"by_name,omitempty"`
}

type PutResponse struct{}

type Thing struct {
	Name       string
	Children   []Thing
	ID         string
	HTTPServer string
	UserID     string
}

func TestRPC(t *testing.T) {
	var buf bytes.Buffer

	if err := gen.RPC(&buf, "store", reflect.TypeFor[Store]()); err != nil {
		t.Fatal(err)
	}

	src := buf.Bytes()

	for _, want := range []string{
		"package store\n",
		"Put(context.Context, PutRequest) (*PutResponse, error)",
		"func RegisterStore(server *Server, store Store)",
		"Thing  *Thing             `json:\"thing\"`",
		"ByName map[string][]Thing `json:\"by_name,omitempty\"`",
		"type PutResponse struct {\n\tError string `json:\"error,omitempty\"`\n}",
		"Name       string  `json:\"name\"`",
		"Children   []Thing `json:\"children\"`",
		"ID         string  `json:\"id\"`",
		"HTTPServer string  `json:\"httpServer\"`",
		"UserID     string  `json:\"userID\"`",
	} {
		if !bytes.Contains(src, []byte(want)) {
			t.Errorf("generated code lacks %q\n%s", want, src)
		}
	}
}
