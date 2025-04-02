package sort

import "strings"

type UserOpts struct {
	Field string
	Order string
}

func BuildClause(opts *UserOpts) string {
	field := "sum_points"
	order := "desc"

	if opts != nil {
		if val, ok := allowedSortUserFields[opts.Field]; ok {
			field = val
		}
		if opts.Order == "asc" || opts.Order == "desc" {
			order = opts.Order
		}
	}

	return field + " " + strings.ToUpper(order)
}

var allowedSortUserFields = map[string]string{
	"sum_points":      "(twitter_points + telegram_points)",
	"created_unix":    "created_unix",
	"twitter_points":  "twitter_points",
	"telegram_points": "telegram_points",
}
