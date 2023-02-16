package vmm

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	"github.com/iwanhae/monolithcloud/pkg/vmm/cloudinit"
)

var _ Message = &StartVMMessage{}

type StartVMMessage struct {
	VMID string
}

func (*StartVMMessage) MetaString() string {
	return "Start VM"
}

func (v *VirtualMachineManager) StartVM(ctx context.Context, msg *StartVMMessage) error {
	// Load VM Meta
	vm, err := v.GetVMMeta(msg.VMID)
	if err != nil {
		return fmt.Errorf("failed to load VM: %w", err)
	}
	vm.SocketPath = path.Join(v.SocketDir, fmt.Sprintf("%s.sock", vm.ID))

	// Prepare CloudConfig
	cloudConfigPath := path.Join(v.DataDir, vm.ID, "cloudconfig.iso")
	if err := cloudinit.GenerateCloudConfigDisk(ctx, vm.CloudConfig, cloudConfigPath); err != nil {
		return err
	}

	// Prepare Config
	CNIConf := &firecracker.CNIConfiguration{NetworkName: v.CNINetworkName, IfName: "eth0"}
	if vm.IPAddress != "" {
		CNIConf.Args = [][2]string{{"IP", vm.IPAddress}, {"IgnoreUnknown", "True"}}
	}
	c := firecracker.Config{
		SocketPath:      vm.SocketPath,
		KernelImagePath: v.getVMLinuxPath(vm.VMLinux),
		KernelArgs:      vm.KernelArgs,
		MachineCfg: models.MachineConfiguration{
			VcpuCount:  firecracker.Int64(vm.VcpuCount),
			MemSizeMib: firecracker.Int64(vm.MemSizeMib),
		},
		Drives: []models.Drive{
			{
				DriveID:      firecracker.String("1"),
				PathOnHost:   firecracker.String(v.getRootFSPath(vm.ID)),
				IsReadOnly:   firecracker.Bool(false),
				IsRootDevice: firecracker.Bool(true),
			},
			{
				DriveID:      firecracker.String("2"),
				PathOnHost:   firecracker.String(cloudConfigPath),
				IsReadOnly:   firecracker.Bool(true),
				IsRootDevice: firecracker.Bool(false),
			},
		},
		LogLevel:          "error",
		NetworkInterfaces: firecracker.NetworkInterfaces{{CNIConfiguration: CNIConf}},
	}

	fcBin, err := exec.LookPath("firecracker")
	if err != nil {
		return err
	}

	// Open Log Files
	stdout, err := os.OpenFile(
		path.Join(v.LogDir, fmt.Sprintf("%s_stdout.log", vm.ID)),
		os.O_WRONLY|os.O_CREATE, 0755,
	)
	if err != nil {
		return fmt.Errorf("fail to open stdout file: %w", err)
	}
	stderr, err := os.OpenFile(
		path.Join(v.LogDir, fmt.Sprintf("%s_stderr.log", vm.ID)),
		os.O_WRONLY|os.O_CREATE, 0755,
	)
	if err != nil {
		return fmt.Errorf("fail to open stderr file: %w", err)
	}

	cmd := firecracker.VMCommandBuilder{}.
		WithBin(fcBin).
		WithSocketPath(vm.SocketPath).
		WithStdout(stdout).
		WithStderr(stderr).
		Build(ctx)
	m, err := firecracker.NewMachine(ctx, c, firecracker.WithProcessRunner(cmd))
	if err != nil {
		return fmt.Errorf("failed creating machine: %w", err)
	}

	if err := m.Start(ctx); err != nil {
		return err
	}

	go func() {
		m.Wait(ctx)
		stdout.Close()
		stderr.Close()
	}()

	return nil
}

func (v *VirtualMachineManager) getFCSocketPath(id string) string {
	return path.Join(
		v.SocketDir,
		fmt.Sprintf("%s.sock", id))
}
