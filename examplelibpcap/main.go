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

type pkthdr struct {
	ts     syscall.Timeval
	caplen uint32
	len    uint32
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

	lib, err := nocgo.Open("libpcap.so")
	if err != nil {
		log.Fatalln("Couldn't load libpcap: ", err)
	}

	var pcapOpenLive func(device []byte, snaplen int32, promisc int32, toMS int32, errbuf []byte) uintptr
	if err := lib.Func("pcap_open_live", &pcapOpenLive); err != nil {
		log.Fatalln("Couldn't get pcap_open_live: ", err)
	}

	var pcapNextEx func(p uintptr, hdr **pkthdr, data *unsafe.Pointer) int32
	if err := lib.Func("pcap_next_ex", &pcapNextEx); err != nil {
		log.Fatalln("Couldn't load pcap_next_ex: ", err)
	}

	errbuf := make([]byte, 512)
	pcapHandle := pcapOpenLive(nocgo.MakeCString(*dev), int32(*snaplen), p, 100, errbuf)

	if pcapHandle == 0 {
		log.Fatalf("Couldn't open %s: %s\n", *dev, nocgo.MakeGoStringFromSlice(errbuf))
	}

	var hdr *pkthdr
	var dataptr unsafe.Pointer

	for {
		if ret := pcapNextEx(pcapHandle, &hdr, &dataptr); ret != 1 {
			log.Fatalln("Unexpected error code ", ret)
		}

		log.Printf("Received packet at %s length %d\n", time.Unix(hdr.ts.Sec, hdr.ts.Usec), hdr.len)
		packet := gopacket.NewPacket((*[1 << 30]byte)(dataptr)[:hdr.caplen], layers.LayerTypeEthernet, gopacket.Default)
		log.Println(packet.Dump())
	}
}
