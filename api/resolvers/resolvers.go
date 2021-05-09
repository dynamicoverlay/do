package resolvers

import (
	"context"
	"fmt"

	"github.com/jinzhu/gorm"
	"imabad.dev/do/api/handlers"
	"imabad.dev/do/lib/models"
)

type RootResolver struct {
	Db *gorm.DB
}

func (r *RootResolver) Users(ctx context.Context) (*[]*models.UserResolver, error) {
	if ctx.Value(handlers.ContextKey("UserID")) == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	var list []*models.UserResolver
	var users []models.User
	r.Db.Find(&users)
	if &users != nil {
		for _, user := range users {
			list = append(list, &models.UserResolver{U: user})
		}
	}
	return &list, nil
}
