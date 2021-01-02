package models

// JSON models

//
// typical generic json
//

type SuccessJson struct {
	Success bool   `json:"success"`
	Message string `json:"status"`
}

type Error struct {
	Success bool  `json:"success"`
	// Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

//
// auth json models
//

// Create a struct to read the username and password from the request body
type Credentials struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	// Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUserJson struct {
	Success bool    	`json:"success"`
	Message string  	`json:"message"`
	//Id 		uint    	`json:"id"`
	Token   string  	`json:"token"`
	User    User `json:"user"`
}

type CreateUserJson struct {
	Success bool    	`json:"success"`
	Message string  	`json:"message"`
	UserId  uint    	`json:"userId"`
	Name    string  	`json:"name"`
	Email   string  	`json:"email"`
	Token   string  	`json:"token"`
	User    User 		`json:"user"`
}


//
// track json models
//


type TrackIndexJson struct {
	UserID    uint	`json:"userID"`
}

//type TrackIndexJson struct {
//	Success bool        `json:"success"`
//	Message  string		`json:"message"`
//	TrackList [] Track  `json:"trackList"`
//	User    User 		`json:"user"`
//}

type CreateLocalTrackJson struct {
	//Track   Track  	    `json:"track"`
	UserID    uint      `json:"userID"`
	Title     string    `json:"title"`
	Desc      string    `json:"desc"`
	//User    User 		`json:"user"`
}

type CreateLocalTrackResponseJson struct {
	Success   bool    	`json:"success"`
	Message   string  	`json:"message"`
	Track     Track  	`json:"track"`
	//Title   string  	`json:"title"`
	//Desc    string  	`json:"desc"`
	//UserID  uint    	`json:"userID"`
	//User    User 		`json:"user"`
}

type OneTrackJson struct {
	TrackID    uint	  `json:"trackID"`
}



//type OneTrackJson struct {
//	Success  bool       `json:"success"`
//	Message  string		`json:"message"`
//	Track     Track     `json:"track"`
//}