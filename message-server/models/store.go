package models

import (
	"context"
	"time"

	"github.com/kaleb-white/letthemknow/message-server/utils"
)

type WriteEnabledStore[T any] interface {
	Init(context.Context) error
	Write(context.Context, *T) (uint64, utils.StoreError)
	// The environment configured timeout. Pass getenv function.
	WriteTimeout(func (string) string) time.Duration	
}

type ReadEnabledStore[T any] interface {
	Init(context.Context) error
	Read(context.Context, uint64) (T, utils.StoreError)
	// The environment configured timeout. Pass getenv function.
	ReadTimeout(func (string) string) time.Duration
}

type DeleteEnabledStore[T any] interface {
	Init(context.Context) error
	Delete(context.Context, uint64) (uint64, utils.StoreError)
	// The environment configured timeout. Pass getenv function.
	DeleteTimeout(func (string) string) time.Duration
}
