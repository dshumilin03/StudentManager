package domain

type Group struct {
	Id          int64  `json:"id"`
	GroupNumber string `json:"group_number" env-required:"true"`
}
