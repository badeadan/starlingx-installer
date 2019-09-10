package lab

func DefaultOamInfo() OamInfo {
	return OamInfo{
		Network: "10.10.10.0/24",
		Gateway: "10.10.10.1",
		FloatAddr: "10.10.10.2",
		Controller0: "10.10.10.3",
		Controller1: "10.10.10.4",
	}
}

func DefaultAioSxLab() AioSxLab {
	return AioSxLab{
		Name: "aiosx",
		SystemMode: "simplex",
		NatNet: "nat2",
		LoopBackPrefix: "127.0.2",
		IntNetPrefix: "intnet",
		Oam: DefaultOamInfo(),
		Cpus: 8,
		Memory: 24,
		DiskCount: 2,
		DiskSize: 520,
	}
}

func DefaultAioDxLab() AioDxLab {
	return AioDxLab{
		Name: "aiodx",
		SystemMode: "duplex",
		NatNet: "nat3",
		LoopBackPrefix: "127.0.3",
		IntNetPrefix: "intnet",
		Oam: DefaultOamInfo(),
		Cpus: 8,
		Memory: 24,
		DiskCount: 2,
		DiskSize: 520,
	}
}

func DefaultStandardLab() StandardLab {
	return StandardLab{
		Name: "standard",
		SystemMode: "standard",
		NatNet: "nat4",
		LoopBackPrefix: "127.0.4",
		IntNetPrefix: "intnet",
		Oam: DefaultOamInfo(),
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
		SystemMode: "duplex",
		NatNet: "nat5",
		LoopBackPrefix: "127.0.5",
		IntNetPrefix: "intnet",
		Oam: DefaultOamInfo(),
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
