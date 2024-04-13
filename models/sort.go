package models

type Sort struct {
	Order int64  `validate:"required,oneof=1 -1"`
	By    string `validate:"required,oneof=price year"`
}
