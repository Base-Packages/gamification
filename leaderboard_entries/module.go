package leaderboard_entries

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
	Controller *LeaderboardEntryController
	Service    *LeaderboardEntryService
	Logger     *logger.Logger
	Storage    *storage.ActiveStorage
}

func NewLeaderboardEntryModule(db *gorm.DB, router *gin.RouterGroup, log logger.Logger, emitter *emitter.Emitter, storage *storage.ActiveStorage) module.Module {

	service := NewLeaderboardEntryService(db, emitter, storage, log)
	controller := NewLeaderboardEntryController(service, storage)

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
	return m.DB.AutoMigrate(&models.LeaderboardEntry{})
}

func (m *Module) GetModels() []interface{} {
	return []interface{}{&models.LeaderboardEntry{}}
}
