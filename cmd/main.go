package main

import (
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"

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

var id int

type Contact struct {
	Name  string
	Email string
	Id    int
}

func newContact(name, email string) Contact {
	id++

	return Contact{
		Name:  name,
		Email: email,
		Id:    id,
	}
}

type Contacts = []Contact

type Data struct {
	Contacts Contacts
}

func (d *Data) hasEmail(email string) bool {
	for i := 0; i < len(d.Contacts); i++ {
		if c := d.Contacts[i]; c.Email == email {
			return true
		}
	}

	return false
}

func (d *Data) deleteContact(id int) {
	for i := 0; i < len(d.Contacts); i++ {
		if c := d.Contacts[i]; c.Id == id {
			d.Contacts = append(d.Contacts[:i], d.Contacts[i+1:]...)
			return
		}
	}
}

func newData() Data {
	return Data{
		Contacts: Contacts{
			newContact("Jonh", "jd@email.com"),
			newContact("Clara", "cd@email.com"),
		},
	}
}

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

func newFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

type Page struct {
	Form FormData
	Data Data
}

func newPage() Page {
	return Page{
		Data: newData(),
		Form: newFormData(),
	}
}

func main() {
	e := echo.New()

	e.Renderer = NewTemplate()
	e.Static("/images", "images")
	e.Static("/css", "css")

	// Middlewares
	e.Use(middleware.Logger())

	// Data
	page := newPage()

	// Routes

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", page)
	})

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		if page.Data.hasEmail(email) {
			formData := newFormData()
			formData.Values["name"] = name
			formData.Values["email"] = email
			formData.Errors["email"] = "Email already exists"

			return c.Render(http.StatusUnprocessableEntity, "form", formData)
		}

		contact := newContact(name, email)
		page.Data.Contacts = append(page.Data.Contacts, contact)

		c.Render(http.StatusOK, "form", newFormData())
		return c.Render(http.StatusOK, "oob-contact", contact)
	})

	e.DELETE("/contacts/:id", func(c echo.Context) error {
		time.Sleep(2 * time.Second)

		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid id")
		}

		page.Data.deleteContact(id)

		return c.NoContent(http.StatusOK)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
