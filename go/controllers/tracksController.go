package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/steven7/go-createmusic/go/context"
	"github.com/steven7/go-createmusic/models"
	"github.com/steven7/go-createmusic/views"
	"net/url"
	"time"

	//"io"
	"log"
	"net/http"
	//"os"
	"strconv"
	"strings"
	"sync"
)

const (
	IndexTracks    = "index_tracks"
	ShowTrack      = "show_track"
	CreateTrack    = "create_track"
	EditTrack      = "edit_track"
	PlayTrack      = "play_track"
	// EditTrackLocal = "edit_tracks"
	// EditTrackDJ    = "edit_tracks"
	// maxMultipartMem = 1 << 20 // 1 megabyte
)

//func NewTracksController(ts models.TrackService, is models.ImageService, mfs models.MusicFileService, r *mux.Router) *TrackController {

func NewTracksController(ts models.TrackService, fs models.FileService, r *mux.Router) *TrackController {
	return &TrackController{
		ChooseTypeView:   	   views.NewView("bootstrap", "tracks/chooseCreateTrackType"),
		// dj
		ChooseDJOptionsView:   views.NewView("bootstrap", "tracks/chooseDJOptions"),
		CreateDJWorkingView:   views.NewView("bootstrap", "tracks/createDJWorking"),
		CreateDJCompleteView:  views.NewView("bootstrap", "tracks/createDJComplete"),
		// local
		CreateLocalView:   	   views.NewView("bootstrap", "tracks/createLocalTrack"),
		EditLocalView:     	   views.NewView("bootstrap", "tracks/editLocalTrack"),
		//EditView:              views.NewView("bootstrap", "tracks/editTracks"),
		IndexView:         	   views.NewView("bootstrap", "tracks/index"),
		ShowView:          	   views.NewView("bootstrap", "tracks/showTracks"),
		PlayView:			   views.NewView("bootstrap", "tracks/playTrack"),
		ts:                	   ts,
		fs: 				   fs,
		//is:                	   is,
		//mfs: 				   mfs,
		r:                 	   r,
	}
}

type TrackController struct {
	ChooseTypeView   	  *views.View
	ChooseDJOptionsView	  *views.View
	CreateDJWorkingView   *views.View
	CreateDJCompleteView  *views.View
	CreateLocalView  	  *views.View
	EditLocalView     	  *views.View
	//EditDJCreatedView     *views.View
	//EditView              *views.View
	IndexView         	  *views.View
	ShowView           	  *views.View
	PlayView           	  *views.View
	ts                	  models.TrackService
	fs                	  models.FileService
	//is                	  models.ImageService
	//mfs					  models.MusicFileService
	r                 	  *mux.Router
}

type TrackForm struct {
	Title      string `schema:"title"`
	Musicfile  string   `schema:"musicfile"`
	CoverImage string   `schema:"image"`
}

// GET /tracks
func (t *TrackController) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	tracks, err := t.ts.ByUserID(user.ID)
	if err != nil {
		// We could attempt to display the index page with
		// no galleries and an error message, but that isn't
		// really more useful than a generic error message so
		// I didn't change this.
		// Regardless of what page we display, we should try
		// to make sure the error is logged so we can debug it
		// later.
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	var vd views.Data
	vd.Yield = tracks
	t.IndexView.Render(w, r, vd)
}

// GET /galleries/:id
func (t *TrackController) Play(w http.ResponseWriter, r *http.Request) {
	track, err := t.trackByID(w, r)
	if err != nil {
		// The galleryByID method will already render the error
		// for us, so we just need to return here.
		return
	}
	fmt.Println("play ")
	fmt.Println(track.ID)
	fmt.Println(track.MusicFile)
	fmt.Println("play " + track.MusicFile.Filename)
	fmt.Println("play " + track.MusicFile.MusicPath())
	// track.MusicFile.Filename = "Tchaikovsky-op19-no3-fueillet-d-album.mp3"
	var vd views.Data
	vd.Yield = track
	t.PlayView.Render(w, r, vd)
}


// POST /tracks
func (t *TrackController) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form TrackForm
	fmt.Println(form)

	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		t.ChooseTypeView.Render(w, r, vd)
		return
	}
	user := context.User(r.Context())

	fmt.Println(user)

	track := models.Track {
		Title:  form.Title,
		UserID: user.ID,
	}

	fmt.Println(track)

	if err := t.ts.Create(&track); err != nil {
		vd.SetAlert(err)
		t.ChooseTypeView.Render(w, r, vd)
		return
	}

	//url, err := g.r.Get(ShowGallery).URL("id",
	//	strconv.Itoa(int(gallery.ID)))
	url, err := t.r.Get(CreateTrack).URL("id",
		strconv.Itoa(int(track.ID)))

	fmt.Println(url)
	fmt.Println(url.Path)

	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

	// POST /tracks/createWithDJ
func (t *TrackController) ChooseDJOptions(w http.ResponseWriter, r *http.Request) {

	var vd views.Data

	t.ChooseDJOptionsView.Render(w, r, vd)

}

// POST /track
func (t *TrackController) CreateWithDJ(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form TrackForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		t.ChooseDJOptionsView.Render(w, r, vd)
		return
	}
	user := context.User(r.Context())
	track := models.Track{
		Title:  form.Title,
		UserID: user.ID,
	}
	//if err := t.ts.Create(&track); err != nil {
	//	vd.SetAlert(err)
	//	fmt.Println(err)
	//	t.ChooseDJOptionsView.Render(w, r, vd)
	//	return
	//}

	//url, err := g.r.Get(ShowGallery).URL("id",
	//	strconv.Itoa(int(gallery.ID)))
	url, err := t.r.Get(IndexTracks).URL("id",
		strconv.Itoa(int(track.ID)))
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

func (t *TrackController) CreateWithDJWorking(w http.ResponseWriter, r *http.Request) {


	var vd views.Data



	t.CreateDJWorkingView.Render(w, r, vd)

	// timer



	// wait (time is in seconds)
	t1 := time.NewTimer(5 * time.Second)
	//<- t1.C
	//fmt.Println("Timer expired")
	//
	//http.Redirect(w, r, "/tracks/createWithDJ/Complete", http.StatusFound)

	go func() {
		<-t1.C
		fmt.Println("Timer fired")
		//t.CreateDJCompleteView.Render(w, r, vd)
		//url, err := t.r.Get(IndexTracks).URL("id",
		//	strconv.Itoa(int(track.ID)))
		//http.Redirect(w, r, "/tracks/createWithDJ/Complete", http.StatusFound)

		form := url.Values{}
		// req, err :=
			http.NewRequest("POST", "/tracks/createWithDJ/Complete", strings.NewReader(form.Encode()))

		t.CreateDJCompleteView.Render(w, r, vd)
	}()


	//track, err := t.trackByID(w, r)
	//if err != nil {
	//
	//	return
	//}

	//vd.Yield = track

}

func (t *TrackController) CreateWithDJComplete(w http.ResponseWriter, r *http.Request) {

	//track, err := t.trackByID(w, r)
	//if err != nil {
	//
	//	return
	//}

	fmt.Println("CreateWithDJComplete")
	var vd views.Data

	user := context.User(r.Context())
	track := models.Track {
		Title:  "Cool DJ Song!!",
		UserID: user.ID,
	}

	vd.Yield = track

	t.CreateDJCompleteView.Render(w, r, vd)
}

	// POST /tracks/newLocal
// GET /galleries/:id/edit
/*
func (t *TrackController) CreateLocal(w http.ResponseWriter, r *http.Request) {
//func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	track, err := t.trackByID(w, r)
	if err != nil {
		// The galleryByID method will already render the error
		// for us, so we just need to return here.
		return
	}
	// A user needs logged in to access this page, so we can
	// assume that the RequireUser middleware has run and
	// set the user for us in the request context.
	user := context.User(r.Context())
	if track.UserID != user.ID {
		http.Error(w, "You do not have permission to edit "+
			"this track", http.StatusForbidden)
		return
	}
	var vd views.Data
	vd.Yield = track
	fmt.Println("creat local on  ", track)
	t.NewLocal.Render(w, r, vd)
}
*/

// GET /galleries/:id
func (t *TrackController) CreateLocal(w http.ResponseWriter, r *http.Request) {
	//track, err := t.trackByID(w, r)
	//if err != nil {
	//	// The galleryByID method will already render the error
	//	// for us, so we just need to return here.
	//	return
	//}
	//fmt.Println("play ")
	//fmt.Println(track.ID)
	//fmt.Println(track.MusicFile)
	//fmt.Println("play " + track.MusicFile.Filename)
	//fmt.Println("play " + track.MusicFile.MusicPath())
	//// track.MusicFile.Filename = "Tchaikovsky-op19-no3-fueillet-d-album.mp3"
	var vd views.Data
	//vd.Yield = track
	t.CreateLocalView.Render(w, r, vd)
}

func (t *TrackController) CreateLocalComplete(w http.ResponseWriter, r *http.Request) {

	var vd views.Data
	var form TrackForm
	// fmt.Println("CreateLocal  ", form)

	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		fmt.Println("CreateLocal  ", err)
		t.CreateLocalView.Render(w, r, vd)
		return
	}
	user := context.User(r.Context())

	track := models.Track{
		Title:  form.Title,
		UserID: user.ID,
	}


	if err := t.ts.Create(&track); err != nil {
		vd.SetAlert(err)
		fmt.Println("CreateLocal  ", err)
		t.CreateLocalView.Render(w, r, vd)
		return
	}


	//
	// Media files
	//

	//
	// music file
	//
	for _, h := range r.MultipartForm.File["musicfile"] {

		musicfile, err := h.Open()
		if err != nil {
			vd.SetAlert(err)
			t.CreateLocalView.Render(w, r, vd)
			return
		}
		defer musicfile.Close()
		err = t.fs.Create(track.ID, musicfile, h.Filename, models.FileTypeMusic)
		if err != nil {
			vd.SetAlert(err)
			t.CreateLocalView.Render(w, r, vd)
			return
		}
	}

	//
	// cover image
	//
	for _, h := range r.MultipartForm.File["image"] {
		imagefile, err := h.Open()
		if err != nil {
			vd.SetAlert(err)
			t.CreateLocalView.Render(w, r, vd)
			return
		}
		defer imagefile.Close()
		err = t.fs.Create(track.ID, imagefile, h.Filename, models.FileTypeImage)
		if err != nil {
			vd.SetAlert(err)
			t.CreateLocalView.Render(w, r, vd)
			return
		}
	}



	//url, err := g.r.Get(ShowGallery).URL("id",
	//	strconv.Itoa(int(gallery.ID)))
	url, err := t.r.Get(IndexTracks).URL("id",
		strconv.Itoa(int(track.ID)))

	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// GET /galleries/:id/edit
func (t *TrackController) EditDJ(w http.ResponseWriter, r *http.Request) {
	track, err := t.trackByID(w, r)
	if err != nil {
		// The galleryByID method will already render the error
		// for us, so we just need to return here.
		return
	}
	// A user needs logged in to access this page, so we can
	// assume that the RequireUser middleware has run and
	// set the user for us in the request context.
	user := context.User(r.Context())
	if track.UserID != user.ID {
		http.Error(w, "You do not have permission to edit "+
			"this gallery", http.StatusForbidden)
		return
	}
	var vd views.Data
	vd.Yield = track
	t.ChooseDJOptionsView.Render(w, r, vd)
}

// GET /galleries/:id/edit
func (t *TrackController) EditLocal(w http.ResponseWriter, r *http.Request) {
	track, err := t.trackByID(w, r)
	if err != nil {
		// The galleryByID method will already render the error
		// for us, so we just need to return here.
		return
	}
	// A user needs logged in to access this page, so we can
	// assume that the RequireUser middleware has run and
	// set the user for us in the request context.
	user := context.User(r.Context())
	if track.UserID != user.ID {
		http.Error(w, "You do not have permission to edit "+
			"this gallery", http.StatusForbidden)
		return
	}
	var vd views.Data
	vd.Yield = track
	t.EditLocalView.Render(w, r, vd)
	//t.IndexView.Render(w, r, vd)
}

// POST /tracks/:id/music
func (t *TrackController) MusicUpload(w http.ResponseWriter, r *http.Request) {
	track, err := t.trackByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if track.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = track
	err = r.ParseMultipartForm(maxMultipartMem)
	if err != nil {
		// If we can't parse the form just render an error alert on the
		// edit gallery page.
		vd.SetAlert(err)
		t.EditLocalView.Render(w, r, vd)
		return
	}

	// Iterate over uploaded files to process them.
	files := r.MultipartForm.File["musicfile"]
	for _, f := range files {
		fmt.Println(f.Filename)
		// Open the uploaded file
		file, err := f.Open()
		if err != nil {
			vd.SetAlert(err)
			t.EditLocalView.Render(w, r, vd)
			return
		}
		defer file.Close()

		fmt.Println("mfs create")
		// Create a musicfile object which also creates destination file on the specifed location
		// fmt.Println("trackC - music upload - " + string(file.) )
		fmt.Println("trackC - music upload - " + string(track.ID) )
		fmt.Println(track.ID)
		fmt.Println("trackC - music upload - " + f.Filename)

		//err = t.mfs.Create(track.ID, file, f.Filename)
		err = t.fs.Create(track.ID, file, f.Filename, models.FileTypeMusic)
		if err != nil {
			vd.SetAlert(err)
			t.EditLocalView.Render(w, r, vd)
			return
		}
	}

	url, err := t.r.Get(EditTrack).URL("id", fmt.Sprintf("%v", track.ID))
	if err != nil {
		http.Redirect(w, r, "/tracks", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// POST /tracks/:id/images
func (t *TrackController) ImageUpload(w http.ResponseWriter, r *http.Request) {
	track, err := t.trackByID(w, r)
	fmt.Println("ImageUpload  ", track)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	fmt.Println("ImageUpload  ", user)
	if track.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}

	var vd views.Data
	vd.Yield = track
	err = r.ParseMultipartForm(maxMultipartMem)
	if err != nil {
		// If we can't parse the form just render an error alert on the
		// edit gallery page.
		vd.SetAlert(err)
		t.EditLocalView.Render(w, r, vd)
		return
	}

	// Iterate over uploaded files to process them.
	files := r.MultipartForm.File["images"]
	for _, f := range files {
		// Open the uploaded file
		file, err := f.Open()
		if err != nil {
			vd.SetAlert(err)
			t.EditLocalView.Render(w, r, vd)
			return
		}
		defer file.Close()

		fmt.Println("ImageUpload  ", f)
		// Create a image which also creates destination file
		// err = t.is.Create(track.ID, file, f.Filename)
		err = t.fs.Create(track.ID, file, f.Filename, models.FileTypeImage)
		if err != nil {
			vd.SetAlert(err)
			fmt.Println("ImageUpload  ", err)
			t.EditLocalView.Render(w, r, vd)
			return
		}
	}

	//vd.Alert = &views.Alert{
	//	Level:   views.AlertLvlSuccess,
	//	Message: "Images successfully uploaded",
	//}
	//g.EditView.Render(w, r, vd)

	// Remove the code used to create a success alert and
	// render the EditView and replace it with the following.
	url, err := t.r.Get(EditTrack).URL("id", fmt.Sprintf("%v", track.ID))
	fmt.Println("ImageUpload  ", url)
	if err != nil {
		fmt.Println("ImageUpload  ", err)
		http.Redirect(w, r, "/tracks", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// follows add new local track
// POST /galleries/:id/create
func (t *TrackController) CreateLocalSongWithDB(w http.ResponseWriter, r *http.Request) {
	track, err := t.trackByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if track.UserID != user.ID {
		http.Error(w, "Track not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = track
	var form TrackForm
	if err := parseForm(r, &form); err != nil {
		// If there is an error we are going to render the
		// EditView again with an alert message.
		vd.SetAlert(err)
		t.EditLocalView.Render(w, r, vd)
		return
	}
	track.Title = form.Title
	// Persist this gallery change in the DB after
	// we add an Update method to our GalleryService in the
	// models package.

	err = t.ts.Update(track)
	// If there is an error our alert will be an error. Otherwise
	// we will still render an alert, but instead it will be
	// a success message.
	if err != nil {
		vd.SetAlert(err)
	} else {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: "Track successfully updated!",
		}
	}
	// Error or not, we are going to render the EditView with
	// our updated information.
	t.EditLocalView.Render(w, r, vd)
}

// follows edit
// POST /galleries/:id/update
func (t *TrackController) EditLocalSongComplete(w http.ResponseWriter, r *http.Request) {
	track, err := t.trackByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if track.UserID != user.ID {
		http.Error(w, "Track not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = track
	var form TrackForm
	if err := parseForm(r, &form); err != nil {
		// If there is an error we are going to render the
		// EditView again with an alert message.
		vd.SetAlert(err)
		t.EditLocalView.Render(w, r, vd)
		return
	}
	track.Title = form.Title
	// Persist this gallery change in the DB after
	// we add an Update method to our GalleryService in the
	// models package.

	err = t.ts.Update(track)
	// If there is an error our alert will be an error. Otherwise
	// we will still render an alert, but instead it will be
	// a success message.
	if err != nil {
		vd.SetAlert(err)
	} else {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: "Track successfully updated!",
		}
	}
	// Error or not, we are going to render the EditView with
	// our updated information.
	t.EditLocalView.Render(w, r, vd)
}

// POST /galleries/:id/link
func (t *TrackController) ImageViaLink(w http.ResponseWriter, r *http.Request) {
	track, err := t.trackByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if track.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = track
	if err := r.ParseForm(); err != nil {
		// If we can't parse the form just render an error alert on the
		// edit gallery page.
		vd.SetAlert(err)
		t.EditLocalView.Render(w, r, vd)
		return
	}
	files := r.PostForm["files"]

	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, fileURL := range files {
		// do this with go routine
		go func(url string) {
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				log.Println("Failed to download the image form:", url)
				return
			}
			defer resp.Body.Close()
			pieces := strings.Split(url, "/")
			filename := pieces[len(pieces) - 1]
			//if err := t.is.Create(track.ID, resp.Body, filename); err != nil {
			//	panic(err)
			//}
			if err := t.fs.Create(track.ID, resp.Body, filename, models.FileTypeImage); err != nil {
				panic(err)
			}
		} (fileURL) // () meanes execute go routine
	}
	wg.Wait()

	//err = r.ParseMultipartForm(maxMultipartMem)
	//if err != nil {
	//	// If we can't parse the form just render an error alert on the
	//	// edit gallery page.
	//	vd.SetAlert(err)
	//	g.EditView.Render(w, r, vd)
	//	return
	//}
	//
	//// Iterate over uploaded files to process them.
	//files := r.MultipartForm.File["images"]
	//for _, f := range files {
	//	// Open the uploaded file
	//	file, err := f.Open()
	//	if err != nil {
	//		vd.SetAlert(err)
	//		g.EditView.Render(w, r, vd)
	//		return
	//	}
	//	defer file.Close()
	//
	//	// Create a image which also creates destination file
	//	err = g.is.Create(gallery.ID, file, f.Filename)
	//	if err != nil {
	//		vd.SetAlert(err)
	//		g.EditView.Render(w, r, vd)
	//		return
	//	}
	//}

	//Remove the code used to create a success alert and
	//render the EditView and replace it with the following.
	url, err := t.r.Get(EditGallery).URL("id", fmt.Sprintf("%v", track.ID))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/galleries", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// GET /tracks/:id/
func (t *TrackController) trackByID(w http.ResponseWriter, r *http.Request) (*models.Track , error){
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return nil, err
	}
	// create gallery
	track, err := t.ts.ByID(uint(id))
	// if error choose way to respond
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Gallery not found", http.StatusNotFound)
		default:
			http.Error(w, "Whoops! Something went wrong.", http.StatusInternalServerError)
		}
		return nil, err
	}

	// music file
	musicfile, _ := t.fs.ByTrackID(track.ID, models.FileTypeMusic)
	track.MusicFile = musicfile

	// cover method
	image, _ := t.fs.ByTrackID(track.ID, models.FileTypeImage)
	track.CoverImage = image

	// list method
	//imagelist, _ := t.is.ListByTrackID(track.ID)
	//track.Images = imagelist

	// ^^^ pick one

	return track, nil
}

// play the song
// GET /tracks/:id/update
/*
func (t *TrackController) PlayTrack(w http.ResponseWriter, r *http.Request) {

	track, err := t.trackByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if track.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	fmt.Println("music upload")
	var vd views.Data
	vd.Yield = track
	err = r.ParseMultipartForm(maxMultipartMem)
	if err != nil {
		// If we can't parse the form just render an error alert on the
		// edit gallery page.
		vd.SetAlert(err)
		t.PLayView.Render(w, r, vd)
		return
	}



	url, err := t.r.Get(PlayTrack).URL("id",
		strconv.Itoa(int(track.ID)))

	fmt.Println("CreateLocal  ", url)
	fmt.Println("CreateLocal  ", url.Path)

	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}
*/