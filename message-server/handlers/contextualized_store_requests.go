package handlers

import (
	"context"

	"github.com/kaleb-white/letthemknow/message-server/models"
	"github.com/kaleb-white/letthemknow/message-server/utils"
)

func ReadFromStoreWithContext[
	N any, 
	S models.ReadEnabledStore[N],
](
	id uint64,
	store S,
	getenv func(string) string,
) (N, utils.StoreError) {
	timeoutDur := store.ReadTimeout(getenv)	
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDur)
	defer cancel()
	return store.Read(ctx, id)
}

func DeleteFromStoreWithContext[
	N any, 
	S models.ReadEnabledStore[N],
](
	id uint64,
	store S,
	getenv func(string) string,
) (N, utils.StoreError) {
	timeoutDur := store.ReadTimeout(getenv)	
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDur)
	defer cancel()
	return store.Read(ctx, id)
}

func WriteToStoreWithContext[
	N any, 
	S models.WriteEnabledStore[N],
](
	obj *N,
	store S,
	getenv func(string) string,
) (uint64, utils.StoreError) {
	timeoutDur := store.WriteTimeout(getenv)	
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDur)
	defer cancel()
	return store.Write(ctx, obj)
}

