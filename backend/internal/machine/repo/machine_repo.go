package machine_repo

import (
	"context"
	machine "rent-computing/internal/machine/data"
)

type MachineRepository interface {
	CreateMachine(ctx context.Context, machine machine.Machine) error
	ListMachines(ctx context.Context) ([]machine.Machine, error)
	ListMachinesByOwner(ctx context.Context, ownerId string) ([]machine.Machine, error)
	GetMachine(ctx context.Context, id string) (*machine.Machine, error)
	UpdateMachine(ctx context.Context, machine machine.Machine) error
	DeleteMachine(ctx context.Context, id string) error
	DeleteMachinesByOwner(ctx context.Context, ownerId string) error
}
