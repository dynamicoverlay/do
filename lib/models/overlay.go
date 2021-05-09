package models

import "github.com/jinzhu/gorm"

type Overlay struct {
	gorm.Model
	Identifier string  `gorm:"type:varchar(255);unique_index"`
	Pin        string  `gorm:"type:varchar(255)"`
	Name       string  `gorm:"type:varchar(255)"`
	Users      []*User `gorm:"many2many:user_overlays;"`
	Modules    []*OverlayModule
}

type OverlayResolver struct {
	O Overlay
}

func (o *OverlayResolver) Identifier() string {
	return o.O.Identifier
}

func (o *OverlayResolver) Pin() string {
	return o.O.Pin
}

func (o *OverlayResolver) Name() string {
	return o.O.Name
}

func (o *OverlayResolver) Modules() []*OverlayModuleResolver {
	var list []*OverlayModuleResolver
	for _, module := range o.O.Modules {
		list = append(list, &OverlayModuleResolver{O: *module})
	}
	return list
}
