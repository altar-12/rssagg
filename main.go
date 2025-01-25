package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/altar-12/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("No PORT present in the environment")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("No DB_URL present in the environment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Unable to connect to the database")
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}
	go scraper(apiCfg.DB, 10, time.Minute)

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/error", handlerError)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.authenticate(apiCfg.handlerFetchUserWithApiKey))
	//v1Router.Delete("/users/{id}", apiCfg.handlerDeleteUser)
	v1Router.Post("/feed", apiCfg.authenticate(apiCfg.handlerCreateFeed))
	v1Router.Get("/feed", apiCfg.handlerGetFeeds)
	v1Router.Post("/feed_follows", apiCfg.authenticate(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.authenticate(apiCfg.handlerFetchUserFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.authenticate(apiCfg.handlerDeleteFeedFollow))
	v1Router.Get("/posts", apiCfg.authenticate(apiCfg.handlerFetchUserPosts))

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	fmt.Println("Server starting on port", portString)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
