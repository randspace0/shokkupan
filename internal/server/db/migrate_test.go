package db

import (
	"testing"

	"github.com/randspace0/shokkupan/internal/models"
)

// Users that existed before roles were added must stay admins, or an upgrade
// would lock everyone out of user management.
func TestExistingUsersBecomeAdmins(t *testing.T) {
	d, err := Init(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	// simulate the pre-roles schema and data
	d.Exec("DROP TABLE users")
	d.Exec("CREATE TABLE users (id integer primary key autoincrement, created_at datetime, updated_at datetime, deleted_at datetime, name text unique, source text, password_hash blob, totp_enabled numeric, totp_secret text, recovery_code text)")
	d.Exec("INSERT INTO users (name, source) VALUES ('old', 'github')")

	if err := d.AutoMigrate(&models.User{}); err != nil {
		t.Fatal(err)
	}
	var u models.User
	if err := d.Where("name = ?", "old").Take(&u).Error; err != nil {
		t.Fatal(err)
	}
	if u.Role != "admin" {
		t.Fatalf("pre-existing user role = %q, want admin", u.Role)
	}
}
