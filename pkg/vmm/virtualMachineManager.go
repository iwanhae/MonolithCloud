package vmm

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
)

const (
	DefaultKernelArgs = "ro console=ttyS0 noapic reboot=k panic=1 pci=off"
)

type VirtualMachineManager struct {
	SocketDir      string
	TemplateDir    string
	DataDir        string
	CNINetworkName string
	KernelArgs     string

	events chan Message
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

func (v *VirtualMachineManager) Start(ctx context.Context) {
	log.Info().Msg("VMM is started")
	defer log.Info().Msg("VMM is terminated")
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
