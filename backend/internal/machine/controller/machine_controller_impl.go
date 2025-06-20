package machine_controller

import (
	"context"
	"errors"
	machine "rent-computing/internal/machine/data"
	machine_manager "rent-computing/internal/machine/manager"
	machine_repo "rent-computing/internal/machine/repo"
)

type MachineControllerImpl struct {
	machineRepo    machine_repo.MachineRepository
	machineManager machine_manager.MachineManager
}

type Options struct {
	MachineRepo    machine_repo.MachineRepository
	MachineManager machine_manager.MachineManager
}

func New(options Options) (*MachineControllerImpl, error) {
	return &MachineControllerImpl{
		machineRepo:    options.MachineRepo,
		machineManager: options.MachineManager,
	}, nil
}

func (c *MachineControllerImpl) Create(ctx context.Context, ownerId string, vnc bool, username, password string) (*MachineRepresentation, error) {
	machine, err := machine.New(machine.Options{OwnerId: ownerId, Vnc: vnc, Username: username, Password: password})
	if err != nil {
		return nil, err
	}

	if vnc {
		addr, err := c.machineManager.StartVNC(ctx, *machine)
		if err != nil {
			return nil, err
		}
		machine.Address = addr
	} else {
		addr, err := c.machineManager.StartShell(ctx, *machine)
		if err != nil {
			return nil, err
		}
		machine.Address = addr
	}

	err = c.machineRepo.CreateMachine(ctx, *machine)
	if err != nil {
		return nil, err
	}

	return &MachineRepresentation{
		Id:      machine.Id,
		OwnerId: machine.OwnerId,
		Address: machine.Address,
		Vnc:     machine.Vnc,
	}, nil
}

func (c *MachineControllerImpl) Delete(ctx context.Context, id, ownerId string) error {
	machine, err := c.machineRepo.GetMachine(ctx, id)
	if err != nil {
		return err
	}

	if machine.OwnerId != ownerId {
		return errors.New("not authorized to access this data")
	}

	return c.machineRepo.DeleteMachine(ctx, id)
}

func (c *MachineControllerImpl) Stop(ctx context.Context, id, ownerId string) error {
	machine, err := c.machineRepo.GetMachine(ctx, id)
	if err != nil {
		return err
	}

	if machine.OwnerId != ownerId {
		return errors.New("not authorized to access this data")
	}

	return c.machineManager.Stop(ctx, *machine)
}

func (c *MachineControllerImpl) Start(ctx context.Context, id, ownerId string) (*MachineRepresentation, error) {
	machine, err := c.machineRepo.GetMachine(ctx, id)
	if err != nil {
		return nil, err
	}

	if machine.OwnerId != ownerId {
		return nil, errors.New("not authorized to access this data")
	}

	if machine.Vnc {
		addr, err := c.machineManager.StartVNC(ctx, *machine)
		if err != nil {
			return nil, err
		}
		machine.Address = addr
	} else {
		addr, err := c.machineManager.StartShell(ctx, *machine)
		if err != nil {
			return nil, err
		}
		machine.Address = addr
	}

	return &MachineRepresentation{
		Id:      machine.Id,
		OwnerId: machine.OwnerId,
		Address: machine.Address,
		Vnc:     machine.Vnc,
	}, nil
}

func (c *MachineControllerImpl) GetMachines(ctx context.Context, ownerId string) ([]MachineRepresentation, error) {
	machines, err := c.machineRepo.ListMachinesByOwner(ctx, ownerId)
	if err != nil {
		return nil, err
	}

	var machinesRet []MachineRepresentation
	for _, machine := range machines {
		if machine.OwnerId != ownerId {
			return nil, errors.New("not authorized to access this data")
		}
		machinesRet = append(machinesRet, MachineRepresentation{
			Id:      machine.Id,
			OwnerId: machine.OwnerId,
			Address: machine.Address,
			Vnc:     machine.Vnc,
		})
	}

	return machinesRet, nil
}
