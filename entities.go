package main

var standard_routes = []string{
	"explore",
	"search",
}

type Post struct {
	Id          int    `json:"id"`
	Content     string `json:"content"`
	Username    string `json:"username"`
	Created_At  string `json:"created_at"`
	Updated_At  string `json:"updated_at"`
	Like_Num    int    `json:"like_num"`
	Dislike_Num int    `json:"dislike_num"`
}

// ? Following and Followers are a list of user id's (ints) of those they are following/followed by
type User struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Posts      []Post `json:"posts"`
	Feed       []Post `json:"feed"`
	Following  []int  `json:"following"`
	Followers  []int  `json:"followers"`
	Bio        string `json:"bio"`
	Phone_Num  string `json:"phone_num"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"error_message"`
	ErrorType    string `json:"error_type"`
}

type SuccessResponse struct {
    Data string `json:"data"`
    Message  string `json:"message"`
}