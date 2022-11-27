package httpserver

import (
	"fmt"
	"io"
	"net/http"
	"pasto/checksum"
	. "pasto/logger"
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
	Logger.Info().Msgf("Received forms: %v\n", r.PostForm)
	receivedPageHash := template.HTMLEscapeString(r.FormValue("page_hash"))
	if len(receivedPageHash) < 1 {
		fmt.Fprintf(w, "Error, empty message\r\n")
		return nil
	}

	fmt.Fprintf(w, "%s\n", "")
	return nil
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	Logger.Info().Msgf("Incoming %s connection from %s with size %d bytes",
		r.Method, r.RemoteAddr, r.ContentLength)
	Logger.Debug().Msgf("Incoming %s [%v] connection from %s with size %d bytes",
		r.Method, r.Header, r.RemoteAddr, r.ContentLength)

	switch r.Method {
	case "GET":
		if err := LoadPage(w, r); err != nil {
			Logger.Error().Err(err).Msg("Loading page failed")
			return
		}
	case "POST":
		r.Body = http.MaxBytesReader(w, r.Body, maxBodySize)
		defer r.Body.Close()

		if err := r.ParseMultipartForm(maxBodySize); err != nil {
			Logger.Error().Err(err).Msgf("Parsing form failed")
			userErrMsg := "Couldn't parse form or mandatory value(s) missing"
			fmt.Fprint(w, userErrMsg+"\r\n")
			return
		}

		formFile, formFileHeaders, err := r.FormFile("file")
		if err != nil {
			Logger.Error().Err(err).Msg("Missing file")
			fmt.Fprint(w, "Missing 'file' parameter\r\n")
			return
		}

		binFile, err := io.ReadAll(formFile)
		if err != nil {
			Logger.Error().Err(err).
				Str("filename", formFileHeaders.Filename).
				Msgf("Reading file failed")
			fmt.Fprint(w, "Error while reading file binary\r\n")
			return
		}
		postHash := checksum.Sha256Sum(binFile)
		Logger.Info().
			Str("filename", formFileHeaders.Filename).
			Str("file_hash", postHash).
			Msg("File hashed")
	}
}
