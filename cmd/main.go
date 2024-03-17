package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
  templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data any, c echo.Context) error {
  return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() *Template {
  return &Template{
    templates: template.Must(template.ParseGlob("views/*.html")),
  }
}

type Count struct {
  Count int
}

func main() {
	e := echo.New()
  e.Renderer = NewTemplate()

  // middlewares
  e.Use(middleware.Logger())

  // Routes
  var count Count

  e.GET("/", func(c echo.Context) error {
    return c.Render(http.StatusOK, "index", count)
  })

  e.POST("/count", func(c echo.Context) error {
    count.Count++
    return c.Render(http.StatusOK, "count", count)
  })

  e.Logger.Fatal(e.Start(":1323"))
}
