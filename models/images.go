package models

// We will need some of these imports later
import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

type ImageService interface {
	Create(trackID uint, r io.Reader, filename string) error
	ByTrackID(galleryID uint) ([]Image, error)
	Delete(i *Image) error
	CoverByTrackID(trackID uint) (Image, error)
	ListByTrackID(trackID uint) ([]Image, error)
}

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct {}


func (is *imageService) Create (trackID uint, r io.Reader, filename string) error {
	//fmt.Println("Create new img")
	path, err := is.mkImageDir(trackID)
	if err != nil {
		return err
	}
	// Clear directory before creating new one. We only need one file at a time.
	is.ClearFiles(path)
	// Create a destination file
	dst, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		return err
	}
	defer dst.Close()
	// Copy reader data to the destination file
	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}
	return nil
}

func (is *imageService) ClearFiles (path string) {
	//fmt.Println(path)
	files, err := filepath.Glob(filepath.Join(path, "*"))
	if err == nil {
		fmt.Println(files)
		for _, imgStr := range files {
			//fmt.Println(imgStr)
			os.Remove(imgStr)
		}
	} else {
		//fmt.Println("lol file error")
		fmt.Println(err)
	}
}

func (is *imageService) ByTrackID(galleryID uint) ([]Image, error) {
	path := is.imageDir(galleryID)
	strings, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return nil, err
	}
	// Setup the Image slice we are returning
	ret := make([]Image, len(strings))
	for i, imgStr := range strings {
		ret[i] = Image{
			Filename: filepath.Base(imgStr),
			TrackID:  galleryID,
		}
	}
	return ret, nil
}

func (is *imageService) CoverByTrackID(trackID uint) (Image, error) {
	path := is.imageDir(trackID)
	fmt.Println("images.go -  " + path)
	strings, err := filepath.Glob(filepath.Join(path, "*"))
	for _, s := range strings {
		fmt.Println("images.go -  " + s )
	}
	//fmt.Println("images.go -  " + strings)
	if err != nil {
		return Image{}, err
	}

	var filename string
	if len(strings) > 0 {
		filename = filepath.Base(strings[0])
	} else {
		filename = ""
	}

	ret := Image{
		Filename: filename,
		TrackID:  trackID,
	}
	return ret, nil
}

func (is *imageService) ListByTrackID(trackID uint) ([] Image, error) {
	path := is.imageDir(trackID)
	strings, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return []Image{}, err
	}
	// Setup the Image slice we are returning
	ret := make([]Image, len(strings))
	for i, imgStr := range strings {
		ret[i] = Image{
			Filename:  filepath.Base(imgStr),
			TrackID: trackID,
		}
	}

	return ret, nil
}

// need this to know when a path is already made
// aka imagePath
func (is *imageService) imageDir(galleryID uint) string {
	return filepath.Join("userfiles", "tracks", fmt.Sprintf("%v", galleryID), "cover")
}

func (is *imageService) mkImageDir(galleryID uint) (string, error) {
	// filepath.Join will return a path like:
	//   images/galleries/123
	// We use filepath.Join instead of building the path
	// manually because the slashes and other characters
	// could vary between operating systems.
	galleryPath := filepath.Join("userfiles", "tracks",
		fmt.Sprintf("%v", galleryID), "cover")
	// Create our directory (and any necessary parent dirs)
	// using 0755 permissions.
	err := os.MkdirAll(galleryPath, 0755)
	if err != nil {
		return "", err
	}
	return galleryPath, nil
}


func (is *imageService) Delete(i *Image) error {
	return os.Remove(i.RelativePath())
}

// Image is used to represent images stored in a Gallery.
// Image is NOT stored in the database, and instead
// references data stored on disk.
type Image struct {
	TrackID   uint   `json:"trackId"`
	Filename  string `json:"filename"`
	// ImageURL  string 	  `gorm:"-"; 	  		  json:"imageURL"`  /// I think this is correct with client side render
}

// Path is used to build the absolute path used to reference this image
// via a web request.
func (i *Image) Path() string {
	temp := url.URL{
		Path: "/" + i.RelativePath(),
	}
	return temp.String()
}

// RelativePath is used to build the path to this image on our local
// disk, relative to where our Go application is run from.

func (i *Image) RelativePath() string {
	// Convert the gallery ID to a string
	trackID := fmt.Sprintf("%v", i.TrackID)
	return filepath.ToSlash(filepath.Join("userfiles", "tracks", trackID, "cover", i.Filename))
}
