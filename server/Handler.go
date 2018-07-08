package server

import (
	"mime/multipart"
	"net/http"
	"fmt"
	"io"
	"strconv"
	"github.com/gorilla/mux"
	"mime"
	"path/filepath"
	"bytes"
	"github.com/nu7hatch/gouuid"
	"path"
)

func (server *Server) FileHandler(w http.ResponseWriter, r *http.Request) {
	var (
		status int
		err    error
	)
	defer func() {
		if nil != err {
			http.Error(w, err.Error(), status)
		}
	}()
	// parse request with maximum memory of _24Kilobits
	const _24K = (1 << 20) * 24
	if err = r.ParseMultipartForm(_24K); nil != err {
		fmt.Println(err)
		status = http.StatusInternalServerError
		return
	}
	for _, fheaders := range r.MultipartForm.File {
		for _, hdr := range fheaders {

			// open uploaded
			var infile multipart.File
			if infile, err = hdr.Open(); nil != err {
				status = http.StatusInternalServerError
				return
			}
			filename := path.Clean(path.Base(hdr.Filename))
			uid, err := uuid.NewV5(uuid.NamespaceURL, []byte(filename))

			err = server.storage.Put(uid.String() + "/" + filename,infile,100)

			returnStr := "http://" + r.Host + "/" + uid.String() + "/" + filename + "\n"

			w.Write([]byte(returnStr))
			if err != nil {

			}
		}
	}
}


func (server *Server) Download(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	token := vars["token"]
	filename := vars["filename"]

	contentType := mime.TypeByExtension(filepath.Ext(filename))

	reader, err := server.storage.Get(token + "/" + filename)

	if err != nil {

	}

	// TODO: Write this to meta data maybe? kills me now if I have to do this for every request
	buf := &bytes.Buffer{}
	contentLength, err := io.Copy(buf, reader)

	if err != nil {
		fmt.Print(err)
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Connection", "close")
	w.Write(buf.Bytes())
}

func (server *Server) HomePage(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("<h1> Home </h1>"))

	return
}

func (server *Server) DownloadMultiStage(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	token := vars["token"]
	filename := vars["filename"]

	contentType := mime.TypeByExtension(filepath.Ext(filename))

	reader, err := server.storage.Get(token + "/" + filename)

	if err != nil {

	}

	// TODO: Write this to meta data maybe? kills me now if I have to do this for every request
	buf := &bytes.Buffer{}
	contentLength, err := io.Copy(buf, reader)

	if err != nil {
		fmt.Print(err)
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Connection", "close")
	w.Write(buf.Bytes())
}
