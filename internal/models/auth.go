package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name   string `gorm:"unique"`
	Source string
	// Role is admin, developer or viewer. The DB default keeps every user that
	// existed before roles were added an admin.
	Role    string `gorm:"default:admin"`
	SSHKeys []SSHKey

	PasswordHash []byte
	TotpEnabled  bool
	TotpSecret   string
	RecoveryCode string
}

type SSHKey struct {
	gorm.Model
	UserID   uint
	GithubID int64 `gorm:"unique"`
	Key      string
}

// AppAccess grants a non-admin user access to one app.
type AppAccess struct {
	gorm.Model
	UserName string `gorm:"uniqueIndex:idx_user_app"`
	AppName  string `gorm:"uniqueIndex:idx_user_app"`
}
