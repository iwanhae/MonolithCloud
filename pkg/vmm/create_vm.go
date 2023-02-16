package vmm

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"

	"github.com/iwanhae/monolithcloud/pkg/vmm/cloudinit"
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

	// Create Base Dir for VM
	dataPath := path.Join(v.DataDir, id)
	if err := os.Mkdir(dataPath, 0755); err != nil {
		return fmt.Errorf("failed to generate a base dir: %w", err)
	}

	// COPY rootfs to datadir
	rootfsPath := path.Join(v.DataDir, id, ROOTFS_FILE_NAME)
	if _, err := copy(
		path.Join(v.TemplateDir, ROOTFS_DIR_NAME, msg.RootfsTemplate, ROOTFS_FILE_NAME),
		rootfsPath,
	); err != nil {
		return fmt.Errorf("copy failed: %w", err)
	}

	// Create Cloud Config
	vm := VirtualMachine{
		ID:         id,
		Name:       msg.Name,
		MemSizeMib: msg.MemSizeMib,
		VcpuCount:  msg.VcpuCount,
		VMLinux:    msg.VmlinuxTemplate,
		Status:     VMStatus_Stopped,
		CloudConfig: cloudinit.CLoudConfig{
			Hostname:         msg.Name,
			DisableRoot:      false,
			PreserveHostname: false,
			SystemInfo: cloudinit.SystemInfoConfig{
				Distro: "ubuntu",
			},
			Users: []cloudinit.UserCoinfig{
				{
					Name: "root",
					SSHAuthorizedKeys: []string{
						"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDPyp+d7AKdQaV8zoO+iPV8kTlXMCSpcD7BR0kPc/S8UAyj8p2uVQWPCpJNo6qkK0g+hGZ1W4J9MihBAOCvL9DAd+XbUJr4bkkc7es3AkQ8OgmP+PR2zGR76cDyZ64tdiyYanWICm3spHAf2yusaJl9ETEletwSSrHLu25aeXD06Hu1y3rxqF1jl95Mu9HuBG7mksIRrlYExtlytl48NeO+xNFlslcgmoomowbrppnjgAcO7MfgjPnOcYjnYOhea+7cYG/ZFCwwstJ9ZP3zX2HY+4D+tuQeTNZsdFpu/JhaYZEshLopF1hB183kzpp6zrv+NisyHOv51cjZa64LfSZyvW8+tBrH+2oUjF5pMB4BXkMBVjpZjhlpDPgkmAcNV9ViE3Rp+Ov19Q3BUI195kw9yjdwlBH8lmxVs+f4loXaOz1cwS87zcAUbTZP2FLzoxNim+V1mqYCtmIoghpx0aeDoqknYpIWxlMBzvGylwMWBDUDDplsV8XN5ZHMjnBdoOc= puppy@LENOVO-YOGA",
					},
				},
			},
			GrowPartition: cloudinit.GrowPartitionConfig{
				Mode:    cloudinit.GrowPartitionMode_Auto,
				Devices: []string{"/"},
			},
		},
		KernelArgs: DefaultKernelArgs,
	}

	if err := v.SetVMMeta(vm); err != nil {
		return fmt.Errorf("failed to update vm meta: %w", err)
	}

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
