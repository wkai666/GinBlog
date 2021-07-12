package routers

import (
	_ "ginApp/docs"
	"ginApp/middleware/jwt"
	"ginApp/pkg/export"
	"ginApp/pkg/setting"
	"ginApp/pkg/upload"
	"ginApp/routers/api"
	v1 "ginApp/routers/api/v1"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func InitRouter() *gin.Engine {

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.AppSetting.RunMode)

	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))

	r.GET("/" ,api.Hello)
	r.GET("/auth", api.GetAuth)
	r.POST("/upload", api.UploadImage)

	url := ginSwagger.URL("http://127.0.0.1:9099/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT())
	{
		// 获取标签
		apiv1.GET("/tags", v1.GetTags)
		// 新建标签
		apiv1.POST("/tags", v1.AddTag)
		// 更新标签
		apiv1.PUT("/tag/:id", v1.EditTag)
		// 删除标签
		apiv1.DELETE("/tag/:id", v1.DelTag)
		// 导出标签
		r.POST("/tags/export", v1.ExportTags)
		// 导入标签
		r.POST("/tags/import", v1.ImportTag)

		// 获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		// 获取单篇文章
		apiv1.GET("/article/:id", v1.GetArticle)
		// 新增文章
		apiv1.POST("/article/:id", v1.AddArticle)
		// 编辑更新指定文章
		apiv1.PUT("article/:id", v1.EditArticle)
		// 删除文章
		apiv1.DELETE("/article/:id", v1.DeleteArticle)
		// 导出文章
		r.POST("/articles/export", v1.ExportArticle)
		// 导入文章
		r.POST("/articles/import", v1.ImportArticle)

	}

	return  r
}
