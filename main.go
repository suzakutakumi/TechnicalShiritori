package main

import (
	"TechnicalShiritori/room"
	"errors"
	"html/template"
	"io"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
)

const NumOfDigits = 4

var rooms map[int]*room.Room

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("docs/*.html")),
	}
	rand.Seed(time.Now().UnixNano())
	e := echo.New()
	e.Renderer = t
	e.GET("/", index)
	e.GET("/rooms", roomList)
	e.POST("/room", create)
	e.GET("/room", enter)
	e.Logger.Fatal(e.Start(":8080"))
}
func index(c echo.Context) error {
	_, err := getUser(c)
	if err != nil {
		_, err = setUser(c)
		if err != nil {
			return err
		}
	}
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{})
}
func roomList(c echo.Context) error {
	rs := []room.Room{}
	for _, v := range rooms {
		rs = append(rs, *v)
	}
	return c.JSON(http.StatusOK, rs)
}
func create(c echo.Context) error {
	var id int
	for {
		id = rand.Intn(int(math.Pow10(NumOfDigits)))
		if _, ok := rooms[id]; !ok {
			break
		}
	}

	user, err := getUser(c)
	if err != nil {
		return err
	}

	rooms[id] = &room.Room{
		id,
		[]string{user},
		[]room.Word{},
	}

	return c.JSON(http.StatusOK, rooms[id])
}
func enter(c echo.Context) error {
	number, err := strconv.Atoi(c.Param("n"))
	if err != nil {
		return err
	}
	_, ok := rooms[number]
	if !ok {
		return errors.New("部屋が存在しません")
	}
	user, err := getUser(c)
	if err != nil {
		return err
	}
	rooms[number].Users = append(rooms[number].Users, user)
	return c.JSON(http.StatusOK, rooms[number])
}

func setUser(c echo.Context) (string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	user := new(http.Cookie)
	user.Name = "user"
	user.Value = u.String()
	user.Expires = time.Now().Add(2 * time.Hour)
	c.SetCookie(user)
	return user.Value, nil
}

func getUser(c echo.Context) (string, error) {
	user, err := c.Cookie("user")
	if err != nil {
		return "", err
	}
	return user.Value, nil
}
