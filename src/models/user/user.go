package models

type Usr struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Password       string `json:"-"`
	ReferalCode    string `json:"referal_code"`
	TwitterPoints  int64  `json:"twitter_points"`
	TelegramPoints int64  `json:"telegram_points"`
	ReferalPoints  int64  `json:"referal_points"`
	SumPoints      int64  `json:"sum_points"`
	CreatedUnix    int64  `json:"created_unix"`
	UpdatedUnix    int64  `json:"updated_unix"`
}

type PublicUsr struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	TwitterPoints  int64  `json:"twitter_points"`
	TelegramPoints int64  `json:"telegram_points"`
	ReferalPoints  int64  `json:"referal_points"`
	SumPoints      int64  `json:"sum_points"`
	CreatedUnix    int64  `json:"created_unix"`
}

// UserToPublicUser todo move to builder
func UserToPublicUser(user *Usr) *PublicUsr {
	return &PublicUsr{
		ID:             user.ID,
		Name:           user.Name,
		TwitterPoints:  user.TwitterPoints,
		TelegramPoints: user.TelegramPoints,
		ReferalPoints:  user.ReferalPoints,
		SumPoints:      user.SumPoints,
		CreatedUnix:    user.CreatedUnix,
	}
}
