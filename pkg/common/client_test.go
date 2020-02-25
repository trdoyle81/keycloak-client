package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	jsoniter "github.com/json-iterator/go"

	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	"github.com/stretchr/testify/assert"
)

const (
	RealmsGetPath                     = "/auth/admin/realms/%s"
	RealmsCreatePath                  = "/auth/admin/realms"
	RealmsDeletePath                  = "/auth/admin/realms/%s"
	UserCreatePath                    = "/auth/admin/realms/%s/users"
	UserDeletePath                    = "/auth/admin/realms/%s/users/%s"
	GroupGetPath                      = "/auth/admin/realms/%s/groups"
	GroupCreatePath                   = "/auth/admin/realms/%s/groups"
	GroupGetDefaults                  = "/auth/admin/realms/%s/default-groups"
	GroupMakeDefaultPath              = "/auth/admin/realms/%s/default-groups/%s"
	GroupCreateClientRole             = "/auth/admin/realms/%s/groups/%s/role-mappings/clients/%s"
	GroupGetClientRoles               = "/auth/admin/realms/%s/groups/%s/role-mappings/clients/%s"
	GroupGetAvailableClientRoles      = "/auth/admin/realms/%s/groups/%s/role-mappings/clients/%s/available"
	AuthenticationFlowUpdateExecution = "/auth/admin/realms/%s/authentication/flows/%s/executions"
	TokenPath                         = "/auth/realms/master/protocol/openid-connect/token" // nolint
)

func getDummyRealm() *v1alpha1.KeycloakRealm {
	return &v1alpha1.KeycloakRealm{
		Spec: v1alpha1.KeycloakRealmSpec{
			Realm: &v1alpha1.KeycloakAPIRealm{
				ID:          "dummy",
				Realm:       "dummy",
				Enabled:     false,
				DisplayName: "dummy",
			},
		},
	}
}

func getDummyUser() *v1alpha1.KeycloakAPIUser {
	return &v1alpha1.KeycloakAPIUser{
		ID:            "dummy",
		UserName:      "dummy",
		FirstName:     "dummy",
		LastName:      "dummy",
		EmailVerified: false,
		Enabled:       false,
	}
}

func TestClient_CreateRealm(t *testing.T) {
	// given
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, RealmsCreatePath, req.URL.Path)
		w.WriteHeader(201)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := Client{
		requester: server.Client(),
		URL:       server.URL,
		token:     "dummy",
	}

	realm := getDummyRealm()

	// when
	err := client.CreateRealm(realm)

	// then
	// no error expected
	// correct path expected on httptest server
	assert.NoError(t, err)
}

func TestClient_DeleteRealmRealm(t *testing.T) {
	// given
	realm := getDummyRealm()

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, fmt.Sprintf(RealmsDeletePath, realm.Spec.Realm.Realm), req.URL.Path)
		w.WriteHeader(204)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := Client{
		requester: server.Client(),
		URL:       server.URL,
		token:     "dummy",
	}

	// when
	err := client.DeleteRealm(realm.Spec.Realm.Realm)

	// then
	// correct path expected on httptest server
	assert.NoError(t, err)
}

func TestClient_CreateUser(t *testing.T) {
	// given
	user := getDummyUser()
	realm := getDummyRealm()

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, fmt.Sprintf(UserCreatePath, realm.Spec.Realm.Realm), req.URL.Path)
		w.WriteHeader(201)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := Client{
		requester: server.Client(),
		URL:       server.URL,
		token:     "dummy",
	}

	// when
	err := client.CreateUser(user, realm.Spec.Realm.Realm)

	// then
	// correct path expected on httptest server
	assert.NoError(t, err)
}

func TestClient_DeleteUser(t *testing.T) {
	// given
	user := getDummyUser()
	realm := getDummyRealm()

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, fmt.Sprintf(UserDeletePath, realm.Spec.Realm.Realm, user.ID), req.URL.Path)
		w.WriteHeader(204)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := Client{
		requester: server.Client(),
		URL:       server.URL,
		token:     "dummy",
	}

	// when
	err := client.DeleteUser(user.ID, realm.Spec.Realm.Realm)

	// then
	// correct path expected on httptest server
	assert.NoError(t, err)
}

func TestClient_GetRealm(t *testing.T) {
	// given
	realm := getDummyRealm()

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, fmt.Sprintf(RealmsGetPath, realm.Spec.Realm.Realm), req.URL.Path)
		assert.Equal(t, req.Method, http.MethodGet)
		json, err := jsoniter.Marshal(realm.Spec.Realm)
		assert.NoError(t, err)

		size, err := w.Write(json)
		assert.NoError(t, err)
		assert.Equal(t, size, len(json))

		w.WriteHeader(204)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := Client{
		requester: server.Client(),
		URL:       server.URL,
		token:     "dummy",
	}

	// when
	newRealm, err := client.GetRealm(realm.Spec.Realm.Realm)

	// then
	// correct path expected on httptest server
	assert.NoError(t, err)

	// returned realm must equal dummy realm
	assert.Equal(t, realm.Spec.Realm.Realm, newRealm.Spec.Realm.Realm)
}

func TestClient_ListRealms(t *testing.T) {
	// given
	realm := getDummyRealm()

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, RealmsCreatePath, req.URL.Path)
		assert.Equal(t, req.Method, http.MethodGet)
		var list []*v1alpha1.KeycloakRealm
		list = append(list, realm)
		json, err := jsoniter.Marshal(list)
		assert.NoError(t, err)

		size, err := w.Write(json)
		assert.NoError(t, err)
		assert.Equal(t, size, len(json))

		w.WriteHeader(204)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := Client{
		requester: server.Client(),
		URL:       server.URL,
		token:     "dummy",
	}

	// when
	realms, err := client.ListRealms()

	// then
	// correct path expected on httptest server
	assert.NoError(t, err)

	// exactly one realms must be returned
	assert.Len(t, realms, 1)
}

func TestClient_FindGroupByName(t *testing.T) {
	const (
		existingGroupName string = "group"
		existingGroupID   string = "12345"
	)
	realm := getDummyRealm()

	handle := withPathAssertionBody(
		t,
		200,
		fmt.Sprintf(GroupGetPath, realm.Spec.Realm.Realm),
		[]*Group{
			&Group{
				ID:   existingGroupID,
				Name: existingGroupName,
			},
		},
	)

	request := func(c *Client) {
		// when the group exists
		foundGroup, err := c.FindGroupByName(existingGroupName, realm.Spec.Realm.Realm)
		// then return the group instance
		assert.NoError(t, err)
		assert.NotNil(t, foundGroup)
		assert.Equal(t, existingGroupID, foundGroup.ID)

		// when the group doesn't exist
		notFoundGroup, err := c.FindGroupByName("not-existing", "dummy")
		// then return `nil`
		assert.NoError(t, err)
		assert.Nil(t, notFoundGroup)
	}

	testClientHTTPRequest(handle, request)
}

func TestClient_CreateGroup(t *testing.T) {
	realm := getDummyRealm()
	const (
		createdGroupName string = "dummy-group"
		createdGroupID   string = "12345"
	)

	handle := withMethodSelection(t, map[string]http.HandlerFunc{
		// When the client requests the groups to find the newly
		// created one
		http.MethodGet: withJSON(t, []*Group{
			&Group{
				ID:   createdGroupID,
				Name: createdGroupName,
			},
		}, 200),

		// When the client requests to create the group
		http.MethodPost: withPathAssertion(t, 201, fmt.Sprintf(GroupCreatePath, realm.Spec.Realm.Realm)),
	})

	request := func(c *Client) {
		groupID, err := c.CreateGroup(createdGroupName, realm.Spec.Realm.Realm)
		assert.NoError(t, err)
		assert.Equal(t, createdGroupID, groupID)
	}

	testClientHTTPRequest(handle, request)
}

func TestClient_MakeGroupDefault(t *testing.T) {
	const groupID string = "12345"
	realm := getDummyRealm()

	handle := withMethodSelection(t, map[string]http.HandlerFunc{
		http.MethodGet: withPathAssertionBody(t, 200, fmt.Sprintf(GroupGetDefaults, realm.Spec.Realm.Realm), []*Group{}),
		http.MethodPut: withPathAssertion(t, 200, fmt.Sprintf(GroupMakeDefaultPath, realm.Spec.Realm.Realm, groupID)),
	})

	request := func(c *Client) {
		err := c.MakeGroupDefault(groupID, realm.Spec.Realm.Realm)
		assert.NoError(t, err)
	}

	testClientHTTPRequest(handle, request)
}

func TestClient_CreateGroupClientRole(t *testing.T) {
	realm := getDummyRealm()
	const (
		groupID  string = "12345"
		clientID string = "client-12345"
	)

	with := withPathAssertion(t, 201, fmt.Sprintf(GroupCreateClientRole, realm.Spec.Realm.Realm, groupID, clientID))
	when := func(c *Client) {
		err := c.CreateGroupClientRole(&v1alpha1.KeycloakUserRole{}, realm.Spec.Realm.Realm, clientID, groupID)
		assert.NoError(t, err)
	}

	testClientHTTPRequest(with, when)
}

func TestClient_ListGroupClientRoles(t *testing.T) {
	realm := getDummyRealm()
	const (
		groupID  = "group12345"
		clientID = "client12345"
	)

	testClientHTTPRequest(
		withPathAssertion(t, 200, fmt.Sprintf(GroupGetClientRoles, realm.Spec.Realm.Realm, groupID, clientID)),
		func(c *Client) {
			_, err := c.ListGroupClientRoles(
				realm.Spec.Realm.Realm, clientID, groupID)

			assert.NoError(t, err)
		},
	)
}

func TestClient_ListAvailableGroupClientRoles(t *testing.T) {
	realm := getDummyRealm()
	const (
		clientID = "client12345"
		groupID  = "group12345"
	)

	testClientHTTPRequest(
		withPathAssertion(t, 200, fmt.Sprintf(GroupGetAvailableClientRoles, realm.Spec.Realm.Realm, clientID, groupID)),
		func(c *Client) {
			_, err := c.ListAvailableGroupClientRoles(realm.Spec.Realm.Realm, groupID, clientID)
			assert.NoError(t, err)
		},
	)
}

func TestClient_UpdateAuthenticationExecutionForFlow(t *testing.T) {
	realm := getDummyRealm()

	const (
		flowAlias string = "test flow"
	)

	requestPath := fmt.Sprintf(AuthenticationFlowUpdateExecution, realm.Spec.Realm.Realm, flowAlias)

	testClientHTTPRequest(
		withPathAssertion(t, 200, requestPath),
		func(c *Client) {
			err := c.UpdateAuthenticationExecutionForFlow(flowAlias, realm.Spec.Realm.Realm, &v1alpha1.AuthenticationExecutionInfo{})
			assert.NoError(t, err)
		},
	)
}

// Utility function to create a test server, register a given handler and perform
// a client function to be tested
func testClientHTTPRequest(
	testHandler http.HandlerFunc,
	testRequest func(c *Client),
) {
	handler := http.HandlerFunc(testHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	client := Client{
		requester: server.Client(),
		URL:       server.URL,
		token:     "not set",
	}

	testRequest(&client)
}

func respondWithJSON(body interface{}, w http.ResponseWriter) (int, error) {
	json, err := json.Marshal(body)
	if err != nil {
		return 0, err
	}

	return w.Write(json)
}

func withJSON(t *testing.T, body interface{}, status int) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		_, err := respondWithJSON(body, w)
		w.WriteHeader(status)
		assert.NoError(t, err)
	}
}

func withPathAssertion(t *testing.T, status int, expectedPath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, expectedPath, req.URL.Path)
		w.WriteHeader(status)
	}
}

func withMethodSelection(t *testing.T, byMethod map[string]http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		methodFunc, ok := byMethod[req.Method]
		assert.True(t, ok)

		methodFunc(w, req)
	}
}

func withPathAssertionBody(t *testing.T, status int, expectedPath string, body interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, expectedPath, req.URL.Path)
		withJSON(t, body, status)(w, req)
	}
}

func testClientMethod(t *testing.T, method func(*Client, ...[]interface{}) (T, error)) func(*Client) {
	return func(c *Client) {
		_, err := method(c)
		assert.NoError(t, err)
	}
}

func TestClient_login(t *testing.T) {
	// given
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, TokenPath, req.URL.Path)
		assert.Equal(t, req.Method, http.MethodPost)

		response := v1alpha1.TokenResponse{
			AccessToken: "dummy",
		}

		json, err := jsoniter.Marshal(response)
		assert.NoError(t, err)

		size, err := w.Write(json)
		assert.NoError(t, err)
		assert.Equal(t, size, len(json))

		w.WriteHeader(204)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	client := Client{
		requester: server.Client(),
		URL:       server.URL,
		token:     "not set",
	}

	// when
	err := client.login("dummy", "dummy")

	// then
	// token must be set on the client now
	assert.NoError(t, err)
	assert.Equal(t, client.token, "dummy")
}
