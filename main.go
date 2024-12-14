package main

import (
	"database/sql"
	// "fmt"
	"log"
	"os"
	"slices"

	// "log"
	"net/http"
	// "strconv"

	// "strings"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)

	const db_file string = "db/database.db"

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
		first_name := c.PostForm("first_name")
		last_name := c.PostForm("last_name")
		username := c.PostForm("username")
		email := c.PostForm("email")
		password := c.PostForm("password")

		c.SetCookie("username", username, 3600, "/", "localhost", false, true)

		var user_exists bool
		var email_exists bool
		var err error

		// ? Check if user exists
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&user_exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				ErrorMessage: "Internal server error",
				ErrorType:    "internal_error",
			})
			return
		}
		if user_exists {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				ErrorMessage: "Username is already in use",
				ErrorType:    "username_exists",
			})
			return
		}

		// ? Check if email exists
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&email_exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				ErrorMessage: "Internal server error",
				ErrorType:    "internal_error",
			})
			return
		}
		if email_exists {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				ErrorMessage: "Email is already in use",
				ErrorType:    "email_exists",
			})
			return
		}

		query := `INSERT INTO users (username, first_name, last_name, email, password) VALUES (?, ?, ?, ?, ?)`
		_, err = db.Exec(query, username, first_name, last_name, email, password)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"ErrorMessage": "Either username or email is already taken",
			})
			return
		}

		c.JSON(http.StatusOK, SuccessResponse {
			Data: username,
			Message: "Signup Sucessful",
		})
		// Redirection to user account happens in sign-up.html with js
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
				c.HTML(http.StatusNotFound, "error.html", gin.H{
					"ErrorMessage": "The username you entered is reserved for system use. Please choose a different username.",
				})
				return
			}
		}

		//? Check if the user exists
		var userInfo User
		query := `SELECT id, username, email, password FROM users WHERE username = ?`
		row := db.QueryRow(query, username)
		err := row.Scan(&userInfo.Id, &userInfo.Username, &userInfo.Email, &userInfo.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				c.HTML(http.StatusNotFound, "error.html", gin.H{
					"ErrorMessage": "The username you entered does not exist. Please check the spelling or register a new account.",
				})
				return
			} else {
				c.HTML(http.StatusNotFound, "error.html", gin.H{
					"ErrorMessage": "Error with database connection",
				})
				return
			}
		}

		// ? Get the users posts
		posts, err := getUserPosts(db, username)
		slices.Reverse(posts)
		// TODO Improve error handling
		if err != nil {
			log.Fatal("Failed to get the users posts", err)
		}

		// ? Commented lines are for data that isnt used in user.html
		c.HTML(200, "user.html", gin.H{
			// "Id":        userInfo.Id,
			"Username": userInfo.Username,
			"Posts":    posts,
			// "Feed":      userInfo.Feed,
			// "Following": userInfo.Following,
			// "Followers": userInfo.Followers,
		})
		log.Println("username : ", username)

	})

	// ! Create post here cant access the username from the post. That is a bad way to do this and should use cookies instead
	// ! Start from sign up and log in and implement cookuies properly from there. Then implement other functiuonality
	router.POST("/create-post", func(c *gin.Context) {

		username := c.PostForm("username")
		contents := c.PostForm("contents")

		query := `INSERT INTO posts (content, username, like_num, dislike_num) VALUES (?, ?, 0, 0)`
		_, err = db.Exec(query, contents, username)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"ErrorMessage": "Failed to create the post.",
			})
			return
		}
		log.Println("UsernAMAME : ", username)
		c.Redirect(http.StatusFound, "/"+username)
	})

	// TODO
	router.POST("/delete-post", func(c *gin.Context) {
		// username := c.PostForm("username")
		id := c.PostForm("id")

		query := `DELETE FROM posts WHERE id = ?`
		_, err := db.Exec(query, id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to delete the post.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Post deleted successfully.",
		})

		// c.Redirect(http.StatusFound, "/"+username)
	})

	router.GET("/test", func(c *gin.Context) {
		cookie, err := c.Cookie("username")

		if err != nil {
			cookie = "NotSet"
			c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
		}

		// fmt.Printf("Cookie value: %s \n", cookie)
		type POS struct {
			X int
			Y int
			Z int
		}

		// positions := []POS{
		// 	{X: 1, Y: 2, Z: 3},
		// 	{X: 4, Y: 5, Z: 6},
		// 	{X: 7, Y: 8, Z: 9},
		// }

		c.HTML(200, "test.html", gin.H{
			"Positions": cookie,
		})

		log.Println(cookie)
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

func getUserPosts(db *sql.DB, username string) ([]Post, error) {

	query := `SELECT * FROM posts WHERE username = ?`
	rows, err := db.Query(query, username)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var postID int
		var postContent string
		var postUsername string
		var postCreatedAt string
		var postUpdatedAt string
		var postLikeNum int
		var postDislikeNum int

		if err := rows.Scan(&postID, &postContent, &postUsername, &postCreatedAt, &postUpdatedAt, &postLikeNum, &postDislikeNum); err != nil {
			return nil, err
		}
		posts = append(posts, Post{
			Id:          postID,
			Content:     postContent,
			Username:    postUsername,
			Created_At:  postCreatedAt,
			Updated_At:  postUpdatedAt,
			Like_Num:    postLikeNum,
			Dislike_Num: postDislikeNum,
		})
	}
	return posts, nil
}
