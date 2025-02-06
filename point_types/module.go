package point_types

import (
	"base/core/emitter"
	"base/core/logger"
	"base/core/module"
	"base/core/packages/gamification/models"
	"base/core/storage"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Module struct {
	module.DefaultModule
	DB         *gorm.DB
	Controller *PointTypeController
	Service    *PointTypeService
	Logger     *logger.Logger
	Storage    *storage.ActiveStorage
}

func NewPointTypeModule(db *gorm.DB, router *gin.RouterGroup, log logger.Logger, emitter *emitter.Emitter, storage *storage.ActiveStorage) module.Module {

	service := NewPointTypeService(db, emitter, storage, log)
	controller := NewPointTypeController(service, storage)

	m := &Module{
		DB:         db,
		Service:    service,
		Controller: controller,
		Logger:     &log,
		Storage:    storage,
	}

	return m
}

func (m *Module) Routes(router *gin.RouterGroup) {
	m.Controller.Routes(router)
}

func (m *Module) Migrate() error {
	return m.DB.AutoMigrate(&models.PointType{})
}

func (m *Module) GetModels() []interface{} {
	return []interface{}{&models.PointType{}}
}
