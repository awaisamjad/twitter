package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// "log"
	"net/http"
	// "strconv"

	// "strings"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var standard_routes = []string{
	"explore",
	"search",
}

type Post struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	User_ID     int    `json:"user_id"`
	Created_At  string `json:"created_at"`
	Updated_At  string `json:"updated_at"`
	Like_Num    int    `json:"like_num"`
	Dislike_Num int    `json:"dislike_num"`
}

type User struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Posts     []Post `json:"posts"`
	Feed      []Post `json:"feed"`
	Following []int  `json:"following"`
	Followers []int  `json:"followers"`
}

func NewUser(id int, username, email, password string, posts, feed []Post, following []int, followers []int) User {
	return User{
		Id:        id,
		Username:  username,
		Email:     email,
		Password:  password,
		Posts:     posts,
		Feed:      feed,
		Following: following,
		Followers: followers,
	}
}

func main() {
	// gin.SetMode(gin.ReleaseMode)

	const db_file string = "/home/awaisamjad/code/go/twitter/backend/src/db/database.db"

	db, err := sql.Open("sqlite3", db_file)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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
		email := c.PostForm("email")
		password := c.PostForm("password")

		query := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
		_, err := db.Exec(query, username, email, password)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "404.html", gin.H{
				"ErrorMessage": "Either username or email is already taken",
			})
			return
		}
		c.Redirect(http.StatusFound, "/"+username)
	})

	// TODO
	// router.GET("/log-in", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "log-in.html", gin.H{})
	// })
	// TODO
	// router.POST("/log-in", func(c *gin.Context) {
	// 	username := c.PostForm("username")
	// 	user := NewUser(20, username, nil, nil, nil, nil)
	// 	users = append(users, user)
	// 	c.Redirect(http.StatusFound, "/"+username)
	// })

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

		query := `SELECT id, username, email, password FROM users WHERE username = ?`
		row := db.QueryRow(query, username)
		err := row.Scan(&userInfo.Id, &userInfo.Username, &userInfo.Email, &userInfo.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				c.HTML(http.StatusNotFound, "404.html", gin.H{
					"ErrorMessage": "The username you entered does not exist. Please check the spelling or register a new account.",
				})
				return
			} else {
				log.Fatal(err)
			}
		}

		// ? Commented lines are for data that isnt used in user.html
		c.HTML(200, "user.html", gin.H{
			// "Id":        userInfo.Id,
			"Username":  userInfo.Username,
			"Posts":     userInfo.Posts,
			// "Feed":      userInfo.Feed,
			// "Following": userInfo.Following,
			// "Followers": userInfo.Followers,
		})
	})

	router.POST("/create-post", func(c *gin.Context) {

		username := c.PostForm("username")
		title := c.PostForm("title")
		contents := c.PostForm("contents")

		// Check if the user exists
		var userId int
		query := `SELECT id FROM users WHERE username = ?`
		err := db.QueryRow(query, username).Scan(&userId)
		if err != nil {
			if err == sql.ErrNoRows {
				c.HTML(http.StatusNotFound, "404.html", gin.H{
					"ErrorMessage": "The username you entered does not exist.",
				})
				return
			} else {
				log.Fatal(err)
			}
		}

		// Insert the new post into the posts table
		query = `INSERT INTO posts (title, content, user_id, like_num, dislike_num) VALUES (?, ?, ?, 0, 0)`
		_, err = db.Exec(query, title, contents, userId)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "404.html", gin.H{
				"ErrorMessage": "Failed to create the post.",
			})
			return
		}

		

		c.Redirect(http.StatusFound, "/"+username)
	})

	// router.POST("/delete-post", func(c *gin.Context) {

	// 	username := c.PostForm("username")
	// 	log.Println(c.Request.URL.Path)
	// 	title := c.PostForm("title")

	// 	for i, user := range users {
	// 		if user.Username == username {
	// 			// Remove the post with the given title and contents
	// 			for j, post := range user.Posts {
	// 				if post.Title == title {
	// 					// this removes the post somehow
	// 					users[i].Posts = append(users[i].Posts[:j], users[i].Posts[j+1:]...)
	// 					break
	// 				}
	// 			}
	// 			break
	// 		}
	// 	}

	// 	c.Redirect(http.StatusFound, "/"+username)
	// })

	router.GET("/test", func(c *gin.Context) {
		cookie, err := c.Cookie("gin_cookie")

		if err != nil {
			cookie = "NotSet"
			c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
		}

		fmt.Printf("Cookie value: %s \n", cookie)
	})

	port := os.Getenv("PORT")

	if port == "" {
		log.Println("Could not find port in .env")
		port = "3000"
	}

	router.Run("0.0.0.0:" + port)
}

func getUserFollowing(db *sql.DB, userId int) ([]int, error) {
	rows, err := db.Query("SELECT following_id FROM user_relationships WHERE user_id = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var following []int
	for rows.Next() {
		var followingId int
		if err := rows.Scan(&followingId); err != nil {
			return nil, err
		}
		following = append(following, followingId)
	}
	return following, nil
}

func getUserFollowers(db *sql.DB, userId int) ([]int, error) {
	rows, err := db.Query("SELECT user_id FROM user_relationships WHERE following_id = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []int
	for rows.Next() {
		var followerId int
		if err := rows.Scan(&followerId); err != nil {
			return nil, err
		}
		followers = append(followers, followerId)
	}
	return followers, nil
}
