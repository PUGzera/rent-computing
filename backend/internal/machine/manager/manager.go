package machine_manager

import (
	"context"
	machine "rent-computing/internal/machine/data"
)

type MachineManager interface {
	StartShell(ctx context.Context, machine machine.Machine) (string, error)
	StartVNC(ctx context.Context, machine machine.Machine) (string, error)
	Stop(ctx context.Context, machine machine.Machine) error
}
