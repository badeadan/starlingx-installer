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

type AioSxLab struct {
	Name           string
	SystemMode     string
	NatNet         string
	LoopBackPrefix string
	IntNetPrefix   string
	Oam            lab.OamInfo
	Cpus           uint
	Memory         uint
	DiskSize       uint
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
	sx := AioSxLab{SystemMode: "simplex"}
	flag.StringVar(&sx.Name, "name", "aiosx", "group name")
	flag.StringVar(&sx.NatNet, "nat-net", "nat2", "nat network name")
	flag.StringVar(&sx.LoopBackPrefix, "loop-prefix", "127.0.2", "nat loopback prefix")
	flag.StringVar(&sx.Oam.Network, "oam-network", "10.10.10.0/24", "oam network address")
	flag.StringVar(&sx.Oam.Gateway, "oam-gateway", "10.10.10.1", "oam gateway")
	flag.StringVar(&sx.Oam.FloatAddr, "oam-float", "10.10.10.3", "oam floating ip")
	flag.StringVar(&sx.IntNetPrefix, "intnet-prefix", "intnet", "internal network  prefix")
	flag.UintVar(&sx.Cpus, "cpus", 8, "controller cpu count")
	flag.UintVar(&sx.Memory, "memory", 16, "controller ram size")
	flag.UintVar(&sx.DiskSize, "disk-size", 520, "controller disk size")

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
	err := t.ExecuteTemplate(buf, "vbox/lab/aiosx", sx)
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
		fmt.Sprintf("%s/vbox-setup.sh", sx.Name),
		0700, buf)

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "vbox/prepare-bootimage", sx)
	if err != nil {
		log.Fatal(err)
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/prepare-bootimage.sh", sx.Name),
		0700, buf)

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "install/lab/aiosx", sx)
	if err != nil {
		log.Fatal(err)
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/install.sh", sx.Name),
		0700, buf)

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "vbox/vmctl", vbox)
	if err != nil {
		log.Fatal(err)
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/vmctl.sh", sx.Name),
		0700, buf)

	err = tw.Close()
	if err != nil {
		log.Fatal(err)
	}
}
