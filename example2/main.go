package main

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	"pagination"
)

func main() {

	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")
	r.GET("/", func(c *gin.Context) {
		//创建一个分页器，100条数据，每页10条
		pagination := pagination.Initialize(c.Request, 100, 10)
		//传到模板中需要转换成template.HTML类型，否则html代码会被转义
		c.HTML(http.StatusOK, "index.html", gin.H{
			"paginate": template.HTML(pagination.Pages()),
		})
	})
	r.Run(":8080")
}
