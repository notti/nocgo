/*
Package nocgo provides a dlopen wrapper that doesn't need cgo.

WARNING! This also supports cgo! So if you want to ensure no cgo, you have to set the environment variable CGO_ENABLED=0 for the build process.

Usage

Don't use this in production! This is meant as a PoC for golang and subject to changes.

So far only floats, integers, and pointers are supported (No structs, no complex types, and no callbacks, but that wouldn't be hard to implement).
See example, examplelibpcap, and the documentation below for examples.

Go to C type mappings

In the type specification struct go types, which will be mapped to C types, have to be used.

	Go type                                    C Type
	=========================================================
	int8, byte                                 char
	uint8, bool                                unsigned char
	int16                                      short
	uint16                                     unsigned short
	int32                                      int
	uint32                                     unsigned int
	int64                                      long
	uint64                                     unsigned long
	float32                                    float
	float64                                    double
	[], uintptr, reflect.UnsafePointer, *      *

The last line means that slices and pointers are mapped to pointers in C. Pointers to structs are possible, but the structs must follow C alignment rules!
Be carefull with the mappings - there is NO type checking, which is actually impossible since the imported library doesn't know these things.

go int is not supported to prevent confusion: int sizes in go for 32bit and 64 bit differ, while in C (cdecl) they do not!

Helperfunctions for converting between C strings and go strings are provided (see below in the function documentation).

Argument Specifications

Arguments to functions must be specified via a struct. If the return value of the function is needed, a field with the tag `nocgo:"ret"` can be added.
This field can appear anywhere in the struct (doesn't need to be at the end). Only one return field can be provided.

Example for pcap_open_live (libpcap):

C declaration:
	pcap_t *pcap_open_live(const char *device, int snaplen, int promisc, int to_ms, char *errbuf);

nocgo declaration:
	type pcapOpenLiveArgs struct {
		device  []byte
		snaplen int32
		promisc int32
		toMS    int32
		errbuf  []byte
		ret     uintptr `nocgo:"ret"`
	}

ret uses uintptr as type since pcap_t* is just passed to every other libpcap function and we don't care or need to know whats actually in there or where it's pointing to.
*/
package nocgo
