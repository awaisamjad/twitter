package main

import (
	"fmt"
	"log"
	// "log"
	"net/http"
	// "strconv"

	// "strings"

	"github.com/gin-gonic/gin"
)

var standard_routes = []string{
	"explore",
	"search",
}

type Post struct {
	Username     string
	Id           int
	Title        string
	Contents     string
	Like_Num     int
	Disklike_Num int
}

type User struct {
	Id        int
	Username  string
	Posts     []Post
	Feed      []Post
	Following []User
	Follwers  []User
}

func NewUser(id int, username string, posts []Post, feed []Post, following []User, followers []User) User {
	return User{
		Id:        id,
		Username:  username,
		Posts:     posts,
		Feed:      feed,
		Following: following,
		Follwers:  followers,
	}
}

var users = []User{
	NewUser(
		1,
		"awais",
		[]Post{
			{Username: "awais", Title: "Post 1", Contents: "Content 1", Like_Num: 0, Disklike_Num: 0},
			{Username: "awais", Title: "Post 2", Contents: "Content 2", Like_Num: 0, Disklike_Num: 0},
		},
		[]Post{
			{Title: "Feed 1", Contents: "Content 1", Like_Num: 0, Disklike_Num: 0},
			{Title: "Feed 2", Contents: "Content 2", Like_Num: 0, Disklike_Num: 0},
		},
		[]User{
			NewUser(3, "awais", nil, nil, nil, nil),
		},
		[]User{
			NewUser(4, "amjad", nil, nil, nil, nil),
		}),

	NewUser(
		2,
		"sameer",
		[]Post{
			{Title: "Post 1", Contents: "Content 1", Like_Num: 0, Disklike_Num: 0},
			{Title: "Post 2", Contents: "Content 2", Like_Num: 0, Disklike_Num: 0},
		},
		[]Post{
			{Title: "Feed 1", Contents: "Content 1", Like_Num: 0, Disklike_Num: 0},
			{Title: "Feed 2", Contents: "Content 2", Like_Num: 0, Disklike_Num: 0},
		},
		[]User{
			NewUser(5, "awais", nil, nil, nil, nil),
		},
		[]User{
			NewUser(6, "amjad", nil, nil, nil, nil),
		}),
}

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.Static("/templates", "./templates")

	for _, route := range standard_routes {
		switch route {
		case "explore":
			router.GET("/"+"explore", func(c *gin.Context) {
				c.HTML(http.StatusOK, "explore.html", gin.H{})
			})
		case "search":
			router.GET("/"+"search", func(c *gin.Context) {
				c.HTML(http.StatusOK, "search.html", gin.H{})
			})
		}
	}

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "feed.html", gin.H{
			"title": "Main website",
		})
	})

	router.GET("/sign-up", func(c *gin.Context) {
		c.HTML(http.StatusOK, "sign-up.html", gin.H{})
	})

	router.POST("/sign-up", func(c *gin.Context) {
		username := c.PostForm("username")
		user := NewUser(20, username, nil, nil, nil, nil)
		users = append(users, user)
		c.Redirect(http.StatusFound, "/"+username)
	})

	router.GET("/:username", func(c *gin.Context) {
		username := c.Param("username")
		for _, notAllowed := range standard_routes {
			if username == notAllowed {
				c.HTML(http.StatusNotFound, "404.html", gin.H{
					"ErrorMessage": "The username you entered is reserved for system use. Please choose a different username.",
				})
				return
			}
		}

		// Check if the user exists
		var userInfo User
		userExists := false
		for _, user := range users {
			if user.Username == username {
				userInfo = user
				userExists = true
				break
			}
		}
		if !userExists {
			c.HTML(http.StatusNotFound, "404.html", gin.H{
				"ErrorMessage": "The username you entered does not exist. Please check the spelling or register a new account.",
			})
			return
		}

		c.HTML(200, "user.html", gin.H{
			"Id":        userInfo.Id,
			"Username":  userInfo.Username,
			"Posts":     userInfo.Posts,
			"Feed":      userInfo.Feed,
			"Following": userInfo.Following,
			"Followers": userInfo.Follwers,
		})
	})

	router.POST("/create-post", func(c *gin.Context) {
		// id := c.PostForm("id")
		// idInt, err := strconv.Atoi(id)
		// if err != nil {
		// 	c.HTML(http.StatusBadRequest, "404.html", gin.H{
		// 		"ErrorMessage": "Invalid ID format.",
		// 	})
		// 	return
		// }

		username := c.PostForm("username")
		title := c.PostForm("title")
		contents := c.PostForm("contents")

		for i, user := range users {
			if user.Username == username {
				newPost := Post{
					Username: username,
					// Id:           idInt,
					Title:        title,
					Contents:     contents,
					Like_Num:     0,
					Disklike_Num: 0,
				}
				users[i].Posts = append(users[i].Posts, newPost)
				break
			}
		}

		c.Redirect(http.StatusFound, "/"+username)
	})

	router.POST("/delete-post", func(c *gin.Context) {

		username := c.PostForm("username")
		log.Println(c.Request.URL.Path)
		title := c.PostForm("title")

		for i, user := range users {
			if user.Username == username {
				// Remove the post with the given title and contents
				for j, post := range user.Posts {
					if post.Title == title {
						// this removes the post somehow
						users[i].Posts = append(users[i].Posts[:j], users[i].Posts[j+1:]...)
						break
					}
				}
				break
			}
		}

		c.Redirect(http.StatusFound, "/"+username)
	})

	router.GET("/test", func(c *gin.Context) {
		cookie, err := c.Cookie("gin_cookie")

		if err != nil {
			cookie = "NotSet"
			c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
		}

		fmt.Printf("Cookie value: %s \n", cookie)
	})

	router.Run(":8080")
}
