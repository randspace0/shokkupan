package models

import "gorm.io/gorm"

// AppStatusEvent records a change in whether an app is running.
type AppStatusEvent struct {
	gorm.Model
	AppName string `gorm:"index"`
	Running bool
}
