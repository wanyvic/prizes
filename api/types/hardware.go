package types

type Hardware struct {
	CPUType         string `json:"cpu_type,omitempty"`
	CPUThread       int64  `json:"cpu_thread,omitempty"`
	MemoryType      string `json:"memory_type,omitempty"`
	MemoryCount     int64  `json:"memory_count,omitempty"` //GB
	GPUType         string `json:"gpu_type,omitempty"`
	GPUCount        int64  `json:"gpu_count,omitempty"`
	PersistentStore string `json:"persistent_store,omitempty"`
}
