package db

import (
	"github.com/glebarez/sqlite"
	"github.com/randspace0/shokkupan/internal/models"
	"gorm.io/gorm"
)

func Init(dsn string) (*gorm.DB, error) {
	dbCfg := &gorm.Config{
		Logger: Logger{},
	}

	/*if cfg.DebugMode == false {
		dbCfg.Logger = dbCfg.Logger.LogMode(logger.Silent)
	}*/

	db, err := gorm.Open(sqlite.Open(dsn), dbCfg)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&models.Server{},
		&models.ServerSecrets{},
		&models.App{},
		&models.Service{},
		&models.User{},
		&models.SSHKey{},
		&models.AppAccess{},
		&models.GithubApp{},
		&models.AppSetupConfig{},
		&models.AppStatusEvent{},
		&models.AppStatsSample{},
		&models.ServerSample{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
