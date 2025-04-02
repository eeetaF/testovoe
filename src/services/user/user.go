package user_service

import (
	"database/sql"
	"fmt"
	"time"

	models "testovoe/src/models/user"
	"testovoe/src/utils"
	"testovoe/src/utils/sort"
)

type FindUserOpts struct {
	Name string
	ID   int64
}

func FindPublicUsers(db *sql.DB, opts *sort.UserOpts, limit, offset int) ([]*models.PublicUsr, error) {
	orderClause := sort.BuildClause(opts)

	rows, err := db.Query(`
		SELECT id, name, twitter_points, telegram_points, referal_points, created_unix
		FROM usr
		ORDER BY `+orderClause+`
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.PublicUsr
	for rows.Next() {
		var u models.PublicUsr
		err = rows.Scan(
			&u.ID,
			&u.Name,
			&u.TwitterPoints,
			&u.TelegramPoints,
			&u.ReferalPoints,
			&u.CreatedUnix,
		)
		if err != nil {
			return nil, err
		}
		u.SumPoints = u.TwitterPoints + u.TelegramPoints + u.ReferalPoints
		result = append(result, &u)
	}
	return result, nil
}

func FindUser(db *sql.DB, opts FindUserOpts) (*models.Usr, error) {
	var row *sql.Row
	if opts.ID > 0 {
		row = db.QueryRow(`
		SELECT id, name, password, referal_code, twitter_points, telegram_points, referal_points, created_unix, updated_unix
		FROM usr
		WHERE id = $1
	`, opts.ID)
	} else if len(opts.Name) > 0 {
		row = db.QueryRow(`
		SELECT id, name, password, referal_code, twitter_points, telegram_points, referal_points, created_unix, updated_unix
		FROM usr
		WHERE name = $1
	`, opts.Name)
	} else {
		return nil, fmt.Errorf("user not found")
	}

	var user models.Usr
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Password,
		&user.ReferalCode,
		&user.TwitterPoints,
		&user.TelegramPoints,
		&user.ReferalPoints,
		&user.CreatedUnix,
		&user.UpdatedUnix,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, err
	}

	user.SumPoints = user.TelegramPoints + user.TwitterPoints + user.ReferalPoints
	return &user, nil
}

func RegisterUser(db *sql.DB, name, rawPassword string) (*models.Usr, error) {
	if len(name) < 5 {
		return nil, fmt.Errorf("name too short")
	}
	if len(rawPassword) < 8 {
		return nil, fmt.Errorf("password too short")
	}

	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM usr WHERE name = $1)", name).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("user already exists")
	}

	referalCode, err := generateUniqueReferalCode(db, 7, 10)
	if err != nil {
		return nil, err
	}

	// todo improve algorithm of generating password
	hashedPassword, err := utils.HashPassword(rawPassword)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()

	var id int64
	err = db.QueryRow(`
		INSERT INTO usr (name, password, referal_code, created_unix, updated_unix)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, name, hashedPassword, referalCode, now, now).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &models.Usr{
		ID:          id,
		Name:        name,
		ReferalCode: referalCode,
		CreatedUnix: now,
		UpdatedUnix: now,
	}, nil
}

func generateUniqueReferalCode(db *sql.DB, length int, maxAttempts int) (string, error) {
	for attempts := 0; attempts < maxAttempts; attempts++ {
		code, err := utils.GenerateReferalCode(length)
		if err != nil {
			return "", err
		}

		var exists bool
		err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM usr WHERE referal_code = $1)", code).Scan(&exists)
		if err != nil {
			return "", err
		}
		if !exists {
			return code, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique referal code")
}
