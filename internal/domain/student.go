package domain

type Student struct {
	Id          int64  `json:"id"`
	FullName    string `json:"full_name" env-required:"true"`
	Age         int    `json:"age" env-required:"true"`
	GroupNumber string `json:"group_number"`
	Email       string `json:"email" env-required:"true"`
}
