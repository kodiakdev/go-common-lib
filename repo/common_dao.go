package commonrepo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommonAudit struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedBy primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedBy primitive.ObjectID `bson:"updatedBy,omitempty" json:"updatedBy"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
