package lab

func DefaultAioSxLab() AioSxLab {
	return AioSxLab{
		Name: "aiosx",
		Hypervisor: "VirtualBox,LibVirt",
		SystemMode: "simplex",
		NatNet: "nat2",
		LoopBackPrefix: "127.0.2",
		IntNetPrefix: "intnet",
		Oam: OamInfo{
			Network: "10.10.12.0/24",
			Gateway: "10.10.12.1",
			FloatAddr: "10.10.12.2",
			Controller0: "10.10.12.3",
			Controller1: "10.10.12.4",
		},
		Cpus: 8,
		Memory: 24,
		DiskCount: 2,
		DiskSize: 520,
	}
}

func DefaultAioDxLab() AioDxLab {
	return AioDxLab{
		Name: "aiodx",
		Hypervisor: "VirtualBox,LibVirt",
		SystemMode: "duplex",
		NatNet: "nat3",
		LoopBackPrefix: "127.0.3",
		IntNetPrefix: "intnet",
		Oam: OamInfo{
			Network: "10.10.13.0/24",
			Gateway: "10.10.13.1",
			FloatAddr: "10.10.13.2",
			Controller0: "10.10.13.3",
			Controller1: "10.10.13.4",
		},
		Cpus: 8,
		Memory: 24,
		DiskCount: 2,
		DiskSize: 520,
	}
}

func DefaultStandardLab() StandardLab {
	return StandardLab{
		Name: "standard",
		Hypervisor: "VirtualBox,LibVirt",
		SystemMode: "standard",
		NatNet: "nat4",
		LoopBackPrefix: "127.0.4",
		IntNetPrefix: "intnet",
		Oam: OamInfo{
			Network: "10.10.14.0/24",
			Gateway: "10.10.14.1",
			FloatAddr: "10.10.14.2",
			Controller0: "10.10.14.3",
			Controller1: "10.10.14.4",
		},
		ControllerCpus: 4,
		ControllerMemory: 16,
		ControllerDiskCount: 2,
		ControllerDiskSize: 520,
		ComputeCount: 2,
		ComputeCpus: 4,
		ComputeMemory: 10,
		ComputeDiskCount: 2,
		ComputeDiskSize: 520,
	}
}

func DefaultStorageLab() StorageLab {
	return StorageLab{
		Name: "storage",
		Hypervisor: "VirtualBox,LibVirt",
		SystemMode: "duplex",
		NatNet: "nat5",
		LoopBackPrefix: "127.0.5",
		IntNetPrefix: "intnet",
		Oam: OamInfo{
			Network: "10.10.15.0/24",
			Gateway: "10.10.15.1",
			FloatAddr: "10.10.15.2",
			Controller0: "10.10.15.3",
			Controller1: "10.10.15.4",
		},
		ControllerCpus: 4,
		ControllerMemory: 16,
		ControllerDiskCount: 2,
		ControllerDiskSize: 520,
		ComputeCount: 2,
		ComputeCpus: 4,
		ComputeMemory: 10,
		ComputeDiskCount: 2,
		ComputeDiskSize: 520,
		StorageCount: 2,
		StorageCpus: 2,
		StorageMemory: 8,
		StorageDiskCount: 4,
		StorageDiskSize: 520,
	}
}
