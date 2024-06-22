package repository

import (
	"gorm.io/gorm"
	"main/ex01/internal/domain/models"
)

type ArticleRepository interface {
	CreateArticle(article *models.Article) error
	GetArticles(offset, limit int) ([]models.Article, error)
	CountArticles() (int, error)
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db}
}

func (r *articleRepository) CreateArticle(article *models.Article) error {
	return r.db.Create(article).Error
}

func (r *articleRepository) GetArticles(offset, limit int) ([]models.Article, error) {
	var articles []models.Article
	err := r.db.Order("id desc").Limit(limit).Offset(offset).Find(&articles).Error
	return articles, err
}

func (r *articleRepository) CountArticles() (int, error) {
	var count int64
	err := r.db.Model(&models.Article{}).Count(&count).Error
	return int(count), err
}
