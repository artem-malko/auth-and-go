package users

import (
	"net/http"

	"github.com/artem-malko/auth-and-go/api/response"
)

type stubResponse struct {
	Data string `json:"data"`
}

// Stub just stub handler
func (h *handlers) Stub(w http.ResponseWriter, _ *http.Request) {
	response.OK(w, &stubResponse{
		Data: "Just stub!",
	})
}
