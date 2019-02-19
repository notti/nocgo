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
	golang, C, Cgo string
}

var mappings = map[datatype]datatypeMap{
	void:      {"", "void", ""},
	i8:        {"int8", "char", "char"},
	u8:        {"uint8", "unsigned char", "uchar"},
	i16:       {"int16", "short", "short"},
	u16:       {"uint16", "unsigned short", "ushort"},
	i32:       {"int32", "int", "int"},
	u32:       {"uint32", "unsigned int", "uint"},
	i64:       {"int64", "long long", "longlong"},
	u64:       {"uint64", "unsigned long long", "ulonglong"},
	f32:       {"float32", "float", "float"},
	f64:       {"float64", "double", "double"},
	byteSlice: {"[]byte", "char *", "unsafe.Pointer"},
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
	return a.Name + " " + mappings[a.Type].golang
}

func (a Argument) C() string {
	return mappings[a.Type].C + " " + a.Name
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
		args[i] = fmt.Sprintf("C.%s(%s)", mappings[arg.Type].Cgo, arg.Name)
	}
	return strings.Join(args, ", ")
}

func (a Arguments) Value() string {
	args := make([]string, len(a))
	for i, arg := range a {
		if arg.GoData != nil {
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
	t.Log({{.Name}}Func)
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
		{
			"int5",
			Value{u8, "a+b", "44"},
			Arguments{
				{"a", Value{u8, "100", nil}},
				{"b", Value{u8, "200", nil}},
			},
		},
		{
			"int6",
			Value{u64, "a", "100"},
			Arguments{
				{"a", Value{u8, "100", nil}},
			},
		},
		{
			"int7",
			Value{u8, "a", "100"},
			Arguments{
				{"a", Value{u64, "100", nil}},
			},
		},
		{
			"intBig",
			Value{u64, "81985529216486895", "uint64(81985529216486895)"},
			Arguments{},
		},
		{
			"float1",
			Value{f32, "10.5", "10.5"},
			Arguments{},
		},
		{
			"float2",
			Value{f64, "10.5", "10.5"},
			Arguments{},
		},
		{
			"stackSpill1",
			Value{i8, "a+b+c+d+e+f+g+h", "8"},
			Arguments{
				{"a", Value{i8, "1", nil}},
				{"b", Value{i8, "1", nil}},
				{"c", Value{i8, "1", nil}},
				{"d", Value{i8, "1", nil}},
				{"e", Value{i8, "1", nil}},
				{"f", Value{i8, "1", nil}},
				{"g", Value{i8, "1", nil}},
				{"h", Value{i8, "1", nil}},
			},
		},
		{
			"stackSpill2",
			Value{f32, "a+b+c+d+e+f+g+h+i+j", "10"},
			Arguments{
				{"a", Value{f32, "1", nil}},
				{"b", Value{f32, "1", nil}},
				{"c", Value{f32, "1", nil}},
				{"d", Value{f32, "1", nil}},
				{"e", Value{f32, "1", nil}},
				{"f", Value{f32, "1", nil}},
				{"g", Value{f32, "1", nil}},
				{"h", Value{f32, "1", nil}},
				{"i", Value{f32, "1", nil}},
				{"j", Value{f32, "1", nil}},
			},
		},
		{
			"stackSpill3",
			Value{i8, "ia+ib+ic+id+ie+f+ig+ih+fa+fb+fc+fd+fe+ff+fg+fh+fi+fj", "18"},
			Arguments{
				{"ia", Value{i8, "1", nil}},
				{"ib", Value{i8, "1", nil}},
				{"ic", Value{i8, "1", nil}},
				{"id", Value{i8, "1", nil}},
				{"ie", Value{i8, "1", nil}},
				{"f", Value{i8, "1", nil}},
				{"ig", Value{i8, "1", nil}},
				{"ih", Value{i8, "1", nil}},
				{"fa", Value{f32, "1", nil}},
				{"fb", Value{f32, "1", nil}},
				{"fc", Value{f32, "1", nil}},
				{"fd", Value{f32, "1", nil}},
				{"fe", Value{f32, "1", nil}},
				{"ff", Value{f32, "1", nil}},
				{"fg", Value{f32, "1", nil}},
				{"fh", Value{f32, "1", nil}},
				{"fi", Value{f32, "1", nil}},
				{"fj", Value{f32, "1", nil}},
			},
		},
		{
			"stackSpill4",
			Value{i8, "ia+ib+ic+id+ie+f+ig+ih+fa+fb+fc+fd+fe+ff+fg+fh+fi+fj", "18"},
			Arguments{
				{"ia", Value{i8, "1", nil}},
				{"fa", Value{f32, "1", nil}},
				{"ib", Value{i8, "1", nil}},
				{"fb", Value{f32, "1", nil}},
				{"ic", Value{i8, "1", nil}},
				{"fc", Value{f32, "1", nil}},
				{"id", Value{i8, "1", nil}},
				{"fd", Value{f32, "1", nil}},
				{"ie", Value{i8, "1", nil}},
				{"fe", Value{f32, "1", nil}},
				{"f", Value{i8, "1", nil}},
				{"ff", Value{f32, "1", nil}},
				{"ig", Value{i8, "1", nil}},
				{"fg", Value{f32, "1", nil}},
				{"ih", Value{i8, "1", nil}},
				{"fh", Value{f32, "1", nil}},
				{"fi", Value{f32, "1", nil}},
				{"fj", Value{f32, "1", nil}},
			},
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
