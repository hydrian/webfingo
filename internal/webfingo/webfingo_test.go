package webfingo

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockDatabase implements the minimal Database interface needed for testing
type MockDatabase struct{}

func (m MockDatabase) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	// Return a mock user for testing
	if email == "john@example.com" {
		return &User{
			ID:        "user-1",
			Username:  "johndoe",
			Email:     "john@example.com",
			RealmID:   "realm-1",
			RealmName: "master",
		}, nil
	}
	return nil, nil
}

func (m MockDatabase) Close() error {
	return nil
}

func TestHandleWebfingerRequest(t *testing.T) {
	// Create a mock database
	db := MockDatabase{}

	// Create a mock KeycloakConfig
	keycloakConfig := KeycloakConfig{
		KeycloakHost: "example.com",
	}

	// Create a request with the required query parameter
	req, err := http.NewRequest("GET", "/.well-known/webfinger?resource=acct:john@example.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function with the KeycloakConfig
	HandleWebfingerRequest(rr, req, db, keycloakConfig)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the content type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, "application/json")
	}

	// Parse the response body
	var response WebFingerResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("Error parsing response body: %v", err)
	}

	// Verify the response content
	expectedSubject := "acct:john@example.com"
	if response.Subject != expectedSubject {
		t.Errorf("handler returned unexpected subject: got %v want %v", response.Subject, expectedSubject)
	}

	// Verify the links
	if len(response.Links) != 1 {
		t.Errorf("handler returned unexpected number of links: got %v want %v", len(response.Links), 1)
	}

	expectedRel := "http://openid.net/specs/connect/1.0/issuer"
	if response.Links[0].Rel != expectedRel {
		t.Errorf("handler returned unexpected link rel: got %v want %v", response.Links[0].Rel, expectedRel)
	}

	expectedHref := "https://example.com/realms/master"
	if response.Links[0].Href != expectedHref {
		t.Errorf("handler returned unexpected link href: got %v want %v", response.Links[0].Href, expectedHref)
	}
}
