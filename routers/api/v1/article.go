package v1

import (
	"ginApp/pkg/app"
	"ginApp/pkg/e"
	"ginApp/pkg/export"
	"ginApp/pkg/logging"
	"ginApp/pkg/qrcode"
	"ginApp/pkg/setting"
	"ginApp/pkg/util"
	"ginApp/service/article_service"
	"ginApp/service/tag_service"
	"github.com/astaxie/beego/validation"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)


// @Summary Get a single article
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failed 500 {object} app.Response
// @Router /api/v1/article/{id} [get]

func GetArticle(c *gin.Context) {
	appG := app.Gin{c}

	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("标签必须大于 0")

	if valid.HasErrors() {
		app.MarkError(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
	}

	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()

	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}

	if ! exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	article, err := articleService.Get()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_ARTICLE_FAIL, nil)
	}

	appG.Response(http.StatusOK, e.SUCCESS, article)

}

func GetArticles(c *gin.Context) {
	appG := app.Gin{C: c}

	data := make(map[string]interface{})
	valid := validation.Validation{}

	var state = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只能为 0 或 1")
	}

	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		valid.Min(tagId, 1, "tag_id").Message("标签 ID 必须大于 0")
	}

	articleService := article_service.Article{
		State: state,
		TagID: tagId,

		PageNum: util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}

	articles, err := articleService.GetAll()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}

	count, err := articleService.Count()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}
	data["list"] = articles
	data["count"] = count
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

type AddArticleForm struct {
	TagID int `form:"tag_id" valid:"Required;Min(1)"`
	Title string `form:"title" valid:"Required;MaxSize(255)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	Desc string `form:"desc" valid:"Required;MaxSize(255)"`
	Content string `form:"content" valid:"Required;MaxSize(65535)"`
	CreatedBy string `form:"created_by" valid:"Required;MaxSize(100)"`
	State int `form:"state" valid:"Range(0,1)"`
}

// @Summary Add article
// @Produce json
// @Param tag_id body string false "TagID"
// @Param title body string false "Title"
// @Param cover_image_url body string false "CoverImageUrl"
// @Param desc body string false "Desc"
// @Param content body string false "Content"
// @Param created_by body string false "createdBy"
// @Param state body string false "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/article/{id} [post]

func AddArticle(c *gin.Context) {

	var (
		appG = app.Gin{C: c}
		form AddArticleForm
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{ID: form.TagID}
	exist, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	articleService := article_service.Article{
		TagID: form.TagID,
		Title: form.Title,
		CoverImageUrl: form.CoverImageUrl,
		Desc: form.Desc,
		Content: form.Content,
		CreatedBy: form.CreatedBy,
		State: form.State,
	}

	if err := articleService.Add(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditArticleForm struct {
	ID int `form:"id" valid:"Required;Min(1)"`
	TagID int `form:"tag_id" valid:"Required;Min(1)"`
	Title string `form:"title" valid:"Required;MaxSize(255)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	Desc string `form:"desc" valid:"Required;MaxSize(255)"`
	Content string `form:"content" valid:"Required;MaxSize(65535)"`
	ModifiedBy string `form:"modified_by" valid:"Required;MaxSize(100)"`
	State int `form:"state" valid:"Required;Range(0,1)"`
}

// @Summary Update article
// @Produce json
// @Param id path int true "ID"
// @Param tag_id body string false "TagID"
// @Param title body string false "Title"
// @Param cover_image_url body string false "CoverImageUrl"
// @Param desc body string false "Desc"
// @Param content body string false "Content"
// @Param modified_by body string false "ModifiedBy"
// @Param state body string false "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/article/{id} [put]

func EditArticle(c *gin.Context) {

	var (
		appG = app.Gin{C:c}
		form = EditArticleForm{ID: com.StrTo(c.Param("id")).MustInt()}
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	articleService := article_service.Article{
		ID:	form.ID,
		TagID: form.TagID,
		Title: form.Title,
		CoverImageUrl: form.CoverImageUrl,
		Desc: form.Desc,
		Content: form.Content,
		ModifiedBy: form.ModifiedBy,
		State: form.State,
	}

	exist, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
	}

	tagService := tag_service.Tag{ID: form.TagID}
	exist, err = tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
	}

	if err = articleService.Edit(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, nil)
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary Delete article
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/article/{id} [delete]

func DeleteArticle(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID 必须大于 0")

	if valid.HasErrors() {
		app.MarkError(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{ID: id}
	exist, err := articleService.ExistByID()

	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}

	if !exist {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	err = articleService.Delete()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)

}

func ExportArticle(c *gin.Context)  {
	appG := app.Gin{C: c}

	id := com.StrTo(c.PostForm("id")).MustInt()
	tagId := com.StrTo(c.PostForm("tag_id")).MustInt()

	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(c.PostForm("state")).MustInt()
	}

	articleService := article_service.Article{
		ID: id,
		TagID: tagId,
		State: state,
	}
	fileName, err := articleService.Export()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXPORT_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"export_url": export.GetExcelFullPath() + fileName,
		"export_save_url": export.GetExcelFullUrl(fileName),
	})
}

func ImportArticle(c *gin.Context) {
	appG := app.Gin{C: c}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{}
	err = articleService.Import(file)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_IMPORT_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK,e.SUCCESS,nil)
}

// 二维码信息
const QRCODE_URL = "hshshshhsh"

func GenerateArticlePoster(c *gin.Context)  {
	appG := app.Gin{C: c}
	article := &article_service.Article{}
	qrc := qrcode.NewQrCode(QRCODE_URL, 300, 300, qr.M, qr.Auto)
	posterName := article_service.GetPosterFlag() + "-" + qrcode.GetQrCodeFileName(qrc.URL) + qrc.GetCodeExt()
	articlePoster := article_service.NewArticlePoster(posterName, article, qrc)
	articlePosterBgService := article_service.NewArticlePosterBg(
			"bg.jpg",
			articlePoster,
			&article_service.Rect{
				X0: 0,
				Y0: 0,
				X1: 555,
				Y1: 700,
			},
			&article_service.Pt{
				X: 125,
				Y: 298,
			},
		)

	_, filePath, err := articlePosterBgService.Generate()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GEN_ARTICLE_POSTER_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"post_url": qrcode.GetQrCodeFullUrl(posterName),
		"post_save_url": filePath + posterName,
	})
}