package installers

import (
	"archive/tar"
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/badeadan/starlingx-vbox-installer/pkg/lab"
	"github.com/gobuffalo/packr/v2"
	"gopkg.in/yaml.v2"
	"io"
	"net"
	"text/template"
)

func MakeAioDxLibvirtInstaller(dx lab.AioDxLab, out io.Writer) error {
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
	box := packr.New("LibvirtTemplates", "./templates/libvirt")
	t = template.Must(lab.DiscoverTemplates(box, "libvirt", t))
	box = packr.New("InstallTemplates", "./templates/install")
	t = template.Must(lab.DiscoverTemplates(box, "install", t))
	tw := &TarWriter{tar.NewWriter(out)}

	virt := lab.LibvirtLab{}
	buf := &bytes.Buffer{}
	err := t.ExecuteTemplate(buf, "libvirt/lab/aiodx", dx)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf.Bytes(), &virt)
	if err != nil {
		return err
	}

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "libvirt/setup", virt)
	if err != nil {
		return err
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/libvirt-setup.sh", dx.Name),
		0700, buf)

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "install/prepare-bootimage", dx)
	if err != nil {
		return err
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/prepare-bootimage.sh", dx.Name),
		0700, buf)

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "install/lab/aiodx", dx)
	if err != nil {
		return err
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/install.sh", dx.Name),
		0700, buf)

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "libvirt/vmctl", virt)
	if err != nil {
		return err
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/vmctl.sh", dx.Name),
		0700, buf)

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "install/readme", dx)
	if err != nil {
		return err
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/README.txt", dx.Name),
		0700, buf)

	return tw.Close()
}
