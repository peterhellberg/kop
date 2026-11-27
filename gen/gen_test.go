package gen_test

import (
	htmltemplate "html/template"
	"io"
	"reflect"
	"strings"
	"testing"
	"text/template"

	"github.com/peterhellberg/kop/gen"
	"github.com/peterhellberg/kop/gen/internal/rpctest/defs"
)

type (
	Req  struct{}
	Resp struct{}
	Name string

	TwoParams    interface{ M(Req, Req) Resp }
	NoResult     interface{ M(Req) }
	PointerParam interface{ M(*Req) Resp }
	NamedParam   interface{ M(Name) Resp }
	ForeignParam interface{ M(defs.AddRequest) Resp }

	UnexportedReq struct{ a int }
	UnexportedSvc interface{ M(UnexportedReq) Resp }
	EmbeddedReq   struct{ Req }
	EmbeddedSvc   interface{ M(EmbeddedReq) Resp }
	FuncReq       struct{ F func() }
	FuncSvc       interface{ M(FuncReq) Resp }
	ChanReq       struct{ C []chan int }
	ChanSvc       interface{ M(ChanReq) Resp }
	AnonReq       struct{ S struct{ A int } }
	AnonSvc       interface{ M(AnonReq) Resp }
	IfaceReq      struct{ R io.Reader }
	IfaceSvc      interface{ M(IfaceReq) Resp }
	NamedReq      struct{ N map[string]Name }
	NamedSvc      interface{ M(NamedReq) Resp }
	ConflictReq   struct {
		T *template.Template
		H *htmltemplate.Template
	}
	ConflictSvc interface{ M(ConflictReq) Resp }
	Valid       interface{ M(Req) Resp }
)

func TestRPCErrors(t *testing.T) {
	for _, tt := range []struct {
		name     string
		services []reflect.Type
		want     string
	}{
		{"not an interface", []reflect.Type{reflect.TypeFor[Req]()}, "must be a named interface"},
		{"two params", []reflect.Type{reflect.TypeFor[TwoParams]()}, "TwoParams.M: method must have the form"},
		{"no result", []reflect.Type{reflect.TypeFor[NoResult]()}, "method must have the form"},
		{"pointer param", []reflect.Type{reflect.TypeFor[PointerParam]()}, "*gen_test.Req must be a struct"},
		{"named param", []reflect.Type{reflect.TypeFor[NamedParam]()}, "gen_test.Name must be a struct"},
		{"foreign param", []reflect.Type{reflect.TypeFor[ForeignParam]()}, "defs.AddRequest must be a struct"},
		{"unexported field", []reflect.Type{reflect.TypeFor[UnexportedSvc]()}, "UnexportedReq.a: fields must be exported"},
		{"embedded field", []reflect.Type{reflect.TypeFor[EmbeddedSvc]()}, "EmbeddedReq.Req: fields must be exported and not embedded"},
		{"func field", []reflect.Type{reflect.TypeFor[FuncSvc]()}, "FuncReq.F: func(): type not supported"},
		{"chan field", []reflect.Type{reflect.TypeFor[ChanSvc]()}, "ChanReq.C: chan int: type not supported"},
		{"anonymous struct", []reflect.Type{reflect.TypeFor[AnonSvc]()}, "AnonReq.S: struct { A int }: type not supported"},
		{"interface field", []reflect.Type{reflect.TypeFor[IfaceSvc]()}, "IfaceReq.R: io.Reader: type not supported"},
		{"named non-struct", []reflect.Type{reflect.TypeFor[NamedSvc]()}, "NamedReq.N: gen_test.Name: only structs can be declared"},
		{"import conflict", []reflect.Type{reflect.TypeFor[ConflictSvc]()}, "package name template conflicts with"},
		{"mixed packages", []reflect.Type{reflect.TypeFor[Valid](), reflect.TypeFor[defs.List]()}, "services must be in the same package"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := gen.RPC(io.Discard, "x", tt.services...)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("err = %v, want containing %q", err, tt.want)
			}
		})
	}
}
