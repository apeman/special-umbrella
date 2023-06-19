package main

import(
	"net/http"
	"fmt"
	"time"
	"crypto/rand"
	"io/ioutil"
	"mime"
	"os"
	"log"
    "strconv"
	"path/filepath"
	"strings"
	
    "github.com/julienschmidt/httprouter"
	
)

const maxUploadSize = 2 * 1024 * 1024 // 2 mb
var uploadPath = "./uploads"
//var tmpl = template.Must(template.ParseGlob("templates/*.html"))



type Post struct {
	PostId string
	Title string
	UserName string
	Caption string
	Location string
	Date string
	Tags string
	NSFW string
	Access string
	PicCount int
	Pics []string
	LoggedIn string
}



func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("path", r.URL.Path)
	if r.Method == "GET" {
		m := rdxHgetall("newpost")
			for k, v := range m {
				fmt.Println(">>", k)
				parts := strings.Split(v, ":::")
				fmt.Println(parts[0])
				fmt.Println(parts[1])
			}
			fmt.Println(m)
		if isLoggedIn(r) {
			tmpl.ExecuteTemplate(w, "head.html", nil)
			tmpl.ExecuteTemplate(w, "nav.html", LoginStatus{LoggedIn: "true"})
			tmpl.ExecuteTemplate(w, "index.html", m)
			tmpl.ExecuteTemplate(w, "footer.html", nil)
		} else {
			tmpl.ExecuteTemplate(w, "head.html", nil)
			tmpl.ExecuteTemplate(w, "nonloginhome.html", nil)
			tmpl.ExecuteTemplate(w, "footer.html", nil)
		}
	}
}



func NewPhotoPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("path", r.URL.Path)
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			fmt.Printf("Could not parse multipart form: %v\n", err)
			renderError(w, "CANT_PARSE_FORM", http.StatusInternalServerError)
			return
		}

var access string
var picsStr string
	
	r.ParseForm()
	slice:=[]string{"public","protected","private"}
	for _, v := range slice {
    	if v == r.Form.Get("access") {
        	fmt.Println(v)
			access = v
    	}
	}
		t := time.Now()
        token := t.Format("20060102150405")
		username := "guest"
		date := token[:8]
 		title := r.Form.Get("title")
 		nsfw := r.Form.Get("nsfw")
		caption := r.Form.Get("caption")
		location := r.Form.Get("location")
		tags := r.Form.Get("tags")
		fmt.Println("NSFW : ", nsfw) // empty or 1
		fmt.Println("Caption : ", caption) 
		fmt.Println("Tags : ", tags) 
	
	
	var pics []string
		files := r.MultipartForm.File["imgfile"]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer file.Close()

		// Get and print out file size
		fileSize := fileHeader.Size
		fmt.Printf("File size (bytes): %v\n", fileSize)
		// validate file size
		if fileSize > maxUploadSize {
			renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
			return
		}
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}

		// check file type, detectcontenttype only needs the first 512 bytes
		detectedFileType := http.DetectContentType(fileBytes)
		switch detectedFileType {
		case "image/png", "image/jpg", "image/jpeg":
			break
		default:
			renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
			return
		}
		fileName := randToken(12)
		fileEndings, err := mime.ExtensionsByType(detectedFileType)
		if err != nil {
			renderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
			return
		}
	if fileEndings[0] == ".jfif" {fileEndings[0] = ".jpg"}
		newFileName := fileName + fileEndings[0]

		newPath := filepath.Join(uploadPath, newFileName)
		fmt.Printf("FileType: %s, File: %s\n", detectedFileType, newPath)

		// write file
		newFile, err := os.Create(newPath)
		if err != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		defer newFile.Close() // idempotent, okay to call twice
		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		

		pics = append(pics,"/img/"+newFileName)
		}
		
	fmt.Println(pics)
	picount := len(pics)
	picsStr = strings.Join(pics, ":::")
	rdxHset("newpost", token, title + ":::" + caption + ":::" + location + ":::" + tags + ":::" + nsfw + ":::" + access) 
	rdxHset("photos", token, username + ":::" + date + ":::" + string(picount) + ":::" + picsStr) 
	http.Redirect(w, r, "/view/"+token, http.StatusSeeOther)
//		w.Write([]byte(fmt.Sprintf("SUCCESS - use /files/%v to access the file", newFileName)))

	}


func NewPhotoGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("path", r.URL.Path)
	switch r.Method {
		case "GET" : {
			if isLoggedIn(r) {
		t := time.Now()
        token := t.Format("20060102150405")
		tok := Tok{Token: token}
				tmpl.ExecuteTemplate(w, "head.html", nil)
				tmpl.ExecuteTemplate(w, "nav.html", LoginStatus{LoggedIn: "true"})
				tmpl.ExecuteTemplate(w, "upload.html", tok)
				tmpl.ExecuteTemplate(w, "footer.html", nil)
			} 	else {
				tmpl.ExecuteTemplate(w, "head.html", nil)
				tmpl.ExecuteTemplate(w, "nav.html", LoginStatus{LoggedIn: "false"})
				tmpl.ExecuteTemplate(w, "nonloginhome.html", nil)
				tmpl.ExecuteTemplate(w, "footer.html", nil)
			}
		}
		default : {
			fmt.Println("Method Not allowed")
		}
	}
}

func ViewPost(w http.ResponseWriter, r *http.Request, postid httprouter.Params) {
	fmt.Println("path", r.URL.Path)
	switch r.Method {
		case "GET" : {
			if isLoggedIn(r) {
				fmt.Println("path: ", r.URL.Path)		
				token := postid.ByName("postid")
				postDetails := rdxHget("newpost", token)
				postPhotos := rdxHget("photos", token)
				fmt.Println(postDetails)
				fmt.Println(postPhotos)
				title, caption, location, tags, nsfw, access := readString(postDetails)
				username, date, piclen, picsStr := readPhotoString(postPhotos)
	resPost := Post {	PostId : token, UserName: username,Title: title,Caption: caption, Location: location, Tags: tags, Date: date, NSFW: nsfw, Access: access, Pics: picsStr, PicCount: piclen,LoggedIn: "true"}
				tmpl.ExecuteTemplate(w, "head.html", nil)
				tmpl.ExecuteTemplate(w, "nav.html", LoginStatus{LoggedIn: "true"})
				tmpl.ExecuteTemplate(w, "post.html", resPost)
				tmpl.ExecuteTemplate(w, "footer.html", nil)
			} 	else {
				fmt.Println("path: ", r.URL.Path)		
				token := postid.ByName("postid")
				postDetails := rdxHget("newpost", token)
				postPhotos := rdxHget("photos", token)
				fmt.Println(postDetails)
				fmt.Println(postPhotos)
				title, caption, location, tags, nsfw, access := readString(postDetails)
				username, date, piclen, picsStr := readPhotoString(postPhotos)
	resPost := Post {	PostId : token, UserName: username,Title: title,Caption: caption, Location: location, Tags: tags, Date: date, NSFW: nsfw, Access: access, Pics: picsStr, PicCount: piclen,LoggedIn: "false"}
				tmpl.ExecuteTemplate(w, "head.html", nil)
				tmpl.ExecuteTemplate(w, "nav.html", LoginStatus{LoggedIn: "false"})
				tmpl.ExecuteTemplate(w, "post.html", resPost)
				tmpl.ExecuteTemplate(w, "footer.html", nil)
			}
		}
		default : {
			fmt.Println("Method Not allowed")
		}
	}
}



func EditPost(w http.ResponseWriter, r *http.Request, postid httprouter.Params) {
	fmt.Println("path", r.URL.Path)
	switch r.Method {
		case "GET" : {
			if isLoggedIn(r) {
				fmt.Println("path: ", r.URL.Path)		
				token := postid.ByName("postid")
				postDetails := rdxHget("newpost", token)
				postPhotos := rdxHget("photos", token)
				fmt.Println(postDetails)
				fmt.Println(postPhotos)
				title, caption, location, tags, nsfw, access := readString(postDetails)
				username, date, piclen, picsStr := readPhotoString(postPhotos)
	resPost := Post {	PostId : token,Title: title,Caption: caption, Location: location, Tags: tags, Date: date, UserName: username, NSFW: nsfw, Access: access, Pics: picsStr, PicCount: piclen,LoggedIn: "true"}
				tmpl.ExecuteTemplate(w, "head.html", nil)
				tmpl.ExecuteTemplate(w, "nav.html", LoginStatus{LoggedIn: "true"})
				tmpl.ExecuteTemplate(w, "editpost.html", resPost)
				tmpl.ExecuteTemplate(w, "footer.html", nil)
			} 	else {
				tmpl.ExecuteTemplate(w, "head.html", nil)
				tmpl.ExecuteTemplate(w, "nonloginhome.html", nil)
				tmpl.ExecuteTemplate(w, "footer.html", nil)
			}
		}
		case "POST" : {
	fmt.Println("path", r.URL.Path)
	var access string
	
	r.ParseForm()
	slice:=[]string{"public","protected","private"}
	for _, v := range slice {
    	if v == r.Form.Get("access") {
        	fmt.Println(v)
			access = v
    	}
	}
        token := postid.ByName("postid")
 		title := r.Form.Get("title")
 		nsfw := r.Form.Get("nsfw")
		caption := r.Form.Get("caption")
		location := r.Form.Get("location")
		tags := r.Form.Get("tags")
		fmt.Println("NSFW : ", nsfw) // empty or 1
		fmt.Println("Caption : ", caption) 
		fmt.Println("Tags : ", tags) 
		rdxHset("newpost", token, title + ":::" + caption + ":::" + location + ":::" + tags + ":::" + nsfw + ":::" + access) 
		http.Redirect(w, r, "/view/"+token, http.StatusSeeOther)
	}
		default : {
			fmt.Println("Method Not allowed")
		}
	}
}

func DeletePost(w http.ResponseWriter, r *http.Request, postid httprouter.Params) {
	fmt.Println("path", r.URL.Path)
	if r.Method  == "POST" {
		if isLoggedIn(r) {
			println(rdxHdel("newpost", postid.ByName("postid")))
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} 	else {
			tmpl.ExecuteTemplate(w, "head.html", nil)
			tmpl.ExecuteTemplate(w, "nonloginhome.html", nil)
			tmpl.ExecuteTemplate(w, "footer.html", nil)
		}
	}
}


func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}



type Tok struct {
	Token string
}

type LoginStatus struct {
	LoggedIn string
	LoggedOut string
}

func readString(str string) (string, string, string, string, string, string){
	parts := strings.Split(str, ":::")
	return parts[0], parts[1], parts[2], parts[3], parts[4], parts[5]
}

func readPhotoString(str string) (string, string, int, []string){
	parts := strings.Split(str, ":::")
	//fmt.Println(parts[0])
	piclen, _ := strconv.Atoi(parts[2])
	var picArray []string
	fmt.Println(len(parts))
	for i := 3; i < len(parts); i++ {
		picArray = append(picArray, parts[i])
	}	
	return parts[0], parts[1], piclen, picArray
}


func AllFiles(w http.ResponseWriter, r *http.Request, postid httprouter.Params) {

    entries, err := os.ReadDir("./uploads")
    if err != nil {
        log.Fatal(err)
    }
 var newarr []string
    for _, e := range entries {
    newarr = append(newarr, e.Name())
    }

	tmpl.ExecuteTemplate(w, "res.html", newarr)
}


func DeletePhoto(w http.ResponseWriter, r *http.Request, postid httprouter.Params) {
	os.Remove("./uploads/"+postid.ByName("photoid"))
}
