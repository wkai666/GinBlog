package models

import "github.com/jinzhu/gorm"

type Article struct {
	Model

	TagID int `json:"tag_id" gorm:"index"`
	Tag   Tag `json:"tag"`

	Title string `json:"title"`
	CoverImageUrl string `json:"cover_image_url"`
	Desc string `json:"desc"`
	Content string `json:"content"`
	CreatedBy string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State int `json:"state"`
}

func ExistArticleById(id int) (bool, error) {
	var article Article

	 err := db.Select("id").Where("id = ? AND deleted_on = ?", id, 0).First(&article).Error

	 if err != nil && err != gorm.ErrRecordNotFound {
	 	return false, err
	 }

	if article.ID > 0 {
		return true, nil
	}

	return false, nil
}

func GetArticleTotal(maps interface{}) (count int) {
	db.Model(&Article{}).Where(maps).Count(&count)

	return
}

func GetArticle(id int) (*Article, error) {

	var article Article

	err := db.Where("id = ? and deleted_on = ?", id, 0).First(&article).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	err = db.Model(&article).Related(&article.Tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &article, nil
}

func GetArticles(pageNum int, pageSize int, maps interface{}) (articles []Article) {
	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)

	return
}

func EditArticle(id int, data interface{}) error {
	err := db.Model(&Article{}).Where("id = ? and deleted_on = ?", id, 0).Updates(data).Error
	if err != nil {
		return err
	}

	return nil
}

// AddArticle add a single article
func AddArticle(data map[string]interface{}) error {
	article := Article{
		TagID: data["tag_id"].(int),
		Title: data["title"].(string),
		CoverImageUrl: data["cover_image_url"].(string),
		Desc: data["desc"].(string),
		Content: data["content"].(string),
		CreatedBy: data["created_by"].(string),
		State: data["state"].(int),
	}

	if err := db.Create(&article).Error; err != nil {
		return err
	}

	return nil
}

// DeleteArticle delete a single article
func DeleteArticle(id int) error {

	err := db.Where("id = ?", id).Delete(Article{}).Error
	if err != nil {
		return err
	}

	return nil
}

func CleanAllArticle() bool {
	db.Unscoped().Where("deleted_on != ?", 0).Delete(&Article{})

	return true
}