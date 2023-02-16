package vmm

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

const (
	DefaultKernelArgs = "ro console=ttyS0 noapic reboot=k panic=1 pci=off network-config=disabled"
)

type VirtualMachineManager struct {
	SocketDir      string
	TemplateDir    string
	DataDir        string
	LogDir         string
	CNINetworkName string
	KernelArgs     string

	events chan Message
	vmMap  map[ /*VMID */ string]VirtualMachine
}

type MessageHandler func(ctx context.Context, m Message) error

type Message interface {
	MetaString() string
}

func (v *VirtualMachineManager) Request(m Message) error {
	if v.events == nil {
		v.events = make(chan Message, 10)
	}
	v.events <- m
	return nil
}

func (v *VirtualMachineManager) SetVMMeta(vm VirtualMachine) error {
	v.vmMap[vm.ID] = vm
	f, err := os.OpenFile(
		path.Join(v.DataDir, vm.ID, METADATA_FILENAME),
		os.O_WRONLY|os.O_CREATE,
		0755,
	)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()
	if err := yaml.NewEncoder(f).Encode(vm); err != nil {
		return fmt.Errorf("failed to marshal vm meta: %w", err)
	}
	return nil
}

func (v *VirtualMachineManager) loadVMMeta(VMID string) (VirtualMachine, error) {
	vm := VirtualMachine{}

	f, err := os.OpenFile(
		path.Join(v.DataDir, VMID, METADATA_FILENAME),
		os.O_RDONLY, 0755)
	if err != nil {
		return vm, err
	}

	if err := yaml.NewDecoder(f).Decode(&vm); err != nil {
		return vm, err
	}

	// Default
	vm.Status = VMStatus_Stopped

	v.vmMap[VMID] = vm
	return vm, nil
}

func (v *VirtualMachineManager) GetVMMeta(VMID string) (VirtualMachine, error) {
	vm, ok := v.vmMap[VMID]
	if !ok {
		return vm, fmt.Errorf("%q VM not found", VMID)
	}
	return vm, nil
}

func (v *VirtualMachineManager) Start(ctx context.Context) {
	log.Info().Msg("VMM is started")
	defer log.Info().Msg("VMM is terminated")

	if v.vmMap == nil {
		v.vmMap = make(map[string]VirtualMachine)
	}

	for event := range v.events {
		log.Info().Str("event", event.MetaString()).Msg("new event received")
		err := func(event Message) error {
			switch evt := event.(type) {
			case *CreateVMMessage:
				return v.CreateVM(ctx, evt)
			case *StartVMMessage:
				return v.StartVM(ctx, evt)
			default:
				return fmt.Errorf("unhandled message")
			}
		}(event)

		if err != nil {
			log.Error().Err(err).Send()
		} else {
			log.Info().Msg("event handled")
		}
	}
}
