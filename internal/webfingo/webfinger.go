package webfingo

import (
	"net/http"
	"strings"
	"net/url"
	"encoding/json"
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

func HandleWebfingerRequest (w http.ResponseWriter, r *http.Request, db Database) {
	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Get resource parameter
	resource := r.URL.Query().Get("resource")
	if resource == "" {
		http.Error(w, "Resource parameter is required", http.StatusBadRequest)
		return
	}

	// Parse the resource parameter (expected format: acct:user@domain)
	if !strings.HasPrefix(resource, "acct:") {
		http.Error(w, "Invalid resource format", http.StatusBadRequest)
		return
	}

	// Extract email from resource
	email := strings.TrimPrefix(resource, "acct:")
	email, err := url.QueryUnescape(email)
	if err != nil {
		http.Error(w, "Invalid resource format", http.StatusBadRequest)
		return
	}

	// Get user from database
	user, err := db.GetUserByEmail(r.Context(), email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Construct base URL from request
	baseURL := "https://" + r.Host

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
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
