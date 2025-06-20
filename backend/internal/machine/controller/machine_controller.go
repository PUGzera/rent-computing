package machine_controller

import "context"

type MachineController interface {
	Create(ctx context.Context, ownerId string, vnc bool, username, password string) (*MachineRepresentation, error)
	Delete(ctx context.Context, id, ownerId string) error
	Start(ctx context.Context, id, ownerId string) (*MachineRepresentation, error)
	Stop(ctx context.Context, id, ownerId string) error
	GetMachines(ctx context.Context, id string) ([]MachineRepresentation, error)
}

type MachineRepresentation struct {
	Id      string `json:"id"`
	OwnerId string `json:"ownerid"`
	Vnc     bool   `json:"vnc"`
	Address string `json:"address"`
}
