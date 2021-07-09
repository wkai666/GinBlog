package app

import (
	"ginApp/pkg/e"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// BindAndValid binds and validate data
func BindAndValid(c *gin.Context, form interface{}) (int, int) {

	err := c.Bind(form)
	if err != nil {
		return http.StatusBadRequest, e.INVALID_PARAMS
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)

	log.Println("ssssssss--- ", form, err)

	if err != nil {
		return http.StatusInternalServerError, e.ERROR
	}

	if !check {
		MarkError(valid.Errors)
		return http.StatusBadRequest, e.INVALID_PARAMS
	}

	return http.StatusOK, e.SUCCESS
}
