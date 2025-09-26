package service

import (
	"testing"
)

func TestComparePassword_Admin(t *testing.T) {
	password := "admin123"
	hash := "$2a$10$N9qo8uLOickgx2ZMRZoMye.IjdQvOQ5eqGStBUKx6XgKnrQvp.Fl6"
	if err := (&Service{}).comparePassword(password, hash); err != nil {
		t.Fatalf("expected match for admin, got error: %v", err)
	}
}

func TestComparePassword_User(t *testing.T) {
	password := "user123"
	hash := "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi"
	if err := (&Service{}).comparePassword(password, hash); err != nil {
		t.Fatalf("expected match for user, got error: %v", err)
	}
}

func TestComparePassword_Demo(t *testing.T) {
	password := "demo123"
	hash := "$2a$10$TKh8H1.PfQx37YgCzwiKb.KjNyWgaHb9cbcoQgdIVFlYg7B77UdFm"
	if err := (&Service{}).comparePassword(password, hash); err != nil {
		t.Fatalf("expected match for demo, got error: %v", err)
	}
}



