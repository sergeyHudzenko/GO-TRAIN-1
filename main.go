package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/sergeyHudzenko/go-rss-aggregator/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}


func main() { 
	// feed, err := urlToFeed("https://www.targaresources.com/rss/news-releases.xml")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(feed)

  godotenv.Load(".env")
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")
	

	if port == "" {
		log.Fatal("PORT is not found")
	}
	
	if dbURL == "" {
		log.Fatal("DB_URL is not found")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err !=nil {
		log.Fatal("Can't connect to database:", err)
	}
 
	db := database.New(conn)
	apiCfg := apiConfig{
		DB: db,
	}
	
	go startScraping(db, 10, time.Minute)

	router := chi.NewRouter();

	router.Use(cors.Handler (cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300, 
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeed)

	v1Router.Post("/feed_follow", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follow", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Delete("/feed_follow/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollows))

	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))
	


	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr: ":"+ port,
	}

	log.Printf("Server starting on port %v", port)

	errSvr := srv.ListenAndServe()
	if errSvr != nil {
		log.Fatal(errSvr)
	}

	fmt.Println("PORT:", port)
}