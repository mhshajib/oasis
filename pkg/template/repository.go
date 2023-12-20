package cli_template

var Repository string = `
package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"{{.ModuleName}}/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TimeStamp struct {
	CreatedAt time.Time  ` + "`json:\"created_at\"` " + `
	UpdatedAt time.Time  ` + "`json:\"updated_at\"` " + `
	DeletedAt *time.Time ` + "`json:\"deleted_at\"` " + `
}

type {{.UcFirstName}} struct {
	ID       primitive.ObjectID ` + "`bson:\"_id,omitempty\"` " + `
	FieldOne string             ` + "`bson:\"field_one\"` " + `
	TimeStamp
}

// {{.UcFirstName}}Mongo represents mongo implementation of {{.SmallName}} repository contract
type {{.UcFirstName}}Mongo struct {
	db *mongo.Database
	c  *mongo.Collection
}

// New{{.UcFirstName}}Mongo return a mongo implementation of {{.SmallName}} storage repository
func New{{.UcFirstName}}Mongo(db *mongo.Database) domain.{{.UcFirstName}}Repository {
	return &{{.UcFirstName}}Mongo{
		db: db,
		c:  db.Collection("{{.SmallPluralName}}"),
	}
}

// Store insert a new {{.SmallName}} to mongodb
func (r *{{.UcFirstName}}Mongo) Store(ctx context.Context, {{.SmallName}} *domain.{{.UcFirstName}}) error {
	{{.SmallName}}Data := {{.UcFirstName}}{
		
	}
	result, err := r.c.InsertOne(ctx, {{.SmallName}})
	if err != nil {
		return fmt.Errorf("repository:mongo: failed to Store {{.SmallName}}: %v", err)
	}
	{{.SmallName}}.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

// Fetch list {{.SmallName}} from mongodb based on criteria
func (r *{{.UcFirstName}}Mongo) Fetch(ctx context.Context, ctr *domain.{{.UcFirstName}}Criteria) ([]*domain.{{.UcFirstName}}, error) {
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

	var {{.SmallName}}List = make([]*domain.{{.UcFirstName}}, 0)
	cursor, err := r.c.Find(ctx, filter, &opts)
	if err != nil {
		return nil, fmt.Errorf("repository:mongo: failed to Fetch {{.SmallName}}: %v", err)
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &{{.SmallName}}List); err != nil {
		return nil, fmt.Errorf("repository:mongo: failed to decode {{.SmallName}}: %v", err)
	}

	return {{.SmallName}}List, nil
}

// Count return the total {{.SmallName}} count from mongodb based on criteria
func (r *{{.UcFirstName}}Mongo) Count(ctx context.Context, ctr *domain.{{.UcFirstName}}Criteria) (int64, error) {
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
		return 0, fmt.Errorf("repository:mongo: failed to Count {{.SmallName}}: %v", err)
	}

	return count, nil
}

// FetchOne fetch a {{.SmallName}} based on criteria
func (r *{{.UcFirstName}}Mongo) FetchOne(ctx context.Context, ctr *domain.{{.UcFirstName}}Criteria) (*domain.{{.UcFirstName}}, error) {
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

	var {{.SmallName}} domain.{{.UcFirstName}}
	if err := r.c.FindOne(ctx, filter).Decode(&{{.SmallName}}); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.Err{{.UcFirstName}}NotFound
		}
		return nil, fmt.Errorf("repository:mongo: failed to FetchOne {{.SmallName}}: %v", err)
	}
	return &{{.SmallName}}, nil
}

// Update update a {{.SmallName}} record
func (r *{{.UcFirstName}}Mongo) Update(ctx context.Context, {{.SmallName}} *domain.{{.UcFirstName}}) error {
	if {{.SmallName}}.ID.Hex() == "" {
		return errors.New("repository:mongodb: Update failed: {{.SmallName}} id required")
	}
	filter := bson.M{"_id": {{.SmallName}}.ID}
	bb, err := bson.Marshal({{.SmallName}})
	if err != nil {
		return fmt.Errorf("repository:mongodb: Update failed: %v", err)
	}

	var update bson.M
	if err := bson.Unmarshal(bb, &update); err != nil {
		return fmt.Errorf("repository:mongodb: Update failed: %v", err)
	}

	result, err := r.c.UpdateOne(ctx, filter, bson.D{{"{{"}}Key: "$set", Value: update{{"}}"}})
	if err != nil {
		return fmt.Errorf("repository:mongodb: Update failed: %v", err)
	}

	if result.ModifiedCount == 0 {
		return errors.New("repository:mongodb: Update failed: 0 document modified")
	}

	return nil
}

// Delete soft delete a {{.SmallName}} record
func (r *{{.UcFirstName}}Mongo) Delete(ctx context.Context, ctr *domain.{{.UcFirstName}}Criteria) error {
	filter := bson.M{}
	if ctr.ID != nil {
		objectID, _ := primitive.ObjectIDFromHex(*ctr.ID)
		filter["_id"] = objectID
	}
	if ctr.FieldOne != nil {
		filter["field_one"] = ctr.FieldOne
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"timestamp.deleted_at": &now,
		},
	}

	result, err := r.c.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("repository:mongo: failed to Delete {{.SmallName}}: %v", err)
	}

	if result.ModifiedCount == 0 {
		return errors.New("repository:mongo: failed to Delete {{.SmallName}}, 0 document modified")
	}

	return nil
}
`
