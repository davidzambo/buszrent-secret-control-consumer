package main

import (
	"buszrent-secret-control-consumer/internal/tokens"
	"log"
	"net/http"
)

type Handler struct {
	client *http.Client
	config Config
	ssoUrl string
}

func startHttpServer(h *Handler) error {
	http.HandleFunc("/get_tokens", h.handleRequestTokens)

	log.Println("Server is running on port :" + h.config.ApiPort)
	return http.ListenAndServe(":"+h.config.ApiPort, nil)
}

func (h *Handler) handleRequestTokens(w http.ResponseWriter, r *http.Request) {
	if err := tokens.FetchNewTokens(r.URL.Query().Get("RefreshToken"), h.config.ClientID, h.ssoUrl); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
