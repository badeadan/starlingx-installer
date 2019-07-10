package main

import (
	"archive/tar"
	"bytes"
	"flag"
	"fmt"
	"github.com/Masterminds/sprig"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"text/template"
	"net"
)

type StorageLab struct {
	Name               string
	Iso                string
	NatNet             string
	IntNetPrefix       string
	Oam                OamInfo
	ControllerCpus     uint
	ControllerMemory   uint
	ControllerDiskSize uint
	ComputeCount       uint
	ComputeCpus        uint
	ComputeMemory      uint
	ComputeDiskSize    uint
	StorageCount       uint
	StorageCpus        uint
	StorageMemory      uint
	StorageDiskCount   uint
	StorageDiskSize    uint
}

func main() {
	stor := StorageLab{}
	flag.StringVar(&stor.Name, "name", "storage", "group name")
	flag.StringVar(&stor.Iso, "iso", "./bootimage.iso", "install iso path")
	flag.StringVar(&stor.NatNet, "nat-net", "nat0", "nat network name")
	flag.StringVar(&stor.Oam.Network, "oam-network", "10.10.10.0/24", "oam network address")
	flag.StringVar(&stor.Oam.Gateway, "oam-gateway", "10.10.10.1", "oam gateway")
	flag.StringVar(&stor.Oam.FloatAddr, "oam-float", "10.10.10.2", "oam floating ip")
	flag.StringVar(&stor.Oam.Controller0, "oam-ctrl-0", "10.10.10.3", "oam controller-0 ip")
	flag.StringVar(&stor.Oam.Controller1, "oam-ctrl-1", "10.10.10.4", "oam controller-1 ip")
	flag.StringVar(&stor.IntNetPrefix, "intnet-prefix", "intnet", "internal network  prefix")
	flag.UintVar(&stor.ControllerCpus, "controller-cpus", 8, "controller cpu count")
	flag.UintVar(&stor.ControllerMemory, "controller-memory", 12, "controller ram size")
	flag.UintVar(&stor.ControllerDiskSize, "controller-disk", 520, "controller disk size")
	flag.UintVar(&stor.ComputeCount, "compute-count", 2, "number of compute hosts")
	flag.UintVar(&stor.ComputeCpus, "compute-cpus", 8, "compute cpu count")
	flag.UintVar(&stor.ComputeMemory, "compute-memory", 12, "compute ram size")
	flag.UintVar(&stor.ComputeDiskSize, "compute-disk", 520, "compute disk size")
	flag.UintVar(&stor.StorageCount, "storage-count", 2, "number of storage hosts")
	flag.UintVar(&stor.StorageCpus, "sorage-cpus", 8, "sorage cpu count")
	flag.UintVar(&stor.StorageMemory, "sorage-memory", 12, "sorage ram size")
	flag.UintVar(&stor.StorageDiskCount, "sorage-disk-count", 4, "number of storage disks per host")
	flag.UintVar(&stor.StorageDiskSize, "sorage-disk", 520, "sorage disk size")

	flag.Parse()

	t := template.New("")
	t = t.Funcs(sprig.TxtFuncMap())
	t = t.Funcs(template.FuncMap{
		"include": func(name string, data interface{}) (string, error) {
			buf := &bytes.Buffer{}
			err := t.ExecuteTemplate(buf, name, data)
			return buf.String(), err
		},
		"NetCidrMask": func(cidr string) (string, error) {
			_, n, err := net.ParseCIDR(cidr)
			mask := ""
			if err == nil {
				mask = fmt.Sprintf("%d.%d.%d.%d",
					n.Mask[0], n.Mask[1], n.Mask[2], n.Mask[3])
			}
			return mask, err
		},
	})
	t = template.Must(t.ParseGlob("./templates/*.tmpl"))

	lab := Lab{}
	buf := &bytes.Buffer{}
	err := t.ExecuteTemplate(buf, "lab-storage-config", stor)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(buf.Bytes(), &lab)
	if err != nil {
		log.Fatal(err)
	}
	err = t.ExecuteTemplate(os.Stdout, "vbox-setup", lab)
	if err != nil {
		log.Fatal(err)
	}
}
