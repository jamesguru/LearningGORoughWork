package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Cat struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func Welcome(c echo.Context) error {
	return c.String(http.StatusOK, "Yall from the web side!")
}

func getCats(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")
	return c.String(http.StatusOK, fmt.Sprintf("Your cat name is: %s\n and his type is %s\n", catName, catType))
}

func addCart(c echo.Context) error {
	cat := Cat{}
	b, err := ioutil.ReadAll(c.Request().Body)

	if err != nil {
		log.Printf("Failed reading the request body: %s, err", err)
		return c.String(http.StatusInternalServerError, "")
	}

	err = json.Unmarshal(b, &cat)

	if err != nil {
		log.Printf("Failed unmarshing in addCarts: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	log.Printf("This is your cat: %#v", cat)
	return c.String(http.StatusOK, "we got your cat!")
}

func mainAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "You are main admin page")
}

// BasicAuth middleware function for authentication
func basicAuth(username, password string, c echo.Context) (bool, error) {
	// Simple hardcoded authentication (replace with database check in production)
	if username == "admin" && password == "password" {
		c.JSON(http.StatusOK, "Authenticated")
		return true, nil
	}
	return false, nil
}

func mainCookie(c echo.Context) error {
	return c.String(http.StatusOK, "You are on the secret cookie page!")
}

func login(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	// Check username and password against DB after hashing the password

	if username == "admin" && password == "1234" {
		cookie := new(http.Cookie)

		cookie.Name = "sessionID"
		cookie.Value = "some_string"
		cookie.Expires = time.Now().Add(48 * time.Hour)
		c.SetCookie(cookie)

		return c.String(http.StatusOK, "You were logged in!")
	} else {
		return c.String(http.StatusUnauthorized, "Your username or password were wrong!")
	}

}

func checkCookie(next echo.HandlerFunc) echo.HandlerFunc {

	return func(c echo.Context) error {
		cookie, err := c.Cookie("sessionID")

		if err != nil {
			if strings.Contains(err.Error(), "named cookie not present") {
				return c.String(http.StatusUnauthorized, "You don't have any cookie")
			}
			log.Println(err)
			return err
		}
		if cookie.Value == "some_string" {
			return next(c)
		}
		return c.String(http.StatusUnauthorized, "You dont have the right cookie, cookie")
	}
}

func main() {
	fmt.Println("Welcome to the server")
	e := echo.New()
	adminGroup := e.Group("/admin")
	cookieGroup := e.Group("/cookie")
	// this logs the server interaction
	adminGroup.Use(middleware.BasicAuth(basicAuth))

	//Admin Route
	adminGroup.GET("/main", mainAdmin)

	// Cookie route
	cookieGroup.Use(checkCookie)
	cookieGroup.GET("/main", mainCookie)

	//login route
	e.GET("/login", login)

	e.GET("/", Welcome)
	e.GET("/cats", getCats)
	e.POST("/cats", addCart)
	e.Start(":8000")
}
