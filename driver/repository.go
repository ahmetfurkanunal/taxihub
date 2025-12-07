package driver

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	Collection *mongo.Collection
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		Collection: client.Database("taxihub").Collection("drivers"),
	}
}

func (r *Repository) Create(ctx context.Context, d Driver) (primitive.ObjectID, error) {
	res, err := r.Collection.InsertOne(ctx, d)
	if err != nil {
		return primitive.NilObjectID, err
	}

	id := res.InsertedID.(primitive.ObjectID)
	return id, nil
}

func (r *Repository) Update(ctx context.Context, id primitive.ObjectID, d Driver) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": d}

	_, err := r.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *Repository) List(ctx context.Context) ([]Driver, error) {
	cur, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return []Driver{}, err
	}

	list := make([]Driver, 0)
	if err := cur.All(ctx, &list); err != nil {
		return []Driver{}, err
	}

	return list, nil
}

func (r *Repository) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
