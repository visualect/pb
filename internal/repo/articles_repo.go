package repo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/visualect/pb/internal/models"
)

// CreateArticleRequest represents fields for creating new article
//
//	@Description	Creat article structure
type CreateArticleRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// UpdateArticleRequest represents updates to an article
//
//	@Description	Partial article update structure
type UpdateArticleRequest struct {
	Title *string `json:"title"`
	Body  *string `json:"body"`
}

type ArticlesRepository interface {
	GetArticles(ctx context.Context) ([]models.Article, error)
	GetArticle(ctx context.Context, id int) (models.Article, error)
	CreateArticle(ctx context.Context, data CreateArticleRequest) (models.Article, error)
	DeleteArticle(ctx context.Context, id int) error
	UpdateArticle(ctx context.Context, id int, data UpdateArticleRequest) error
}

type articlesRepo struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) ArticlesRepository {
	return &articlesRepo{db}
}

func (r *articlesRepo) GetArticles(ctx context.Context) ([]models.Article, error) {
	fn := "articlesRepo.GetArticles"
	rows, err := r.db.Query(ctx, "SELECT * FROM articles")
	if err != nil {
		return []models.Article{}, fmt.Errorf("error in %s: %s", fn, err)
	}
	articles, err := pgx.CollectRows(rows, pgx.RowToStructByPos[models.Article])
	if err != nil {
		return []models.Article{}, fmt.Errorf("error in %s: %s", fn, err)
	}
	return articles, nil
}

func (r *articlesRepo) GetArticle(ctx context.Context, id int) (models.Article, error) {
	fn := "articlesRepo.GetArticle"
	var a models.Article
	row := r.db.QueryRow(ctx, "SELECT * FROM articles WHERE id = $1", id)
	err := row.Scan(&a.ID, &a.Title, &a.Body, &a.CreatedAt)
	if errors.Is(pgx.ErrNoRows, err) {
		return models.Article{}, nil
	}
	if err != nil {
		return models.Article{}, fmt.Errorf("error in %s: %s", fn, err)
	}
	return a, nil
}

func (r *articlesRepo) CreateArticle(ctx context.Context, data CreateArticleRequest) (models.Article, error) {
	fn := "articlesRepo.CreateArticle"
	var a models.Article
	row := r.db.QueryRow(ctx, "INSERT INTO articles (title, body) VALUES ($1, $2) RETURNING id, title, body, created_at", data.Title, data.Body)
	err := row.Scan(&a.ID, &a.Title, &a.Body, &a.CreatedAt)
	if err != nil {
		return models.Article{}, fmt.Errorf("error in %s: %s", fn, err)
	}
	return a, nil
}

func (r *articlesRepo) DeleteArticle(ctx context.Context, id int) error {
	fn := "articlesRepo.DeleteArticle"
	_, err := r.db.Exec(ctx, "DELETE FROM articles WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error in %s: %s", fn, err)
	}
	return nil
}

func (r *articlesRepo) UpdateArticle(ctx context.Context, id int, data UpdateArticleRequest) error {
	fn := "articlesRepo.UpdateArticle"
	var clauses []string
	var args []any
	argIdx := 1

	if data.Title != nil {
		clauses = append(clauses, fmt.Sprintf("title = $%d", argIdx))
		args = append(args, *data.Title)
		argIdx++
	}

	if data.Body != nil {
		clauses = append(clauses, fmt.Sprintf("body = $%d", argIdx))
		args = append(args, *data.Body)
		argIdx++
	}

	if len(clauses) == 0 {
		return fmt.Errorf("error in %s: at least one field required", fn)
	}
	args = append(args, id)

	query := fmt.Sprintf("UPDATE articles SET %s WHERE id = $%d", strings.Join(clauses, ", "), argIdx)
	_, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error in %s: %s", fn, err)
	}
	return nil
}
