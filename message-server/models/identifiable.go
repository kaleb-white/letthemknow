package models

type Identifiable interface {
	SetId(uint64)
	GetId() uint64
}
