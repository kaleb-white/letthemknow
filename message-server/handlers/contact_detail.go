package handlers

import (
	"net/http"

	"github.com/kaleb-white/letthemknow/message-server/models"
)

const SOURCE_CONTACT_DETAIL_HANDLER string = "contact_detail_handler"

func HandleReadContactDetail(contactDetailStore *models.ContactDetailStore, getenv func(string) string) http.Handler {
	return  genericReadHandler(contactDetailStore, getenv, SOURCE_CONTACT_DETAIL_HANDLER, "ContactDetail")
}
