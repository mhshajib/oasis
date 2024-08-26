
package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"oasis/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TimeStamp struct {
	CreatedAt time.Time  `bson:"created_at"` 
	UpdatedAt time.Time  `bson:"updated_at"` 
	DeletedAt *time.Time `bson:"deleted_at"` 
}

type  struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"` 
	FieldOne string             `bson:"field_one"` 
	TimeStamp
}

// Mongo represents mongo implementation of  repository contract
type Mongo struct {
	db *mongo.Database
	c  *mongo.Collection
}

// NewMongo return a mongo implementation of  storage repository
func NewMongo(db *mongo.Database) domain.Repository {
	return &Mongo{
		db: db,
		c:  db.Collection(""),
	}
}

// Store insert a new  to mongodb
func (r *Mongo) Store(ctx context.Context,  *domain.) (*domain., error) {
	Data := {
		FieldOne:     .FieldOne,
		TimeStamp: TimeStamp{
			CreatedAt: .CreatedAt,
			UpdatedAt: .UpdatedAt,
			DeletedAt: .DeletedAt,
		},
	}
	result, err := r.c.InsertOne(ctx, Data)
	if err != nil {
		return nil, fmt.Errorf("repository:mongo: failed to Store : %v", err)
	}
	.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return , nil
}

// Fetch list  from mongodb based on criteria
func (r *Mongo) Fetch(ctx context.Context, ctr *domain.Criteria) ([]*domain., error) {
	filter := bson.M{}
	if ctr.ID != nil {
		objectID, _ := primitive.ObjectIDFromHex(*ctr.ID)
		filter["_id"] = objectID
	}
	if ctr.FieldOne != nil {
		filter["field_one"] = ctr.FieldOne
	}

	if ctr.WithDeleted != nil {
		if *ctr.WithDeleted {
			filter["timestamp.deleted_at"] = bson.M{"$ne": nil}
		} else {
			filter["timestamp.deleted_at"] = bson.M{"$eq": nil}
		}
	}

	opts := options.FindOptions{}
	if ctr.Offset != nil {
		opts.SetSkip(*ctr.Offset)
	}
	if ctr.Limit != nil {
		opts.SetLimit(*ctr.Limit)
	}

	if ctr.SortAsc {
		opts.SetSort(bson.M{"timestamp.created_at": 1})
	} else {
		opts.SetSort(bson.M{"timestamp.created_at": -1})
	}

	var List = make([]*, 0)
	cursor, err := r.c.Find(ctx, filter, &opts)
	if err != nil {
		return nil, fmt.Errorf("repository:mongo: failed to Fetch : %v", err)
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &List); err != nil {
		return nil, fmt.Errorf("repository:mongo: failed to decode : %v", err)
	}

	return convertList(List), nil
}

// Count return the total  count from mongodb based on criteria
func (r *Mongo) Count(ctx context.Context, ctr *domain.Criteria) (int64, error) {
	filter := bson.M{}
	if ctr.ID != nil {
		objectID, _ := primitive.ObjectIDFromHex(*ctr.ID)
		filter["_id"] = objectID
	}
	if ctr.FieldOne != nil {
		filter["field_one"] = ctr.FieldOne
	}

	if ctr.WithDeleted != nil {
		if *ctr.WithDeleted {
			filter["timestamp.deleted_at"] = bson.M{"$ne": nil}
		} else {
			filter["timestamp.deleted_at"] = bson.M{"$eq": nil}
		}
	}

	opts := options.FindOptions{}
	if ctr.Offset != nil {
		opts.SetSkip(*ctr.Offset)
	}
	if ctr.Limit != nil {
		opts.SetLimit(*ctr.Limit)
	}

	count, err := r.c.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("repository:mongo: failed to Count : %v", err)
	}

	return count, nil
}

// FetchOne fetch a  based on criteria
func (r *Mongo) FetchOne(ctx context.Context, ctr *domain.Criteria) (*domain., error) {
	filter := bson.M{}
	if ctr.ID != nil {
		objectID, _ := primitive.ObjectIDFromHex(*ctr.ID)
		filter["_id"] = objectID
	}
	if ctr.FieldOne != nil {
		filter["field_one"] = ctr.FieldOne
	}

	if ctr.WithDeleted != nil {
		if *ctr.WithDeleted {
			filter["timestamp.deleted_at"] = bson.M{"$ne": nil}
		} else {
			filter["timestamp.deleted_at"] = bson.M{"$eq": nil}
		}
	}

	var  
	if err := r.c.FindOne(ctx, filter).Decode(&); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("repository:mongo: failed to FetchOne : %v", err)
	}
	return convert(&), nil
}

// Update update a  record
func (r *Mongo) Update(ctx context.Context,  *domain.) (*domain., error) {
	if .ID == "" {
		return nil, errors.New("repository:mongodb: Update failed:  id required")
	}
	objectId, err := primitive.ObjectIDFromHex(.ID)
	if err != nil {
		return nil, errors.New("repository:mongodb: Update failed:  primary id unable to convert from hex to ObjectID")
	}

	filter := bson.M{"_id": objectId}
	Data := {
		ID:           objectId,
		FieldOne:     .FieldOne,
		TimeStamp: TimeStamp{
			CreatedAt: .CreatedAt,
			UpdatedAt: .UpdatedAt,
			DeletedAt: .DeletedAt,
		},
	}

	bb, err := bson.Marshal(Data)
	if err != nil {
		return nil, fmt.Errorf("repository:mongodb: Update failed: %v", err)
	}

	var update bson.M
	if err := bson.Unmarshal(bb, &update); err != nil {
		return nil, fmt.Errorf("repository:mongodb: Update failed: %v", err)
	}

	result, err := r.c.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return nil, fmt.Errorf("repository:mongodb: Update failed: %v", err)
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("repository:mongodb: Update failed: 0 document modified")
	}

	return , nil
}

// Delete soft delete a  record
func (r *Mongo) Delete(ctx context.Context, ctr *domain.Criteria) error {
	if ctr.ID == nil {
		return errors.New("repository:mongodb: Delte failed:  primary id required")
	}

	objectId, err := primitive.ObjectIDFromHex(*ctr.ID)
	if err != nil {
		return errors.New("repository:mongodb: Delete failed:  primary id unable to convert from hex to ObjectID")
	}

	filter := bson.M{"_id": objectId}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"timestamp.deleted_at": &now,
		},
	}

	result, err := r.c.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("repository:mongo: failed to Delete : %v", err)
	}

	if result.ModifiedCount == 0 {
		return errors.New("repository:mongo: failed to Delete , 0 document modified")
	}

	return nil
}

// Convert ...
func convert(c *) *domain. {
	return &domain.{
		ID:           c.ID.Hex(),
		FieldOne:     c.FieldOne,
		TimeStamp: domain.TimeStamp{
			CreatedAt: c.TimeStamp.CreatedAt,
			UpdatedAt: c.TimeStamp.UpdatedAt,
			DeletedAt: c.TimeStamp.DeletedAt,
		},
	}
}

// ConvertList ...
func convertList(cl []*) []*domain. {
	list := make([]*domain., 0)
	for _, c := range cl {
		list = append(list, convert(c))
	}
	return list
}
