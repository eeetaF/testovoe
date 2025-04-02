package route

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"

	models "testovoe/src/models/user"
	"testovoe/src/services/auth"
	user_service "testovoe/src/services/user"
	"testovoe/src/utils/sort"
)

// TODO move requests and responses to api package
type usrRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type registerResponse struct {
	AccessToken string `json:"access_token"`
}

type taskCompleteRequest struct {
	Type    int    `json:"type"`              // 0: referal, 1: telegram, 2: twitter
	Referal string `json:"referal,omitempty"` // only for referal
}

// UsersHandler handles /users/{id}/...
func UsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Trim(r.URL.Path, "/")
		parts := strings.Split(path, "/")

		if len(parts) < 2 || parts[0] != "users" {
			http.NotFound(w, r)
			return
		}

		userID, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			http.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}

		// Routing
		switch {
		case len(parts) == 3 && parts[2] == "status":
			UsersStatusHandler(db, userID).ServeHTTP(w, r)
			return

		case len(parts) == 4 && parts[2] == "task" && parts[3] == "complete":
			UsersTaskCompleteHandler(db, userID).ServeHTTP(w, r)
			return

		default:
			http.NotFound(w, r)
		}
	}
}

// UsersRegisterHandler POST /users/register
func UsersRegisterHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req usrRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		user, err := user_service.RegisterUser(db, req.Name, req.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		token, err := auth.GenerateToken(user.ID)
		if err != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&registerResponse{AccessToken: token})
	}
}

// UsersSignInHandler POST /users/sign_in
func UsersSignInHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req usrRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		user, err := user_service.FindUser(db, user_service.FindUserOpts{Name: req.Name})
		if err != nil {
			http.Error(w, "not found user with such credentials", http.StatusUnauthorized)
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			http.Error(w, "not found user with such credentials", http.StatusUnauthorized)
			return
		}

		token, err := auth.GenerateToken(user.ID)
		if err != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&registerResponse{AccessToken: token})
	}
}

// UsersStatusHandler GET /users/{id}/status
func UsersStatusHandler(db *sql.DB, userID int64) http.HandlerFunc {
	return auth.Middleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		user, err := user_service.FindUser(db, user_service.FindUserOpts{ID: userID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if user.ID != r.Context().Value("userID") {
			json.NewEncoder(w).Encode(models.UserToPublicUser(user))
			return
		}

		json.NewEncoder(w).Encode(user)
	})
}

// UsersLeaderboardHandler GET /users/leaderboard
func UsersLeaderboardHandler(db *sql.DB) http.HandlerFunc {
	return auth.Middleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// todo change offset to page, add paginator
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		limit := 10
		offset := 0

		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
		if o, err := strconv.Atoi(offsetStr); err == nil && o > 0 {
			offset = o
		}

		users, err := user_service.FindPublicUsers(db, &sort.UserOpts{Field: "sum_points", Order: "desc"}, limit, offset)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		if users == nil {
			users = []*models.PublicUsr{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string][]*models.PublicUsr{
			"users": users,
		})
	})
}

// UsersTaskCompleteHandler POST /users/{id}/task/complete
func UsersTaskCompleteHandler(db *sql.DB, userID int64) http.HandlerFunc {
	return auth.Middleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userIDFromToken, ok := r.Context().Value("userID").(int64)
		if !ok || userID != userIDFromToken {
			http.Error(w, "cannot complete task for another user", http.StatusForbidden)
			return
		}

		var req taskCompleteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		pointsRewarded, err := user_service.CompleteTask(db, userIDFromToken, &user_service.TaskComplete{
			Type:    req.Type,
			Referal: req.Referal,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"message":         "task completed",
			"points_rewarded": pointsRewarded,
		})
	})
}
