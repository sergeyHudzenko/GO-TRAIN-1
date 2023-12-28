package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/sergeyHudzenko/go-rss-aggregator/internal/database"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"` 
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithErr(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: user.ID,
		FeedID: params.FeedID,
	})

	if err != nil {
		respondWithErr(w, 400, fmt.Sprintf("Err with creating feed follow: %s", err))
		return
	}
	
	respondWithJSON(w, 201, databaseFeedFollowToFeedFollow(feedFollow))
} 

func (apiCfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := apiCfg.DB.GetFeedFollow(r.Context(), user.ID);
	if err != nil {
		respondWithErr(w, 404, fmt.Sprintf("Feed follows not found: %s", err))
		return
	}

	respondWithJSON(w, 200, databaseFeedFollowsToFeedFollows(feedFollows))
}

func (apiCfg *apiConfig) handlerDeleteFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	 feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	 feedFollowID, err := uuid.Parse(feedFollowIDStr)

	 if err != nil {
		respondWithErr(w, 400, fmt.Sprintf("Couldn't parse feed follow ID: %v", err))
		return
	 }

	 err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	 })

	 if err != nil {
		respondWithErr(w, 400, fmt.Sprintf("Couldn't delete feed follow: %v", err))
		return
	 }

	 respondWithJSON(w, 200, struct{}{})
}