package vmm

import (
	"context"
	"fmt"
	"path"
)

var _ Message = &StartVMMessage{}

type StartVMMessage struct {
	VMID string
}

func (*StartVMMessage) MetaString() string {
	return "Start VM"
}

func (v *VirtualMachineManager) StartVM(ctx context.Context, msg *StartVMMessage) error {

	return nil
}

func (v *VirtualMachineManager) getFCSocketPath(id string) string {
	return path.Join(
		v.SocketDir,
		fmt.Sprintf("%s.sock", id))
}
