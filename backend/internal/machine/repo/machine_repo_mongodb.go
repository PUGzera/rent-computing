package machine_repo

import (
	"context"
	machine "rent-computing/internal/machine/data"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MachineRepoMongoDB struct {
	collection *mongo.Collection
}

type OptionsMongoDB struct {
	Collection *mongo.Collection
}

func NewMongoDB(OptionsMongoDB OptionsMongoDB) (*MachineRepoMongoDB, error) {
	collection := OptionsMongoDB.Collection
	return &MachineRepoMongoDB{
		collection: collection,
	}, nil
}

//CRUD Operations

// Create
func (u *MachineRepoMongoDB) CreateMachine(ctx context.Context, machine machine.Machine) error {
	_, err := u.collection.InsertOne(ctx, machine)
	return err
}

// Read

func (u *MachineRepoMongoDB) ListMachines(ctx context.Context) ([]machine.Machine, error) {
	cursor, err := u.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var machines []machine.Machine

	for cursor.Next(ctx) {
		var machine machine.Machine
		if err := cursor.Decode(&machine); err != nil {
			return nil, err
		}
		machines = append(machines, machine)
	}

	return machines, nil
}

func (u *MachineRepoMongoDB) ListMachinesByOwner(ctx context.Context, ownerId string) ([]machine.Machine, error) {
	cursor, err := u.collection.Find(ctx, bson.M{"ownerid": ownerId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var machines []machine.Machine

	for cursor.Next(ctx) {
		var machine machine.Machine
		if err := cursor.Decode(&machine); err != nil {
			return nil, err
		}
		machines = append(machines, machine)
	}

	return machines, nil
}

func (u *MachineRepoMongoDB) GetMachine(ctx context.Context, id string) (*machine.Machine, error) {
	var machine machine.Machine
	err := u.collection.FindOne(ctx, bson.M{"id": id}).Decode(&machine)
	if err != nil {
		return nil, err
	}

	return &machine, nil
}

// Update
func (u *MachineRepoMongoDB) UpdateMachine(ctx context.Context, machine machine.Machine) error {
	_, err := u.collection.UpdateOne(ctx, bson.M{"id": machine.Id}, bson.M{"$set": machine})
	return err
}

// Delete

func (u *MachineRepoMongoDB) DeleteMachine(ctx context.Context, id string) error {
	_, err := u.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (u *MachineRepoMongoDB) DeleteMachinesByOwner(ctx context.Context, ownerId string) error {
	_, err := u.collection.DeleteMany(ctx, bson.M{"ownerId": ownerId})
	return err
}
