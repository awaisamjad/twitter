package main

import (
	"database/sql"
	"time"
	// "fmt"
	"log"
	"os"
	"slices"

	// "log"
	"net/http"
	// "strconv"
	// "strings"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var store *sessions.CookieStore


func main() {
	// gin.SetMode(gin.ReleaseMode)

	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	initSessionStore()

	go func() {
        for {
            time.Sleep(24 * time.Hour)
            rotateSessionKey()
        }
    }()

	const db_file string = "db/database.db"

	db, err := sql.Open("sqlite3", db_file)
	if err != nil {
		log.Fatal("DB Error :", err)
	}
	defer db.Close()

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	// router.Static("/templates", "./templates")

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
		session, err := store.Get(c.Request, "current-session")
		session.Options = &sessions.Options{
			Path: "/",
			Secure: true,
			HttpOnly: true,
		}
		if err != nil {
			
			log.Panic("Failed to create session")
		}
		username, ok := session.Values["username"].(string)
		if !ok || username == "" {
			log.Println("Show generic feed")
		} else {
			log.Println("Show custom feed for", username)
		}
		c.Redirect(http.StatusFound, "/feed")
	})

	router.GET("/feed", func(c *gin.Context) {
		c.HTML(http.StatusFound, "feed.html", gin.H{})
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
		password, err := HashPassword(password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				ErrorMessage: "Internal server error",
				ErrorType:    "internal_error",
			})
			return
		}

		// ? Validating inputs
		if !IsNameValid(first_name) || !IsNameValid(last_name) {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				ErrorMessage: "Invalid First or Last name. Ensure they are between 2 and 50 characters long, contain no special characters, and have no spaces.",
				ErrorType:    "name_invalid",
			})
			return
		}

		if !IsUsernameValid(username) {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				ErrorMessage: "Invalid Username. Ensure it is between 3 and 20 characters long, and has no spaces.",
				ErrorType:    "username_invalid",
			})
			return
		}

		if !IsEmailValid(email) {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				ErrorMessage: "Invalid email",
				ErrorType:    "email_invalid",
			})
			return
		}

		if !IsPasswordValid(password) {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				ErrorMessage: "Password does not meet the required criteria.\nEnsure it is at least 8 characters long, includes at least one uppercase letter, one lowercase letter, one digit, and one special character (e.g., !@#$%^&*...).",
				ErrorType:    "password_invalid",
			})
			return
		}

		// ? Check if username already exists in dB
		var user_exists bool
		var email_exists bool

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

		// ? Check if email already exists in dB
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

		c.JSON(http.StatusOK, SuccessResponse{
			Data:    username,
			Message: "Signup Sucessful",
		})
		//* Redirection to log in page happens in sign-up.html with js
	})

	// TODO
	router.GET("/log-in", func(c *gin.Context) {
		c.HTML(http.StatusOK, "log-in.html", gin.H{})
	})

	router.POST("/log-in", func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")

		var storedPassword, username string
		//? Check if email exists and get the stored password and username
		err = db.QueryRow("SELECT password, username FROM users WHERE email = ?", email).Scan(&storedPassword, &username)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusBadRequest, ErrorResponse{
					ErrorMessage: "Email is not registered",
					ErrorType:    "email_does_not_exist",
				})
			} else {
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					ErrorMessage: "Internal server error",
					ErrorType:    "internal_error",
				})
				log.Println("POST Method, Log in")
			}
			return
		}

		//? Check if the password matches
		if !CheckPasswordAgainstPasswordHash(password, storedPassword) {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				ErrorMessage: "Incorrect password",
				ErrorType:    "incorrect_password",
			})
			return
		}
		
		session, _ := store.Get(c.Request, "current-session")
        session.Values["username"] = username
        session.Save(c.Request, c.Writer)
		log.Println("Session username", username)
		c.JSON(http.StatusOK, SuccessResponse{
			Data:    username,
			Message: "Log in Sucessful",
		})
		//* Redirection to user account happens in log-in.html with js
	})

	router.GET("/:username", func(c *gin.Context) {
		c.HTML(http.StatusOK, "user.html", gin.H{})
	})

	router.POST("/api/posts/:username", func(c *gin.Context) {
		session, _ := store.Get(c.Request, "current-session")
		username := c.Param("username")
		session_username := session.Values["username"]
		log.Println(session_username)
		if session_username != username {
			c.HTML(http.StatusFound, "error.html", gin.H{
				"ErrorMessage": "User is not registered to view account",
			})
			return
		}

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
		err = row.Scan(&userInfo.Id, &userInfo.Username, &userInfo.Email, &userInfo.Password)
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
		// c.HTML(200, "user.html", gin.H{
		
		// 	// "Id":        userInfo.Id,
		// 	"Username": userInfo.Username,
		// 	"Posts":    posts,
		// 	// "Feed":      userInfo.Feed,
		// 	// "Following": userInfo.Following,
		// 	// "Followers": userInfo.Followers,
		// })

		c.JSON(http.StatusOK, gin.H{
			"Username" : userInfo.Username,
			"Posts" : posts,
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
		session, _ := store.Get(c.Request, "current-session")
		for key, value := range session.Values {
			log.Printf("%s: %v\n", key, value)
		}

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
