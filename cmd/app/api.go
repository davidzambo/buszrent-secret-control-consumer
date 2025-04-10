package main

import (
	"buszrent-secret-control-consumer/internal/login"
	"buszrent-secret-control-consumer/internal/tokens"
	"encoding/json"
	"fmt"
	"io"
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
	http.HandleFunc("/", methodHandler(http.MethodGet, h.handleRequestTokens))
	http.HandleFunc("/set-token", methodHandler(http.MethodPost, h.handleUpdateTokens))
	http.HandleFunc("/login", methodHandler(http.MethodPost, h.handleLogin))

	log.Println("Server is running on port :" + h.config.ApiPort)
	return http.ListenAndServe(":"+h.config.ApiPort, nil)
}

func (h *Handler) handleUpdateTokens(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Failed to read token update body: %v\n", err)
		w.WriteHeader(400)
		return
	}

	var content = map[string]string{}

	unmarshallErr := json.Unmarshal(body, &content)

	if unmarshallErr != nil {
		fmt.Printf("Failed to unmarshall token body: %v\n", err)
		w.WriteHeader(400)
		return
	}

	tokens.SetAccessToken(content["AccessToken"])
	fmt.Printf("New accessToken arrived: %s\n", content["refreshToken"])

	tokens.SetRefreshToken(content["RefreshToken"])
	fmt.Printf("New refreshToken arrived: %s\n", content["refreshToken"])

	w.WriteHeader(201)
}

func (h *Handler) handleRequestTokens(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("New login callback called with refreshToken: %s\n", r.URL.Query().Get("RefreshToken"))
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
	statusCode, err := login.ToWebFlotta(h.client, login.Config(h.config), webFlottaSso)
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
