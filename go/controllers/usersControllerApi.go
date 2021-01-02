package controllers

import (
	"fmt"
	"github.com/steven7/go-createmusic/go/email"
	"github.com/steven7/go-createmusic/go/models"
	"net/http"
)

// GET /signup
//func (u *Users) New(w http.ResponseWriter, r *http.Request) {
//	var form SignupForm
//	parseURLParams(r, &form)
//	u.NewView.Render(w, r, form)
//}

func NewUsersAPI(us models.UserService, emailer *email.Client) *Users {
	return &Users{
		us:        	  us,
		emailer:   	  emailer,
	}
}

// Create is used to process the signup form when a user
// tries to create a new user account.
//
// POST /create
func (u *Users) CreateWithAPI(w http.ResponseWriter, r *http.Request) {

	//
	// Parse parameters
	//
	var credentials models.Credentials

	ParseJSONParameters(w, r, &credentials)

	//
	// Create user in database
	//
	user := models.User{
		Name:     credentials.Name,
		Email:    credentials.Email,
		Password: credentials.Password,
	}
	if err := u.us.Create(&user); err != nil {

		errorData := models.Error {
			Success: false,
			Title:  "Could not create user",
			Detail: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		WriteJson(w, errorData)

		return
	}

	//
	//Create new JWT token for the newly registered account
	//
	tokenString, err := CreateJWTToken(user.ID)
	if err != nil {
		errorData := models.Error {
			Title:  "Could not create user",
			Detail: err.Error(),
		}
		WriteJson(w, errorData)
		return
	}

	createUserJson := models.CreateUserJson {
		Message: "Account has been created",
		Success: true,
		UserId:  user.ID,
		Name:    credentials.Name,
		Email:   credentials.Email,
		Token:   tokenString,
		//User:
	}

	// return response
	WriteJson(w, createUserJson)
}

// Authenticate/login is used to sign the given user in via cookies
func (u *Users) AuthenticateWithAPI(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials

	ParseJSONParameters(w, r, &credentials)

	email := credentials.Email
	pw := credentials.Password

	user, err := u.us.Authenticate(email, pw)
	if err == models.ErrInvalidCredentials {
		errorData := models.Error {
			Success: false,
			Title:  "ErrCredentials",
			Detail: err.Error(),
		}
		fmt.Println(err.Error())
		WriteJson(w, errorData)
		return
	} else if err == models.ErrPasswordIncorrect {
		errorData := models.Error {
			Success: false,
			Title:  "ErrPassword",
			Detail: err.Error(),
		}
		fmt.Println(err.Error())
		WriteJson(w, errorData)
		return
	} else if err != nil {

		errorData := models.Error{
			Success: false,
			Title:  "Could not authenticate",
			Detail: err.Error(),
		}
		fmt.Println("Could not authenticate")
		fmt.Println(err.Error())
		//fmt.Println(err.Error())
		WriteJson(w, errorData)
		return
	}

	//
	// Authentication successful
	// Create jwt
	//
	token, err := CreateJWTToken(user.ID)
	if err != nil {
		//c.JSON(http.StatusUnprocessableEntity, err.Error())
		errorData := models.Error {
			Title:  "Token error",
			Detail: "Token could not be created",
		}
		fmt.Println("Token could not be created")
		WriteJson(w, errorData)
		return
	}

	user.Token = token


	loginUserJson := models.LoginUserJson{
		Message: "User successfully validated!!",
		Success: true,
		Token:   token,
		//Id: 	 id,
		User:    *user,
	}

	fmt.Println("User successfully validated!!")

	WriteJson(w, loginUserJson)
}

//
// Logout will happen mostly on the client side for now
//
// POST /logout
/*
func (u *Users) LogoutWithAPI(w http.ResponseWriter, r *http.Request) {
	// First expire the user's cookie”
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	// Then we update the user with a new remember token
	user := context.User(r.Context())
	// We are ignoring errors for now because they are
	// unlikely, and even if they do occur we can't recover
	// now that the user doesn't have a valid cookie
	token, _ := rand.RememberToken()
	user.Remember = token
	u.us.Update(user)
	// Finally send the user to the home page
	http.Redirect(w, r, "/", http.StatusFound)
}
*/
