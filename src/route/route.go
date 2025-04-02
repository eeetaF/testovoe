package route

import (
	"database/sql"
	"log"
	"net/http"
)

func InitRoutes(database *sql.DB) {
	http.HandleFunc("/users/register", UsersRegisterHandler(database))
	http.HandleFunc("/users/sign_in", UsersSignInHandler(database))
	http.HandleFunc("/users/leaderboard", UsersLeaderboardHandler(database))
	http.HandleFunc("/users/", UsersHandler(database))

	log.Println("Server is running on 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Can't run server: %v", err)
	}
}
