package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/visualect/pb/models"
)

type ArticlesRepository interface {
	GetArticles(ctx context.Context) ([]models.Article, error)
	GetArticle(ctx context.Context, id int) (models.Article, error)
	CreateArticle(ctx context.Context) (models.Article, error)
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
	articles, err := pgx.CollectRows(rows, pgx.RowTo[models.Article])
	if err != nil {
		return []models.Article{}, fmt.Errorf("error in %s: %s", fn, err)
	}
	return articles, nil
}

func (r *articlesRepo) GetArticle(ctx context.Context, id int) (models.Article, error) {
	fn := "articlesRepo.GetArticle"
	var article models.Article
	row := r.db.QueryRow(ctx, "SELECT * FROM articles WHERE id = $1", id)
	err := row.Scan(&article)
	if errors.Is(pgx.ErrNoRows, err) {
		return models.Article{}, nil
	}
	if err != nil {
		return models.Article{}, fmt.Errorf("error in %s: %s", fn, err)
	}
	return article, nil
}

func (r *articlesRepo) CreateArticle(ctx context.Context, id int) (models.Article, error) {
	fn := "articlesRepo.CreateArticle"
	var article models.Article
	_, err := r.db.Exec(ctx, "SELECT * FROM articles WHERE id = $1", id)

	// err := row.Scan(&article)
	// if errors.Is(pgx.ErrNoRows, err) {
	// 	return models.Article{}, nil
	// }
	// if err != nil {
	// 	return models.Article{}, fmt.Errorf("error in %s: %s", fn, err)
	// }
	// return article, nil
}
