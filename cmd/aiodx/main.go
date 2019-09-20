package main

import (
	"flag"
	"github.com/badeadan/starlingx-vbox-installer/pkg/lab"
	"github.com/badeadan/starlingx-vbox-installer/pkg/lab/installers"
	"github.com/gobuffalo/packr/v2"
	"log"
	"os"
)

func main() {
	dx := lab.AioDxLab{SystemMode: "duplex"}
	_default := lab.DefaultAioDxLab()
	flag.StringVar(&dx.Hypervisor, "hypervisor", "virtualbox", "hypervisor")
	flag.StringVar(&dx.Name, "name", _default.Name, "group name")
	flag.StringVar(&dx.NatNet, "nat-net", _default.NatNet, "nat network name")
	flag.StringVar(&dx.LoopBackPrefix, "loop-prefix", _default.LoopBackPrefix, "nat loopback prefix")
	flag.StringVar(&dx.Oam.Network, "oam-network", _default.Oam.Network, "oam network address")
	flag.StringVar(&dx.Oam.Gateway, "oam-gateway", _default.Oam.Gateway, "oam gateway")
	flag.StringVar(&dx.Oam.FloatAddr, "oam-float", _default.Oam.FloatAddr, "oam floating ip")
	flag.StringVar(&dx.Oam.Controller0, "oam-ctrl-0", _default.Oam.Controller0, "oam controller-0 ip")
	flag.StringVar(&dx.Oam.Controller1, "oam-ctrl-1", _default.Oam.Controller1, "oam controller-1 ip")
	flag.StringVar(&dx.IntNetPrefix, "intnet-prefix", _default.IntNetPrefix, "internal network  prefix")
	flag.UintVar(&dx.Cpus, "cpus", _default.Cpus, "controller cpu count")
	flag.UintVar(&dx.Memory, "memory", _default.Memory, "controller ram size")
	flag.UintVar(&dx.DiskSize, "disk-size", _default.DiskSize, "controller disk size")
	flag.UintVar(&dx.DiskCount, "disk-count", _default.DiskCount, "number of extra controller disks")

	flag.Parse()

	// Force packr template discovery
	_ = packr.New("VboxTemplates", "./templates/vbox")
	_ = packr.New("LibvirtTemplates", "./templates/libvirt")
	_ = packr.New("InstallTemplates", "./templates/install")

	if dx.Hypervisor == "libvirt" {
		dx.LoopBackPrefix = ""
		err := installers.MakeAioDxLibvirtInstaller(dx, os.Stdout)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := installers.MakeAioDxInstaller(dx, os.Stdout)
		if err != nil {
			log.Fatal(err)
		}
	}
}
