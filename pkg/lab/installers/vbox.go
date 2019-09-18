package installers

import (
	"archive/tar"
	"bytes"
	"fmt"
	"github.com/badeadan/starlingx-vbox-installer/pkg/lab"
	"github.com/gobuffalo/packr/v2"
	"gopkg.in/yaml.v2"
	"io"
	"text/template"
)

func MakeAioSxInstaller(sx lab.AioSxLab, out io.Writer) error {
	t := lab.NewTxtTemplate()
	box := packr.New("VboxTemplates", "./templates/vbox")
	t = template.Must(lab.DiscoverTemplates(box, "vbox", t))
	box = packr.New("InstallTemplates", "./templates/install")
	t = template.Must(lab.DiscoverTemplates(box, "install", t))
	tw := &TarWriter{tar.NewWriter(out)}

	vbox := lab.Lab{}
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
		fmt.Sprintf("%s/setup.sh", sx.Name),
		0700, buf)

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "install/prepare-bootimage", sx)
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

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "install/readme", sx)
	if err != nil {
		return err
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/README.txt", sx.Name),
		0700, buf)

	return tw.Close()
}

func MakeAioDxInstaller(dx lab.AioDxLab, out io.Writer) error {
	t := lab.NewTxtTemplate()
	box := packr.New("VboxTemplates", "./templates/vbox")
	t = template.Must(lab.DiscoverTemplates(box, "vbox", t))
	box = packr.New("InstallTemplates", "./templates/install")
	t = template.Must(lab.DiscoverTemplates(box, "install", t))
	tw := &TarWriter{tar.NewWriter(out)}

	vbox := lab.Lab{}
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
		fmt.Sprintf("%s/setup.sh", dx.Name),
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
	err = t.ExecuteTemplate(buf, "vbox/vmctl", vbox)
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

func MakeStandardInstaller(sl lab.StandardLab, out io.Writer) error {
	t := lab.NewTxtTemplate()
	box := packr.New("VboxTemplates", "./templates/vbox")
	t = template.Must(lab.DiscoverTemplates(box, "vbox", t))
	box = packr.New("InstallTemplates", "./templates/install")
	t = template.Must(lab.DiscoverTemplates(box, "install", t))
	tw := &TarWriter{tar.NewWriter(out)}

	vbox := lab.Lab{}
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
		fmt.Sprintf("%s/setup.sh", sl.Name),
		0700, buf)

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "install/prepare-bootimage", sl)
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

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "install/readme", sl)
	if err != nil {
		return err
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/README.txt", sl.Name),
		0700, buf)

	return tw.Close()
}

func MakeStorageInstaller(sl lab.StorageLab, out io.Writer) error {
	t := lab.NewTxtTemplate()
	box := packr.New("VboxTemplates", "./templates/vbox")
	t = template.Must(lab.DiscoverTemplates(box, "vbox", t))
	box = packr.New("InstallTemplates", "./templates/install")
	t = template.Must(lab.DiscoverTemplates(box, "install", t))
	tw := &TarWriter{tar.NewWriter(out)}

	vbox := lab.Lab{}
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
		fmt.Sprintf("%s/setup.sh", sl.Name),
		0700, buf)

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "install/prepare-bootimage", sl)
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

	buf = &bytes.Buffer{}
	err = t.ExecuteTemplate(buf, "install/readme", sl)
	if err != nil {
		return err
	}
	tw.WriteFileBytes(
		fmt.Sprintf("%s/README.txt", sl.Name),
		0700, buf)

	return tw.Close()
}
