package models

import (
	"context"

	"github.com/kaleb-white/letthemknow/message-server/utils"
)

type WriteEnabledStore[T any] interface {
	Init(context.Context) error
	Write(context.Context, *T) (uint64, utils.StoreError)
}

type ReadEnabledStore[T any] interface {
	Init(context.Context) error
	Read(context.Context, uint64) (T, utils.StoreError)
}

type DeleteEnabledStore[T any] interface {
	Init(context.Context) error
	Delete(context.Context, uint64) (uint64, utils.StoreError)
}


