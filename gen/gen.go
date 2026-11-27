// Package gen generates JSON-over-HTTP RPC servers and clients from Go
// interfaces whose methods have the form Method(Request) Response.
package gen

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"go/format"
	"io"
	"maps"
	"path"
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
// client for each of the services, given as interface types from a single
// package, and every struct they reach in that package. Types from other
// packages are imported.
func RPC(w io.Writer, pkg string, services ...reflect.Type) error {
	g := generator{
		objects: map[string]*object{},
		imports: map[string]string{}, // package name -> path
	}

	for _, p := range []string{"bytes", "context", "encoding/json", "errors", "fmt", "io", "net/http", "time"} {
		g.imports[path.Base(p)] = p
	}

	data := struct {
		Package  string
		Imports  map[string]string // path -> package name
		Services []service
		Objects  []*object
	}{Package: pkg, Imports: map[string]string{}}

	for _, s := range services {
		svc, err := g.service(s)
		if err != nil {
			return err
		}

		data.Services = append(data.Services, svc)
	}

	for name, p := range g.imports {
		data.Imports[p] = name
	}

	for _, name := range slices.Sorted(maps.Keys(g.objects)) {
		data.Objects = append(data.Objects, g.objects[name])
	}

	var buf bytes.Buffer

	tmpl := template.Must(template.New("").Funcs(template.FuncMap{"lowerFirst": lowerFirst, "base": path.Base}).Parse(rpcTemplate))

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

type generator struct {
	pkgPath string
	objects map[string]*object
	imports map[string]string
}

func (g *generator) service(s reflect.Type) (service, error) {
	if s.Kind() != reflect.Interface || s.Name() == "" {
		return service{}, fmt.Errorf("%s: service must be a named interface", s)
	}

	if g.pkgPath == "" {
		g.pkgPath = s.PkgPath()
	} else if s.PkgPath() != g.pkgPath {
		return service{}, fmt.Errorf("%s: services must be in the same package as %s", s, g.pkgPath)
	}

	svc := service{Name: s.Name()}

	for m := range s.Methods() {
		in, out, err := g.signature(m.Type)
		if err != nil {
			return svc, fmt.Errorf("%s.%s: %w", s, m.Name, err)
		}

		g.objects[out].Output = true

		svc.Methods = append(svc.Methods, method{m.Name, in, out})
	}

	return svc, nil
}

func (g *generator) signature(m reflect.Type) (in, out string, err error) {
	if m.NumIn() != 1 || m.NumOut() != 1 {
		return "", "", errors.New("method must have the form Method(Request) Response")
	}

	for _, t := range []reflect.Type{m.In(0), m.Out(0)} {
		if t.Kind() != reflect.Struct || t.PkgPath() != g.pkgPath {
			return "", "", fmt.Errorf("%s must be a struct declared in %s", t, g.pkgPath)
		}

		if err := g.object(t); err != nil {
			return "", "", err
		}
	}

	return m.In(0).Name(), m.Out(0).Name(), nil
}

// object adds the struct t, declared in the definitions package, along
// with every type its fields reach.
func (g *generator) object(t reflect.Type) error {
	if g.objects[t.Name()] != nil {
		return nil
	}

	o := &object{Name: t.Name()}
	g.objects[o.Name] = o

	for f := range t.Fields() {
		if !f.IsExported() || f.Anonymous {
			return fmt.Errorf("%s.%s: fields must be exported and not embedded", t, f.Name)
		}

		typ, err := g.typeName(f.Type)
		if err != nil {
			return fmt.Errorf("%s.%s: %w", t, f.Name, err)
		}

		json := f.Tag.Get("json")
		if json == "" {
			json = lowerFirst(f.Name)
		}

		o.Fields = append(o.Fields, field{f.Name, typ, json})
	}

	return nil
}

// typeName returns the Go source for t in the generated package, adding
// objects and imports as needed.
func (g *generator) typeName(t reflect.Type) (string, error) {
	switch k := t.Kind(); {
	case k == reflect.Func, k == reflect.Chan, k == reflect.UnsafePointer, k == reflect.Interface && t.NumMethod() > 0:
		return "", fmt.Errorf("%s: type not supported", t)
	case t.Name() != "" && t.PkgPath() == g.pkgPath:
		if t.Kind() != reflect.Struct {
			return "", fmt.Errorf("%s: only structs can be declared in %s", t, g.pkgPath)
		}

		return t.Name(), g.object(t)
	case t.Name() != "" && t.PkgPath() != "":
		name := strings.TrimSuffix(t.String(), "."+t.Name())

		if p, ok := g.imports[name]; ok && p != t.PkgPath() {
			return "", fmt.Errorf("%s: package name %s conflicts with %s", t, name, p)
		}

		g.imports[name] = t.PkgPath()

		return t.String(), nil
	case t.Name() != "":
		return t.Name(), nil
	}

	switch t.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Array:
		elem, err := g.typeName(t.Elem())

		return strings.TrimSuffix(t.String(), t.Elem().String()) + elem, err
	case reflect.Map:
		key, err := g.typeName(t.Key())
		if err != nil {
			return "", err
		}

		elem, err := g.typeName(t.Elem())

		return "map[" + key + "]" + elem, err
	case reflect.Interface:
		return "any", nil
	}

	return "", fmt.Errorf("%s: type not supported", t)
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
