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
	golang, C, Cgo, Cgop string
}

var mappings = map[datatype]datatypeMap{
	void:      {"", "void", "", ""},
	i8:        {"int8", "char", "C.char", ""},
	u8:        {"uint8", "unsigned char", "C.uchar", ""},
	i16:       {"int16", "short", "C.short", ""},
	u16:       {"uint16", "unsigned short", "C.ushort", ""},
	i32:       {"int32", "int", "C.int", ""},
	u32:       {"uint32", "unsigned int", "C.uint", ""},
	i64:       {"int64", "long long", "C.longlong", ""},
	u64:       {"uint64", "unsigned long long", "C.ulonglong", ""},
	f32:       {"float32", "float", "C.float", ""},
	f64:       {"float64", "double", "C.double", ""},
	byteSlice: {"[]byte", "char *", "unsafe.Pointer", "*C.char"},
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
		if mappings[arg.Type].Cgo == "unsafe.Pointer" {
			args[i] = fmt.Sprintf("(%s)(%s(&%s[0]))", mappings[arg.Type].Cgop, mappings[arg.Type].Cgo, arg.Name)
		} else {
			args[i] = fmt.Sprintf("%s(%s)", mappings[arg.Type].Cgo, arg.Name)
		}
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
	Name            string
	Ret             Value
	Arguments       Arguments
	CgoPre, CgoPost string
	NOCGOPre, NOCGOPost string
	Bench           bool
	Multitest       bool
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

func (t Test) BenchmarkName() string {
	return "Benchmark" + strings.Title(t.Name)
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

import "unsafe"

var _ = unsafe.Sizeof(0)

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
	{{.CgoPre}}
	{{if .Ret.Void}}{{.Name}}({{.Arguments.Value}}){{else}}ret := {{.Name}}({{.Arguments.Value}})
	if ret != {{.Ret.GoData}} {
		t.Fatalf("Expected %v, but got %v\n", {{.Ret.GoData}}, ret)
	}{{end}}
	{{.CgoPost}}
}
{{if .Bench}}
func {{.BenchmarkName}}(b *testing.B) {
	{{.CgoPre}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		{{.Name}}({{.Arguments.Value}})
	}
}
{{end}}
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

{{if .Multitest}}
func {{.TestName}}Multi(t *testing.T) {
	for i:=0; i < 100; i++ {
		t.Run("{{.TestName}}Multi", func(t *testing.T) {
			t.Parallel()
			{{.NOCGOPre}}
			arg := &{{.Name}}Spec{ {{.DataInit}} }
			t.Log({{.Name}}Func)
			{{.Name}}Func.Call(unsafe.Pointer(arg)){{if not .Ret.Void}}
			if arg.ret != {{.Ret.GoData}} {
				t.Fatalf("Expected %v, but got %v\n", {{.Ret.GoData}}, arg.ret)
			}{{end}}
			{{.NOCGOPost}}
		})
	}
}
{{else}}
func {{.TestName}}(t *testing.T) {
	{{.NOCGOPre}}
	arg := &{{.Name}}Spec{ {{.DataInit}} }
	t.Log({{.Name}}Func)
	{{.Name}}Func.Call(unsafe.Pointer(arg)){{if not .Ret.Void}}
	if arg.ret != {{.Ret.GoData}} {
		t.Fatalf("Expected %v, but got %v\n", {{.Ret.GoData}}, arg.ret)
	}{{end}}
	{{.NOCGOPost}}
}
{{end}}
{{if .Bench}}
func {{.BenchmarkName}}(b *testing.B) {
	{{.NOCGOPre}}
	arg := &{{.Name}}Spec{ {{.DataInit}} }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		{{.Name}}Func.Call(unsafe.Pointer(arg))
	}
}
{{end}}
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
			Name:  "empty",
			Ret:   Value{void, nil, nil},
			Bench: true,
		},
		{
			Name: "int1",
			Ret:  Value{i8, "10", "10"},
		},
		{
			Name: "int2",
			Ret:  Value{i8, "-10", "-10"},
		},
		{
			Name: "int3",
			Ret:  Value{u8, "10", "10"},
		},
		{
			Name: "int4",
			Ret:  Value{u8, "-10", "246"},
		},
		{
			Name: "int5",
			Ret:  Value{u8, "a+b", "44"},
			Arguments: Arguments{
				{"a", Value{u8, "100", nil}},
				{"b", Value{u8, "200", nil}},
			},
		},
		{
			Name: "int6",
			Ret:  Value{u64, "a", "100"},
			Arguments: Arguments{
				{"a", Value{u8, "100", nil}},
			},
		},
		{
			Name: "int7",
			Ret:  Value{u8, "a", "100"},
			Arguments: Arguments{
				{"a", Value{u64, "100", nil}},
			},
		},
		{
			Name: "intBig1",
			Ret:  Value{u64, "81985529216486895", "uint64(81985529216486895)"},
		},
		{
			Name: "intBig2",
			Ret:  Value{u64, "a", "uint64(81985529216486895)"},
			Arguments: Arguments{
				{"a", Value{u64, "81985529216486895", nil}},
			},
		},
		{
			Name: "float1",
			Ret:  Value{f32, "10.5", "10.5"},
		},
		{
			Name:  "float2",
			Ret:   Value{f64, "10.5", "10.5"},
			Bench: true,
		},
		{
			Name: "stackSpill1",
			Ret:  Value{i8, "a+b+c+d+e+f+g+h", "8"},
			Arguments: Arguments{
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
			Name: "stackSpill2",
			Ret:  Value{f32, "a+b+c+d+e+f+g+h+i+j", "10"},
			Arguments: Arguments{
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
			Name: "stackSpill3",
			Ret:  Value{i8, "ia+ib+ic+id+ie+f+ig+ih+fa+fb+fc+fd+fe+ff+fg+fh+fi+fj", "18"},
			Arguments: Arguments{
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
			Bench: true,
		},
		{
			Name: "stackSpill4",
			Ret:  Value{i8, "ia+ib+ic+id+ie+f+ig+ih+fa+fb+fc+fd+fe+ff+fg+fh+fi+fj", "18"},
			Arguments: Arguments{
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
		{
			Name: "funcall1",
			Ret:  Value{i32, `sprintf(s, "test from C: %d %1.1f %s\n", a, b, c)`, "27"},
			Arguments: Arguments{
				{"s", Value{byteSlice, nil, "buf"}},
				{"a", Value{i8, "-1", nil}},
				{"b", Value{f32, "1.5", nil}},
				{"c", Value{byteSlice, nil, `[]byte("gotest\000")`}},
			},
			CgoPre: "buf := make([]byte, 1024)",
			CgoPost: `	if string(buf[:ret]) != "test from C: -1 1.5 gotest\n" {
		t.Fatalf("Expected \"test from C: -1 1.5 gotest\n\", but got \"%s\"", string(buf[:ret]))
	}`,
			NOCGOPre: "buf := make([]byte, 1024)",
			NOCGOPost: `	if string(buf[:arg.ret]) != "test from C: -1 1.5 gotest\n" {
		t.Fatalf("Expected \"test from C: -1 1.5 gotest\n\", but got \"%s\"", string(buf[:arg.ret]))
	}`,
			Multitest: true,
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
