package article_service

import (
	"encoding/json"
	"ginApp/models"
	"ginApp/pkg/export"
	"ginApp/pkg/gredis"
	"ginApp/pkg/logging"
	"ginApp/service/cache_service"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/tealeg/xlsx"
	"io"
	"strconv"
	"time"
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

func (a *Article) Count() (int, error) {
	count, err := models.GetArticleTotal(a.GetMaps())
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (a *Article) GetAll() ([]models.Article, error) {
	articles, err := models.GetArticles(a.PageNum, a.PageSize, a.GetMaps())
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (a *Article) GetMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0

	if a.ID > 0 {
		maps["id"] = a.ID
	}

	if a.TagID > 0 {
		maps["tag_id"] = a.TagID
	}

	if a.State >= 0{
		maps["state"] = a.State
	}

	return maps
}

func (a *Article) ExistByID() (bool, error) {
	return models.ExistArticleById(a.ID)
}

func (a *Article) Delete() (error) {
	return models.DeleteArticle(a.ID)
}

func (a *Article) Export() (string, error) {
	articles, err := a.GetAll()
	if err != nil {
		return "", err
	}

	file := xlsx.NewFile()
	sheet, err := file.AddSheet("文章")
	if err != nil {
		return "", err
	}

	titles := []string{"文章ID", "标签", "文章标题", "封面", "简述", "内容", "创建者", "创建时间", "修改者", "修改时间"}
	row := sheet.AddRow()

	var cell *xlsx.Cell
	for _, v := range titles {
		cell = row.AddCell()
		cell.Value = v
	}

	for _, article := range articles {
		value := []string{
			strconv.Itoa(article.ID),
			article.Tag.Name,
			article.Title,
			article.CoverImageUrl,
			article.Desc,
			article.Content,
			article.CreatedBy,
			article.ModifiedBy,
		}

		row = sheet.AddRow()
		for _, v := range value {
			cell = row.AddCell()
			cell.Value = v
		}
	}

	timeStr := strconv.Itoa(int(time.Now().Unix()))
	fileName := "articles_" + timeStr + ".xlsx"

	fullPath := export.GetExcelFullPath() + fileName

	err = file.Save(fullPath)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (a *Article) Import(r io.Reader) error {
	xls, err := excelize.OpenReader(r)
	if err != nil {
		return err
	}

	rows := xls.GetRows("文章信息")
	for i, row := range rows {
		if i > 0 {
			data := make(map[string]interface{})

			data["tag_id"],_ = strconv.Atoi(row[0])
			data["title"] = row[1]
			data["cover_image_url"] = row[2]
			data["desc"] = row[3]
			data["content"] = row[4]
			data["created_by"] = row[5]
			data["state"],_ = strconv.Atoi(row[6])

			models.AddArticle(data)
		}
	}

	return nil
}