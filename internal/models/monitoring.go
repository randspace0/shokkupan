package models

import "gorm.io/gorm"

// AppStatusEvent records a change in whether an app is running.
type AppStatusEvent struct {
	gorm.Model
	AppName string `gorm:"index"`
	Running bool
}

// AppStatsSample is one poll's summed container usage for an app.
type AppStatsSample struct {
	gorm.Model
	AppName  string `gorm:"index"`
	CPU      float64
	MemBytes int64
	MemLimit int64
}
