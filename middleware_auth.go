package main

import (
	"fmt"
	"net/http"

	"github.com/sergeyHudzenko/go-rss-aggregator/internal/auth"
	"github.com/sergeyHudzenko/go-rss-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIkey(r.Header)
		if err != nil {
			respondWithErr(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

	 user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
	 if err != nil {
		respondWithErr(w, 400, fmt.Sprintf("Can't get user: %v", err))
		return
	 }

	 handler(w, r, user)
	}
}