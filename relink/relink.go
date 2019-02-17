package main

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

type def struct {
	interp string
	libs   []string
}

var defaults = map[string]def{
	"amd64_linux": def{
		interp: "/lib64/ld-linux-x86-64.so.2",
		libs:   []string{"libc.so.6", "libpthread.so", "libdl.so"},
	},
}

type elfFile struct {
	f         *os.File
	e         *elf.File
	phoff     uint64
	phentsize uint64
	shoff     uint64
	shentsize uint64
	shstrndx  uint64
	buffer    [1024]byte
}

func openElfFile(name string) (*elfFile, error) {
	ret := &elfFile{}
	var err error
	ret.f, err = os.OpenFile(name, os.O_RDWR, 0777)
	if err != nil {
		return nil, err
	}
	ret.e, err = elf.NewFile(ret.f)
	if err != nil {
		return nil, err
	}

	_, err = ret.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	switch ret.e.Class {
	case elf.ELFCLASS32:
		var hdr elf.Header32
		err = ret.ReadData(&hdr)
		if err != nil {
			return nil, err
		}
		ret.phoff = uint64(hdr.Phoff)
		ret.phentsize = uint64(hdr.Phentsize)
		ret.shoff = uint64(hdr.Shoff)
		ret.shentsize = uint64(hdr.Shentsize)
		ret.shstrndx = uint64(hdr.Shstrndx)
	case elf.ELFCLASS64:
		var hdr elf.Header64
		err = ret.ReadData(&hdr)
		if err != nil {
			return nil, err
		}
		ret.phoff = hdr.Phoff
		ret.phentsize = uint64(hdr.Phentsize)
		ret.shoff = hdr.Shoff
		ret.shentsize = uint64(hdr.Shentsize)
		ret.shstrndx = uint64(hdr.Shstrndx)
	}

	return ret, nil
}

func (e *elfFile) Seek(where int64, whence int) (int64, error) {
	return e.f.Seek(where, whence)
}

func (e *elfFile) ReadData(data interface{}) error {
	return binary.Read(e.f, e.e.ByteOrder, data)
}

func (e *elfFile) WriteData(data interface{}) error {
	return binary.Write(e.f, e.e.ByteOrder, data)
}

func (e *elfFile) Write(b []byte) error {
	_, err := e.f.Write(b)
	return err
}

func (e *elfFile) WriteAt(b []byte, off uint64) error {
	_, err := e.f.WriteAt(b, int64(off))
	return err
}

func (e *elfFile) Read8() (uint8, error) {
	_, err := e.f.Read(e.buffer[:1])
	return e.buffer[0], err
}

func (e *elfFile) Read16() (uint16, error) {
	_, err := e.f.Read(e.buffer[:2])
	return e.e.ByteOrder.Uint16(e.buffer[:]), err
}

func (e *elfFile) Read32() (uint32, error) {
	_, err := e.f.Read(e.buffer[:4])
	return e.e.ByteOrder.Uint32(e.buffer[:]), err
}

func (e *elfFile) Read64() (uint64, error) {
	_, err := e.f.Read(e.buffer[:8])
	return e.e.ByteOrder.Uint64(e.buffer[:]), err
}

func (e *elfFile) Read8At(where int64) (uint8, error) {
	_, err := e.f.ReadAt(e.buffer[:1], where)
	return e.buffer[0], err
}

func (e *elfFile) Read16At(where int64) (uint16, error) {
	_, err := e.f.ReadAt(e.buffer[:2], where)
	return e.e.ByteOrder.Uint16(e.buffer[:]), err
}

func (e *elfFile) Read32At(where int64) (uint32, error) {
	_, err := e.f.ReadAt(e.buffer[:4], where)
	return e.e.ByteOrder.Uint32(e.buffer[:]), err
}

func (e *elfFile) Read64At(where int64) (uint64, error) {
	_, err := e.f.ReadAt(e.buffer[:8], where)
	return e.e.ByteOrder.Uint64(e.buffer[:]), err
}

func (e *elfFile) Write8(data uint8) error {
	_, err := e.f.Write([]byte{data})
	return err
}

func (e *elfFile) Write16(data uint16) error {
	e.e.ByteOrder.PutUint16(e.buffer[:], data)
	_, err := e.f.Write(e.buffer[:2])
	return err
}

func (e *elfFile) Write32(data uint32) error {
	e.e.ByteOrder.PutUint32(e.buffer[:], data)
	_, err := e.f.Write(e.buffer[:4])
	return err
}

func (e *elfFile) Write64(data uint64) error {
	e.e.ByteOrder.PutUint64(e.buffer[:], data)
	_, err := e.f.Write(e.buffer[:8])
	return err
}

func (e *elfFile) Write8At(data uint8, where uint64) error {
	return e.WriteAt([]byte{data}, where)
}

func (e *elfFile) Write16At(data uint16, where uint64) error {
	e.e.ByteOrder.PutUint16(e.buffer[:], data)
	return e.WriteAt(e.buffer[:2], where)
}

func (e *elfFile) Write32At(data uint32, where uint64) error {
	e.e.ByteOrder.PutUint32(e.buffer[:], data)
	return e.WriteAt(e.buffer[:4], where)
}

func (e *elfFile) Write64At(data uint64, where uint64) error {
	e.e.ByteOrder.PutUint64(e.buffer[:], data)
	return e.WriteAt(e.buffer[:8], where)
}

func (e *elfFile) WriteSections() error {
	if _, err := e.f.Seek(int64(e.shoff), io.SeekStart); err != nil {
		return err
	}
	shstrtab := 0
	shstrtabpos := e.shoff
	names := make([]uint64, len(e.e.Sections))
	for i, sec := range e.e.Sections {
		names[i] = e.shoff - shstrtabpos
		if err := e.Write(append([]byte(sec.Name), 0)); err != nil {
			return err
		}
		e.shoff += uint64(len(sec.Name)) + 1
		if sec.Name == ".shstrtab" {
			shstrtab = i
		}
	}
	e.e.Sections[shstrtab].Offset = shstrtabpos
	e.e.Sections[shstrtab].FileSize = e.shoff - shstrtabpos

	for i, sec := range e.e.Sections {
		if err := e.WriteSection(sec.SectionHeader, names[i]); err != nil {
			return err
		}
	}
	switch e.e.Class {
	case elf.ELFCLASS32:
		if err := e.Write32At(uint32(e.shoff), 0x20); err != nil {
			return err
		}
		if err := e.Write16At(uint16(len(e.e.Sections)), 0x30); err != nil {
			return err
		}
	case elf.ELFCLASS64:
		if err := e.Write64At(e.shoff, 0x28); err != nil {
			return err
		}
		if err := e.Write16At(uint16(len(e.e.Sections)), 0x3C); err != nil {
			return err
		}
	}
	return nil
}

func (e *elfFile) WriteSection(sh elf.SectionHeader, name uint64) error {
	switch e.e.Class {
	case elf.ELFCLASS32:
		hdr := elf.Section32{
			Name:      uint32(name),
			Type:      uint32(sh.Type),
			Flags:     uint32(sh.Flags),
			Addr:      uint32(sh.Addr),
			Off:       uint32(sh.Offset),
			Size:      uint32(sh.FileSize),
			Link:      sh.Link,
			Info:      sh.Info,
			Addralign: uint32(sh.Addralign),
			Entsize:   uint32(sh.Entsize),
		}
		return e.WriteData(hdr)
	case elf.ELFCLASS64:
		hdr := elf.Section64{
			Name:      uint32(name),
			Type:      uint32(sh.Type),
			Flags:     uint64(sh.Flags),
			Addr:      sh.Addr,
			Off:       sh.Offset,
			Size:      sh.FileSize,
			Link:      sh.Link,
			Info:      sh.Info,
			Addralign: sh.Addralign,
			Entsize:   sh.Entsize,
		}
		return e.WriteData(hdr)
	}
	// compression header not handeled
	return errors.New("Unknown elf bit size")
}

func (e *elfFile) WritePrograms() error {
	if _, err := e.Seek(int64(e.phoff), io.SeekStart); err != nil {
		return err
	}
	for _, prog := range e.e.Progs {
		if err := e.WriteProgram(prog.ProgHeader); err != nil {
			return err
		}
	}
	switch e.e.Class {
	case elf.ELFCLASS32:
		if err := e.Write32At(uint32(e.phoff), 0x1C); err != nil {
			return err
		}
		if err := e.Write16At(uint16(len(e.e.Progs)), 0x2C); err != nil {
			return err
		}
	case elf.ELFCLASS64:
		if err := e.Write64At(e.phoff, 0x20); err != nil {
			return err
		}
		if err := e.Write16At(uint16(len(e.e.Progs)), 0x38); err != nil {
			return err
		}
	}
	return nil
}

func (e *elfFile) WriteProgram(ph elf.ProgHeader) error {
	switch e.e.Class {
	case elf.ELFCLASS32:
		hdr := elf.Prog32{
			Type:   uint32(ph.Type),
			Flags:  uint32(ph.Flags),
			Off:    uint32(ph.Off),
			Vaddr:  uint32(ph.Vaddr),
			Paddr:  uint32(ph.Paddr),
			Filesz: uint32(ph.Filesz),
			Memsz:  uint32(ph.Memsz),
			Align:  uint32(ph.Align),
		}
		return e.WriteData(hdr)
	case elf.ELFCLASS64:
		hdr := elf.Prog64{
			Type:   uint32(ph.Type),
			Flags:  uint32(ph.Flags),
			Off:    ph.Off,
			Vaddr:  ph.Vaddr,
			Paddr:  ph.Paddr,
			Filesz: ph.Filesz,
			Memsz:  ph.Memsz,
			Align:  ph.Align,
		}
		return e.WriteData(hdr)
	}
	return errors.New("Unknown elf bit size")
}

func (e *elfFile) Read(b []byte) error {
	_, err := e.f.Read(b)
	return err
}

func (e *elfFile) Copy(src, dst int64, length int) error {
	buffer := make([]byte, length)
	_, err := e.Seek(src, io.SeekStart)
	if err != nil {
		return err
	}
	err = e.Read(buffer)
	if err != nil {
		return err
	}
	_, err = e.Seek(dst, io.SeekStart)
	if err != nil {
		return err
	}
	err = e.Write(buffer)
	if err != nil {
		return err
	}
	return nil
}

func (e *elfFile) Close() error {
	return e.f.Close()
}

// Dyn contains a single entry of the dynamic table
type Dyn struct {
	Tag elf.DynTag
	Val uint64
}

func padding(addr, align uint64) uint64 {
	align1 := align - 1
	return (align - (addr & align1)) & align1
}

// DynSymbol represents a dynamic symbol
type DynSymbol struct {
	Name    string
	Value   uint64
	Size    uint64
	Bind    elf.SymBind
	Type    elf.SymType
	Vis     elf.SymVis
	Section int
}

func (e *elfFile) makeDynsym(elements []DynSymbol) (dynsym, dynstr []byte) {
	sym := &bytes.Buffer{}
	str := &bytes.Buffer{}
	for _, elem := range elements {
		namei := str.Len()
		str.Write(append([]byte(elem.Name), 0))
		switch e.e.Class {
		case elf.ELFCLASS32:
			binary.Write(sym, e.e.ByteOrder, elf.Sym32{
				Name:  uint32(namei),
				Value: uint32(elem.Value),
				Size:  uint32(elem.Size),
				Info:  byte(elem.Bind)<<4 | byte(elem.Type)&0x0f,
				Other: byte(elem.Vis) & 0x03,
				Shndx: uint16(elem.Section),
			})
		case elf.ELFCLASS64:
			binary.Write(sym, e.e.ByteOrder, elf.Sym64{
				Name:  uint32(namei),
				Value: uint64(elem.Value),
				Size:  uint64(elem.Size),
				Info:  byte(elem.Bind)<<4 | byte(elem.Type)&0x0f,
				Other: byte(elem.Vis) & 0x03,
				Shndx: uint16(elem.Section),
			})
		}
	}
	if str.Len() == 0 {
		str.WriteByte(0)
	}
	return sym.Bytes(), str.Bytes()
}

func (e *elfFile) makeDynsec(elements []Dyn) []byte {
	ret := &bytes.Buffer{}
	switch e.e.Class {
	case elf.ELFCLASS32:
		var secs []elf.Dyn32
		for _, sec := range elements {
			secs = append(secs, elf.Dyn32{
				Tag: int32(sec.Tag),
				Val: uint32(sec.Val),
			})
		}
		binary.Write(ret, e.e.ByteOrder, secs)
	case elf.ELFCLASS64:
		var secs []elf.Dyn64
		for _, sec := range elements {
			secs = append(secs, elf.Dyn64{
				Tag: int64(sec.Tag),
				Val: uint64(sec.Val),
			})
		}
		binary.Write(ret, e.e.ByteOrder, secs)
	}
	return ret.Bytes()
}

// RelSymbol represents a symbol in need of relocation
type RelSymbol struct {
	Off   uint64
	SymNo uint64
}

func (e *elfFile) makeDynRel(symbols []RelSymbol) ([]byte, bool, uint64) {
	ret := &bytes.Buffer{}
	var rela bool
	var relt uint64
	switch e.e.Machine {
	case elf.EM_386:
		rela = false
		relt = uint64(elf.R_386_JMP_SLOT)
	case elf.EM_X86_64:
		rela = true
		relt = uint64(elf.R_X86_64_JMP_SLOT)
	default:
		log.Fatal("Unknown machine type ", e.e.Machine)
	}

	var relsz uint64

	switch e.e.Class {
	case elf.ELFCLASS32:
		if rela {
			for _, symbol := range symbols {
				binary.Write(ret, e.e.ByteOrder, elf.Rela32{
					Off:  uint32(symbol.Off),
					Info: uint32(symbol.SymNo<<8 | relt),
				})
			}
			relsz = 12
		} else {
			for _, symbol := range symbols {

				binary.Write(ret, e.e.ByteOrder, elf.Rel32{
					Off:  uint32(symbol.Off),
					Info: uint32(symbol.SymNo<<8 | relt),
				})
			}
			relsz = 8
		}
	case elf.ELFCLASS64:
		if rela {
			for _, symbol := range symbols {
				binary.Write(ret, e.e.ByteOrder, elf.Rela64{
					Off:  symbol.Off,
					Info: symbol.SymNo<<32 | relt,
				})
			}
			relsz = 24
		} else {
			for _, symbol := range symbols {
				binary.Write(ret, e.e.ByteOrder, elf.Rel64{
					Off:  symbol.Off,
					Info: symbol.SymNo<<32 | relt,
				})
			}
			relsz = 16
		}
	}
	return ret.Bytes(), rela, relsz
}

func main() {
	libsArg := flag.String("libs", "", "Load comma separated list of given libs instead of defaults")
	interp := flag.String("interp", "", "Use given interp instead of default")

	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatal("Need a static golang binary as argument")
	}

	libs := strings.Split(*libsArg, ",")

	f, err := openElfFile(flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}

	/*
		        KEEP EXEC (we are not dyn after all)
		        try to put new program headers, dyn, interp into first 4k
		        0  -+-----------------------------------+--
		            | ELF                               |
		            +-----------------------------------+
		            | program headers                   |
		            +-----------------------------------+
		            | interp                            |
		            +-----------------------------------+
		should be   | dyn stuff                         |
		 below 4k ->+-----------------------------------+
		            | other stuff that needs relocation |
		            +-----------------------------------+<- ensure mapping until here
		            +-----------------------------------+
		   entry -> | Everything else (e.g., text)      |
		            +-----------------------------------+
		            | .shstrtab                         |
		            +-----------------------------------+
		            | Section headers                   |
		            +-----------------------------------+
	*/

	// First some sanity checks - and checks if we can do our meddling, after all we don't support everything in this POC

	if f.e.Type != elf.ET_EXEC {
		log.Fatal("only static binaries not using an interp supported")
	}

	var base uint64
	var baseProg int

	for i, prog := range f.e.Progs {
		if prog.Type == elf.PT_INTERP || prog.Type == elf.PT_DYNAMIC {
			log.Fatal("only static binaries not using an interp supported")
		}
		if prog.Type == elf.PT_LOAD {
			if base == 0 {
				base = prog.Vaddr
				baseProg = i
			} else if prog.Vaddr < base {
				base = prog.Vaddr
				baseProg = i
			}
		}
	}

	if uint64(f.phoff+f.phentsize*uint64(len(f.e.Progs))) > f.e.Entry {
		log.Fatal("Not enough space before entry point")
	}

	symbolList, err := f.e.Symbols()
	if err != nil {
		log.Fatal(err)
	}

	for _, sym := range symbolList {
		if strings.HasPrefix(sym.Name, "_rt0_") {
			if d, ok := defaults[sym.Name[5:]]; ok {
				libs = d.libs
				*interp = d.interp
				break
			}
		}
	}

	interpProg := len(f.e.Progs)

	f.e.Progs = append(f.e.Progs, &elf.Prog{
		ProgHeader: elf.ProgHeader{
			Type:   elf.PT_INTERP,
			Flags:  elf.PF_R,
			Off:    0, // fill later
			Vaddr:  0, // fill later
			Paddr:  0, // fill later
			Filesz: 0, // fill later
			Memsz:  0, // fill later
			Align:  1,
		}})

	dynsecProg := len(f.e.Progs)

	f.e.Progs = append(f.e.Progs, &elf.Prog{
		ProgHeader: elf.ProgHeader{
			Type:   elf.PT_DYNAMIC,
			Flags:  elf.PF_R | elf.PF_W,
			Off:    0, // fill later
			Vaddr:  0, // fill later
			Paddr:  0, // fill later
			Filesz: 0, // fill later
			Memsz:  0, // fill later
			Align:  8,
		}})

	interpPos := f.phoff + f.phentsize*uint64(len(f.e.Progs))
	interpB := append([]byte(*interp), 0)
	interpLen := uint64(len(interpB))

	f.e.Progs[interpProg].Off = interpPos
	f.e.Progs[interpProg].Vaddr = interpPos + base
	f.e.Progs[interpProg].Paddr = interpPos + base
	f.e.Progs[interpProg].Filesz = interpLen
	f.e.Progs[interpProg].Memsz = interpLen

	hashPos := interpPos + interpLen
	hashPos += padding(hashPos, 8)
	hash := make([]byte, 8*4) // Empty 64bit DT_HASH
	hashLen := uint64(len(hash))

	var relList []RelSymbol

	var symsection int

	var symdefs []DynSymbol

	symdefs = append(symdefs, DynSymbol{
		Name:    "",
		Value:   0,
		Size:    0,
		Bind:    elf.STB_LOCAL,
		Type:    elf.STT_NOTYPE,
		Vis:     elf.STV_DEFAULT,
		Section: int(elf.SHN_UNDEF),
	})

	xCgoInit := uint64(0)
	cgoInit := uint64(0)
	cgoSize := uint64(0)

	for _, sym := range symbolList {
		if strings.HasSuffix(sym.Name, "__dynload") {
			parts := strings.Split(sym.Name, ".")
			name := parts[len(parts)-1]
			dynsym := name[:len(name)-9]

			symsection = int(sym.Section)
			relList = append(relList, RelSymbol{
				Off:   sym.Value,
				SymNo: uint64(len(symdefs)),
			})
			symdefs = append(symdefs, DynSymbol{
				Name:    dynsym,
				Value:   0,
				Size:    0,
				Bind:    elf.STB_GLOBAL,
				Type:    elf.STT_FUNC,
				Vis:     elf.STV_DEFAULT,
				Section: int(elf.SHN_UNDEF),
			})
		}
		if sym.Name == "x_cgo_init" {
			xCgoInit = sym.Value
		}
		if sym.Name == "_cgo_init" {
			sec := f.e.Sections[sym.Section]
			cgoInit = sym.Value - sec.Addr + sec.Offset
			cgoSize = sym.Size
		}
	}

	if xCgoInit != 0 && cgoInit != 0 && cgoSize != 0 {
		switch cgoSize {
		case 4:
			f.Write32At(uint32(xCgoInit), cgoInit)
		case 8:
			f.Write64At(xCgoInit, cgoInit)
		default:
			log.Fatalln("Unknown symbol size", cgoSize)
		}
	}

	dynsym, dynstr := f.makeDynsym(symdefs)

	var libOffsets []uint64

	for _, l := range libs {
		libOffsets = append(libOffsets, uint64(len(dynstr)))
		dynstr = append(dynstr, []byte(l)...)
		dynstr = append(dynstr, 0)
	}

	dynsymLocal := 0
	dynstrPos := hashPos + hashLen
	dynstrLen := uint64(len(dynstr))

	dynsymPos := dynstrPos + dynstrLen
	dynsymPos += padding(dynsymPos, 8)
	dynsymLen := uint64(len(dynsym))

	// TODO: DT_BIND_NOW?

	dynrel, rela, relsz := f.makeDynRel(relList)
	dynrelPos := dynsymPos + dynsymLen
	dynrelPos += padding(dynrelPos, 8)
	dynrelLen := uint64(len(dynrel))

	var dynsecs []Dyn
	for _, offset := range libOffsets {
		dynsecs = append(dynsecs, Dyn{Tag: elf.DT_NEEDED, Val: uint64(offset)})
	}

	if rela {
		dynsecs = append(dynsecs, Dyn{Tag: elf.DT_RELA, Val: uint64(base + dynrelPos)})
		dynsecs = append(dynsecs, Dyn{Tag: elf.DT_RELASZ, Val: uint64(dynrelLen)})
		dynsecs = append(dynsecs, Dyn{Tag: elf.DT_RELAENT, Val: uint64(relsz)})
	} else {
		dynsecs = append(dynsecs, Dyn{Tag: elf.DT_REL, Val: uint64(base + dynrelPos)})
		dynsecs = append(dynsecs, Dyn{Tag: elf.DT_RELSZ, Val: uint64(dynrelLen)})
		dynsecs = append(dynsecs, Dyn{Tag: elf.DT_RELENT, Val: uint64(relsz)})
	}

	dynsecs = append(dynsecs, []Dyn{
		{Tag: elf.DT_STRTAB, Val: base + dynstrPos},
		{Tag: elf.DT_STRSZ, Val: dynstrLen},
		{Tag: elf.DT_SYMTAB, Val: base + dynsymPos},
		{Tag: elf.DT_SYMENT, Val: dynsymLen},
		{Tag: elf.DT_HASH, Val: hashPos + base},
		{Tag: elf.DT_BIND_NOW, Val: 0},
		{Tag: elf.DT_NULL, Val: 0},
	}...)

	dynsec := f.makeDynsec(dynsecs)
	dynsecPos := dynrelPos + dynrelLen
	dynsecPos += padding(dynsecPos, 8)
	dynsecLen := uint64(len(dynsec))

	f.e.Progs[dynsecProg].Off = dynsecPos
	f.e.Progs[dynsecProg].Vaddr = dynsecPos + base
	f.e.Progs[dynsecProg].Paddr = dynsecPos + base
	f.e.Progs[dynsecProg].Filesz = dynsecLen
	f.e.Progs[dynsecProg].Memsz = dynsecLen

	afterDynsec := dynsecPos + dynsecLen

	relPos := afterDynsec
	var torelocate []*elf.Section
	relocated := make(map[int]bool)

	for {
		var newRelocate []*elf.Section
		for i, sec := range f.e.Sections {
			if sec.Type == elf.SHT_NULL {
				continue
			}
			if sec.Offset < relPos && !relocated[i] {
				newRelocate = append(newRelocate, sec)
				relocated[i] = true
			}
		}
		if len(newRelocate) == 0 {
			break
		}
		torelocate = append(torelocate, newRelocate...)

		sort.Slice(torelocate, func(i, j int) bool { return torelocate[i].Offset < torelocate[j].Offset })
		relPos = afterDynsec
		for _, sec := range torelocate {
			relPos += sec.Size
			if sec.Addralign > 1 {
				relPos += padding(relPos, sec.Addralign)
			}
		}
	}

	for _, sec := range torelocate {
		data := make([]byte, sec.Size)
		if _, err := f.f.ReadAt(data, int64(sec.Offset)); err != nil {
			log.Fatal(err)
		}
		if sec.Addralign > 1 {
			afterDynsec += padding(afterDynsec, sec.Addralign)
		}
		if err := f.WriteAt(data, afterDynsec); err != nil {
			log.Fatal(err)
		}
		for _, prog := range f.e.Progs {
			if prog.Off == sec.Offset {
				prog.Off = afterDynsec
			}
			if prog.Vaddr == sec.Offset+base {
				prog.Vaddr = afterDynsec + base
				prog.Paddr = afterDynsec + base
			}
		}

		sec.Addr += afterDynsec - sec.Offset // or base + offset
		sec.Offset, afterDynsec = afterDynsec, afterDynsec+sec.Offset
	}

	if afterDynsec > f.e.Entry {
		log.Fatal("not enough space before entry point")
	}

	if f.e.Progs[baseProg].Filesz < afterDynsec {
		f.e.Progs[baseProg].Filesz = afterDynsec
		f.e.Progs[baseProg].Memsz = afterDynsec
	}

	if err := f.WritePrograms(); err != nil {
		log.Fatal(err)
	}

	if err := f.WriteAt(interpB, interpPos); err != nil {
		log.Fatal(err)
	}

	if err := f.WriteAt(hash, hashPos); err != nil {
		log.Fatal(err)
	}

	if err := f.WriteAt(dynstr, dynstrPos); err != nil {
		log.Fatal(err)
	}

	if err := f.WriteAt(dynsym, dynsymPos); err != nil {
		log.Fatal(err)
	}

	if err := f.WriteAt(dynrel, dynrelPos); err != nil {
		log.Fatal(err)
	}

	if err := f.WriteAt(dynsec, dynsecPos); err != nil {
		log.Fatal(err)
	}

	f.e.Sections = append(f.e.Sections, &elf.Section{
		SectionHeader: elf.SectionHeader{
			Name:      ".interp",
			Type:      elf.SHT_PROGBITS,
			Flags:     elf.SHF_ALLOC,
			Addr:      base + interpPos,
			Offset:    interpPos,
			FileSize:  interpLen,
			Addralign: 1,
		}})

	dynstrI := len(f.e.Sections)

	f.e.Sections = append(f.e.Sections, &elf.Section{
		SectionHeader: elf.SectionHeader{
			Name:      ".dynstr",
			Type:      elf.SHT_STRTAB,
			Flags:     elf.SHF_ALLOC,
			Addr:      base + dynstrPos,
			Offset:    dynstrPos,
			FileSize:  dynstrLen,
			Addralign: 1,
		}})

	entSize := uint64(24)
	if f.e.Class == elf.ELFCLASS32 {
		entSize = 16
	}

	dynsymSec := len(f.e.Sections)

	f.e.Sections = append(f.e.Sections, &elf.Section{
		SectionHeader: elf.SectionHeader{
			Name:      ".dynsym",
			Type:      elf.SHT_DYNSYM,
			Flags:     elf.SHF_ALLOC,
			Addr:      base + dynsymPos,
			Offset:    dynsymPos,
			FileSize:  dynsymLen,
			Addralign: 8,
			Link:      uint32(dynstrI),
			Entsize:   entSize,
			Info:      uint32(dynsymLocal + 1),
		}})

	entSize = uint64(16)
	if f.e.Class == elf.ELFCLASS32 {
		entSize = 8
	}

	f.e.Sections = append(f.e.Sections, &elf.Section{
		SectionHeader: elf.SectionHeader{
			Name:      ".dynamic",
			Type:      elf.SHT_DYNAMIC,
			Flags:     elf.SHF_ALLOC | elf.SHF_WRITE,
			Addr:      base + dynsecPos,
			Offset:    dynsecPos,
			FileSize:  dynsecLen,
			Addralign: 8,
			Link:      uint32(dynstrI),
			Entsize:   entSize,
		}})

	dynname := ".rel"
	if rela {
		dynname = ".rela"
	}
	dynname += f.e.Sections[symsection].Name

	shtype := elf.SHT_REL
	if rela {
		shtype = elf.SHT_RELA
	}

	f.e.Sections = append(f.e.Sections, &elf.Section{
		SectionHeader: elf.SectionHeader{
			Name:      dynname,
			Type:      shtype,
			Flags:     elf.SHF_ALLOC,
			Addr:      base + dynrelPos,
			Offset:    dynrelPos,
			FileSize:  dynrelLen,
			Addralign: 8,
			Link:      uint32(dynsymSec),
			Info:      uint32(symsection),
			Entsize:   relsz,
		}})

	f.e.Sections = append(f.e.Sections, &elf.Section{
		SectionHeader: elf.SectionHeader{
			Name:      ".hash",
			Type:      elf.SHT_HASH,
			Flags:     elf.SHF_ALLOC,
			Addr:      base + hashPos,
			Offset:    hashPos,
			FileSize:  hashLen,
			Addralign: 8,
			Link:      uint32(dynsymSec),
		}})

	shoff, err := f.f.Seek(0, io.SeekEnd)
	if err != nil {
		log.Fatal(err)
	}
	f.shoff = uint64(shoff)

	if err := f.WriteSections(); err != nil {
		log.Fatal(err)
	}

	f.Close()
}
