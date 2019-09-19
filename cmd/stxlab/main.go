// Produces a form like:
//
package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/badeadan/starlingx-vbox-installer/pkg/lab"
	"github.com/badeadan/starlingx-vbox-installer/pkg/lab/installers"
	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/schema"
	"github.com/joncalhoun/form"
	"html/template"
	"io"
	"net/http"
	"time"
)

func redirect2AioSx(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/aiosx", 301)
}

type AioSxForm struct {
	Name           string `form:"label=Lab Name"`
	Hypervisor     string `form:"label=Hypervisor;type=select"`
	NatNet         string `form:"label=NAT Network"`
	LoopBackPrefix string `form:"label=Loopback Prefix"`
	IntNetPrefix   string `form:"label=Internal Network Prefix"`
	Network        string `form:"label=OAM Network prefix & mask"`
	Gateway        string `form:"label=OAM Gateway"`
	FloatAddr      string `form:"label=OAM Controller IP address"`
	Cpus           uint   `form:"label=Number of CPUs"`
	Memory         uint   `form:"label=Memory size (GB)"`
	DiskSize       uint   `form:"label=Disk size (GB)"`
	DiskCount      uint   `form:"label=Number of extra controller disks"`
}

func AioSxForm2Lab(form AioSxForm) lab.AioSxLab {
	return lab.AioSxLab{
		Name:           form.Name,
		Hypervisor:     form.Hypervisor,
		NatNet:         form.NatNet,
		LoopBackPrefix: form.LoopBackPrefix,
		IntNetPrefix:   form.IntNetPrefix,
		Oam: lab.OamInfo{
			Network:   form.Network,
			Gateway:   form.Gateway,
			FloatAddr: form.FloatAddr,
		},
		Cpus:      form.Cpus,
		Memory:    form.Memory,
		DiskSize:  form.DiskSize,
		DiskCount: form.DiskCount,
	}
}

func AioSxLab2Form(l lab.AioSxLab) AioSxForm {
	return AioSxForm{
		Name:           l.Name,
		Hypervisor:     l.Hypervisor,
		NatNet:         l.NatNet,
		LoopBackPrefix: l.LoopBackPrefix,
		IntNetPrefix:   l.IntNetPrefix,
		Network:        l.Oam.Network,
		Gateway:        l.Oam.Gateway,
		FloatAddr:      l.Oam.FloatAddr,
		Cpus:           l.Cpus,
		Memory:         l.Memory,
		DiskSize:       l.DiskSize,
		DiskCount:      l.DiskCount,
	}
}

func handleAioSx(w http.ResponseWriter, r *http.Request) {
	box := packr.New("WebTemplates", "./templates/web")
	inputTpl, err := box.FindString("input.tmpl")
	if err != nil {
		panic(err)
	}
	tpl := template.Must(template.New("").Funcs(sprig.FuncMap()).Funcs(form.FuncMap()).Parse(inputTpl))
	fb := form.Builder{
		InputTemplate: tpl,
	}
	pageTpl, err := box.FindString("page.tmpl")
	if err != nil {
		panic(err)
	}
	tpl = template.Must(template.New("").Funcs(sprig.FuncMap()).Funcs(fb.FuncMap()).Parse(pageTpl))
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "text/html")
		data := struct {
			Type   string
			Form   AioSxForm
			Errors []error
		}{
			Type: "AioSX",
			Form: AioSxLab2Form(lab.DefaultAioSxLab()),
		}
		err := tpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
		return
	case http.MethodPost:
		r.ParseForm()
		dec := schema.NewDecoder()
		dec.IgnoreUnknownKeys(true)
		var form AioSxForm
		err := dec.Decode(&form, r.PostForm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		modtime := time.Now()
		var content bytes.Buffer
		if form.Hypervisor == "LibVirt" {
			err = installers.MakeAioSxLibvirtInstaller(
				AioSxForm2Lab(form), io.Writer(&content))
		} else {
			err = installers.MakeAioSxInstaller(
				AioSxForm2Lab(form), io.Writer(&content))
		}
		const name = "stx-lab-aiosx.tar"
		w.Header().Add("Content-Disposition",
			fmt.Sprintf("Attachment; filename=%s", name))
		http.ServeContent(w, r, name, modtime, bytes.NewReader(content.Bytes()))
	default:
		http.NotFound(w, r)
		return
	}
}

type AioDxForm struct {
	Name           string `form:"label=Lab Name"`
	Hypervisor     string `form:"label=Hypervisor;type=select"`
	NatNet         string `form:"label=NAT Network"`
	LoopBackPrefix string `form:"label=Loopback Prefix"`
	IntNetPrefix   string `form:"label=Internal Network Prefix"`
	Network        string `form:"label=OAM Network prefix & mask"`
	Gateway        string `form:"label=OAM Gateway"`
	FloatAddr      string `form:"label=OAM Floating IP address"`
	Controller0    string `form:"label=OAM Controller-0 IP address"`
	Controller1    string `form:"label=OAM Controller-1 IP address"`
	Cpus           uint   `form:"label=Number of CPUs"`
	Memory         uint   `form:"label=Memory size (GB)"`
	DiskSize       uint   `form:"label=Disk size (GB)"`
	DiskCount      uint   `form:"label=Number of extra controller disks"`
}

func AioDxForm2Lab(form AioDxForm) lab.AioDxLab {
	return lab.AioDxLab{
		Name:           form.Name,
		Hypervisor:     form.Hypervisor,
		NatNet:         form.NatNet,
		LoopBackPrefix: form.LoopBackPrefix,
		IntNetPrefix:   form.IntNetPrefix,
		Oam: lab.OamInfo{
			Network:     form.Network,
			Gateway:     form.Gateway,
			FloatAddr:   form.FloatAddr,
			Controller0: form.Controller0,
			Controller1: form.Controller1,
		},
		Cpus:      form.Cpus,
		Memory:    form.Memory,
		DiskSize:  form.DiskSize,
		DiskCount: form.DiskCount,
	}
}

func AioDxLab2Form(l lab.AioDxLab) AioDxForm {
	return AioDxForm{
		Name:           l.Name,
		Hypervisor:     l.Hypervisor,
		NatNet:         l.NatNet,
		LoopBackPrefix: l.LoopBackPrefix,
		IntNetPrefix:   l.IntNetPrefix,
		Network:        l.Oam.Network,
		Gateway:        l.Oam.Gateway,
		FloatAddr:      l.Oam.FloatAddr,
		Controller0:    l.Oam.Controller0,
		Controller1:    l.Oam.Controller1,
		Cpus:           l.Cpus,
		Memory:         l.Memory,
		DiskSize:       l.DiskSize,
		DiskCount:      l.DiskCount,
	}
}

func handleAioDx(w http.ResponseWriter, r *http.Request) {
	box := packr.New("WebTemplates", "./templates/web")
	inputTpl, err := box.FindString("input.tmpl")
	if err != nil {
		panic(err)
	}
	tpl := template.Must(template.New("").Funcs(sprig.FuncMap()).Funcs(form.FuncMap()).Parse(inputTpl))
	fb := form.Builder{
		InputTemplate: tpl,
	}
	pageTpl, err := box.FindString("page.tmpl")
	if err != nil {
		panic(err)
	}
	tpl = template.Must(template.New("").Funcs(sprig.FuncMap()).Funcs(fb.FuncMap()).Parse(pageTpl))
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "text/html")
		data := struct {
			Type   string
			Form   AioDxForm
			Errors []error
		}{
			Type: "AioDX",
			Form: AioDxLab2Form(lab.DefaultAioDxLab()),
		}
		err := tpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
		return
	case http.MethodPost:
		r.ParseForm()
		dec := schema.NewDecoder()
		dec.IgnoreUnknownKeys(true)
		var form AioDxForm
		err := dec.Decode(&form, r.PostForm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		modtime := time.Now()
		var content bytes.Buffer
		if form.Hypervisor == "LibVirt" {
			err = installers.MakeAioDxLibvirtInstaller(
				AioDxForm2Lab(form),
				io.Writer(&content))
		} else {
			err = installers.MakeAioDxInstaller(
				AioDxForm2Lab(form),
				io.Writer(&content))
		}
		const name = "stx-lab-aiodx.tar"
		w.Header().Add("Content-Disposition",
			fmt.Sprintf("Attachment; filename=%s", name))
		http.ServeContent(w, r, name, modtime, bytes.NewReader(content.Bytes()))
	default:
		http.NotFound(w, r)
		return
	}
}

type StandardForm struct {
	Name                string `form:"label=Lab Name"`
	Hypervisor          string `form:"label=Hypervisor;type=select"`
	NatNet              string `form:"label=NAT Network"`
	LoopBackPrefix      string `form:"label=Loopback Prefix"`
	IntNetPrefix        string `form:"label=Internal Network Prefix"`
	Network             string `form:"label=OAM Network prefix & mask"`
	Gateway             string `form:"label=OAM Gateway"`
	FloatAddr           string `form:"label=OAM Floating IP address"`
	Controller0         string `form:"label=OAM Controller-0 IP address"`
	Controller1         string `form:"label=OAM Controller-1 IP address"`
	ControllerCpus      uint   `form:"label=Controller Number of CPUs"`
	ControllerMemory    uint   `form:"label=Controller Memory size (GB)"`
	ControllerDiskSize  uint   `form:"label=Controller Disk size (GB)"`
	ControllerDiskCount uint   `form:"label=Number of extra controller disks"`
	ComputeCount        uint   `form:"label=Number of Computes"`
	ComputeCpus         uint   `form:"label=Compute Number of CPUs"`
	ComputeMemory       uint   `form:"label=Compute Memory size (GB)"`
	ComputeDiskSize     uint   `form:"label=Compute Disk size (GB)"`
	ComputeDiskCount    uint   `form:"label=Number of extra compute disks"`
}

func StandardForm2Lab(form StandardForm) lab.StandardLab {
	return lab.StandardLab{
		Name:           form.Name,
		Hypervisor:     form.Hypervisor,
		NatNet:         form.NatNet,
		LoopBackPrefix: form.LoopBackPrefix,
		IntNetPrefix:   form.IntNetPrefix,
		Oam: lab.OamInfo{
			Network:     form.Network,
			Gateway:     form.Gateway,
			FloatAddr:   form.FloatAddr,
			Controller0: form.Controller0,
			Controller1: form.Controller1,
		},
		ControllerCpus:      form.ControllerCpus,
		ControllerMemory:    form.ControllerMemory,
		ControllerDiskSize:  form.ControllerDiskSize,
		ControllerDiskCount: form.ControllerDiskCount,
		ComputeCount:        form.ComputeCount,
		ComputeCpus:         form.ComputeCpus,
		ComputeMemory:       form.ComputeMemory,
		ComputeDiskSize:     form.ComputeDiskSize,
		ComputeDiskCount:    form.ComputeDiskCount,
	}
}

func StandardLab2Form(l lab.StandardLab) StandardForm {
	return StandardForm{
		Name:                l.Name,
		Hypervisor:          l.Hypervisor,
		NatNet:              l.NatNet,
		LoopBackPrefix:      l.LoopBackPrefix,
		IntNetPrefix:        l.IntNetPrefix,
		Network:             l.Oam.Network,
		Gateway:             l.Oam.Gateway,
		FloatAddr:           l.Oam.FloatAddr,
		Controller0:         l.Oam.Controller0,
		Controller1:         l.Oam.Controller1,
		ControllerCpus:      l.ControllerCpus,
		ControllerMemory:    l.ControllerMemory,
		ControllerDiskSize:  l.ControllerDiskSize,
		ControllerDiskCount: l.ControllerDiskCount,
		ComputeCount:        l.ComputeCount,
		ComputeCpus:         l.ComputeCpus,
		ComputeMemory:       l.ComputeMemory,
		ComputeDiskSize:     l.ComputeDiskSize,
		ComputeDiskCount:    l.ComputeDiskCount,
	}
}

func handleStandard(w http.ResponseWriter, r *http.Request) {
	box := packr.New("WebTemplates", "./templates/web")
	inputTpl, err := box.FindString("input.tmpl")
	if err != nil {
		panic(err)
	}
	tpl := template.Must(template.New("").Funcs(sprig.FuncMap()).Funcs(form.FuncMap()).Parse(inputTpl))
	fb := form.Builder{
		InputTemplate: tpl,
	}
	pageTpl, err := box.FindString("page.tmpl")
	if err != nil {
		panic(err)
	}
	tpl = template.Must(template.New("").Funcs(sprig.FuncMap()).Funcs(fb.FuncMap()).Parse(pageTpl))
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "text/html")
		data := struct {
			Type   string
			Form   StandardForm
			Errors []error
		}{
			Type: "Standard",
			Form: StandardLab2Form(lab.DefaultStandardLab()),
		}
		err := tpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
		return
	case http.MethodPost:
		r.ParseForm()
		dec := schema.NewDecoder()
		dec.IgnoreUnknownKeys(true)
		var form StandardForm
		err := dec.Decode(&form, r.PostForm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		modtime := time.Now()
		var content bytes.Buffer
		if form.Hypervisor == "LibVirt" {
			err = installers.MakeStandardLibvirtInstaller(
				StandardForm2Lab(form),
				io.Writer(&content))
		} else {

			err = installers.MakeStandardInstaller(
				StandardForm2Lab(form),
				io.Writer(&content))
		}
		const name = "stx-lab-standard.tar"
		w.Header().Add("Content-Disposition",
			fmt.Sprintf("Attachment; filename=%s", name))
		http.ServeContent(w, r, name, modtime, bytes.NewReader(content.Bytes()))
	default:
		http.NotFound(w, r)
		return
	}
}

type StorageForm struct {
	Name                string `form:"label=Lab Name"`
	Hypervisor          string `form:"label=Hypervisor;type=select"`
	NatNet              string `form:"label=NAT Network"`
	LoopBackPrefix      string `form:"label=Loopback Prefix"`
	IntNetPrefix        string `form:"label=Internal Network Prefix"`
	Network             string `form:"label=OAM Network prefix & mask"`
	Gateway             string `form:"label=OAM Gateway"`
	FloatAddr           string `form:"label=OAM Floating IP address"`
	Controller0         string `form:"label=OAM Controller-0 IP address"`
	Controller1         string `form:"label=OAM Controller-1 IP address"`
	ControllerCpus      uint   `form:"label=Controller Number of CPUs"`
	ControllerMemory    uint   `form:"label=Controller Memory size (GB)"`
	ControllerDiskSize  uint   `form:"label=Controller Disk size (GB)"`
	ControllerDiskCount uint   `form:"label=Number of extra controller disks"`
	ComputeCount        uint   `form:"label=Number of Compute nodes"`
	ComputeCpus         uint   `form:"label=Compute Number of CPUs"`
	ComputeMemory       uint   `form:"label=Compute Memory size (GB)"`
	ComputeDiskSize     uint   `form:"label=Compute Disk size (GB)"`
	ComputeDiskCount    uint   `form:"label=Number of extra compute disks"`
	StorageCount        uint   `form:"label=Number of Storage nodes"`
	StorageCpus         uint   `form:"label=Storage Number of CPUs"`
	StorageMemory       uint   `form:"label=Storage Memory size (GB)"`
	StorageDiskSize     uint   `form:"label=Storage Disk size (GB)"`
	StorageDiskCount    uint   `form:"label=Number of storage disks (OSDs) per host"`
}

func StorageForm2Lab(form StorageForm) lab.StorageLab {
	return lab.StorageLab{
		Name:           form.Name,
		Hypervisor:     form.Hypervisor,
		NatNet:         form.NatNet,
		LoopBackPrefix: form.LoopBackPrefix,
		IntNetPrefix:   form.IntNetPrefix,
		Oam: lab.OamInfo{
			Network:     form.Network,
			Gateway:     form.Gateway,
			FloatAddr:   form.FloatAddr,
			Controller0: form.Controller0,
			Controller1: form.Controller1,
		},
		ControllerCpus:      form.ControllerCpus,
		ControllerMemory:    form.ControllerMemory,
		ControllerDiskSize:  form.ControllerDiskSize,
		ControllerDiskCount: form.ControllerDiskCount,
		ComputeCount:        form.ComputeCount,
		ComputeCpus:         form.ComputeCpus,
		ComputeMemory:       form.ComputeMemory,
		ComputeDiskSize:     form.ComputeDiskSize,
		ComputeDiskCount:    form.ComputeDiskCount,
		StorageCount:        form.StorageCount,
		StorageCpus:         form.StorageCpus,
		StorageMemory:       form.StorageMemory,
		StorageDiskSize:     form.StorageDiskSize,
		StorageDiskCount:    form.StorageDiskCount,
	}
}

func StorageLab2Form(l lab.StorageLab) StorageForm {
	return StorageForm{
		Name:                l.Name,
		Hypervisor:          l.Hypervisor,
		NatNet:              l.NatNet,
		LoopBackPrefix:      l.LoopBackPrefix,
		IntNetPrefix:        l.IntNetPrefix,
		Network:             l.Oam.Network,
		Gateway:             l.Oam.Gateway,
		FloatAddr:           l.Oam.FloatAddr,
		Controller0:         l.Oam.Controller0,
		Controller1:         l.Oam.Controller1,
		ControllerCpus:      l.ControllerCpus,
		ControllerMemory:    l.ControllerMemory,
		ControllerDiskSize:  l.ControllerDiskSize,
		ControllerDiskCount: l.ControllerDiskCount,
		ComputeCount:        l.ComputeCount,
		ComputeCpus:         l.ComputeCpus,
		ComputeMemory:       l.ComputeMemory,
		ComputeDiskSize:     l.ComputeDiskSize,
		ComputeDiskCount:    l.ComputeDiskCount,
		StorageCount:        l.StorageCount,
		StorageCpus:         l.StorageCpus,
		StorageMemory:       l.StorageMemory,
		StorageDiskSize:     l.StorageDiskSize,
		StorageDiskCount:    l.StorageDiskCount,
	}
}

func handleStorage(w http.ResponseWriter, r *http.Request) {
	box := packr.New("WebTemplates", "./templates/web")
	inputTpl, err := box.FindString("input.tmpl")
	if err != nil {
		panic(err)
	}
	tpl := template.Must(template.New("").Funcs(sprig.FuncMap()).Funcs(form.FuncMap()).Parse(inputTpl))
	fb := form.Builder{
		InputTemplate: tpl,
	}
	pageTpl, err := box.FindString("page.tmpl")
	if err != nil {
		panic(err)
	}
	tpl = template.Must(template.New("").Funcs(sprig.FuncMap()).Funcs(fb.FuncMap()).Parse(pageTpl))
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "text/html")
		data := struct {
			Type   string
			Form   StorageForm
			Errors []error
		}{
			Type: "Storage",
			Form: StorageLab2Form(lab.DefaultStorageLab()),
		}
		err := tpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
		return
	case http.MethodPost:
		r.ParseForm()
		dec := schema.NewDecoder()
		dec.IgnoreUnknownKeys(true)
		var form StorageForm
		err := dec.Decode(&form, r.PostForm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		modtime := time.Now()
		var content bytes.Buffer
		if form.Hypervisor == "LibVirt" {
			err = installers.MakeStorageLibvirtInstaller(
				StorageForm2Lab(form),
				io.Writer(&content))
		} else {
			err = installers.MakeStorageInstaller(
				StorageForm2Lab(form),
				io.Writer(&content))
		}
		const name = "stx-lab-storage.tar"
		w.Header().Add("Content-Disposition",
			fmt.Sprintf("Attachment; filename=%s", name))
		http.ServeContent(w, r, name, modtime, bytes.NewReader(content.Bytes()))
	default:
		http.NotFound(w, r)
		return
	}
}

func main() {
	var port uint = 3000
	flag.UintVar(&port, "port", 3000, "port to listen for HTTP connections")
	flag.Parse()

	// force packr template discovery
	_ = packr.New("VboxTemplates", "./templates/vbox")
	_ = packr.New("InstallTemplates", "./templates/install")
	_ = packr.New("WebTemplates", "./templates/web")

	http.HandleFunc("/", redirect2AioSx)
	http.HandleFunc("/aiosx", handleAioSx)
	http.HandleFunc("/aiodx", handleAioDx)
	http.HandleFunc("/standard", handleStandard)
	http.HandleFunc("/storage", handleStorage)

	fmt.Printf("Listening for HTTP connections on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

type fieldError struct {
	Field string
	Issue string
}

func (fe fieldError) Error() string {
	return fmt.Sprintf("%v: %v", fe.Field, fe.Issue)
}

func (fe fieldError) FieldError() (field, err string) {
	return fe.Field, fe.Issue
}
