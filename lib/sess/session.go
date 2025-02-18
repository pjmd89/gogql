package sess

import "go.mongodb.org/mongo-driver/bson/primitive"

type Session struct {
	UserID primitive.ObjectID
}
