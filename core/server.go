//极简博客

package core

import (
	"embed"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//go:embed dist/*
var s embed.FS

func StartServer(path string, port string, auth string) {

	//记录日志
	f, _ := os.Create(path + "blog.log")
	//拒绝gin记录日志
	//f, _ := os.Create("")
	//调试信息
	gin.DefaultWriter = io.MultiWriter(f)
	//错误信息
	gin.DefaultErrorWriter = io.MultiWriter(f)

	//gin 初始定义
	router := gin.Default()
	//添加压缩协议
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	//静态资源添加
	staticFs, _ := fs.Sub(s, "dist/static")
	router.StaticFS("/static", http.FS(staticFs))

	staticFsScript, _ := fs.Sub(s, "dist/static/scripts")
	router.StaticFS("/scripts", http.FS(staticFsScript))

	//判断是否需要用户名密码认证
	var username, password string
	if auth != "" {
		accountInfo := strings.Split(auth, ":")
		if len(accountInfo) == 2 {
			username = accountInfo[0]
			password = accountInfo[1]
		}
	}

	//根据需要添加http认证
	if username != "" && password != "" {
		accountList := map[string]string{
			username: password,
		}
		authorized := router.Group("/", gin.BasicAuth(accountList))
		authorized.GET("", func(c *gin.Context) {
			indexHTML, _ := s.ReadFile("dist/index.html")
			_, _ = c.Writer.WriteString(string(indexHTML))
		})
	} else {
		router.GET("/", func(c *gin.Context) {
			indexHTML, _ := s.ReadFile("dist/index.html")
			_, _ = c.Writer.WriteString(string(indexHTML))
		})
	}

	//重定向添加
	router.GET("/g/:name", func(c *gin.Context) {
		c.Request.URL.Path = "/"
		router.HandleContext(c)
	})

	//重定向添加
	router.GET("/a/:group/:name", func(c *gin.Context) {
		c.Request.URL.Path = "/"
		router.HandleContext(c)
	})

	//获取菜单项
	router.GET("/api/menuData", func(c *gin.Context) {

		status := "ok"

		//定义菜单数组
		var menuArray []MenuItem

		//添加默认菜单项
		var homepage MenuItem
		homepage.Path = "/g/all"
		homepage.Name = "首页"
		homepage.Component = "./Index"
		menuArray = append(menuArray, homepage)

		//遍历public内所有内容
		err := filepath.Walk(path+"public/", func(path string, f os.FileInfo, err error) error {

			//当遍历的是文件夹同时不是public文件夹时执行操作
			if f.IsDir() && f.Name() != "public" {
				//默认该菜单项是文章列表
				var group MenuItem
				group.Path = "/g/" + f.Name()
				group.Name = f.Name()
				group.Component = "./Index"

				//判断是否包含 index.md
				thisFile := path + "/index.md"
				_, err1 := os.Stat(thisFile)
				if !os.IsNotExist(err1) {
					//当该菜单包含 index.md 认为该菜单需要直接显示文章，无需展示文章列表
					group.Path = "/a/" + f.Name() + "/index"
					group.Name = f.Name()
					group.Component = "./Article"
				}
				//status = thisFile
				menuArray = append(menuArray, group)
			}

			return nil
		})
		if err != nil {
			status = "error"
		}

		data := map[string]interface{}{
			"status": status,
			"data":   menuArray,
		}
		c.JSONP(http.StatusOK, data)
	})

	//获取文章列表
	router.GET("/api/group/:name", func(c *gin.Context) {

		status := "ok"

		name := c.Param("name")

		var articleArray []Article

		if name == "all" {
			name = ""
		}

		err := filepath.Walk(path+"public/"+name, func(path string, f os.FileInfo, err error) error {

			if !f.IsDir() && f.Name() != "index.md" {
				var thisOne Article
				thisOne.Title = strings.Replace(f.Name(), ".md", "", 1)
				thisOne.SubTitle = f.ModTime().Format("2006-01-02")
				dir, _ := filepath.Split(path)
				thisOne.Group = strings.Replace(dir, "public"+string(os.PathSeparator), "", 1)
				thisOne.Group = strings.ReplaceAll(thisOne.Group, string(os.PathSeparator), "")
				if thisOne.Group == "" {
					thisOne.Group = "all"
				}
				thisOne.Avatar = "https://gw.alipayobjects.com/zos/antfincdn/UCSiy1j6jx/xingzhuang.svg"
				articleArray = append(articleArray, thisOne)
			}

			return nil
		})
		if err != nil {
			status = "error"
		}

		data := map[string]interface{}{
			"status": status,
			"data":   articleArray,
		}
		c.JSONP(http.StatusOK, data)
	})

	//获取文章
	router.GET("/api/article/:group/:name", func(c *gin.Context) {

		status := "ok"
		dataContent := "## 没有找到对应的文章"

		group := c.Param("group")
		name := c.Param("name")

		if group == "all" {
			group = ""
		} else {
			group = group + "/"
		}

		thisFile := path + "public/" + group + name + ".md"
		_, err1 := os.Stat(thisFile)
		//检查文章是否存在
		if !os.IsNotExist(err1) {
			var content []byte
			content, _ = os.ReadFile(thisFile)
			dataContent = string(content)
		}

		data := map[string]interface{}{
			"status": status,
			"data":   dataContent,
		}
		c.JSONP(http.StatusOK, data)
	})

	router.Run(":" + port)
}
