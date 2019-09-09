package main

import (
	"archive/tar"
	"bytes"
	"flag"
	"github.com/badeadan/starlingx-vbox-installer/pkg/lab"
	"log"
	"os"
)

type TarWriter struct {
	*tar.Writer
}

func (tw *TarWriter) WriteFileBytes(name string, mode int64, buffer *bytes.Buffer) error {
	err := tw.WriteHeader(&tar.Header{
		Name: name,
		Mode: mode,
		Size: int64(buffer.Len()),
	})
	if err == nil {
		_, err = tw.Write(buffer.Bytes())
	}
	return err
}

func main() {
	sl := lab.StorageLab{SystemMode: "standard"}
	_default := lab.DefaultStorageLab()
	flag.StringVar(&sl.Name, "name", _default.Name, "group name")
	flag.StringVar(&sl.NatNet, "nat-net", _default.NatNet, "nat network name")
	flag.StringVar(&sl.LoopBackPrefix, "loopback-prefix", _default.LoopBackPrefix, "nat loopback prefix")
	flag.StringVar(&sl.Oam.Network, "oam-network", _default.Oam.Network, "oam network address")
	flag.StringVar(&sl.Oam.Gateway, "oam-gateway", _default.Oam.Gateway, "oam gateway")
	flag.StringVar(&sl.Oam.FloatAddr, "oam-float", _default.Oam.FloatAddr, "oam floating ip")
	flag.StringVar(&sl.Oam.Controller0, "oam-ctrl-0", _default.Oam.Controller0, "oam controller-0 ip")
	flag.StringVar(&sl.Oam.Controller1, "oam-ctrl-1", _default.Oam.Controller1, "oam controller-1 ip")
	flag.StringVar(&sl.IntNetPrefix, "intnet-prefix", _default.IntNetPrefix, "internal network  prefix")
	flag.UintVar(&sl.ControllerCpus, "controller-cpus", _default.ControllerCpus, "controller cpu count")
	flag.UintVar(&sl.ControllerMemory, "controller-memory", _default.ControllerMemory, "controller ram size")
	flag.UintVar(&sl.ControllerDiskSize, "controller-disk", _default.ControllerDiskSize, "controller disk size")
	flag.UintVar(&sl.ControllerDiskCount, "controller-disk-count", _default.ControllerDiskCount, "number of extra controller disks")
	flag.UintVar(&sl.ComputeCount, "compute-count", _default.ComputeCount, "number of compute hosts")
	flag.UintVar(&sl.ComputeCpus, "compute-cpus", _default.ComputeCpus, "compute cpu count")
	flag.UintVar(&sl.ComputeMemory, "compute-memory", _default.ComputeMemory, "compute ram size")
	flag.UintVar(&sl.ComputeDiskCount, "compute-disk-count", _default.ComputeDiskCount, "number of extra compute disks")
	flag.UintVar(&sl.ComputeDiskSize, "compute-disk", _default.ComputeDiskSize, "compute disk size")
	flag.UintVar(&sl.StorageCount, "storage-count", _default.StorageCount, "number of storage hosts")
	flag.UintVar(&sl.StorageCpus, "storage-cpus", _default.StorageCpus, "storage cpu count")
	flag.UintVar(&sl.StorageMemory, "storage-memory", _default.StorageMemory, "storage ram size")
	flag.UintVar(&sl.StorageDiskCount, "storage-disk-count", _default.StorageDiskCount, "number of storage disks (OSDs) per host")
	flag.UintVar(&sl.StorageDiskSize, "storage-disk", _default.StorageDiskSize, "storage disk size")

	flag.Parse()
	err := lab.MakeStorageInstaller(sl, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
