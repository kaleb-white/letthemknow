package models

type IdResponse struct {
	Id uint64
}

type AddedIdsResponse struct {
	Added []uint64
}

type RemovedIdsResponse struct {
	Removed []uint64
}

type AddContactsToListRequest struct {
	Ids []uint64
}

type RemoveContactsFromListRequest struct {
	Ids []uint64
}
