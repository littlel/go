/*
 * @Description:
 * @Author:
 * @Date: 2021-12-22 17:03:04
 */
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	web_server()
}

func web_server() {
	router := gin.Default()

	router.Use(Cors())

	router.Static("/assets", "./assets")
	router.StaticFS("/more_static", http.Dir("my_file_system"))
	router.StaticFS("/video", http.Dir("./video"))
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.POST("/upload", uploadFile)
	router.GET("/config", get_config)

	//重定向测试
	router.GET("/redirect", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "http://www.baidu.com/")
	})

	// 监听并服务于 0.0.0.0:8080
	router.Run(":8080")

}

func get_config(c *gin.Context) {
	filepath := "/home/projects/Release/config/IvmsConfig.xml"
	xml_context := read3(filepath)
	c.String(200, xml_context)
}

// ioutil
func read3(filepath string) string {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	f, _ := ioutil.ReadAll(file)
	txt := string(f)
	return txt
}

func uploadFile(c *gin.Context) {
	// FormFile方法会读取参数“upload”后面的文件名，返回值是一个File指针，和一个FileHeader指针，和一个err错误。
	file, header, err := c.Request.FormFile("upload")
	if err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}
	// header调用Filename方法，就可以得到文件名
	filename := header.Filename
	fmt.Println(file, err, filename)
	// 创建一个文件，文件名为filename，这里的返回值out也是一个File指针
	out, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer out.Close()
	// 将file的内容拷贝到out
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	//执行升级，升级完才返回成功

	c.String(http.StatusCreated, "upload successful \n")
}

//支持浏览器简单预检
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("origin") //请求头部
		if len(origin) == 0 {
			origin = c.Request.Header.Get("Origin")
		}
		//接收客户端发送的origin （重要！）
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		//允许客户端传递校验信息比如 cookie (重要)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		//服务器支持的所有跨域请求的方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET")
		//c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		// 设置预验请求有效期为 86400 秒
		//c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		//if c.Request.Method == "OPTIONS" {
		//	c.AbortWithStatus(204)
		//	return
		//}
		c.Next()
	}
}
