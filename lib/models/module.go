package models

import (
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type Module struct {
	gorm.Model
	Name           string `gorm:"type:varchar(255)"`
	Identifier     string `gorm:"type:varchar(255);unique_index"`
	StateFormat    postgres.Jsonb
	SettingsFormat postgres.Jsonb
}

type ModuleResolver struct {
	M Module
}

func (m *ModuleResolver) Identifier() string {
	return m.M.Identifier
}

func (m *ModuleResolver) Name() string {
	return m.M.Name
}

func (m *ModuleResolver) StateFormat() string {
	value, _ := m.M.StateFormat.Value()
	return value.(string)
}

func (m *ModuleResolver) SettingsFormat() string {
	value, _ := m.M.SettingsFormat.Value()
	return value.(string)
}

type OverlayModule struct {
	gorm.Model
	OverlayID uint
	Overlay   Overlay
	ModuleID  uint
	Module    Module
	Enabled   bool
	Settings  postgres.Jsonb
}

type OverlayModuleResolver struct {
	O OverlayModule
}

func (o *OverlayModuleResolver) ID() int32 {
	return int32(o.O.ID)
}

func (o *OverlayModuleResolver) Overlay() *OverlayResolver {
	return &OverlayResolver{O: o.O.Overlay}
}

func (o *OverlayModuleResolver) Module() *ModuleResolver {
	return &ModuleResolver{M: o.O.Module}
}

func (o *OverlayModuleResolver) Enabled() bool {
	return o.O.Enabled
}

func (o *OverlayModuleResolver) Settings() string {
	value, _ := o.O.Settings.Value()
	return value.(string)
}
