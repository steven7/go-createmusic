package models

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
)

type FileType int

const (
	FileTypeImage    FileType = 0
	FileTypeMusic    FileType = 1
)

// experimental idea
// wil do though

type FileService interface {
	Create(trackID uint, r io.Reader, filename string, ft FileType) error
	ByTrackID(trackID uint, ft FileType) (File, error)
	//ListByTrackID(trackID uint) ([]File, error)
	Delete(mf *File, ft FileType) error
}

func NewFileService() FileService {
	return &fileService{}
}

type fileService struct {}


func (fs *fileService) Create (trackID uint, r io.Reader, filename string, ft FileType) error {
	fmt.Println(" create music file")

	//var path string
	//var err error
	//
	//if ft == FileTypeMusic {
	//
	//} else if ft == FileTypeImage {
	//	path, err = mfs.mkImageDir(trackID)
	//}

	path, err := fs.mkDir(trackID, ft)
	if err != nil {
		return err
	}
	fmt.Println(path)
	// Clear directory before creating new one. We only need one file at a time.
	fs.ClearFiles(path)
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

func (fs *fileService) ClearFiles (path string) {
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

func (fs *fileService) ByTrackID(trackID uint, ft FileType) (File, error) {
	path := fs.getDir(trackID, ft)
	// fmt.Println("musicfiles.go -  " + path)
	strings, err := filepath.Glob(filepath.Join(path, "*"))
	for _, s := range strings {
		fmt.Println(" files.go -  " + s )
	}
	//fmt.Println("images.go -  " + strings)
	if err != nil {
		return File{}, err
	}

	var filename string
	if len(strings) > 0 {
		filename = filepath.Base(strings[0])
	} else {
		filename = ""
	}

	ret := File{
		Filename: filename,
		TrackID:  trackID,
	}

	return ret, nil
}

//func (fs *fileService) ListByTrackID(trackID uint) ([]File, error) {
//	panic("implement me")
//}

func (fs *fileService) Delete(f *File, ft FileType) error {
	return os.Remove(f.RelativePath(ft))
}


func (fs *fileService) getDir(trackID uint, ft FileType) string {
	var dir string
	if ft == FileTypeMusic {
		dir = filepath.Join("userfiles", "tracks", fmt.Sprintf("%v", trackID), "music")
	} else if ft == FileTypeImage {
		dir = filepath.Join("userfiles", "tracks", fmt.Sprintf("%v", trackID), "cover")
	}
	return dir
}

func (fs *fileService) mkDir(galleryID uint, ft FileType) (string, error) {
	// filepath.Join will return a path like:
	//   images/galleries/123
	// We use filepath.Join instead of building the path
	// manually because the slashes and other characters
	// could vary between operating systems.
	var galleryPath string
	if ft == FileTypeMusic {
		galleryPath = filepath.Join("userfiles", "tracks",
			fmt.Sprintf("%v", galleryID), "music")
	} else if ft == FileTypeImage {
		galleryPath = filepath.Join("userfiles", "tracks",
			fmt.Sprintf("%v", galleryID), "cover")
	}
	// Create our directory (and any necessary parent dirs)
	// using 0755 permissions.
	err := os.MkdirAll(galleryPath, 0755)
	if err != nil {
		return "", err
	}
	return galleryPath, nil
}

// File is used to represent images stored in a Gallery.
// File is NOT stored in the database, and instead
// references data stored on disk.
type File struct {
	TrackID   uint   `json:"trackId"`
	Filename  string `json:"filename"`
}

// Path is used to build the absolute path used to reference this image
// via a web request.
//func (f *File) Path(ft FileType) string {
//	temp := url.URL{
//		Path: "/" + f.RelativePath(ft),
//	}
//	return temp.String()
//}

func (f *File) MusicPath() string {
	temp := url.URL{
		Path: f.RelativePath(FileTypeMusic),
	}
	return temp.String()
}

func (f *File) ImagePath() string {
	temp := url.URL{
		Path: f.RelativePath(FileTypeImage),
	}
	return temp.String()
}

// RelativePath is used to build the path to this image on our local
// disk, relative to where our Go application is run from.
func (f *File)  RelativePath(ft FileType) string {
	// Convert the gallery ID to a string
	trackID := fmt.Sprintf("%v", f.TrackID)
	var path string
	if ft == FileTypeMusic {
		path = filepath.ToSlash(filepath.Join("userfiles", "tracks", trackID, "music", f.Filename))
	} else if ft == FileTypeImage {
		path = filepath.ToSlash(filepath.Join("userfiles", "tracks", trackID, "cover", f.Filename))
	}
	return path
}