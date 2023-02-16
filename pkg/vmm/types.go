package vmm

import (
	"path"

	"github.com/iwanhae/monolithcloud/pkg/vmm/cloudinit"
)

const (
	ROOTFS_FILE_NAME = "rootfs.ext4"
	ROOTFS_DIR_NAME  = "rootfs"

	VMLINUX_FILENAME = "vmlinux.bin"
	VMLINUX_DIR_NAME = "vmlinux"

	METADATA_FILENAME = "metadata.yaml"
)

func (v *VirtualMachineManager) getVMLinuxPath(vmlinux string) string {
	return path.Join(v.TemplateDir, VMLINUX_DIR_NAME, vmlinux, VMLINUX_FILENAME)
}
func (v *VirtualMachineManager) getRootFSPath(VMID string) string {
	return path.Join(v.DataDir, VMID, ROOTFS_FILE_NAME)
}

type VirtualMachine struct {
	ID        string `json:"id"`
	IPAddress string `json:"ip"`

	Name       string `json:"name"`
	MemSizeMib int64  `json:"memory_size_mb"`
	VcpuCount  int64  `json:"vcpu_count"`
	VMLinux    string `json:"vmlinux"`

	CloudConfig cloudinit.CLoudConfig `json:"cloud_config"`
	KernelArgs  string

	Status     VMStatus `json:"-"`
	SocketPath string   `json:"-"`
}

type VMStatus string

const (
	VMStatus_Running   VMStatus = "Running"
	VMStatus_Stopped   VMStatus = "Stopped"
	VMStatus_Preparing VMStatus = "Preparing"
)
