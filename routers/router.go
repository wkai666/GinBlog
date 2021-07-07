package routers

import (
	"ginApp/middleware/jwt"
	"ginApp/pkg/setting"
	"ginApp/routers/api"
	v1 "ginApp/routers/api/v1"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.RunMode)

	r.GET("/auth", api.GetAuth)

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

	}

	//r.GET("/test", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"message":"test",
	//	})
	//})

	return  r
}
