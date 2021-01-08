package main

import (
	"github.com/gin-gonic/gin"
	//"fmt"
	"fmt"
)

func main() {
	r := gin.Default()
	//r.GET("/",func(c *gin.Context) {
	// http://localhost:8080/
	//c.JSON(200,gin.H{
	//	"Blog":"www.flysnow.org",
	//	"wechat":"flysnow_org",
	//})

	// http://localhost:8080/?wechat=flysnow_org
	//c.String(200,c.Query("wechat"))  // string type
	//c.String(200,c.DefaultQuery("id","0"))    // int type

	// Array
	// http://localhost:8080/?media=blog&media=wechat
	//c.JSON(200,c.QueryArray("media"))
	//})

	// Map
	//r.GET("/map",func(c *gin.Context) {
	//	// http://localhost:8080/map?ids[a]=123&ids[b]=456&ids[c]=789
	//	// a=!b=!c
	//	c.JSON(200,c.QueryMap("ids"))
	//})

	// 接收表单数据
	// terminal: curl -d wechat=flysnow_org http://localhost:8080/
	//r.POST("/",func(c *gin.Context) {
	//	c.String(200,c.PostForm("wechat"))
	//})

	// 分组路由 Group(relativePath string, handlers ...HandlerFunc)
	// 路由中间件 ...HandlerFunc
	//v1Group := r.Group("/v1")
	//v1Group := r.Group("/v1",func(c *gin.Context) {
	//	fmt.Println("/v1中间件")
	//})
	//{
	//	v1Group.GET("/users", func(c *gin.Context) {
	//		c.String(200, "/v1/users")
	//	})
	//
	//	v1Group.GET("/products", func(c *gin.Context) {
	//		c.String(200, "/v1/products")
	//	})
	//}
	//v2Group := r.Group("/v2")
	//v2Group := r.Group("/v2",func(c *gin.Context) {
	//	fmt.Println("/v1中间件")
	//})
	//{
	//	v2Group.GET("/users", func(c *gin.Context) {
	//		c.String(200, "/v2/users")
	//	})
	//
	//	v2Group.GET("/products", func(c *gin.Context) {
	//		c.String(200, "/v2/products")
	//	})
	//}

	// http://localhost:8080/?address=192.168.186.137:3306&address=192.168.186.137:3307&address=192.168.186.138:3306&user=root&password=Abc727364
	r.GET("/",func(c *gin.Context){
		c.JSON(200,c.QueryArray("address"))
		fmt.Println(c.QueryArray("address"))
		//c.String(200,c.Query("user"))
		//c.String(200,c.Query("password"))
		fmt.Println(c.Query("user"))
		fmt.Println(c.Query("password"))

	})

	r.Run(":8080")
}