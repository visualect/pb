package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/visualect/pb/internal/repo"
)

type ArticlesHandler struct {
	repo repo.ArticlesRepository
}

func New(repo repo.ArticlesRepository) *ArticlesHandler {
	return &ArticlesHandler{repo}
}

// GetArticles gets a list of all articles
//
//	@Summary		List all articles
//	@Description	Get all articlesGet all articles
//	@Tags			articles
//	@Produce		json
//	@Success		200	{array}		models.Article
//	@Failure		400	{string}	string	"Bad Request"
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Router			/v1/articles [get]
func (a *ArticlesHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := a.repo.GetArticles(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=300")
	err = json.NewEncoder(w).Encode(articles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetArticle get an article by ID
//
//	@Summary		Get article by ID
//	@Description	Get article by his ID
//	@Tags			articles
//	@Param			id	path	string	true	"Article ID"
//	@Produce		json
//	@Success		200	{object}	models.Article
//	@Failure		400	{string}	string	"Bad Request"
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Router			/v1/articles/{id} [get]
func (a *ArticlesHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	article, err := a.repo.GetArticle(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	err = json.NewEncoder(w).Encode(article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateArticle creates a new article
//
//	@Summary		Create a new article
//	@Description	Create a new article with the provided data
//	@Tags			articles
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		201	{object}	models.Article
//	@Failure		400	{string}	string	"Bad Request"
//	@Failure		403	{string}	string	"Forbidden"
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Router			/v1/articles [post]
func (a *ArticlesHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var b repo.CreateArticleRequest
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if strings.Trim(b.Body, " ") == "" {
		http.Error(w, "body is required", http.StatusBadRequest)
		return
	}

	if strings.Trim(b.Title, " ") == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	article, err := a.repo.CreateArticle(r.Context(), b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteArticle deletes an article by ID
//
//	@Summary		Delete article by ID
//	@Description	Delete article by his ID
//	@Tags			articles
//	@Param			id	path	string	true	"Article ID"
//	@Produce		json
//	@Security		BearerAuth
//	@Success		204	{string}	string	"No Content"
//	@Failure		400	{string}	string	"Bad Request"
//	@Failure		403	{string}	string	"Forbidden"
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Router			/v1/articles/{id} [delete]
func (a *ArticlesHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.repo.DeleteArticle(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UpdateArticle updated an article by ID
//
//	@Summary		Partial update article by ID
//	@Description	Partial update article by his ID
//	@Tags			articles
//	@Param			id	path	string	true	"Article ID"
//	@Accept			json
//	@Produce		json
//	@Success		204	{string}	string	"No Content"
//	@Failure		400	{string}	string	"Bad Request"
//	@Failure		403	{string}	string	"Forbidden"
//	@Failure		500	{string}	string	"Internal Server Error"
//	@Router			/v1/articles/{id} [patch]
func (a *ArticlesHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var data repo.UpdateArticleRequest
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = a.repo.UpdateArticle(r.Context(), id, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
