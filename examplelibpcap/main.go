// Example using libpcap to read packets from a networkdevice and decode them with gopacket.
package main

import (
	"flag"
	"log"
	"syscall"
	"time"
	"unsafe"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/notti/nocgo"
)

type pcapOpenLiveArgs struct {
	device  []byte
	snaplen int32
	promisc int32
	toMS    int32
	errbuf  []byte
	ret     uintptr `nocgo:"ret"`
}

type pkthdr struct {
	ts     syscall.Timeval
	caplen uint32
	len    uint32
}

type pcapNextExArgs struct {
	p    uintptr
	hdr  **pkthdr
	data *unsafe.Pointer
	ret  int32 `nocgo:"ret"`
}

func main() {
	dev := flag.String("dev", "lo", "The device to listen to")
	snaplen := flag.Int("snaplen", 1500, "Maximum capture length")
	promisc := flag.Bool("promisc", true, "Set device to promiscous mode")

	flag.Parse()

	p := int32(0)
	if *promisc {
		p = 1
	}

	openLiveArg := pcapOpenLiveArgs{
		device:  nocgo.MakeCString(*dev),
		snaplen: int32(*snaplen),
		promisc: p,
		toMS:    100,
		errbuf:  make([]byte, 512),
	}

	lib, err := nocgo.Open("libpcap.so")
	if err != nil {
		log.Fatalln("Couldn't load libpcap: ", err)
	}
	pcapOpenLive, err := lib.Func("pcap_open_live", openLiveArg)
	if err != nil {
		log.Fatalln("Couldn't get pcap_open_live: ", err)
	}
	pcapNextEx, err := lib.Func("pcap_next_ex", pcapNextExArgs{})
	if err != nil {
		log.Fatalln("Couldn't load pcap_next_ex: ", err)
	}

	pcapOpenLive.Call(unsafe.Pointer(&openLiveArg))

	if openLiveArg.ret == 0 {
		log.Fatalf("Couldn't open %s: %s\n", *dev, nocgo.MakeGoStringFromSlice(openLiveArg.errbuf))
	}

	var hdr *pkthdr
	var dataptr unsafe.Pointer

	nextExArg := pcapNextExArgs{
		p:    openLiveArg.ret,
		hdr:  &hdr,
		data: &dataptr,
	}

	for {
		pcapNextEx.Call(unsafe.Pointer(&nextExArg))

		if nextExArg.ret != 1 {
			log.Fatalln("Unexpected error code ", nextExArg.ret)
		}

		log.Printf("Received packet at %s length %d\n", time.Unix(hdr.ts.Sec, hdr.ts.Usec), hdr.len)
		packet := gopacket.NewPacket((*[1 << 30]byte)(dataptr)[:hdr.caplen], layers.LayerTypeEthernet, gopacket.Default)
		log.Println(packet.Dump())
	}
}
