package main

import (
	"buszrent-secret-control-consumer/internal/login"
	"buszrent-secret-control-consumer/internal/tokens"
	"encoding/json"
	"log"
	"net/http"
)

type Handler struct {
	client *http.Client
	config Config
	ssoUrl string
}

func methodHandler(method string, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handlerFunc(w, r)
	}
}

func startHttpServer(h *Handler) error {
	http.HandleFunc("/get_tokens", methodHandler(http.MethodGet, h.handleRequestTokens))
	http.HandleFunc("/login", methodHandler(http.MethodPost, h.handleLogin))

	log.Println("Server is running on port :" + h.config.ApiPort)
	return http.ListenAndServe(":"+h.config.ApiPort, nil)
}

func (h *Handler) handleRequestTokens(w http.ResponseWriter, r *http.Request) {
	statusCode, err := tokens.FetchNewTokens(h.config.Slack.DevAPIToken, h.config.Slack.DevChannel, r.URL.Query().Get("RefreshToken"), h.config.ClientID, webFlottaSso)
	if err != nil {
		if statusCode != 0 {
			w.WriteHeader(statusCode)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	w.WriteHeader(statusCode)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	statusCode, err := login.LoginToWebFlotta(h.client, login.Config(h.config), webFlottaSso)
	if err != nil {
		if statusCode != 0 {
			w.WriteHeader(statusCode)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	w.WriteHeader(statusCode)
}
