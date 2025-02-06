package achievement_criteria

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
	Controller *AchievementCriteriaController
	Service    *AchievementCriteriaService
	Logger     *logger.Logger
	Storage    *storage.ActiveStorage
}

func NewAchievementCriteriaModule(db *gorm.DB, router *gin.RouterGroup, log logger.Logger, emitter *emitter.Emitter, storage *storage.ActiveStorage) module.Module {

	service := NewAchievementCriteriaService(db, emitter, storage, log)
	controller := NewAchievementCriteriaController(service, storage)

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
	return m.DB.AutoMigrate(&models.AchievementCriteria{})
}

func (m *Module) GetModels() []interface{} {
	return []interface{}{&models.AchievementCriteria{}}
}
