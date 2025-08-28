package handlers

import (
	"encoding/json"
	"net/http"
)

// GetProfile is a protected endpoint that returns user profile
func GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user info from context (set by AuthMiddleware)
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// In a real app, you'd fetch user from database using userID
	user, exists := users[email]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Verify that the user ID matches
	if user.ID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(user)
}
