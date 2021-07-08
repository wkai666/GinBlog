package models

import (
	"fmt"
	"ginApp/pkg/logging"
	"ginApp/pkg/setting"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var db *gorm.DB

type Model struct {
	ID int `gorm:"primary_key" json:"id"`
	CreatedOn int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn int `json:"deleted_on"`
}

func Setup() {
	var (
		err error
		dbType, dbName, user, passwd, host, tablePrefix string
	)

	dbType = setting.DatabaseSetting.Type
	host = setting.DatabaseSetting.Host
	user = setting.DatabaseSetting.User
	passwd = setting.DatabaseSetting.Password
	dbName = setting.DatabaseSetting.DBName
	tablePrefix = setting.DatabaseSetting.TablePrefix

	db, err = gorm.Open(
		dbType,
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			user, passwd, host, dbName))
	if err != nil {
		logging.Error(err)
	}

	gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}

	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
}

func CloseDB() {
	defer db.Close()
}

func ExistTagByID(id int) bool {
	var tag Tag
	db.Select("id").Where("id = ?", id).First(&tag)
	if tag.ID > 0 {
		return true
	}

	return false
}

func DeleteTag(id int) bool {
	db.Where("id = ?", id).Delete(&Tag{})

	return true
}

func EditTag(id int, data interface{}) bool {
	db.Model(&Tag{}).Where("id = ?", id).Update(data)

	return true
}

// updateTimeStampForCreateCallback will set `createdOn` `modifiedOn` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if ! scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

// updateTimeStampForUpdateCallback will set `ModifyTime` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

// deleteCallback will
func deleteCallback(scope *gorm.Scope) {
	if ! scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn")
		if !scope.Search.Unscoped && hasDeletedOnField {
			scope.Raw(fmt.Sprintf(
					"UPDATE %v SET %v = %v%v%v",
					scope.QuotedTableName(),
					scope.Quote(deletedOnField.DBName),
					scope.AddToVars(time.Now().Unix()),
					addExtraSpaceIfExist(scope.CombinedConditionSql()),
					addExtraSpaceIfExist(extraOption),
				)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
					"DELETE FROM %v%v%v",
					scope.QuotedTableName(),
					addExtraSpaceIfExist(scope.CombinedConditionSql()),
					addExtraSpaceIfExist(extraOption),
				)).Exec()
		}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}