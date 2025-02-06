package gamification

import (
	// MODULE_IMPORT_MARKER - Do not remove this comment because it's used by the CLI to add new module imports

	"base/core/config"
	"base/core/database"
	"base/core/emitter"
	"base/core/logger"
	"base/core/module"
	"base/core/storage"
	"base/packages/gamification/achievement_criteria"
	"base/packages/gamification/achievements"
	"base/packages/gamification/activity_types"
	"base/packages/gamification/challenges"
	"base/packages/gamification/leaderboard_entries"
	"base/packages/gamification/leaderboards"
	"base/packages/gamification/levels"
	"base/packages/gamification/point_types"
	"base/packages/gamification/user_achievements"
	"base/packages/gamification/user_activities"
	"base/packages/gamification/user_challenges"
	"base/packages/gamification/user_levels"
	"base/packages/gamification/user_points"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Gamification struct {
	DB      *gorm.DB
	Router  *gin.Engine
	Log     logger.Logger
	Emitter *emitter.Emitter
	Storage *storage.ActiveStorage
	Modules []module.Module
}

// NewApp creates and initializes a new App instance
func NewApp(cfg *config.Config) (*Gamification, error) {
	// Initialize logger
	logConfig := logger.Config{
		Environment: "development",
		LogPath:     "logs",
		Level:       "debug",
	}
	log, err := logger.NewLogger(logConfig)
	if err != nil {
		return nil, err
	}
	// Initialize router
	router := gin.Default()
	// Initialize emitter
	emitter := &emitter.Emitter{}
	// Initialize database (you'll need to implement this)
	db, err := initDB(cfg)
	if err != nil {
		return nil, err
	}
	// Initialize storage
	storageConfig := storage.Config{
		Provider:  cfg.StorageProvider,
		Path:      cfg.StoragePath,
		BaseURL:   cfg.StorageBaseURL,
		APIKey:    cfg.StorageAPIKey,
		APISecret: cfg.StorageAPISecret,
		Endpoint:  cfg.StorageEndpoint,
		Bucket:    cfg.StorageBucket,
		CDN:       cfg.CDN,
	}
	activeStorage, err := storage.NewActiveStorage(db, storageConfig)
	if err != nil {
		return nil, err
	}
	app := &Gamification{
		DB:      db,
		Router:  router,
		Log:     log,
		Emitter: emitter,
		Storage: activeStorage,
		Modules: make([]module.Module, 0),
	}
	// Initialize modules
	moduleInitializer := &GamificationModuleInitializer{
		DB:      db,
		Router:  router.Group("/api"),
		Logger:  log,
		Emitter: emitter,
		Storage: activeStorage,
	}
	app.Modules = moduleInitializer.InitializeModules(db)
	return app, nil
}

// GamificationModuleInitializer holds all dependencies needed for app module initialization
type GamificationModuleInitializer struct {
	DB      *gorm.DB
	Router  *gin.RouterGroup
	Logger  logger.Logger
	Emitter *emitter.Emitter
	Storage *storage.ActiveStorage
}

// InitializeModules initializes all application modules
func (a *GamificationModuleInitializer) InitializeModules(db *gorm.DB) []module.Module {
	var modules []module.Module
	// Initialize modules
	moduleMap := a.getModules(db)
	// Register and initialize each module
	for name, mod := range moduleMap {

		if err := module.RegisterModule(name, mod); err != nil {
			a.Logger.Error("Failed to register module",
				logger.String("module", name),
				logger.String("error", err.Error()))
			continue
		}
		// Initialize the module
		if err := mod.Init(); err != nil {
			a.Logger.Error("Failed to initialize module",
				logger.String("module", name),
				logger.String("error", err.Error()))
			continue
		}
		// Migrate the module
		if err := mod.Migrate(); err != nil {
			a.Logger.Error("Failed to migrate module",
				logger.String("module", name),
				logger.String("error", err.Error()))
			continue
		}
		// Set up routes for the module
		if routeModule, ok := mod.(interface{ Routes(*gin.RouterGroup) }); ok {
			routeModule.Routes(a.Router)
		}
		modules = append(modules, mod)
	}
	return modules
}

// getModules returns a map of module name to module instance
func (a *GamificationModuleInitializer) getModules(db *gorm.DB) map[string]module.Module {
	modules := make(map[string]module.Module)
	// Define the module initializers directly
	moduleInitializers := map[string]func(*gorm.DB, *gin.RouterGroup, logger.Logger, *emitter.Emitter, *storage.ActiveStorage) module.Module{

		"activity_types": func(db *gorm.DB, router *gin.RouterGroup, log logger.Logger, emitter *emitter.Emitter, activeStorage *storage.ActiveStorage) module.Module {
			return activity_types.NewActivityTypeModule(db, router, log, emitter, activeStorage)
		},

		"user_activities": func(db *gorm.DB, router *gin.RouterGroup, log logger.Logger, emitter *emitter.Emitter, activeStorage *storage.ActiveStorage) module.Module {
			return user_activities.NewUserActivityModule(db, router, log, emitter, activeStorage)
		},

		"point_types": func(db *gorm.DB, router *gin.RouterGroup, log logger.Logger, emitter *emitter.Emitter, activeStorage *storage.ActiveStorage) module.Module {
			return point_types.NewPointTypeModule(db, router, log, emitter, activeStorage)
		},

		"user_points": func(db *gorm.DB, router *gin.RouterGroup, log logger.Logger, emitter *emitter.Emitter, activeStorage *storage.ActiveStorage) module.Module {
			return user_points.NewUserPointModule(db, router, log, emitter, activeStorage)
		},

		"achievements": func(db *gorm.DB, router *gin.RouterGroup, log logger.Logger, emitter *emitter.Emitter, activeStorage *storage.ActiveStorage) module.Module {
			return achievements.NewAchievementModule(db, router, log, emitter, activeStorage)
		},

		"user_achievements": func(db *gorm.DB, router *gin.RouterGroup, log logger.Logger, emitter *emitter.Emitter, activeStorage *storage.ActiveStorage) module.Module {
			return user_achievements.NewUserAchievementModule(db, router, log, emitter, activeStorage)
		},

		"achievement_criteria": func(db *gorm.DB, router *gin.RouterGroup, log logger.Logger, emitter *emitter.Emitter, activeStorage *storage.ActiveStorage) module.Module {
			return achievement_criteria.NewAchievementCriteriaModule(db, router, log, emitter, activeStorage)
		},

		"levels": func(db *gorm.DB, router *gin.RouterGroup, log logger.Logger, emitter *emitter.Emitter, activeStorage *storage.ActiveStorage) module.Module {
			return levels.NewLevelModule(db, router, log, emitter, activeStorage)
		},

		"user_levels": func(db *gorm.DB, router *gin.RouterGroup, log logger.Logger, emitter *emitter.Emitter, activeStorage *storage.ActiveStorage) module.Module {
			return user_levels.NewUserLevelModule(db, router, log, emitter, activeStorage)
		},

		"challenges": func(db *gorm.DB, router *gin.RouterGroup, log logger.Logger, emitter *emitter.Emitter, activeStorage *storage.ActiveStorage) module.Module {
			return challenges.NewChallengeModule(db, router, log, emitter, activeStorage)
		},

		"user_challenges": func(db *gorm.DB, router *gin.RouterGroup, log logger.Logger, emitter *emitter.Emitter, activeStorage *storage.ActiveStorage) module.Module {
			return user_challenges.NewUserChallengeModule(db, router, log, emitter, activeStorage)
		},

		"leaderboards": func(db *gorm.DB, router *gin.RouterGroup, log logger.Logger, emitter *emitter.Emitter, activeStorage *storage.ActiveStorage) module.Module {
			return leaderboards.NewLeaderboardModule(db, router, log, emitter, activeStorage)
		},

		"leaderboard_entries": func(db *gorm.DB, router *gin.RouterGroup, log logger.Logger, emitter *emitter.Emitter, activeStorage *storage.ActiveStorage) module.Module {
			return leaderboard_entries.NewLeaderboardEntryModule(db, router, log, emitter, activeStorage)
		},

		// MODULE_INITIALIZER_MARKER - Do not remove this comment because it's used by the CLI to add new module initializers
	}

	// Initialize and register each module
	for name, initializer := range moduleInitializers {
		modules[name] = initializer(db, a.Router, a.Logger, a.Emitter, a.Storage)
	}

	return modules
}

// initDB initializes the database connection
func initDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := database.InitDB(cfg)
	if err != nil {
		return nil, err
	}
	return db.DB, nil
}
