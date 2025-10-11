package graphdb

import (
	"context"
	"go-graphdb/testenv"
	"testing"
)

func TestSecurity_Enable(t *testing.T) {
	testenv.WithEnv(t, func(url string) {
		client := New(url)

		enabled, err := client.Security().Enabled(context.Background())
		if err != nil {
			t.Errorf("checking if security is enabled failed: %v", err)
		}

		if enabled {
			t.Errorf("security is enabled when it shouldn't be")
		}

		err = client.Security().SetEnabled(context.Background(), true)
		if err != nil {
			t.Errorf("enabling security failed: %v", err)
		}

		enabled, _ = client.Security().Enabled(context.Background())
		if !enabled {
			t.Errorf("security is disabled when it shouldn't be")
		}

		client = New(url, WithBasicAuth("admin", "root"))
		err = client.Security().SetEnabled(context.Background(), false)
		if err != nil {
			t.Errorf("failed to disable security: %v", err)
		}

		enabled, _ = client.Security().Enabled(context.Background())
		if enabled {
			t.Errorf("security should be disabled")
		}

		free, err := client.Security().FreeAccess(context.Background())
		if err != nil {
			t.Errorf("failed to get free access settings: %v", err)
		}
		if free.Enabled {
			t.Errorf("free access should be disabled")
		}

		repoID, err := createRepository(t, client)
		if err != nil {
			t.Errorf("failed to create repository for free access: %v", err)
		}

		_ = client.Security().SetEnabled(context.Background(), true)

		freeAccess := FreeAccess{
			Enabled: true,
			AppSettings: map[string]bool{
				// todo: these are required, unfortunately, graphdb doesn't return anything about this
				"DEFAULT_INFERENCE": true,
			},
		}
		freeAccess.WriteRepos(repoID)

		err = client.Security().ConfigureFreeAccess(context.Background(), freeAccess)
		if err != nil {
			t.Errorf("failed to configure free access: %v", err)
		}

		free, err = client.Security().FreeAccess(context.Background())
		if err != nil {
			t.Errorf("failed to get free access settings: %v", err)
		}
		if !free.Enabled {
			t.Errorf("free access should be enabled")
		}
		if len(free.Authorities) != 1 {
			t.Errorf("expected 1 authority, found %d", len(free.Authorities))
		}

		user := User{
			Username: "provisioner",
			Password: "password",
		}

		if err = client.Security().CreateUser(context.Background(), user); err != nil {
			t.Error("should create user")
		}

		users, err := client.Security().Users(context.Background())
		if err != nil || len(users) != 2 {
			t.Errorf("expected 2 user, found %d", len(users))
		}

		user, err = client.Security().User(context.Background(), "provisioner")
		if err != nil {
			t.Error("admin user should exist")
		}

		if user.Password != "" {
			t.Error("password should not be returned")
		}

		expected := "ROLE_ADMIN"
		user = User{
			Username:           "provisioner",
			GrantedAuthorities: []string{expected},
		}

		if err = client.Security().UpdateUser(context.Background(), user); err != nil {
			t.Errorf("failed to update user: %v", err)
		}

		user, _ = client.Security().User(context.Background(), "provisioner")
		var role string
		for _, authority := range user.GrantedAuthorities {
			if authority == expected {
				role = authority
				break
			}
		}
		if role != expected {
			t.Errorf("expected role '%s', found '%s'", expected, role)
		}

		err = client.Security().DeleteUser(context.Background(), "provisioner")
		if err != nil {
			t.Errorf("failed to delete user: %v", err)
		}

		_, err = client.Security().User(context.Background(), "provisioner")
		if err == nil {
			t.Error("should receive 404")
		}

		token, details, err := client.Security().Login(context.Background(), "admin", "root")
		if err != nil {
			t.Errorf("failed to login: %v", err)
		}

		if token == "" {
			t.Error("token should not be empty")
		}

		if details.Username != "admin" {
			t.Error("username should be admin")
		}
	})

}

func TestSecurity_CustomRoles(t *testing.T) {
	testenv.WithEnv(t, func(url string) {
		client := New(url)
		roles := map[string][]string{
			"custom_pishki": {"admin"},
		}

		err := client.Security().ReplaceCustomRoles(context.Background(), roles)
		if err != nil {
			t.Errorf("failed to set custom roles: %v", err)
		}
	})
}
