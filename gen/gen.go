// Package gen generates JSON-over-HTTP RPC servers and clients from Go
// interfaces whose methods have the form Method(Request) Response.
package gen

import (
	"bytes"
	_ "embed"
	"go/format"
	"io"
	"maps"
	"reflect"
	"slices"
	"strings"
	"text/template"
	"unicode"
)

//go:embed rpc.go.tmpl
var rpcTemplate string

type service struct {
	Name    string
	Methods []method
}

type method struct{ Name, Input, Output string }

type object struct {
	Name   string
	Output bool
	Fields []field
}

type field struct{ Name, Type, JSON string }

// RPC writes Go source for package pkg to w, containing a server and a
// client for each of the services, given as interface types, and every
// struct they reach in the same package. Anything the generated code cannot
// express is left for the compiler to reject.
func RPC(w io.Writer, pkg string, services ...reflect.Type) error {
	objects := map[string]*object{}

	data := struct {
		Package  string
		Services []service
		Objects  []*object
	}{Package: pkg}

	for _, s := range services {
		prefix := strings.TrimSuffix(s.String(), s.Name())

		var add func(t reflect.Type) *object

		add = func(t reflect.Type) *object {
			for t.Kind() == reflect.Slice || t.Kind() == reflect.Pointer || t.Kind() == reflect.Map {
				t = t.Elem()
			}

			if t.Kind() != reflect.Struct || t.PkgPath() != s.PkgPath() || objects[t.Name()] != nil {
				return objects[t.Name()]
			}

			o := &object{Name: t.Name()}
			objects[o.Name] = o

			for f := range t.Fields() {
				json := f.Tag.Get("json")
				if json == "" {
					json = lowerFirst(f.Name)
				}

				o.Fields = append(o.Fields, field{f.Name, strings.ReplaceAll(f.Type.String(), prefix, ""), json})
				add(f.Type)
			}

			return o
		}

		svc := service{Name: s.Name()}

		for m := range s.Methods() {
			in, out := m.Type.In(0), m.Type.Out(0)

			add(in)

			if o := add(out); o != nil {
				o.Output = true
			}

			svc.Methods = append(svc.Methods, method{m.Name, in.Name(), out.Name()})
		}

		data.Services = append(data.Services, svc)
	}

	for _, name := range slices.Sorted(maps.Keys(objects)) {
		data.Objects = append(data.Objects, objects[name])
	}

	var buf bytes.Buffer

	tmpl := template.Must(template.New("").Funcs(template.FuncMap{"lowerFirst": lowerFirst}).Parse(rpcTemplate))

	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	_, err = w.Write(src)

	return err
}

// lowerFirst lowercases the leading capital or initialism of s,
// e.g. "Items" -> "items", "ID" -> "id", "HTTPServer" -> "httpServer".
func lowerFirst(s string) string {
	r := []rune(s)

	for i := range r {
		if !unicode.IsUpper(r[i]) || i > 0 && i+1 < len(r) && unicode.IsLower(r[i+1]) {
			break
		}

		r[i] = unicode.ToLower(r[i])
	}

	return string(r)
}
