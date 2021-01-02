package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/steven7/go-createmusic/go/compose_ai"
	"github.com/steven7/go-createmusic/models"
	"net/http"
	"strconv"
)

func NewTracksAPI(ts models.TrackService, fs models.FileService, r *mux.Router) *TrackController {
	return &TrackController{
		ts:                	   ts,
		fs: 				   fs,
		//is:                	   is,
		//mfs: 				   mfs,
		r:                 	   r,
	}
}

// GET /api/tracks/index
func (t *TrackController) IndexWithAPI(w http.ResponseWriter, r *http.Request) {

	fmt.Println("IndexWithAPI!!")

	//
	// validate jwt token is done with middleware
	//
	var info models.TrackIndexJson
	ParseJSONParameters(w, r, &info)

	uintUserID := info.UserID
	fmt.Println(uintUserID)

	//if err != nil {
	//	//app.serverError(w, err)
	//	fmt.Println(uintUserID)
	//	errorData := models.Error {
	//		Title:  "Error parsing userID. Could not convert to uint",
	//		Detail: err.Error(),
	//	}
	//	WriteJson(w, errorData)
	//	return
	//}

	tracks, err := t.ts.ByUserID(uint(uintUserID))
	fmt.Println("json")
	for _ , track := range tracks {
		fmt.Println(track)
	}

	if err != nil {
		// app.serverError(w, err)
		errorData := models.Error {
			Title:  "Error fetching tracks for user",
			Detail: err.Error(),
		}
		WriteJson(w, errorData)
		return
	}

	WriteJson(w, tracks)

}

// GET /api/tracks/one
func (t *TrackController) GetTrackWithAPI(w http.ResponseWriter, r *http.Request) {

	//
	// validate jwt token is done with middleware
	//
	var info models.OneTrackJson
	ParseJSONParameters(w, r, &info)
	uintTrackID := info.TrackID

	/*
	uintTrackID, err := strconv.ParseUint(info.TrackID, 10, 64)
	if err != nil {
		//app.serverError(w, err)
		errorData := models.Error {
			Title:  "Error parsing userID. Could not convert to uint",
			Detail: err.Error(),
		}
		WriteJson(w, errorData)
		return
	}
	*/
	track, err := t.ts.ByID(uint(uintTrackID))
	// if error choose way to respond
	if err != nil {
		errorData := models.Error {
			Title:  "Error fetching one track",
			Detail: err.Error(),
		}
		WriteJson(w, errorData)
		return
	}

	// music file
	//musicfile, err := t.fs.ByTrackID(track.ID, models.FileTypeMusic)
	//if err != nil {
	//	errorData := models.Error {
	//		Title:  "Error fetching the track file",
	//		Detail: err.Error(),
	//	}
	//	WriteJson(w, errorData)
	//	return
	//}

	// cover image
	//coverimage, err := t.fs.ByTrackID(track.ID, models.FileTypeImage)
	//if err != nil {
	//	errorData := models.Error {
	//		Title:  "Error fetching the cover image file",
	//		Detail: err.Error(),
	//	}
	//	WriteJson(w, errorData)
	//	return
	//}

	WriteJson(w, track)
}


// GET /api/tracks/one/coverimage
func (t *TrackController) GetTrackCoverFileWithAPI(w http.ResponseWriter, r *http.Request) {

	fmt.Println("GetTrackCoverFileWithAPI")

	var info models.OneTrackJson
	ParseJSONParameters(w, r, &info)
	uintTrackID :=  info.TrackID

	fmt.Println(uintTrackID)

	// cover image
	coverimage, err := t.fs.ByTrackID(uint(uintTrackID), models.FileTypeImage)
	if err != nil {
		errorData := models.Error {
			Title:  "Error fetching the cover image file",
			Detail: err.Error(),
		}
		WriteJson(w, errorData)
		return
	}

	WriteFile(w, coverimage.ImagePath())

}

// GET /api/tracks/one/musicfile
func (t *TrackController) GetTrackMusicFileWithAPI(w http.ResponseWriter, r *http.Request) {

	fmt.Println("GetTrackMusicFileWithAPI")

	var info models.OneTrackJson
	ParseJSONParameters(w, r, &info)
	uintTrackID := info.TrackID

	fmt.Println(uintTrackID)

	// music file

	musicfile, err := t.fs.ByTrackID(uint(uintTrackID), models.FileTypeMusic)
	if err != nil {
		errorData := models.Error {
			Title:  "Error fetching the track file",
			Detail: err.Error(),
		}
		WriteJson(w, errorData)
		return
	}

	WriteFile(w, musicfile.MusicPath())

}

// POST /api/tracks/createlocal
func (t *TrackController) CreateLocalWithAPI(w http.ResponseWriter, r *http.Request) {

	//
	// validate jwt token is done with middleware
	//

	r.ParseMultipartForm(128)

	f := r.MultipartForm

	fmt.Println("form data ", f)

	uid := f.Value["userID"][0]
	fmt.Println("user id: ", uid)

	userID, err := strconv.ParseUint(f.Value["userID"][0], 10, 32)
	if err != nil {
		userID = 0
		fmt.Println(err)
	}
	title := f.Value["title"][0]
	artist := f.Value["artist"][0]
	desc := f.Value["desc"][0]

	///
	///
	///


	// Create track object in memory
	track := models.Track {
		Title:  title, //trackData.Track.Title,
		Artist: artist,
		Description: desc, //"Created with create local track api endpoint",
		UserID: uint(userID), //trackData.UserID, //trackData.User.ID,
	}

	//
	//
	// Create track file in database
	//
	//

	fmt.Println("Creating for track id %ui", track.ID)

	if err := t.ts.Create(&track); err != nil {
		// app.serverError(w, err)
		errorData := models.Error {
			Title:  "Error creating text",
			Detail: err.Error(),
		}
		WriteJson(w, errorData)
		return
	}

	//
	// validate input files first
	//

	//
	// cover image
	//
	imgArr := f.File["coverimage"]

	if len(imgArr) > 0 {

		imageHeader := f.File["coverimage"][0]
		fmt.Println(imageHeader)

		imagefile, err := imageHeader.Open()
		if err != nil {
			errorData := models.Error {
				Title:  "Error with cover image",
				Detail: err.Error(),
			}
			WriteJson(w, errorData)
			return
		}
		defer imagefile.Close()

		err = t.fs.Create(track.ID, imagefile, imageHeader.Filename, models.FileTypeImage)
		if err != nil {
			errorData := models.Error {
				Title:  "Error uploading cover image",
				Detail: err.Error(),
			}
			WriteJson(w, errorData)
			return
		}
		track.CoverImageFilename = imageHeader.Filename

	}

	//
	//
	// music file
	//
 	//

	if len(f.File["musicfile"]) == 0 {
		errorData := models.Error {
			Title:  "Error creating track",
			Detail: "A valid music file must be included",
		}
		WriteJson(w, errorData)
		return
	}

	musicHeader := f.File["musicfile"][0]
	fmt.Println(musicHeader)

	musicfile, err := musicHeader.Open()
	if err != nil {
		errorData := models.Error {
			Title:  "Error with uploaded file",
			Detail: err.Error(),
		}
		WriteJson(w, errorData)
		return
	}
	defer musicfile.Close()
	err = t.fs.Create(track.ID, musicfile, musicHeader.Filename, models.FileTypeMusic)
	if err != nil {
		errorData := models.Error {
			Title:  "Error uploading file",
			Detail: err.Error(),
		}
		WriteJson(w, errorData)
		return
	}
	track.MusicFileFilename = musicHeader.Filename


	//fmt.Println("created track ")
	//fmt.Println(track)
	//fmt.Println(track.UserID)
	//fmt.Println(track.Title)
	//fmt.Println(track.ID)

	////
	///
	///
	createLocalTrackResponse := models.CreateLocalTrackResponseJson{
		Success: true,
		Message: "Track successfully created!!",
		Track: track,
	}

	WriteJson(w, createLocalTrackResponse)

}

// POST /tracks/createWithComposeAI
// only work with deep jazz for the momment
func (t *TrackController) CreateWithComposeAI(w http.ResponseWriter, r *http.Request) {

	fmt.Println("CreateWithComposeAI!!")

	//r.ParseMultipartForm(128)
	//
	//f := r.MultipartForm
	//
	//fmt.Println("form data ", f)
	//
	//uid := f.Value["userID"][0]
	//fmt.Println("user id: ", uid)
	//
	//userID, err := strconv.ParseUint(f.Value["userID"][0], 10, 32)
	//if err != nil {
	//	userID = 0
	//	fmt.Println(err)
	//}
	userID := 1
	//title := f.Value["title"][0]
	//artist := f.Value["artist"][0]
	//desc := f.Value["desc"][0]

	///
	///
	///


	// Create track object
	track := models.Track {
		Title:  "compose.ai - jazz", //trackData.Track.Title,
		Artist: "compose.ai",
		Description: "compose.ai created this", //desc, //"Created with create local track api endpoint",
		UserID: uint(userID), //trackData.UserID, //trackData.User.ID,
	}

	err := compose_ai.ComposeJazz()
	if err != nil {

	}

	//
	// Create track file in database if compose ai successful
	//

	//fmt.Println("Creating for track id %ui", track.ID)

	//if err := t.ts.Create(&track); err != nil {
	//	// app.serverError(w, err)
	//	errorData := models.Error {
	//		Title:  "Error creating text",
	//		Detail: err.Error(),
	//	}
	//	WriteJson(w, errorData)
	//	return
	//}

	createLocalTrackResponse := models.CreateLocalTrackResponseJson{
		Success: true,
		Message: "Track successfully tested. This feature is not done yet", //"Track successfully created!!",
		Track: track,
	}

	WriteJson(w, createLocalTrackResponse)
}

// POST /track/create
func (t *TrackController) CreateDJWithAPI(w http.ResponseWriter, r *http.Request) {

	//userId := 2 // default text id
	//trackId := 10 // track 10 is 0x00ss.jpg cover image and Tchaikovsky-Waltz-op39-no8 file


}

// POST /track/create
func (t *TrackController) CreateDJ_Jazz_WithAPI(w http.ResponseWriter, r *http.Request) {

	//userId := 2 // default text id
	//trackId := 10 // track 10 is 0x00ss.jpg cover image and Tchaikovsky-Waltz-op39-no8 file


	//cmd := exec.Command("script.py")
	//
	//cmd := exec.Command("cmd", python, script)
	//out, err := cmd.Output()
	//fmt.Println(string(out))

}


// GET /track/download
func (t *TrackController) DownloadFileWithAPI(w http.ResponseWriter, r *http.Request) {

}

// GET /track/stream
func (t *TrackController) StreamWithAPI(w http.ResponseWriter, r *http.Request) {

}