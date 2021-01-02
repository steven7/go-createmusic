package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"

	"github.com/steven7/go-createmusic/context"
	"github.com/steven7/go-createmusic/models"
	"github.com/steven7/go-createmusic/views"
)

const (
	IndexGalleries  = "index_galleries"
	ShowGallery     = "show_gallery"
	EditGallery     = "edit_gallery"

	maxMultipartMem = 1 << 20 // 1 megabyte
)

func NewGalleries(gs models.GalleryService, is models.ImageService, r *mux.Router) *Galleries {
	return &Galleries{
		New:       views.NewView("bootstrap", "galleries/new"),
		ShowView:  views.NewView("bootstrap", "galleries/show"),
		EditView:  views.NewView("bootstrap", "galleries/edit"),
		IndexView: views.NewView("bootstrap", "galleries/index"),
		gs:        gs,
		is:		   is,
		r:         r,
	}
}

type Galleries struct {
	New       *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	gs        models.GalleryService
	is        models.ImageService
	r         *mux.Router
}

type GalleryForm struct {
	Title string `schema:"title"`
}

// GET /galleries
func (g *Galleries) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	galleries, err := g.gs.ByUserID(user.ID)
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
	vd.Yield = galleries
	g.IndexView.Render(w, r, vd)
}

// GET /galleries/:id
func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		// The galleryByID method will already render the error
		// for us, so we just need to return here.
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(w, r, vd)
}

// GET /galleries/:id/edit
func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		// The galleryByID method will already render the error
		// for us, so we just need to return here.
		return
	}
	// A user needs logged in to access this page, so we can
	// assume that the RequireUser middleware has run and
	// set the user for us in the request context.
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You do not have permission to edit "+
			"this gallery", http.StatusForbidden)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.EditView.Render(w, r, vd)
}

// follows edit
// POST /galleries/:id/update
func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		// If there is an error we are going to render the
		// EditView again with an alert message.
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}
	gallery.Title = form.Title
	// Persist this gallery change in the DB after
	// we add an Update method to our GalleryService in the
	// models package.

	err = g.gs.Update(gallery)
	// If there is an error our alert will be an error. Otherwise
	// we will still render an alert, but instead it will be
	// a success message.
	if err != nil {
		vd.SetAlert(err)
	} else {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: "Gallery successfully updated!",
		}
	}
	// Error or not, we are going to render the EditView with
	// our updated information.
	g.EditView.Render(w, r, vd)
}

// POST /galleries/:id/images
func (g *Galleries) ImageUpload(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}

	var vd views.Data
	vd.Yield = gallery
	err = r.ParseMultipartForm(maxMultipartMem)
	if err != nil {
		// If we can't parse the form just render an error alert on the
		// edit gallery page.
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}

	// Iterate over uploaded files to process them.
	files := r.MultipartForm.File["images"]
	for _, f := range files {
		// Open the uploaded file
		file, err := f.Open()
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(w, r, vd)
			return
		}
		defer file.Close()

		// Create a image which also creates destination file
		err = g.is.Create(gallery.ID, file, f.Filename)
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(w, r, vd)
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
	url, err := g.r.Get(EditGallery).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		http.Redirect(w, r, "/galleries", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// POST /galleries/:id/link
func (g *Galleries) ImageViaLink(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	if err := r.ParseForm(); err != nil {
		// If we can't parse the form just render an error alert on the
		// edit gallery page.
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
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
			if err := g.is.Create(gallery.ID, resp.Body, filename); err != nil {
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
	url, err := g.r.Get(EditGallery).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/galleries", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// POST /galleries/:id/images/:filename/delete
func (g *Galleries) ImageDelete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You do not have permission to edit "+
			"this gallery or image", http.StatusForbidden)
		return
	}

	// Get the filename from the path
	filename := mux.Vars(r)["filename"]
	// Build the Image model
	i := models.Image{
		Filename:  filename,
		TrackID: gallery.ID,
	}
	// Try to delete the image.
	err = g.is.Delete(&i)
	if err != nil {
		// Render the edit page with any errors.
		var vd views.Data
		vd.Yield = gallery
		vd.SetAlert(err)
		g.EditView.Render(w, r, vd)
		return
	}
	// If all goes well, redirect to the edit gallery page.
	url, err := g.r.Get(EditGallery).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		http.Redirect(w, r, "/galleries", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}



// POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, r, vd)
		return
	}
	user := context.User(r.Context())
	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, r, vd)
		return
	}

	//url, err := g.r.Get(ShowGallery).URL("id",
	//	strconv.Itoa(int(gallery.ID)))
	url, err := g.r.Get(EditGallery).URL("id",
		strconv.Itoa(int(gallery.ID)))

	fmt.Println(url)
	fmt.Println(url.Path)

	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// POST /galleries/:id/delete
func (g *Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	// Lookup the gallery using the galleryByID we wrote earlier
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		// If there is an error the galleryByID will have rendered
		// it for us already.
		return
	}
	// We also need to retrieve the user and verify they have
	// permission to delete this gallery. This means we will
	// need to use the RequireUser middleware on any routes
	// mapped to this method.
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You do not have permission to edit "+
			"this gallery", http.StatusForbidden)
		return
	}

	var vd views.Data
	err = g.gs.Delete(gallery.ID)
	if err != nil {
		// If there is an error we want to set an alert and
		// render the edit page with the error. We also need
		// to set the Yield to gallery so that the EditView
		// is rendered correctly.
		vd.SetAlert(err)
		vd.Yield = gallery
		g.EditView.Render(w, r,vd)
		return
	}
	// We will want to redirect to the index
	// page that lists all galleries this user owns
	url, err := g.r.Get(IndexGalleries).URL()
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// GET /galleries/:id/
func (g *Galleries) galleryByID(w http.ResponseWriter, r *http.Request) (*models.Gallery , error){
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return nil, err
	}
	// create gallery
	gallery, err := g.gs.ByID(uint(id))
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
	//
	images, _ := g.is.ByTrackID(gallery.ID)
	gallery.Images = images
	return gallery, nil
}