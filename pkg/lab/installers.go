package lab

import (
	"archive/tar"
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/gobuffalo/packr/v2"
	"gopkg.in/yaml.v2"
	"io"
	"net"
	"text/template"
	"time"
)

type TarWriter struct {
	*tar.Writer
}

func (tw *TarWriter) WriteFileBytes(name string, mode int64, buffer *bytes.Buffer) error {
	modtime := time.Now()
	err := tw.WriteHeader(&tar.Header{
		Name:    name,
		Mode:    mode,
		Size:    int64(buffer.Len()),
		ModTime: modtime,
	})
	if err == nil {
		_, err = tw.Write(buffer.Bytes())
	}
	return err
}

func MakeAioSxInstaller(sx AioSxLab, out io.Writer) error {
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
	box := packr.New("VboxTemplates", "./templates/vbox")
	t = template.Must(DiscoverTemplates(box, "vbox", t))
	box = packr.New("InstallTemplates", "./templates/install")
	t = template.Must(DiscoverTemplates(box, "install", t))
	tw := &TarWriter{tar.NewWriter(out)}

	vbox := Lab{}
	buf := &bytes.Buffer{}
	err := t.ExecuteTemplate(buf, "vbox/lab/aiosx", sx)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf.Bytes(), &vbox)
	if err != nil {
		return err
	}

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "vbox/setup", vbox)
	if err != nil {
		return err
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/vbox-setup.sh", sx.Name),
		0700, buf)

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "vbox/prepare-bootimage", sx)
	if err != nil {
		return err
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/prepare-bootimage.sh", sx.Name),
		0700, buf)

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "install/lab/aiosx", sx)
	if err != nil {
		return err
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/install.sh", sx.Name),
		0700, buf)

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "vbox/vmctl", vbox)
	if err != nil {
		return err
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/vmctl.sh", sx.Name),
		0700, buf)

	return tw.Close()
}

func MakeAioDxInstaller(dx AioDxLab, out io.Writer) error {
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
	box := packr.New("VboxTemplates", "./templates/vbox")
	t = template.Must(DiscoverTemplates(box, "vbox", t))
	box = packr.New("InstallTemplates", "./templates/install")
	t = template.Must(DiscoverTemplates(box, "install", t))
	tw := &TarWriter{tar.NewWriter(out)}

	vbox := Lab{}
	buf := &bytes.Buffer{}
	err := t.ExecuteTemplate(buf, "vbox/lab/aiodx", dx)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf.Bytes(), &vbox)
	if err != nil {
		return err
	}

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "vbox/setup", vbox)
	if err != nil {
		return err
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/vbox-setup.sh", dx.Name),
		0700, buf)

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "vbox/prepare-bootimage", dx)
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
	err = t.ExecuteTemplate(buf, "vbox/vmctl", vbox)
	if err != nil {
		return err
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/vmctl.sh", dx.Name),
		0700, buf)

	return tw.Close()
}

func MakeStandardInstaller(sl StandardLab, out io.Writer) error {
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
	box := packr.New("VboxTemplates", "./templates/vbox")
	t = template.Must(DiscoverTemplates(box, "vbox", t))
	box = packr.New("InstallTemplates", "./templates/install")
	t = template.Must(DiscoverTemplates(box, "install", t))
	tw := &TarWriter{tar.NewWriter(out)}

	vbox := Lab{}
	buf := &bytes.Buffer{}
	err := t.ExecuteTemplate(buf, "vbox/lab/standard", sl)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf.Bytes(), &vbox)
	if err != nil {
		return err
	}

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "vbox/setup", vbox)
	if err != nil {
		return err
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/vbox-setup.sh", sl.Name),
		0700, buf)

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "vbox/prepare-bootimage", sl)
	if err != nil {
		return err
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/prepare-bootimage.sh", sl.Name),
		0700, buf)

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "install/lab/standard", sl)
	if err != nil {
		return err
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/install.sh", sl.Name),
		0700, buf)

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "vbox/vmctl", vbox)
	if err != nil {
		return err
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/vmctl.sh", sl.Name),
		0700, buf)

	return tw.Close()
}

func MakeStorageInstaller(sl StorageLab, out io.Writer) error {

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
	box := packr.New("VboxTemplates", "./templates/vbox")
	t = template.Must(DiscoverTemplates(box, "vbox", t))
	box = packr.New("InstallTemplates", "./templates/install")
	t = template.Must(DiscoverTemplates(box, "install", t))
	tw := &TarWriter{tar.NewWriter(out)}

	vbox := Lab{}
	buf := &bytes.Buffer{}
	err := t.ExecuteTemplate(buf, "vbox/lab/storage", sl)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buf.Bytes(), &vbox)
	if err != nil {
		return err
	}

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "vbox/setup", vbox)
	if err != nil {
		return err
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/vbox-setup.sh", sl.Name),
		0700, buf)

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "vbox/prepare-bootimage", sl)
	if err != nil {
		return err
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/prepare-bootimage.sh", sl.Name),
		0700, buf)

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "install/lab/storage", sl)
	if err != nil {
		return err
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/install.sh", sl.Name),
		0700, buf)

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "vbox/vmctl", vbox)
	if err != nil {
		return err
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/vmctl.sh", sl.Name),
		0700, buf)

	return tw.Close()
}
