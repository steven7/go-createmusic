package models


import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

type MusicFileService interface {
	Create(trackID uint, r io.Reader, filename string) error
	ByTrackID(trackID uint) (MusicFile, error)
	ListByTrackID(trackID uint) ([]MusicFile, error)
	Delete(mf *MusicFile) error
}

func NewMusicFileService() MusicFileService {
	return &musicFileService{}
}

type musicFileService struct {}


func (mfs *musicFileService) Create (trackID uint, r io.Reader, filename string) error {
	fmt.Println(" create music file")
	path, err := mfs.mkMusicDir(trackID)
	if err != nil {
		return err
	}
	fmt.Println(path)
	// Clear directory before creating new one. We only need one file at a time.
	mfs.ClearFiles(path)
	// Create a destination file
	dst, err := os.Create(filepath.Join(path, filename))
	fmt.Println(dst)
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

func (mfs *musicFileService) ClearFiles (path string) {
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

func (mfs *musicFileService) musicDir(trackID uint) string {
	return filepath.Join("userfiles", "tracks", fmt.Sprintf("%v", trackID), "music")
}

func (mfs *musicFileService) mkMusicDir(trackID uint) (string, error) {
	// filepath.Join will return a path like:
	//   images/galleries/123
	// We use filepath.Join instead of building the path
	// manually because the slashes and other characters
	// could vary between operating systems.
	trackPath := filepath.Join("userfiles", "tracks",
		fmt.Sprintf("%v", trackID), "music")
	// Create our directory (and any necessary parent dirs)
	// using 0755 permissions.
	err := os.MkdirAll(trackPath, 0755)
	if err != nil {
		return "", err
	}
	return trackPath, nil
}

//// need this to know when a path is already made
//// aka imagePath
//func (mfs *musicFileService) imageDir(galleryID uint) string {
//	return filepath.Join("images", "galleries", fmt.Sprintf("%v", galleryID))
//}
//
//// make dir for image with id
//// aka mkImagePath
//func (mfs *musicFileService) mkImageDir(galleryID uint) (string, error) {
//	// filepath.Join will return a path like:
//	//   images/galleries/123
//	// We use filepath.Join instead of building the path
//	// manually because the slashes and other characters
//	// could vary between operating systems.
//	galleryPath := filepath.Join("images", "galleries",
//		fmt.Sprintf("%v", galleryID))
//	// Create our directory (and any necessary parent dirs)
//	// using 0755 permissions.
//	err := os.MkdirAll(galleryPath, 0755)
//	if err != nil {
//		return "", err
//	}
//	return galleryPath, nil
//}

func (mfs *musicFileService) Delete(i *MusicFile) error {
	panic("implement me")
}

func (mfs *musicFileService) ByTrackID(trackID uint) (MusicFile, error) {
	path := mfs.musicDir(trackID)
	// fmt.Println("musicfiles.go -  " + path)
	strings, err := filepath.Glob(filepath.Join(path, "*"))
	for _, s := range strings {
		fmt.Println("musicfiles.go -  " + s )
	}
	//fmt.Println("images.go -  " + strings)
	if err != nil {
		return MusicFile{}, err
	}

	var filename string
	if len(strings) > 0 {
		filename = filepath.Base(strings[0])
	} else {
		filename = ""
	}

	ret := MusicFile{
		Filename: filename,
		TrackID:  trackID,
	}

	return ret, nil
}



func (mfs *musicFileService) ListByTrackID(trackID uint) ([]MusicFile, error) {
	panic("implement me")
}


// Image is used to represent images stored in a Gallery.
// Image is NOT stored in the database, and instead
// references data stored on disk.
type MusicFile struct {
	TrackID   uint
	Filename  string
}

// Path is used to build the absolute path used to reference this image
// via a web request.
func (mf *MusicFile) Path() string {
	temp := url.URL{
		Path: "/" + mf.RelativePath(),
	}
	return temp.String()
}

// RelativePath is used to build the path to this image on our local
// disk, relative to where our Go application is run from.
func (mf *MusicFile)  RelativePath() string {
	// Convert the gallery ID to a string
	trackID := fmt.Sprintf("%v", mf.TrackID)
	return filepath.ToSlash(filepath.Join("userfiles", "tracks", trackID, "music", mf.Filename))
}
