package httpserver

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"pasto/checksum"
	"pasto/logger"
	"text/template"
)

const (
	maxBodySize int64 = 32<<20 + 512
)

func LoadPage(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error parsing form\r\n")
		return nil
	}
	logger.Infof("Received forms: %v\n", r.PostForm)
	receivedPageHash := template.HTMLEscapeString(r.FormValue("page_hash"))
	if len(receivedPageHash) < 1 {
		fmt.Fprintf(w, "Error, empty message\r\n")
		return nil
	}

	fmt.Fprintf(w, "%s\n", "")
	return nil
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	logger.Infof("Incoming %s connection from %s with size %d bytes",
		r.Method, r.RemoteAddr, r.ContentLength)
	logger.Debugf("Incoming %s [%v] connection from %s with size %d bytes",
		r.Method, r.Header, r.RemoteAddr, r.ContentLength)

	switch r.Method {
	case "GET":
		if err := LoadPage(w, r); err != nil {
			logger.Errorf("%s", err)
			return
		}
	case "POST":
		r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
		defer r.Body.Close()

		if err := r.ParseMultipartForm(maxBodySize); err != nil {
			logger.Errorf("parsing form failed: %s", err)
			userErrMsg := "Couldn't parse form or mandatory value(s) missing"
			fmt.Fprint(w, userErrMsg+"\r\n")
			return
		}

		formFile, formFileHeaders, err := r.FormFile("file")
		if err != nil {
			logger.Error("missing file")
			fmt.Fprint(w, "Missing 'file' parameter\r\n")
			return
		}

		binFile, err := ioutil.ReadAll(formFile)
		if err != nil {
			logger.Errorf("reading file %s failed: %v",
				formFileHeaders.Filename, err)
			fmt.Fprint(w, "Error while reading file binary\r\n")
			return
		}
		postHash := checksum.Sha256Sum(binFile)
		logger.Infof("Filename %s has hash %s",
			formFileHeaders.Filename, postHash)

	}
}
