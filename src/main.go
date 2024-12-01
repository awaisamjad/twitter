package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	// "slices"
)

var standard_routes = []string{
	"explore",
	"search",
}

type Post struct {
	Title        string
	Contents     string
	Like_Num     int
	Disklike_Num int
}

type User struct {
	Username  string
	Posts     []Post
	Feed      []Post
	Following []User
	Follwers  []User
}

func NewUser(username string, posts []Post, feed []Post, following []User, followers []User) User {
	return User{
		Username:  username,
		Posts:     posts,
		Feed:      feed,
		Following: following,
		Follwers:  followers,
	}
}

var users = []User{
	NewUser(
		"awais",
		[]Post{
			{Title: "Post 1", Contents: "Content 1", Like_Num: 0, Disklike_Num: 0},
			{Title: "Post 2", Contents: "Content 2", Like_Num: 0, Disklike_Num: 0},
		},
		[]Post{
			{Title: "Feed 1", Contents: "Content 1", Like_Num: 0, Disklike_Num: 0},
			{Title: "Feed 2", Contents: "Content 2", Like_Num: 0, Disklike_Num: 0},
		},
		[]User{
			NewUser("awais", nil, nil, nil, nil),
		},
		[]User{
			NewUser("amjad", nil, nil, nil, nil),
		}),

	NewUser(
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
			NewUser("awais", nil, nil, nil, nil),
		},
		[]User{
			NewUser("amjad", nil, nil, nil, nil),
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
		fmt.Println(userInfo.Username)
		c.HTML(200, "user.html", gin.H{
			"Username":  userInfo.Username,
			"Posts":     userInfo.Posts,
			"Feed":      userInfo.Feed,
			"Following": userInfo.Following,
			"Followers": userInfo.Follwers,
		})
	})

	router.POST("/create-post", func(c *gin.Context) {
		username := c.PostForm("username")
		title := c.PostForm("title")
		contents := c.PostForm("contents")

		for i, user := range users {
			if user.Username == username {
				newPost := Post{
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

	router.DELETE("/delete-post", func(c *gin.Context) {
		username := c.PostForm("username")
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

	router.Run(":8080")
}
