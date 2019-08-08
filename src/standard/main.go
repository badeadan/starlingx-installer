package main

import (
	"archive/tar"
	"bytes"
	"flag"
	"fmt"
	"github.com/Masterminds/sprig"
	"gopkg.in/yaml.v2"
	"lab"
	"log"
	"net"
	"os"
	"text/template"
)

type StandardLab struct {
	Name           string
	SystemMode     string
	NatNet         string
	LoopBackPrefix string
	IntNetPrefix   string
	Oam            lab.OamInfo
	ControllerCpus     uint
	ControllerMemory   uint
	ControllerDiskSize uint
	ComputeCount       uint
	ComputeCpus        uint
	ComputeMemory      uint
	ComputeDiskSize    uint
}

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
	sl := StandardLab{SystemMode: "duplex"}
	flag.StringVar(&sl.Name, "name", "standard", "group name")
	flag.StringVar(&sl.NatNet, "nat-net", "nat4", "nat network name")
	flag.StringVar(&sl.LoopBackPrefix, "loop-prefix", "127.0.4", "nat loopback prefix")
	flag.StringVar(&sl.Oam.Network, "oam-network", "10.10.10.0/24", "oam network address")
	flag.StringVar(&sl.Oam.Gateway, "oam-gateway", "10.10.10.1", "oam gateway")
	flag.StringVar(&sl.Oam.FloatAddr, "oam-float", "10.10.10.2", "oam floating ip")
	flag.StringVar(&sl.Oam.Controller0, "oam-ctrl-0", "10.10.10.3", "oam controller-0 ip")
	flag.StringVar(&sl.Oam.Controller1, "oam-ctrl-1", "10.10.10.4", "oam controller-1 ip")
	flag.StringVar(&sl.IntNetPrefix, "intnet-prefix", "intnet", "internal network  prefix")
	flag.UintVar(&sl.ControllerCpus, "controller-cpus", 4, "controller cpu count")
	flag.UintVar(&sl.ControllerMemory, "controller-memory", 16, "controller ram size")
	flag.UintVar(&sl.ControllerDiskSize, "controller-disk-size", 520, "controller disk size")
	flag.UintVar(&sl.ComputeCount, "compute-count", 2, "number of compute hosts")
	flag.UintVar(&sl.ComputeCpus, "compute-cpus", 4, "compute cpu count")
	flag.UintVar(&sl.ComputeMemory, "compute-memory", 10, "compute ram size")
	flag.UintVar(&sl.ComputeDiskSize, "compute-disk", 520, "compute disk size")

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
	t = template.Must(lab.DiscoverTemplates("./templates/vbox", t))
	t = template.Must(lab.DiscoverTemplates("./templates/install", t))
	tw := &TarWriter{tar.NewWriter(os.Stdout)}

	vbox := lab.Lab{}
	buf := &bytes.Buffer{}
	err := t.ExecuteTemplate(buf, "vbox/lab/standard", sl)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(buf.Bytes(), &vbox)
	if err != nil {
		log.Fatal(err)
	}

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "vbox/setup", vbox)
	if err != nil {
		log.Fatal(err)
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/vbox-setup.sh", sl.Name),
		0700, buf)

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "vbox/prepare-bootimage", sl)
	if err != nil {
		log.Fatal(err)
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/prepare-bootimage.sh", sl.Name),
		0700, buf)

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "install/lab/standard", sl)
	if err != nil {
		log.Fatal(err)
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/install.sh", sl.Name),
		0700, buf)

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "vbox/vmctl", vbox)
	if err != nil {
		log.Fatal(err)
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/vmctl.sh", sl.Name),
		0700, buf)

	err = tw.Close()
	if err != nil {
		log.Fatal(err)
	}
}
