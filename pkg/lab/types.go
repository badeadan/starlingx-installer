package lab

type Disk struct {
	Size   uint
	Medium string
}

type StorageController struct {
	Name        string
	Type        string
	Chipset     string
	PortCount   uint
	HostIOCache bool
	Bootable    bool
}

type StorageAttachment struct {
	Controller    string
	Port          uint
	Device        uint
	Type          string
	Medium        string
	NonRotational bool
	Discard       bool
}

type Nic struct {
	Mode     string
	Network  string
	Type     string
	BootPrio uint
	Promisc  string
}

type Vm struct {
	Name               string
	Cpus               uint
	Memory             uint
	BootOrder          []string
	Disks              []Disk
	StorageControllers []StorageController
	StorageAttachments []StorageAttachment
	Nics               []Nic
}

type Network struct {
	Name        string
	Mode        string
	Address     string
	LoopbackMap []string
	PortForward []string
}

type OamInfo struct {
	Network     string
	Gateway     string
	FloatAddr   string
	Controller0 string
	Controller1 string
}

type Lab struct {
	Type     string
	Group    string
	BasePath string
	Vms      []Vm
	Networks []Network
}

type AioSxLab struct {
	Name           string
	SystemMode     string
	NatNet         string
	LoopBackPrefix string
	IntNetPrefix   string
	Oam            OamInfo
	Cpus           uint
	Memory         uint
	DiskSize       uint
}

type AioDxLab struct {
	Name           string
	SystemMode     string
	NatNet         string
	LoopBackPrefix string
	IntNetPrefix   string
	Oam            OamInfo
	Cpus           uint
	Memory         uint
	DiskSize       uint
}

type StandardLab struct {
	Name               string
	SystemMode         string
	NatNet             string
	LoopBackPrefix     string
	IntNetPrefix       string
	Oam                OamInfo
	ControllerCpus     uint
	ControllerMemory   uint
	ControllerDiskSize uint
	ComputeCount       uint
	ComputeCpus        uint
	ComputeMemory      uint
	ComputeDiskSize    uint
}

type StorageLab struct {
	Name               string
	SystemMode         string
	NatNet             string
	LoopBackPrefix     string
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