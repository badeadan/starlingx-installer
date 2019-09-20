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
	sx := lab.AioSxLab{SystemMode: "simplex"}
	_default := lab.DefaultAioSxLab()
	flag.StringVar(&sx.Hypervisor, "hypervisor", "virtualbox", "hypervisor")
	flag.StringVar(&sx.Name, "name", _default.Name, "group name")
	flag.StringVar(&sx.NatNet, "nat-net", _default.NatNet, "nat network name")
	flag.StringVar(&sx.LoopBackPrefix, "loop-prefix", _default.LoopBackPrefix, "nat loopback prefix")
	flag.StringVar(&sx.Oam.Network, "oam-network", _default.Oam.Network, "oam network address")
	flag.StringVar(&sx.Oam.Gateway, "oam-gateway", _default.Oam.Gateway, "oam gateway")
	flag.StringVar(&sx.Oam.FloatAddr, "oam-float", _default.Oam.FloatAddr, "oam floating ip")
	flag.StringVar(&sx.IntNetPrefix, "intnet-prefix", _default.IntNetPrefix, "internal network  prefix")
	flag.UintVar(&sx.Cpus, "cpus", _default.Cpus, "controller cpu count")
	flag.UintVar(&sx.Memory, "memory", _default.Memory, "controller ram size")
	flag.UintVar(&sx.DiskSize, "disk-size", _default.DiskSize, "controller disk size")
	flag.UintVar(&sx.DiskCount, "disk-count", _default.DiskCount, "number of extra controller disks")

	flag.Parse()

	// Force packr template discovery
	_ = packr.New("VboxTemplates", "./templates/vbox")
	_ = packr.New("LibvirtTemplates", "./templates/libvirt")
	_ = packr.New("InstallTemplates", "./templates/install")

	if sx.Hypervisor == "libvirt" {
		sx.LoopBackPrefix = ""
		err := installers.MakeAioSxLibvirtInstaller(sx, os.Stdout)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := installers.MakeAioSxInstaller(sx, os.Stdout)
		if err != nil {
			log.Fatal(err)
		}
	}
}
