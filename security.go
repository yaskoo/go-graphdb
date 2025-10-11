package graphdb

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

const (
	PathLogin              = "/rest/login"
	PathSecurity           = "/rest/security"
	PathSecurityFreeAccess = "/rest/security/free-access"
	PathSecurityUsers      = "/rest/security/users"
	PathSecurityUser       = "/rest/security/users/%s"
	PathCustomRoles        = "/rest/security/custom-roles"
	PathCustomRole         = "/rest/security/custom-roles/%s"
)

type User struct {
	Username           string          `json:"username,omitempty"`
	Password           string          `json:"password,omitempty"`
	GrantedAuthorities []string        `json:"grantedAuthorities,omitempty"`
	AppSettings        map[string]bool `json:"appSettings,omitempty"`
	DateCreated        int64           `json:"dateCreated,omitempty"`
	GptThreads         []interface{}   `json:"gptThreads,omitempty"`
}

// todo: UserDetails is pretty much the same as User

type UserDetails struct {
	Username    string          `json:"username,omitempty"`
	Password    string          `json:"password,omitempty"`
	Authorities []string        `json:"authorities,omitempty"`
	AppSettings map[string]bool `json:"appSettings,omitempty"`
	External    bool            `json:"external,omitempty"`
}

type FreeAccess struct {
	Enabled     bool            `json:"enabled,omitempty"`
	Authorities []string        `json:"authorities,omitempty"`
	AppSettings map[string]bool `json:"appSettings,omitempty"`
}

func (f *FreeAccess) WriteRepos(repos ...string) {
	for _, repo := range repos {
		f.Authorities = append(f.Authorities, fmt.Sprintf("WRITE_REPO_%s", strings.ToLower(repo)))
	}
}

func (f *FreeAccess) ReadRepos(repos ...string) {
	for _, repo := range repos {
		f.Authorities = append(f.Authorities, fmt.Sprintf("READ_REPO_%s", strings.ToLower(repo)))
	}
}

type SecurityClient struct {
	client *Client
}

func (s *SecurityClient) Enabled(ctx context.Context, config ...RequestConfig) (bool, error) {
	var enabled bool
	return enabled, s.client.get(ctx, PathSecurity, func(resp *http.Response) error {
		return json.NewDecoder(resp.Body).Decode(&enabled)
	}, config...)
}

func (s *SecurityClient) SetEnabled(ctx context.Context, enabled bool, config ...RequestConfig) error {
	config = append(config, JsonBody(enabled))
	return s.client.post(ctx, PathSecurity, nil, func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			all, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("security: toggling to %t failed, %s", enabled, string(all))
		}
		return nil
	}, config...)
}

func (s *SecurityClient) FreeAccess(ctx context.Context, config ...RequestConfig) (FreeAccess, error) {
	var free FreeAccess
	return free, s.client.get(ctx, PathSecurityFreeAccess, func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			all, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("security/free-access: %s", string(all))
		}
		return json.NewDecoder(resp.Body).Decode(&free)
	}, config...)
}

func (s *SecurityClient) ConfigureFreeAccess(ctx context.Context, access FreeAccess, config ...RequestConfig) error {
	config = append(config, JsonBody(access))
	return s.client.post(ctx, PathSecurityFreeAccess, nil, func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			all, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("security/free-access: %s", string(all))
		}
		return nil
	}, config...)
}

func (s *SecurityClient) Users(ctx context.Context, config ...RequestConfig) ([]User, error) {
	var users []User
	return users, s.client.get(ctx, PathSecurityUsers, func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			all, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("users: %s", string(all))
		}
		return json.NewDecoder(resp.Body).Decode(&users)
	}, config...)
}

func (s *SecurityClient) User(ctx context.Context, username string, config ...RequestConfig) (User, error) {
	var user User
	return user, s.client.get(ctx, fmt.Sprintf(PathSecurityUser, username), func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			all, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("users: %s", string(all))
		}
		return json.NewDecoder(resp.Body).Decode(&user)
	}, config...)
}

func (s *SecurityClient) CreateUser(ctx context.Context, user User, config ...RequestConfig) error {
	config = append(config, JsonBody(user))

	return s.client.post(ctx, fmt.Sprintf(PathSecurityUser, user.Username), nil, func(resp *http.Response) error {
		if resp.StatusCode != http.StatusCreated {
			all, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("user: %s", string(all))
		}
		return nil
	}, config...)
}

func (s *SecurityClient) UpdateUser(ctx context.Context, user User, config ...RequestConfig) error {
	config = append(config, JsonBody(user))

	return s.client.put(ctx, fmt.Sprintf(PathSecurityUser, user.Username), nil, func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			all, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("user: %s", string(all))
		}
		return nil
	}, config...)
}

func (s *SecurityClient) UpdateUserSettings(ctx context.Context, user User, config ...RequestConfig) error {
	config = append(config, JsonBody(user.AppSettings))

	return s.client.patch(ctx, fmt.Sprintf(PathSecurityUser, user.Username), nil, func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			all, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("users: %s", string(all))
		}
		return nil
	}, config...)
}

func (s *SecurityClient) DeleteUser(ctx context.Context, username string, config ...RequestConfig) error {
	return s.client.delete(ctx, fmt.Sprintf(PathSecurityUser, username), nil, func(resp *http.Response) error {
		if resp.StatusCode != http.StatusNoContent {
			all, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("users: %s", string(all))
		}
		return nil
	}, config...)
}

func (s *SecurityClient) Login(ctx context.Context, username, password string, config ...RequestConfig) (string, UserDetails, error) {
	var token string
	var details UserDetails

	config = append(config, JsonBody(map[string]string{
		"username": username,
		"password": password,
	}))
	return token, details, s.client.post(ctx, PathLogin, nil, func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			all, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("users_login: %s", string(all))
		}

		parts := strings.SplitN(resp.Header.Get("authorization"), " ", 2)
		if len(parts) == 2 || parts[0] != "GDB" {
			token = parts[1]
		}
		return json.NewDecoder(resp.Body).Decode(&details)
	}, config...)
}

func (s *SecurityClient) GetCustomRoles(ctx context.Context, config ...RequestConfig) (map[string]string, error) {
	var roles map[string]string
	return roles, s.client.get(ctx, PathCustomRoles, func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			all, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("custom_roles: %s", string(all))
		}
		return errors.Wrap(json.NewDecoder(resp.Body).Decode(&roles), "custom_roles")
	}, config...)
}

func (s *SecurityClient) ReplaceCustomRoles(ctx context.Context, roles map[string][]string, config ...RequestConfig) error {
	config = append(config, JsonBody(roles))
	return s.client.put(ctx, PathCustomRoles, nil, func(resp *http.Response) error {
		if resp.StatusCode != http.StatusOK {
			all, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("custom_roles: %s", string(all))
		}
		return nil
	}, config...)
}

func (s *SecurityClient) GetCustomRoleUsers(ctx context.Context, role string, config ...RequestConfig) ([]string, error) {
	var users []string
	return users, s.client.get(ctx, fmt.Sprintf(PathCustomRole, role), func(resp *http.Response) error {
		all, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("custom_roles: %s", string(all))
		}
		return errors.Wrap(json.NewDecoder(resp.Body).Decode(&users), "custom_roles")
	}, config...)
}

func (s *SecurityClient) ReplaceCustomRoleUsers(ctx context.Context, role string, users []string, config ...RequestConfig) error {
	config = append(config, JsonBody(users))
	return s.client.put(ctx, fmt.Sprintf(PathCustomRole, role), nil, func(resp *http.Response) error {
		all, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("custom_roles: %s", string(all))
		}
		return nil
	}, config...)
}

func (s *SecurityClient) AddCustomRoleUsers(ctx context.Context, role string, users []string, config ...RequestConfig) error {
	config = append(config, JsonBody(users))
	return s.client.post(ctx, fmt.Sprintf(PathCustomRole, role), nil, func(resp *http.Response) error {
		all, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("custom_roles: %s", string(all))
		}
		return nil
	}, config...)
}

func (s *SecurityClient) RemoveCustomRoleUsers(ctx context.Context, role string, users []string, config ...RequestConfig) error {
	config = append(config, JsonBody(users))
	return s.client.delete(ctx, fmt.Sprintf(PathCustomRole, role), nil, func(resp *http.Response) error {
		all, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("custom_roles: %s", string(all))
		}
		return nil
	}, config...)
}
