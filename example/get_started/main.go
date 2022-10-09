package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	"tinyweb/tinyweb"
)

type student struct {
	Name string
	Age  int
}

func main() {
	r := tinyweb.New()

	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": func(t time.Time) string {
			year, month, day := t.Date()
			return fmt.Sprintf("%d-%02d-%02d", year, month, day)
		},
	})

	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static", []string{"/css/style.css"})

	// 计时器中间件
	r.Use(func(c *tinyweb.Context) {
		t := time.Now()
		log.Println("global timer start...")
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Request.RequestURI, time.Since(t))
		log.Println("global timer finish...")
	})

	r.GET("/", func(ctx *tinyweb.Context) {
		ctx.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(c *tinyweb.Context) {
		stu1 := &student{Name: "Peterson", Age: 20}
		stu2 := &student{Name: "Marry", Age: 22}
		c.HTML(http.StatusOK, "arr.tmpl", tinyweb.H{
			"title":  "tiny web",
			"stuArr": [2]*student{stu1, stu2},
		})
	})

	r.GET("/date", func(c *tinyweb.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", tinyweb.H{
			"title": "tiny web",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	r.GET("/hello", func(c *tinyweb.Context) {
		// expect /hello?name=someone
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *tinyweb.Context) {
		c.JSON(http.StatusOK, tinyweb.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	r.GET("/u/:username/:password", func(c *tinyweb.Context) {
		c.JSON(http.StatusOK, tinyweb.H{
			"username": c.Params["username"],
			"password": c.Params["password"],
		})
	})

	v1 := r.Group("/v1")
	v1.Use(func(c *tinyweb.Context) { // 应用在group上的中间件
		t := time.Now()
		log.Println("v1 timer start...")
		c.Next()
		log.Printf("[%d] %s in %v for group v1", c.StatusCode, c.Request.RequestURI, time.Since(t))
		log.Println("v1 timer finish...")
	})
	{
		v1.GET("/:username", func(c *tinyweb.Context) {
			c.HTML(
				http.StatusOK,
				"index.html",
				fmt.Sprintf("<h1>welcome, %s</h1>\n", c.Params["username"]),
			)

		})
	}

	r.GET("/panic", func(c *tinyweb.Context) {
		names := []string{"a string"}
		c.String(http.StatusOK, names[100])
	})

	err := r.Run("localhost:23456")
	if err != nil {
		panic(err)
	}
}
