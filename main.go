package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/visualect/pb/repo"
)

type ArticlesHanlder struct {
	repo repo.ArticlesRepository
}

type CreateArticleRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

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
		log.Fatalf("failed to conect to databse: %s\n", err)
	} else {
		log.Println("successfully connected to database")
	}

	articlesRepo := repo.New(dbpool)
	handler := &ArticlesHanlder{
		repo: articlesRepo,
	}

	http.HandleFunc("GET /v1/getArticles", handler.getArticles)
	http.HandleFunc("GET /v1/getArticle/{id}", handler.getArticle)
	// http.HandleFunc("POST /v1/createtArticle", createArticle)
	// http.HandleFunc("DELETE /v1/deleteArticle/{id}", deleteArticle)
	// http.HandleFunc("PATCH /v1/updateArticle/{id}", updateArticle)

	log.Println("running on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func (a *ArticlesHanlder) getArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := a.repo.GetArticles(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(articles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *ArticlesHanlder) getArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	article, err := a.repo.GetArticle(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
