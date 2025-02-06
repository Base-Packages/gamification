package leaderboards

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
	Controller *LeaderboardController
	Service    *LeaderboardService
	Logger     *logger.Logger
	Storage    *storage.ActiveStorage
}

func NewLeaderboardModule(db *gorm.DB, router *gin.RouterGroup, log logger.Logger, emitter *emitter.Emitter, storage *storage.ActiveStorage) module.Module {

	service := NewLeaderboardService(db, emitter, storage, log)
	controller := NewLeaderboardController(service, storage)

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
	return m.DB.AutoMigrate(&models.Leaderboard{})
}

func (m *Module) GetModels() []interface{} {
	return []interface{}{&models.Leaderboard{}}
}
