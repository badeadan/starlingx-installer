package main

import (
	"flag"
	"github.com/badeadan/starlingx-vbox-installer/pkg/lab"
	"log"
	"os"
)

func main() {
	dx := lab.AioDxLab{SystemMode: "duplex"}
	flag.StringVar(&dx.Name, "name", "aiodx", "group name")
	flag.StringVar(&dx.NatNet, "nat-net", "nat3", "nat network name")
	flag.StringVar(&dx.LoopBackPrefix, "loop-prefix", "127.0.3", "nat loopback prefix")
	flag.StringVar(&dx.Oam.Network, "oam-network", "10.10.10.0/24", "oam network address")
	flag.StringVar(&dx.Oam.Gateway, "oam-gateway", "10.10.10.1", "oam gateway")
	flag.StringVar(&dx.Oam.FloatAddr, "oam-float", "10.10.10.2", "oam floating ip")
	flag.StringVar(&dx.Oam.Controller0, "oam-ctrl-0", "10.10.10.3", "oam controller-0 ip")
	flag.StringVar(&dx.Oam.Controller1, "oam-ctrl-1", "10.10.10.4", "oam controller-1 ip")
	flag.StringVar(&dx.IntNetPrefix, "intnet-prefix", "intnet", "internal network  prefix")
	flag.UintVar(&dx.Cpus, "cpus", 8, "controller cpu count")
	flag.UintVar(&dx.Memory, "memory", 16, "controller ram size")
	flag.UintVar(&dx.DiskSize, "disk-size", 520, "controller disk size")

	flag.Parse()
	err := lab.MakeAioDxInstaller(dx, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
