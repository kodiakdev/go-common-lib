package commonauth

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/emicklei/go-restful"
)

const (
	RequesterUserID = "requesterUserId"
)

//ExtractRequesterID extract userID from token
func ExtractRequesterID(req *restful.Request) primitive.ObjectID {
	iUserID := req.Attribute(RequesterUserID)
	userID := fmt.Sprintf("%v", iUserID)
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		id, _ = primitive.ObjectIDFromHex("")
	}
	return id
}
