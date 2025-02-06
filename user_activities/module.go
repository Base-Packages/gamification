package user_activities

import (
	"base/core/emitter"
	"base/core/logger"
	"base/core/module"
	"base/core/storage"
	"base/packages/gamification/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Module struct {
	module.DefaultModule
	DB         *gorm.DB
	Controller *UserActivityController
	Service    *UserActivityService
	Logger     *logger.Logger
	Storage    *storage.ActiveStorage
}

func NewUserActivityModule(db *gorm.DB, router *gin.RouterGroup, log logger.Logger, emitter *emitter.Emitter, storage *storage.ActiveStorage) module.Module {

	service := NewUserActivityService(db, emitter, storage, log)
	controller := NewUserActivityController(service, storage)

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
	return m.DB.AutoMigrate(&models.UserActivity{})
}

func (m *Module) GetModels() []interface{} {
	return []interface{}{&models.UserActivity{}}
}
