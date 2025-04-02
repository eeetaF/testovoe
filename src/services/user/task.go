package user_service

import (
	"database/sql"
	"fmt"
)

type TaskComplete struct {
	Type    int
	Referal string
}

func CompleteTask(db *sql.DB, userID int64, task *TaskComplete) (int, error) {
	switch task.Type {
	case 1: // Telegram
		_, err := db.Exec(`UPDATE usr SET telegram_points = telegram_points + 2 WHERE id = $1`, userID)
		return 2, err

	case 2: // Twitter
		_, err := db.Exec(`UPDATE usr SET twitter_points = twitter_points + 1 WHERE id = $1`, userID)
		return 1, err

	case 3: // Referal
		if task.Referal == "" {
			return 0, fmt.Errorf("referal code required")
		}

		// 1. check if already used
		var alreadyUsed bool
		err := db.QueryRow(`SELECT EXISTS(SELECT 1 FROM usr WHERE id = $1 AND referal_points > 0)`, userID).Scan(&alreadyUsed)
		if err != nil {
			return 0, err
		}
		if alreadyUsed {
			return 0, fmt.Errorf("referal already used")
		}

		// 2. find refered user
		var referedUserID int64
		err = db.QueryRow(`SELECT id FROM usr WHERE referal_code = $1`, task.Referal).Scan(&referedUserID)
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("invalid referal code")
		}
		if err != nil {
			return 0, err
		}
		if referedUserID == userID {
			return 0, fmt.Errorf("cannot refer yourself")
		}

		// 3. apply +10 referal points
		_, err = db.Exec(`UPDATE usr SET referal_points = referal_points + 10 WHERE id = $1`, userID)
		return 10, err

	default:
		return 0, fmt.Errorf("invalid task type")
	}
}
