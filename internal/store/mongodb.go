package store

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongodb struct {
	client     *mongo.Client
	db         *mongo.Database
	collection *mongo.Collection
}

const collectionName string = "links"

var _ (Store) = (*mongodb)(nil)

func NewMongoDBStore(ctx context.Context, user string, password string, host string, databaseName string) (*mongodb, error) {
	connectionString := fmt.Sprintf("mongodb://%s:%s@%s", user, password, host)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	exists := false
	for _, d := range databases {
		if d == databaseName {
			exists = true
			break
		}
	}
	db := client.Database(databaseName)
	if !exists {
		err = db.CreateCollection(ctx, collectionName)
	}
	if err != nil {
		return nil, err
	}
	collection := db.Collection(collectionName)
	textModel := mongo.IndexModel{Keys: bson.D{{Key: "_id", Value: "text"}, {Key: "description", Value: "text"}}}
	_, err = collection.Indexes().CreateMany(ctx, []mongo.IndexModel{textModel})
	if err != nil {
		return nil, err
	}
	return &mongodb{
		client:     client,
		db:         db,
		collection: collection,
	}, nil
}

func (m *mongodb) Close(ctx context.Context) error {
	if m.client == nil {
		return nil
	}
	return m.client.Disconnect(ctx)
}

// CreateLink implements Store.
func (m *mongodb) CreateLink(ctx context.Context, link Link) error {
	_, err := m.collection.InsertOne(ctx, link)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrIDExists
		}
		return err
	}
	return nil
}

// DisableLink implements Store.
func (m *mongodb) DisableLink(ctx context.Context, name string) error {
	result, err := m.collection.DeleteOne(ctx, bson.M{"_id": name})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrLinkNotFound
	}
	return nil
}

// GetLinkByName implements Store.
func (m *mongodb) GetLinkByName(ctx context.Context, name string) (Link, error) {
	result := m.collection.FindOne(ctx, bson.M{"_id": name})
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return Link{}, ErrLinkNotFound
		}
		return Link{}, result.Err()
	}
	link := Link{}
	err := result.Decode(&link)
	if err != nil {
		return Link{}, err
	}
	return link, nil
}

// GetLinkByURL implements Store.
func (m *mongodb) GetLinkByURL(ctx context.Context, url string) (Link, error) {
	result := m.collection.FindOne(ctx, bson.M{"url": url})
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return Link{}, ErrLinkNotFound
		}
		return Link{}, result.Err()
	}
	var link Link
	err := result.Decode(&link)
	if err != nil {
		return Link{}, err
	}
	return link, nil
}

// GetOwnedLinks implements Store.
func (m *mongodb) GetOwnedLinks(ctx context.Context, email string) ([]Link, error) {
	cursor, err := m.collection.Find(ctx, bson.M{"created_by": email})
	if err != nil {
		return []Link{}, err
	}
	if cursor.RemainingBatchLength() == 0 {
		return []Link{}, nil
	}
	var links []Link
	err = cursor.All(ctx, &links)
	if err != nil {
		return []Link{}, err
	}
	return links, nil
}

// GetPopularLinks implements Store.
func (m *mongodb) GetPopularLinks(ctx context.Context, size int) ([]Link, error) {
	cursor, err := m.collection.Find(ctx, bson.D{}, options.Find().SetSort(bson.D{{Key: "views", Value: -1}}).SetLimit(int64(size)))
	if err != nil {
		return []Link{}, err
	}
	if cursor.RemainingBatchLength() == 0 {
		return []Link{}, nil
	}
	var links []Link
	err = cursor.All(ctx, &links)
	if err != nil {
		return []Link{}, err
	}
	return links, nil
}

// GetRecentLinks implements Store.
func (m *mongodb) GetRecentLinks(ctx context.Context, size int) ([]Link, error) {
	cursor, err := m.collection.Find(ctx, bson.D{}, options.Find().SetSort(bson.D{{Key: "updated_at", Value: -1}}).SetLimit(int64(size)))
	if err != nil {
		return []Link{}, err
	}
	if cursor.RemainingBatchLength() == 0 {
		return []Link{}, nil
	}
	var links []Link
	err = cursor.All(ctx, &links)
	if err != nil {
		return []Link{}, err
	}
	return links, nil
}

// IncrementLinkViews implements Store.
func (m *mongodb) IncrementLinkViews(ctx context.Context, name string) error {
	link, err := m.GetLinkByName(ctx, name)
	if err != nil {
		return err
	}
	link.Views++
	err = m.updateLink(ctx, link)
	if err != nil {
		return err
	}
	return nil
}

// QueryLinks implements Store.
func (m *mongodb) QueryLinks(ctx context.Context, query string) ([]Link, error) {
	filter := bson.D{{Key: "$text", Value: bson.D{{Key: "$search", Value: fmt.Sprintf("\"%s\"", query)}}}}
	cursor, err := m.collection.Find(ctx, filter)
	if err != nil {
		return []Link{}, err
	}
	if cursor.RemainingBatchLength() == 0 {
		return []Link{}, nil
	}
	var links []Link
	err = cursor.All(ctx, &links)
	if err != nil {
		return []Link{}, err
	}
	return links, nil
}

func (m *mongodb) updateLink(ctx context.Context, link Link) error {
	// resp, err := m.collection.UpdateOne(ctx, bson.M{"_id": link.Name}, link)
	resp := m.collection.FindOneAndReplace(ctx, bson.M{"_id": link.Name}, link)
	if resp.Err() != nil {
		if errors.Is(resp.Err(), mongo.ErrNoDocuments) {
			return ErrLinkNotFound
		}
		return resp.Err()
	}
	return nil
}
