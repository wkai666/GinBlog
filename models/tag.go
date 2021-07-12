package models

import "github.com/jinzhu/gorm"

type Tag struct {
	Model

	Name string `json:"name"`
	CreatedBy string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State int `json:"state"`
}

func GetTags(pageNum int, pageSize int, maps interface {}) (tags []Tag, err error) {

	if pageSize > 0 && pageNum > 0 {
		err = db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags).Error
	} else {
		err = db.Where(maps).Find(&tags).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return tags, nil
}

// GetTagTotal count the total number of tags
func GetTagTotal(maps interface {}) (int, error) {
	var count int

	err := db.Model(&Tag{}).Where(maps).Count(&count).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}

	return count, nil
}

// ExistTagByName checks if tag is already existed
func ExistTagByName(name string) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("name = ? AND deleted_on = ?", name, 0).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if tag.ID > 0 {
		return true, nil
	}
	return false, nil
}

// AddTag add a tag
func AddTag(name string, state int, createdBy string) error {
	err := db.Create(&Tag{
		Name: name,
		State: state,
		CreatedBy: createdBy,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

// ExistTagByID whether a Tag exists based on the ID
func ExistTagByID(id int) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("id = ? AND deleted_on = ?", id, 0).First(&tag).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if tag.ID > 0 {
		return true, nil
	}

	return false, nil
}

// DeleteTag delete a tag
func DeleteTag(id int) error {
	err := db.Where("id = ?", id).Delete(&Tag{}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	return nil
}

// EditTag modify a single tag
func EditTag(id int, data interface{}) error {
	err := db.Model(&Tag{}).Where("id = ? AND deleted_on = ?", id, 0).Update(data).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	return nil
}

func CleanAllTag() bool {
	db.Unscoped().Where("deleted_on != ? ", 0).Delete(&Tag{})

	return true
}
