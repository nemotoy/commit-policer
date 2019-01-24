package handler

import (
	"fmt"
	"net/http"
)

// WebhookHandler ...
type WebhookHandler struct {
}

// NewWebhookHandler ...
func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{}
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, fmt.Sprintf("webhook"))
}
