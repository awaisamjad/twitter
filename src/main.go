package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	// "slices"
)

var standard_routes = []string{
	"explore",
	"search",
}

type Post struct {
	title        string
	contents     string
	like_num     int
	disklike_num int
}

type User struct {
	username  string
	posts     []Post
	feed      []Post
	following []User
	follwers  []User
}

func (u User) New(username string, posts []Post, feed []Post, following []User, followers []User) User {
	return User{
		username:  username,
		posts:     posts,
		feed:      feed,
		following: following,
		follwers:  followers,
	}
}

var users = []User{
	User{}.New(
		"awais",
		[]Post{
			Post{title: "Post 1", contents: "Content 1", like_num: 0, disklike_num: 0},
			Post{title: "Post 2", contents: "Content 2", like_num: 0, disklike_num: 0},
		},
		[]Post{
			Post{title: "Feed 1", contents: "Content 1", like_num: 0, disklike_num: 0},
			Post{title: "Feed 2", contents: "Content 2", like_num: 0, disklike_num: 0},
		},
		[]User{
			User{}.New("awais", nil, nil, nil, nil),
		},
		[]User{
			User{}.New("amjad", nil, nil, nil, nil),
		}),

	User{}.New(
		"sameer",
		[]Post{
			Post{title: "Post 1", contents: "Content 1", like_num: 0, disklike_num: 0},
			Post{title: "Post 2", contents: "Content 2", like_num: 0, disklike_num: 0},
		},
		[]Post{
			Post{title: "Feed 1", contents: "Content 1", like_num: 0, disklike_num: 0},
			Post{title: "Feed 2", contents: "Content 2", like_num: 0, disklike_num: 0},
		},
		[]User{
			User{}.New("awais", nil, nil, nil, nil),
		},
		[]User{
			User{}.New("amjad", nil, nil, nil, nil),
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

		userExists := false
		for _, user := range users {
			if user.username == username {
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

		// user_info := getUserInfo(username)

		c.HTML(200, "user.html", gin.H{
			"Username": username,
		})
	})

	router.Run(":8080")
}
