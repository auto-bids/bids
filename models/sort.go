package models

type Sort struct {
	Order string `validate:"omitempty,oneof=desc asc"`
	By    string `validate:"omitempty,oneof=price year engine_capacity power"`
}
