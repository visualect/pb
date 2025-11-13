//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server Petstore server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
// @BasePath	/swagger/*
package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/visualect/pb/docs"
	"github.com/visualect/pb/internal/handlers"
	"github.com/visualect/pb/internal/handlers/middleware"
	"github.com/visualect/pb/internal/repo"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("failed to conect to database: %s\n", err)
	} else {
		log.Println("successfully connected to database")
	}

	articlesRepo := repo.New(dbpool)
	h := handlers.New(articlesRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/articles", h.GetArticles)
	mux.HandleFunc("GET /v1/articles/{id}", h.GetArticle)
	mux.HandleFunc("POST /v1/articles", middleware.AuthRequired(h.CreateArticle))
	mux.HandleFunc("DELETE /v1/articles/{id}", middleware.AuthRequired(h.DeleteArticle))
	mux.HandleFunc("PATCH /v1/articles/{id}", middleware.AuthRequired(h.UpdateArticle))
	mux.HandleFunc("GET /swagger/", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json")))

	c := cors.New(cors.Options{
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPatch},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		MaxAge:         86400,
	})
	handler := c.Handler(mux)

	log.Println("running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
