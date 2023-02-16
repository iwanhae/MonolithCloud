package vmm

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
)

var _ Message = &CreateVMMessage{}

type CreateVMMessage struct {
	Name       string
	MemSizeMib int64
	VcpuCount  int64
	StorageGB  int64

	RootfsTemplate  string
	VmlinuxTemplate string
}

func (*CreateVMMessage) MetaString() string {
	return "Create VM"
}

func (v *VirtualMachineManager) CreateVM(ctx context.Context, msg *CreateVMMessage) error {
	id, err := v.generateVMID()
	if err != nil {
		return fmt.Errorf("failed to generate a new VMID: %w", err)
	}

	dataPath := path.Join(v.DataDir, id)
	if err := os.Mkdir(dataPath, 0755); err != nil {
		return fmt.Errorf("failed to generate a base dir: %w", err)
	}

	vmlinuxPath := path.Join(v.TemplateDir, "vmlinux", msg.VmlinuxTemplate)
	rootfsPath := path.Join(v.DataDir, id, "rootfs.ext4")
	if _, err := copy(
		path.Join(v.TemplateDir, "rootfs", msg.RootfsTemplate, "rootfs.ext4"),
		rootfsPath,
	); err != nil {
		return fmt.Errorf("copy failed: %w", err)
	}

	fmt.Println(vmlinuxPath, rootfsPath)

	return nil
}

func (v *VirtualMachineManager) generateVMID() (string, error) {
	entries, err := os.ReadDir(v.DataDir)
	if err != nil {
		return "", err
	}
	max := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		v, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}
		if max < v {
			max = v
		}
	}
	return fmt.Sprintf("%05d", max+1), nil
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
