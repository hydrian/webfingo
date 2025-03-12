package webfingo

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

// WebFingerResponse represents the JSON structure for WebFinger responses
type WebFingerResponse struct {
	Subject string              `json:"subject"`
	Links   []WebFingerLinkItem `json:"links"`
}

// WebFingerLinkItem represents a link in the WebFinger response
type WebFingerLinkItem struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

func HandleWebfingerRequest(w http.ResponseWriter, r *http.Request, db Database, keycloakConfig KeycloakConfig, logger *Logger) {
	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Get resource parameter
	resource := r.URL.Query().Get("resource")
	if resource == "" {
		logger.Printf("WebFinger request missing required 'resource' parameter from %s", r.RemoteAddr)
		http.Error(w, "Resource parameter is required", http.StatusBadRequest)
		return
	}

	// Parse the resource parameter (expected format: acct:user@domain)
	if !strings.HasPrefix(resource, "acct:") {
		logger.Printf("WebFinger request with invalid resource format: %s from %s", resource, r.RemoteAddr)
		http.Error(w, "Invalid resource format", http.StatusBadRequest)
		return
	}

	// Extract email from resource
	email := strings.TrimPrefix(resource, "acct:")
	email, err := url.QueryUnescape(email)
	if err != nil {
		logger.Printf("Error unescaping email from resource '%s': %v", resource, err)
		http.Error(w, "Invalid resource format", http.StatusBadRequest)
		return
	}

	// Get user from database
	user, err := db.GetUserByEmail(r.Context(), email)
	if err != nil {
		logger.Printf("Error retrieving user by email '%s': %v", email, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Use the configured host
	baseURL := "https://" + keycloakConfig.KeycloakHost

	// Create WebFinger response
	response := WebFingerResponse{
		Subject: resource,
		Links: []WebFingerLinkItem{
			{
				Rel:  "http://openid.net/specs/connect/1.0/issuer",
				Href: baseURL + "/realms/" + user.RealmName,
			},
		},
	}

	// Encode response as JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Printf("Error encoding WebFinger response for '%s': %v", email, err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
