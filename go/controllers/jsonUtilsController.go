package controllers

import (
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

//
// differnt json methods
//

func Message(status bool, message string) (map[string] interface{}) {
	return map[string]interface{} {"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string] interface{}) {
	w.Header().Add("Content-Type", "application/json")
	//w.Header().Add("Access-Control-Allow-Origin", "*")
	//w.Header().Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	//w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	// addHeaders(w)
	json.NewEncoder(w).Encode(data)
}

func WriteJson(w http.ResponseWriter, data interface{}) {

	var jsonData []byte
	jsonData, err := json.Marshal(data)

	if err != nil {
		log.Println(err)
		fmt.Println(err)
	}

	fmt.Println(string(jsonData))
	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonData)
}

func WriteJsonWithStatus(w http.ResponseWriter, data interface{}, status int) {

	var jsonData []byte
	jsonData, err := json.Marshal(data)

	if err != nil {
		log.Println(err)
		fmt.Println(err)
	}

	fmt.Println(string(jsonData))

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	w.Write(jsonData)
}


func WriteFile(w http.ResponseWriter, filename string) {

	//Check if file exists and open
	Openfile, err := os.Open(filename)
	defer Openfile.Close() //Close after function return
	if err != nil {
		//File not found, send 404
		fmt.Println("File not found!!")
		fmt.Println(err)
		http.Error(w, "File not found.", 404)
		return
	}

	//File is found, create and send the correct headers

	//Get the Content-Type of the file
	//Create a buffer to store the header of the file in
	FileHeader := make([]byte, 512)
	//Copy the headers into the FileHeader buffer
	Openfile.Read(FileHeader)
	//Get content type of file
	FileContentType := http.DetectContentType(FileHeader)

	//Get the file size
	FileStat, _ := Openfile.Stat()                     //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	//Send the headers
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	//Send the file
	//We read 512 bytes from the file already, so we reset the offset back to 0
	Openfile.Seek(0, 0)
	io.Copy(w, Openfile) //'Copy' the file to the client

	return

}

func WriteImage(w http.ResponseWriter, filename string) {

	//Check if file exists and open
	Openfile, err := os.Open(filename)
	//(filename)
	defer Openfile.Close() //Close after function return
	if err != nil {
		//File not found, send 404
		fmt.Println("Image not found!!")
		fmt.Println(err)
		http.Error(w, "Image not found.", 404)
		return
	}

	//File is found, create and send the correct headers

	//Get the Content-Type of the file
	//Create a buffer to store the header of the file in
	FileHeader := make([]byte, 512)
	//Copy the headers into the FileHeader buffer
	Openfile.Read(FileHeader)
	//Get content type of file
	FileContentType := http.DetectContentType(FileHeader)

	//Get the file size
	FileStat, _ := Openfile.Stat()                     //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	//Send the headers
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	//Send the file
	//We read 512 bytes from the file already, so we reset the offset back to 0
	//Openfile.Seek(0, 0)
	//io.Copy(w, Openfile) //'Copy' the file to the client

	// create image file
	img, err := jpeg.Decode(Openfile)
	if err != nil {

		fmt.Println("Image could not be processed!!")
		fmt.Println(err)
		http.Error(w, "Image could not be processed.", 404)
		log.Fatal(err)

	}
	Openfile.Close()

	jpeg.Encode(w, img, nil)

	return

}