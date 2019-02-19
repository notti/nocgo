package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

type datatype int

const (
	void datatype = iota
	i8
	u8
	i16
	u16
	i32
	u32
	i64
	u64
	f32
	f64
	byteSlice
)

type datatypeMap struct {
	golang, C string
}

var mappings = map[datatype]datatypeMap{
	void:      {"", "void"},
	i8:        {"int8", "char"},
	u8:        {"uint8", "unsigned char"},
	i16:       {"int16", "short"},
	u16:       {"uint16", "unsigned short"},
	i32:       {"int32", "int"},
	u32:       {"uint32", "unsigned int"},
	i64:       {"int64", "long"},
	u64:       {"uint64", "unsigned long"},
	f32:       {"float32", "float"},
	f64:       {"float64", "double"},
	byteSlice: {"[]byte", "char *"},
}

type Value struct {
	Type   datatype
	CData  interface{}
	GoData interface{}
}

func (v Value) Go() string {
	return mappings[v.Type].golang
}

func (v Value) C() string {
	return mappings[v.Type].C
}

func (v Value) Void() bool {
	return v.Type == void
}

type Argument struct {
	Name string
	Value
}

func (a Argument) Go() string {
	return fmt.Sprint(a.Name, mappings[a.Type].golang)
}

func (a Argument) C() string {
	return fmt.Sprint(a.Name, mappings[a.Type].C)
}

type Arguments []Argument

func (a Arguments) Go() string {
	args := make([]string, len(a))
	for i, arg := range a {
		args[i] = arg.Go()
	}
	return strings.Join(args, ", ")
}

func (a Arguments) C() string {
	args := make([]string, len(a))
	for i, arg := range a {
		args[i] = arg.C()
	}
	return strings.Join(args, ", ")
}

func (a Arguments) Call() string {
	args := make([]string, len(a))
	for i, arg := range a {
		args[i] = arg.Name
	}
	return strings.Join(args, ", ")
}

func (a Arguments) Value() string {
	args := make([]string, len(a))
	for i, arg := range a {
		if arg.GoData == nil {
			args[i] = fmt.Sprint(arg.GoData)
		} else {
			args[i] = fmt.Sprint(arg.CData)
		}
	}
	return strings.Join(args, ", ")
}

type Test struct {
	Name      string
	Ret       Value
	Arguments Arguments
}

func (t Test) NOCGO() string {
	args := make([]string, len(t.Arguments))
	for i, arg := range t.Arguments {
		args[i] = "\t" + arg.Go()
	}
	if !t.Ret.Void() {
		args = append(args, fmt.Sprint("\tret ", t.Ret.Go(), " `nocgo:\"ret\"`"))
	}
	return strings.Join(args, "\n")
}

func (t Test) DataInit() string {
	args := make([]string, len(t.Arguments))
	for i, arg := range t.Arguments {
		if arg.GoData == nil {
			args[i] = fmt.Sprintf("%s: %s", arg.Name, arg.CData)
		} else {
			args[i] = fmt.Sprintf("%s: %s", arg.Name, arg.GoData)
		}
	}
	return strings.Join(args, ", ")
}

func (t Test) TestName() string {
	return "Test" + strings.Title(t.Name)
}

var ccode = template.Must(template.New("ccode").Parse(`#include <stdio.h>

{{range .}}{{.Ret.C}} {{.Name}}({{.Arguments.C}}) {
	return{{with .Ret.CData}} {{.}}{{end}};
}

{{end}}
`))

var bridge = template.Must(template.New("bridge").Parse(`package testlib

{{range .}}// {{.Ret.C}} {{.Name}}({{.Arguments.C}});
{{end}}import "C"

{{range .}}func {{.Name}}({{.Arguments.Go}}) {{.Ret.Go}} {
	{{if not .Ret.Void}}return {{.Ret.Go}}({{end}}C.{{.Name}}({{.Arguments.Call}}){{if not .Ret.Void}}){{end}}
}

{{end}}

`))

var testingCgo = template.Must(template.New("testingCgo").Parse(`package testlib

import (
	"testing"
)

{{range .}}func {{.TestName}}(t *testing.T) {
	{{if .Ret.Void}}{{.Name}}({{.Arguments.Value}}){{else}}ret := {{.Name}}({{.Arguments.Value}})
	if ret != {{.Ret.GoData}} {
		t.Fatalf("Expected %v, but got %v\n", {{.Ret.GoData}}, ret)
	}{{end}}
}

{{end}}
`))

var testingNOCGO = template.Must(template.New("testingNOCGO").Parse(`package testlib

import (
	"log"
	"os"
	"runtime"
	"testing"
	"unsafe"

	"github.com/notti/nocgo"
)

{{range .}}type {{.Name}}Spec struct {
{{.NOCGO}}
}

var {{.Name}}Func nocgo.Spec

func {{.TestName}}(t *testing.T) {
	arg := &{{.Name}}Spec{ {{.DataInit}} }
	{{.Name}}Func.Call(unsafe.Pointer(arg)){{if not .Ret.Void}}
	if arg.ret != {{.Ret.GoData}} {
		t.Fatalf("Expected %v, but got %v\n", {{.Ret.GoData}}, arg.ret)
	}{{end}}
}

{{end}}

func TestMain(m *testing.M) {
	var lib string
	switch runtime.GOARCH {
	case "386":
		lib = "libcalltest32.so.1"
	case "amd64":
		lib = "libcalltest64.so.1"
	default:
		log.Fatalln("Unknown arch ", runtime.GOARCH)
	}

	l, err := nocgo.Open(lib)
	if err != nil {
		log.Fatal(err)
	}

	{{range .}}{{.Name}}Func, err = l.Func("{{.Name}}", {{.Name}}Spec{})
	if err != nil {
		log.Fatal(err)
	}

	{{end}}	

	os.Exit(m.Run())
}
`))

/*
func BenchmarkCall(b *testing.B) {
	arg := &testCall{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, -11, 12}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f2.Call(unsafe.Pointer(arg))
	}
}
*/

func main() {
	cfile, err := os.Create("testlib/test.c")
	if err != nil {
		log.Fatalln("Couldn't open c file")
	}
	bridgefile, err := os.Create("testlib/cgo_bridge.go")
	if err != nil {
		log.Fatalln("Couldn't open bridge file")
	}
	testingCgofile, err := os.Create("testlib/cgo_test.go")
	if err != nil {
		log.Fatalln("Couldn't open bridge file")
	}
	testingNOCGOfile, err := os.Create("nocgo/nocgo_test.go")
	if err != nil {
		log.Fatalln("Couldn't open bridge file")
	}
	tests := []Test{
		{
			"empty",
			Value{void, nil, nil},
			Arguments{},
		},
		{
			"int1",
			Value{i8, "10", "10"},
			Arguments{},
		},
		{
			"int2",
			Value{i8, "-10", "-10"},
			Arguments{},
		},
		{
			"int3",
			Value{u8, "10", "10"},
			Arguments{},
		},
		{
			"int4",
			Value{u8, "-10", "246"},
			Arguments{},
		},
	}
	if err := ccode.Execute(cfile, tests); err != nil {
		log.Fatal(err)
	}
	cfile.Close()
	if err := bridge.Execute(bridgefile, tests); err != nil {
		log.Fatal(err)
	}
	bridgefile.Close()
	if err := testingCgo.Execute(testingCgofile, tests); err != nil {
		log.Fatal(err)
	}
	testingCgofile.Close()
	if err := testingNOCGO.Execute(testingNOCGOfile, tests); err != nil {
		log.Fatal(err)
	}
	testingNOCGOfile.Close()
}
