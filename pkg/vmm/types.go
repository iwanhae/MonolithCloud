package vmm

type VirtualMachine struct {
	ID        string `json:"id"`
	IPAddress string `json:"ip"`

	Name       string `json:"name"`
	MemSizeMib int64  `json:"memory_size_mb"`
	VcpuCount  int64  `json:"vcpu_count"`

	RootfsPath  string   `json:"-"`
	VMLinuxPath string   `json:"-"`
	DataPath    string   `json:"-"`
	Status      VMStatus `json:"-"`
	SocketPath  string   `json:"-"`
}

type VMStatus string

const (
	VMStatus_Running VMStatus = "Running"
	VMStatus_Stopped VMStatus = "Stopped"
)
