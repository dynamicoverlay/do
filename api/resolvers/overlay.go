package resolvers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm/dialects/postgres"
	jsonSchema "github.com/qri-io/jsonschema"
	"imabad.dev/do/api/handlers"
	"imabad.dev/do/lib/messaging"
	"imabad.dev/do/lib/models"
	"imabad.dev/do/lib/utils"
)

type updateStateArgs struct {
	Overlay string
	Key     string
	Value   string
}

type getOverlayArgs struct {
	ID string
}

type createOverlayArgs struct {
	Name string
}

type addModuleArgs struct {
	OverlayID string
	ModuleID  string
}

type removeModuleArgs struct {
	OverlayID string
	ModuleID  int
}

type updateModuleArgs struct {
	OverlayID string
	ModuleID  int
	Enabled   bool
	Settings  string
}

type UpdateStateResponse struct {
	updated bool
}

func (r *UpdateStateResponse) Updated() bool {
	return r.updated
}

func (r *RootResolver) UpdateState(ctx context.Context, args updateStateArgs) (*UpdateStateResponse, error) {
	if ctx.Value(handlers.ContextKey("UserID")) == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	//TODO: Check if they have permission to modify overlay
	var overlay models.Overlay
	err := r.Db.First(&overlay, "identifier = ?", args.Overlay).Error
	if err != nil {
		return nil, err
	}
	if &overlay != nil {
		var newValue map[string]interface{}
		json.Unmarshal([]byte(args.Value), &newValue)
		changeStateRequest := messaging.ChangeStateRequest{
			Overlay: args.Overlay,
			Key:     args.Key,
			Value:   newValue,
		}
		bytes, err := json.Marshal(changeStateRequest)
		if err != nil {
			return nil, err
		}
		messaging.PublishWithName("changeOverlayState", "application/json", bytes)
		return &UpdateStateResponse{
			updated: true,
		}, nil
	}
	return nil, fmt.Errorf("invalid overlay")
}

func (r *RootResolver) Overlay(ctx context.Context, args getOverlayArgs) (*models.OverlayResolver, error) {
	if ctx.Value(handlers.ContextKey("UserID")) == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	var overlay models.Overlay
	err := r.Db.First(&overlay, "identifier = ?", args.ID).Error
	if err != nil {
		return nil, err
	}
	if &overlay != nil {
		return &models.OverlayResolver{O: overlay}, nil
	}
	return nil, fmt.Errorf("invalid overlay")
}

func (r *RootResolver) Overlays(ctx context.Context) (*[]*models.OverlayResolver, error) {
	if ctx.Value(handlers.ContextKey("UserID")) == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	userID := ctx.Value(handlers.ContextKey("UserID")).(*int)
	var user models.User
	err := r.Db.First(&user, utils.IntValue(userID)).Error
	if err != nil {
		return nil, err
	}
	if &user != nil {
		var overlays []models.Overlay
		var list []*models.OverlayResolver
		r.Db.Model(&user).Related(&overlays, "Overlays")
		for _, overlay := range overlays {
			list = append(list, &models.OverlayResolver{O: overlay})
		}
		return &list, nil
	}
	return nil, fmt.Errorf("invalid user")
}

func (r *RootResolver) CreateOverlay(ctx context.Context, args createOverlayArgs) (*models.OverlayResolver, error) {
	if ctx.Value(handlers.ContextKey("UserID")) == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	userID := ctx.Value(handlers.ContextKey("UserID")).(*int)
	var user models.User
	err := r.Db.First(&user, utils.IntValue(userID)).Error
	if err != nil {
		return nil, err
	}
	if &user != nil {
		overlay := models.Overlay{
			Name:       args.Name,
			Identifier: uuid.New().String(),
			Pin:        uuid.New().String(),
		}
		if err := r.Db.Create(&overlay).Error; err != nil {
			return nil, err
		}
		r.Db.Model(&user).Association("Overlays").Append(overlay)
		return &models.OverlayResolver{O: overlay}, nil
	}
	return nil, fmt.Errorf("invalid user")
}

func (r *RootResolver) AddModuleToOverlay(ctx context.Context, args addModuleArgs) (*models.OverlayResolver, error) {
	if ctx.Value(handlers.ContextKey("UserID")) == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	userID := ctx.Value(handlers.ContextKey("UserID")).(*int)
	var user models.User
	err := r.Db.First(&user, utils.IntValue(userID)).Error
	if err != nil {
		return nil, err
	}
	if &user != nil {
		var overlay models.Overlay
		err := r.Db.First(&overlay, "identifier = ?", args.OverlayID).Error
		if err != nil {
			return nil, err
		}
		if &overlay == nil {
			return nil, fmt.Errorf("invalid overlay")
		}
		var module models.Module
		err = r.Db.First(&module, "identifier = ?", args.ModuleID).Error
		if err != nil {
			return nil, err
		}
		if &module == nil {
			return nil, fmt.Errorf("invalid module")
		}
		newModule := models.OverlayModule{
			Overlay:  overlay,
			Module:   module,
			Enabled:  true,
			Settings: postgres.Jsonb{},
		}
		r.Db.Save(&newModule)
		r.Db.Model(&overlay).Association("Modules").Append(&newModule)
		return &models.OverlayResolver{O: overlay}, nil
	}
	return nil, fmt.Errorf("invalid user")
}

func (r *RootResolver) RemoveModuleFromOverlay(ctx context.Context, args removeModuleArgs) (*models.OverlayResolver, error) {
	if ctx.Value(handlers.ContextKey("UserID")) == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	userID := ctx.Value(handlers.ContextKey("UserID")).(*int)
	var user models.User
	err := r.Db.First(&user, utils.IntValue(userID)).Error
	if err != nil {
		return nil, err
	}
	if &user != nil {
		var overlay models.Overlay
		err := r.Db.First(&overlay, "identifier = ?", args.OverlayID).Error
		if err != nil {
			return nil, err
		}
		if &overlay == nil {
			return nil, fmt.Errorf("invalid overlay")
		}
		var module models.OverlayModule
		err = r.Db.First(&module, "id = ?", args.ModuleID).Error
		if err != nil {
			return nil, err
		}
		if &module == nil {
			return nil, fmt.Errorf("invalid module")
		}
		r.Db.Model(&overlay).Association("Modules").Delete(&module)
		r.Db.Delete(&module)
		return &models.OverlayResolver{O: overlay}, nil
	}
	return nil, fmt.Errorf("invalid user")
}
func (r *RootResolver) UpdateOverlayModule(ctx context.Context, args updateModuleArgs) (*models.OverlayModuleResolver, error) {
	if ctx.Value(handlers.ContextKey("UserID")) == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	userID := ctx.Value(handlers.ContextKey("UserID")).(*int)
	var user models.User
	err := r.Db.First(&user, utils.IntValue(userID)).Error
	if err != nil {
		return nil, err
	}
	if &user != nil {
		var overlay models.Overlay
		err := r.Db.First(&overlay, "identifier = ?", args.OverlayID).Error
		if err != nil {
			return nil, err
		}
		if &overlay == nil {
			return nil, fmt.Errorf("invalid overlay")
		}
		var module models.OverlayModule
		err = r.Db.Preload("Module").First(&module, "id = ?", args.ModuleID).Error
		if err != nil {
			return nil, err
		}
		if &module == nil {
			return nil, fmt.Errorf("invalid module")
		}
		module.Enabled = args.Enabled
		rs := &jsonSchema.Schema{}
		val, _ := module.Module.SettingsFormat.Value()
		if err := json.Unmarshal([]byte(val.(string)), rs); err != nil {
			return nil, err
		}
		errs, err := rs.ValidateBytes(ctx, []byte(args.Settings))
		if err != nil {
			return nil, err
		}
		if len(errs) > 0 {
			var errorstrings []string
			for _, err := range errs {
				errorstrings = append(errorstrings, err.Error())
			}
			return nil, fmt.Errorf(strings.Join(errorstrings, ", "))
		}
		module.Settings = postgres.Jsonb{RawMessage: json.RawMessage(args.Settings)}
		r.Db.Save(&module)
		return &models.OverlayModuleResolver{O: module}, nil
	}
	return nil, fmt.Errorf("invalid user")
}
func (r *RootResolver) Modules(ctx context.Context) (*[]*models.ModuleResolver, error) {
	if ctx.Value(handlers.ContextKey("UserID")) == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	var list []*models.ModuleResolver
	var modules []models.Module
	r.Db.Find(&modules)
	if &modules != nil {
		for _, module := range modules {
			list = append(list, &models.ModuleResolver{M: module})
		}
	}
	return &list, nil
}
