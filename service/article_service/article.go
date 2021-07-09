package article_service

import (
	"encoding/json"
	"ginApp/models"
	"ginApp/pkg/gredis"
	"ginApp/pkg/logging"
	"ginApp/service/cache_service"
)

type Article struct {
	ID int
	TagID int
	Title string
	CoverImageUrl string
	Desc string
	Content string
	State int
	CreatedBy string
	ModifiedBy string

	PageNum int
	PageSize int
}

func (a *Article) Add() error {
	article := map[string]interface{}{
		"tag_id": a.TagID,
		"title": a.Title,
		"cover_image_url": a.CoverImageUrl,
		"desc": a.Desc,
		"content": a.Content,
		"state": a.State,
		"created_by":a.CreatedBy,
	}

	if err := models.AddArticle(article); err != nil {
		return err
	}

	return nil
}

func (a *Article) Edit() error {
	return models.EditArticle(a.ID, map[string]interface{}{
		"tag_id": a.TagID,
		"title": a.Title,
		"cover_image_url": a.CoverImageUrl,
		"desc": a.Desc,
		"content": a.Content,
		"state": a.State,
		"modified_by": a.ModifiedBy,
	})
}

func (a *Article) Get() (*models.Article, error) {
	var cacheArticle *models.Article

	cache := cache_service.Article{ID: a.ID}
	key := cache.GetArticleKey()
	if gredis.Exists(key) {
		 if data, err := gredis.Get(key); err != nil {
		 	logging.Info(err)
		 } else {
			json.Unmarshal(data, &cacheArticle)
			return cacheArticle, nil
		 }
	}

	article, err := models.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}
	gredis.Set(key, article, 3600)
	return article, nil
}

func (a *Article) ExistByID() (bool, error) {
	return models.ExistArticleById(a.ID)
}

func (a *Article) Delete() (error) {
	return models.DeleteArticle(a.ID)
}