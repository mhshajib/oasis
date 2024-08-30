package cli_template

var Repository string = `
package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"{{.DomainPath}}"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TimeStamp struct {
	CreatedAt time.Time  ` + "`bson:\"created_at\"` " + `
	UpdatedAt time.Time  ` + "`bson:\"updated_at\"` " + `
	DeletedAt *time.Time ` + "`bson:\"deleted_at\"` " + `
}

type {{.UcFirstName}} struct {
	ID       primitive.ObjectID ` + "`bson:\"_id,omitempty\"` " + ` {{range .Fields}}
    {{.Name}}    {{.Type}}          ` + "`json:\"{{.JsonTag}}\"`" + ` {{end}}
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
func (r *{{.UcFirstName}}Mongo) Store(ctx context.Context, {{.SmallName}} *domain.{{.UcFirstName}}) (*domain.{{.UcFirstName}}, error) {
	{{.SmallName}}Data := {{.UcFirstName}}{ {{range .Fields}}
    	{{.Name}}:    {{.SmallName}}.{{.Name}}, {{end}}
		TimeStamp: TimeStamp{
			CreatedAt: {{.SmallName}}.CreatedAt,
			UpdatedAt: {{.SmallName}}.UpdatedAt,
			DeletedAt: {{.SmallName}}.DeletedAt,
		},
	}
	result, err := r.c.InsertOne(ctx, {{.SmallName}}Data)
	if err != nil {
		return nil, fmt.Errorf("repository:mongo: failed to Store {{.SmallName}}: %v", err)
	}
	{{.SmallName}}.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return {{.SmallName}}, nil
}

// Fetch list {{.SmallName}} from mongodb based on criteria
func (r *{{.UcFirstName}}Mongo) Fetch(ctx context.Context, ctr *domain.{{.UcFirstName}}Criteria) ([]*domain.{{.UcFirstName}}, error) {
	filter := bson.M{}
	if ctr.ID != nil {
		objectID, _ := primitive.ObjectIDFromHex(*ctr.ID)
		filter["_id"] = objectID
	}
	
	{{range .CriteriaFields}}
	if {{if eq .Type "*string"}}ctr.{{.Name}} != nil{{else if or (eq .Type "*int") (eq .Type "*float64") (eq .Type "*float32")}}*ctr.{{.Name}} >= 0{{else if or (eq .Type "[]string") (eq .Type "[]*string") (eq .Type "[]int") (eq .Type "[]*int")}}len(ctr.{{.Name}}) > 0{{end}} {
		filter["{{.JsonTag}}"] = ctr.{{.Name}}
	}
	{{end}}

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

	var {{.SmallName}}List = make([]*{{.UcFirstName}}, 0)
	cursor, err := r.c.Find(ctx, filter, &opts)
	if err != nil {
		return nil, fmt.Errorf("repository:mongo: failed to Fetch {{.SmallName}}: %v", err)
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &{{.SmallName}}List); err != nil {
		return nil, fmt.Errorf("repository:mongo: failed to decode {{.SmallName}}: %v", err)
	}

	return convert{{.UcFirstName}}List({{.SmallName}}List), nil
}

// Count return the total {{.SmallName}} count from mongodb based on criteria
func (r *{{.UcFirstName}}Mongo) Count(ctx context.Context, ctr *domain.{{.UcFirstName}}Criteria) (int64, error) {
	filter := bson.M{}
	if ctr.ID != nil {
		objectID, _ := primitive.ObjectIDFromHex(*ctr.ID)
		filter["_id"] = objectID
	}
	{{range .CriteriaFields}}
	if {{if eq .Type "*string"}}ctr.{{.Name}} != nil{{else if or (eq .Type "*int") (eq .Type "*float64") (eq .Type "*float32")}}*ctr.{{.Name}} >= 0{{else if or (eq .Type "[]string") (eq .Type "[]*string") (eq .Type "[]int") (eq .Type "[]*int")}}len(ctr.{{.Name}}) > 0{{end}} {
		filter["{{.JsonTag}}"] = ctr.{{.Name}}
	}
	{{end}}

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
	{{range .CriteriaFields}}
	if {{if eq .Type "*string"}}ctr.{{.Name}} != nil{{else if or (eq .Type "*int") (eq .Type "*float64") (eq .Type "*float32")}}*ctr.{{.Name}} >= 0{{else if or (eq .Type "[]string") (eq .Type "[]*string") (eq .Type "[]int") (eq .Type "[]*int")}}len(ctr.{{.Name}}) > 0{{end}} {
		filter["{{.JsonTag}}"] = ctr.{{.Name}}
	}
	{{end}}

	if ctr.WithDeleted != nil {
		if *ctr.WithDeleted {
			filter["timestamp.deleted_at"] = bson.M{"$ne": nil}
		} else {
			filter["timestamp.deleted_at"] = bson.M{"$eq": nil}
		}
	}

	var {{.SmallName}} {{.UcFirstName}}
	if err := r.c.FindOne(ctx, filter).Decode(&{{.SmallName}}); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.Err{{.UcFirstName}}NotFound
		}
		return nil, fmt.Errorf("repository:mongo: failed to FetchOne {{.SmallName}}: %v", err)
	}
	return convert{{.UcFirstName}}(&{{.SmallName}}), nil
}

// Update update a {{.SmallName}} record
func (r *{{.UcFirstName}}Mongo) Update(ctx context.Context, {{.SmallName}} *domain.{{.UcFirstName}}) (*domain.{{.UcFirstName}}, error) {
	if {{.SmallName}}.ID == "" {
		return nil, errors.New("repository:mongodb: Update failed: {{.SmallName}} id required")
	}
	objectId, err := primitive.ObjectIDFromHex({{.SmallName}}.ID)
	if err != nil {
		return nil, errors.New("repository:mongodb: Update failed: {{.SmallName}} primary id unable to convert from hex to ObjectID")
	}

	filter := bson.M{"_id": objectId}
	{{.SmallName}}Data := {{.UcFirstName}}{
		ID:           objectId, {{range .Fields}}
    	{{.Name}}:    {{.SmallName}}.{{.Name}}, {{end}}
		TimeStamp: TimeStamp{
			CreatedAt: {{.SmallName}}.CreatedAt,
			UpdatedAt: {{.SmallName}}.UpdatedAt,
			DeletedAt: {{.SmallName}}.DeletedAt,
		},
	}

	bb, err := bson.Marshal({{.SmallName}}Data)
	if err != nil {
		return nil, fmt.Errorf("repository:mongodb: Update failed: %v", err)
	}

	var update bson.M
	if err := bson.Unmarshal(bb, &update); err != nil {
		return nil, fmt.Errorf("repository:mongodb: Update failed: %v", err)
	}

	result, err := r.c.UpdateOne(ctx, filter, bson.D{{"{{"}}Key: "$set", Value: update{{"}}"}})
	if err != nil {
		return nil, fmt.Errorf("repository:mongodb: Update failed: %v", err)
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("repository:mongodb: Update failed: 0 document modified")
	}

	return {{.SmallName}}, nil
}

// Delete soft delete a {{.SmallName}} record
func (r *{{.UcFirstName}}Mongo) Delete(ctx context.Context, ctr *domain.{{.UcFirstName}}Criteria) error {
	if ctr.ID == nil {
		return errors.New("repository:mongodb: Delte failed: {{.SmallName}} primary id required")
	}

	objectId, err := primitive.ObjectIDFromHex(*ctr.ID)
	if err != nil {
		return errors.New("repository:mongodb: Delete failed: {{.SmallName}} primary id unable to convert from hex to ObjectID")
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
		return fmt.Errorf("repository:mongo: failed to Delete {{.SmallName}}: %v", err)
	}

	if result.ModifiedCount == 0 {
		return errors.New("repository:mongo: failed to Delete {{.SmallName}}, 0 document modified")
	}

	return nil
}

// Convert{{.UcFirstName}} ...
func convert{{.UcFirstName}}(c *{{.UcFirstName}}) *domain.{{.UcFirstName}} {
	return &domain.{{.UcFirstName}}{
		ID:           c.ID.Hex(), {{range .Fields}}
    	{{.Name}}:    c.{{.Name}}, {{end}}
		TimeStamp: domain.TimeStamp{
			CreatedAt: c.TimeStamp.CreatedAt,
			UpdatedAt: c.TimeStamp.UpdatedAt,
			DeletedAt: c.TimeStamp.DeletedAt,
		},
	}
}

// Convert{{.UcFirstName}}List ...
func convert{{.UcFirstName}}List(cl []*{{.UcFirstName}}) []*domain.{{.UcFirstName}} {
	list := make([]*domain.{{.UcFirstName}}, 0)
	for _, c := range cl {
		list = append(list, convert{{.UcFirstName}}(c))
	}
	return list
}
`
