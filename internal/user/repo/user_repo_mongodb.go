package user_repo

import (
	"context"
	user "rent-computing/internal/user/data"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepoMongoDB struct {
	collection *mongo.Collection
}

type OptionsMongoDB struct {
	Collection *mongo.Collection
}

func NewMongoDB(OptionsMongoDB OptionsMongoDB) (*UserRepoMongoDB, error) {
	collection := OptionsMongoDB.Collection
	return &UserRepoMongoDB{
		collection: collection,
	}, nil
}

//CRUD Operations

// Create
func (u *UserRepoMongoDB) CreateUser(ctx context.Context, user user.User) error {
	_, err := u.collection.InsertOne(ctx, user)
	return err
}

// Read

func (u *UserRepoMongoDB) ListUsers(ctx context.Context) ([]user.User, error) {
	cursor, err := u.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []user.User

	for cursor.Next(ctx) {
		var user user.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *UserRepoMongoDB) GetUser(ctx context.Context, id string) (*user.User, error) {
	var user user.User
	err := u.collection.FindOne(ctx, bson.M{"id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepoMongoDB) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	var user user.User
	err := u.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepoMongoDB) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	var user user.User
	err := u.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update
func (u *UserRepoMongoDB) UpdateUser(ctx context.Context, user user.User) error {
	_, err := u.collection.UpdateOne(ctx, bson.M{"username": user.Username}, bson.M{"$set": user})
	return err
}

// Delete

func (u *UserRepoMongoDB) DeleteUser(ctx context.Context, id string) error {
	_, err := u.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (u *UserRepoMongoDB) DeleteUserByUsername(ctx context.Context, username string) error {
	_, err := u.collection.DeleteOne(ctx, bson.M{"username": username})
	return err
}

func (u *UserRepoMongoDB) DeleteUserByEmail(ctx context.Context, email string) error {
	_, err := u.collection.DeleteOne(ctx, bson.M{"email": email})
	return err
}
