package v1

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/alifcapital/keycloak_module/conf"
	"github.com/alifcapital/keycloak_module/src/iam"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type KeycloakHTTPClient struct {
	cfg    *conf.Config
	client *http.Client
}

func NewKeycloakHTTPClient(cfg *conf.Config) (*KeycloakHTTPClient, error) {
	ctx := context.Background()
	issuer := fmt.Sprintf("%s/realms/%s", cfg.KeycloakBaseUrl, cfg.KeycloakRealm)
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to build openid provider"), err)
	}

	config := oauth2.Config{
		ClientID:     "security-admin-console",
		ClientSecret: "",
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "",
		Scopes:       nil,
	}

	token, err := config.PasswordCredentialsToken(ctx, cfg.KeycloakUserName, cfg.KeycloakPassword)
	if err != nil {
		return nil, err
	}
	client := config.Client(ctx, token)

	return &KeycloakHTTPClient{
		cfg:    cfg,
		client: client,
	}, nil
}

func (kc *KeycloakHTTPClient) Store(_ context.Context, user *iam.User) error {
	payload, err := kc.marshal(user)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to marshal user: %+v", user), err)
	}

	url := kc.createUserUrl()

	resp, err := kc.client.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return errors.Join(fmt.Errorf("failed to send a requiest to: %s", url), err)
	}
	defer resp.Body.Close()

	var b bytes.Buffer
	if _, err := b.ReadFrom(resp.Body); err != nil {
		return errors.Join(fmt.Errorf("failed to read response body, requiest to: %s", url), err)
	}

	// we have an error
	if resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("failed to execute requiest to: %s, err: %s, http_code: %d", url, b.String(), resp.StatusCode)
	}

	// parse from location url which has user id
	location := resp.Header.Get("location")
	parts := strings.Split(location, "/")
	userID := parts[len(parts)-1]
	user.ID = userID

	return nil
}

func (kc *KeycloakHTTPClient) Get(_ context.Context, userID string) (*iam.User, error) {
	url := kc.userRepresentationUrl(userID)

	resp, err := kc.client.Get(url)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to send a request to: %s", url), err)
	}
	defer resp.Body.Close()

	var b bytes.Buffer
	_, err = b.ReadFrom(resp.Body)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to read response body, requiest to: %s", url), err)
	}

	// we have an error
	if resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("failed to execute requiest to: %s, err: %s, http_code: %d", url, b.String(), resp.StatusCode)
	}

	user := new(iam.User)
	if err := kc.unmarshal(b.Bytes(), user); err != nil {
		return nil, errors.Join(fmt.Errorf("failed to unmarshal response into user"), err)
	}
	return user, nil
}

func (kc *KeycloakHTTPClient) Update(ctx context.Context, user *iam.User) error {
	payload, err := kc.marshal(user)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to marshal user: %+v", user), err)
	}

	url := kc.updateUserUrl(user.ID)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return errors.Join(fmt.Errorf("failed to build request to url: %s", url), err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")

	kc.MoveLDAPUserToOU(user.Username, user.Attributes["newOU"].(string))

	resp, err := kc.client.Do(req)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to send a requiest to: %s", url), err)
	}
	defer resp.Body.Close()

	// we have an error
	if resp.StatusCode >= http.StatusMultipleChoices {
		var b bytes.Buffer
		_, _ = b.ReadFrom(resp.Body)
		return fmt.Errorf("failed to execute requiest to: %s, err: %s, http_code: %d", url, b.String(), resp.StatusCode)
	}
	return nil
}

func (kc *KeycloakHTTPClient) MoveLDAPUserToOU(username, newOU string) error {
	l, err := ldap.DialURL(kc.cfg.LDAPServer, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
	if err != nil {
		return fmt.Errorf("ldap dial failed: %w", err)
	}
	defer l.Close()

	// Bind
	err = l.Bind(kc.cfg.LDAPBindUser, kc.cfg.LDAPBindPassword)
	if err != nil {
		return fmt.Errorf("ldap bind failed: %w", err)
	}

	// Search user DN
	searchRequest := ldap.NewSearchRequest(
		kc.cfg.LDAPBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 1, 0, false,
		fmt.Sprintf("(&(objectClass=user)(sAMAccountName=%s))", username),
		[]string{"dn", "cn"},
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		return fmt.Errorf("search error: %w", err)
	}
	if len(sr.Entries) == 0 {
		return fmt.Errorf("user not found")
	}

	entry := sr.Entries[0]
	oldDN := entry.DN
	cn := entry.GetAttributeValue("cn")

	// Modify DN (move)
	modDN := ldap.NewModifyDNRequest(oldDN, "CN="+cn, true, fmt.Sprintf("%s,%s", newOU, kc.cfg.LDAPBaseDN))
	err = l.ModifyDN(modDN)
	if err != nil {
		return fmt.Errorf("move failed: %w", err)
	}

	return nil
}

func (kc *KeycloakHTTPClient) SetTemporaryPassword(ctx context.Context, userID, tempPassword string) error {
	if tempPassword == "" {
		tempPassword = uuid.NewString()
	}
	newPass := map[string]interface{}{
		"value":     tempPassword,
		"type":      "password",
		"temporary": true,
	}
	payload, err := kc.marshal(newPass)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to marshal user map with password data"), err)
	}

	url := kc.updateUserPasswordUrl(userID)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return errors.Join(fmt.Errorf("failed to build request to url: %s", url), err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")

	resp, err := kc.client.Do(req)
	if err != nil {
		return errors.Join(fmt.Errorf("failed to send a requiest to: %s", url), err)
	}
	defer resp.Body.Close()

	// we have an error
	if resp.StatusCode >= http.StatusMultipleChoices {
		var b bytes.Buffer
		_, _ = b.ReadFrom(resp.Body)
		return fmt.Errorf("failed to execute requiest to: %s, err: %s, http_code: %d", url, b.String(), resp.StatusCode)
	}
	return nil
}

func (kc *KeycloakHTTPClient) marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (kc *KeycloakHTTPClient) unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (kc *KeycloakHTTPClient) createUserUrl() string {
	return fmt.Sprintf("%s/admin/realms/%s/users", kc.cfg.KeycloakBaseUrl, kc.cfg.KeycloakRealm)
}

func (kc *KeycloakHTTPClient) userRepresentationUrl(userID string) string {
	return fmt.Sprintf("%s/admin/realms/%s/users/%s", kc.cfg.KeycloakBaseUrl, kc.cfg.KeycloakRealm, userID)
}

func (kc *KeycloakHTTPClient) updateUserUrl(userID string) string {
	return fmt.Sprintf("%s/admin/realms/%s/users/%s", kc.cfg.KeycloakBaseUrl, kc.cfg.KeycloakRealm, userID)
}

func (kc *KeycloakHTTPClient) updateUserPasswordUrl(userID string) string {
	return fmt.Sprintf("%s/admin/realms/%s/users/%s/reset-password", kc.cfg.KeycloakBaseUrl, kc.cfg.KeycloakRealm, userID)
}
