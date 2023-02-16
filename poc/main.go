package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	firecracker "github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	log "github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("Hello World")
	ctx := context.Background()
	CloudInitMain(ctx)
	// return
	if err := vmMain(ctx); err != nil {
		fmt.Println("ERROR", err)
		panic(err)
	}
}

func vmMain(ctx context.Context) error {

	machineOpts := []firecracker.Opt{
		firecracker.WithLogger(log.NewEntry(log.New())),
	}

	firecrackerBinary, err := exec.LookPath("firecracker")
	if err != nil {
		return err
	}

	cmd := firecracker.VMCommandBuilder{}.
		WithBin(firecrackerBinary).
		WithSocketPath("./tmp2.sock").
		WithStdin(os.Stdin).
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		Build(ctx)

	machineOpts = append(machineOpts, firecracker.WithProcessRunner(cmd))

	f, err := os.OpenFile("log", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	m, err := firecracker.NewMachine(ctx, firecracker.Config{
		SocketPath:      "./tmp2.sock",
		KernelImagePath: "./templates/vmlinux/5.10.156/vmlinux.bin",
		KernelArgs:      "ro console=ttyS0 noapic reboot=k panic=1 pci=off network-config=disabled",
		MachineCfg: models.MachineConfiguration{
			MemSizeMib: firecracker.Int64(2048),
			VcpuCount:  firecracker.Int64(6),
		},
		Drives: []models.Drive{
			{
				DriveID:      firecracker.String("1"),
				PathOnHost:   firecracker.String("./templates/rootfs/ubuntu_2204/rootfs.ext4"),
				IsReadOnly:   firecracker.Bool(false),
				IsRootDevice: firecracker.Bool(true),
			},
			{
				DriveID:      firecracker.String("2"),
				PathOnHost:   firecracker.String("./output.iso"),
				IsReadOnly:   firecracker.Bool(true),
				IsRootDevice: firecracker.Bool(false),
			},
		},
		LogLevel: "error",
		NetworkInterfaces: firecracker.NetworkInterfaces{
			{
				CNIConfiguration: &firecracker.CNIConfiguration{
					NetworkName: "fcnet",
					IfName:      "eth0",
					Args:        [][2]string{{"IP", "192.168.0.150"}, {"IgnoreUnknown", "True"}},
				},
			},
		},
	}, machineOpts...)
	if err != nil {
		return fmt.Errorf("failed creating machine: %w", err)
	}
	if err := m.Start(ctx); err != nil {
		return err
	}

	/*
		if err := m.SetMetadata(ctx, map[string]string{
			"local-hostname": "HelloWorld",
		}); err != nil {
			return err
		}
	*/

	if err := m.Wait(ctx); err != nil {
		return fmt.Errorf("wait returned an error %w", err)
	}

	return nil
}

func installSignalHandlers(ctx context.Context, m *firecracker.Machine) {
	go func() {
		// Clear some default handlers installed by the firecracker SDK:
		signal.Reset(os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

		for {
			switch s := <-c; {
			case s == syscall.SIGTERM || s == os.Interrupt:
				log.Printf("Caught signal: %s, requesting clean shutdown", s.String())
				if err := m.Shutdown(ctx); err != nil {
					log.Errorf("An error occurred while shutting down Firecracker VM: %v", err)
				}
			case s == syscall.SIGQUIT:
				log.Printf("Caught signal: %s, forcing shutdown", s.String())
				if err := m.StopVMM(); err != nil {
					log.Errorf("An error occurred while stopping Firecracker VMM: %v", err)
				}
			}
		}
	}()
}
