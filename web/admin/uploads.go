package admin

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/madappgang/identifo/model"
)

// UploadJWTKeys is for uploading public and private keys used for signing JWTs.
func (ar *Router) UploadJWTKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1024 * 1024 * 1); err != nil {
			ar.Error(w, err, http.StatusBadRequest, fmt.Sprintf("Error parsing a request body as multipart/form-data: %s", err.Error()))
			return
		}

		formKeys := r.MultipartForm.File["keys"]

		keys := &model.JWTKeys{}

		for _, fileHeader := range formKeys {
			f, err := fileHeader.Open()
			if err != nil {
				ar.Error(w, err, http.StatusBadRequest, fmt.Sprintf("Error uploading key: %s", err.Error()))
				return
			}
			defer f.Close()

			switch fileHeader.Filename {
			case "private.pem":
				keys.Private = f
			case "public.pem":
				keys.Public = f
			default:
				ar.Error(w, fmt.Errorf("Invalid key field name '%s'", fileHeader.Filename), http.StatusBadRequest, "")
				return
			}
		}

		if err := ar.configurationStorage.InsertKeys(keys); err != nil {
			ar.Error(w, err, http.StatusInternalServerError, "")
			return
		}
		ar.ServeJSON(w, http.StatusOK, nil)
	}
}

// UploadADDAFile is for uploading Apple Developer Domain Association File.
func (ar *Router) UploadADDAFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1024 * 1024 * 1); err != nil {
			ar.Error(w, err, http.StatusBadRequest, fmt.Sprintf("Error parsing a request body as multipart/form-data: %s", err.Error()))
			return
		}

		formFile, _, err := r.FormFile("file")
		if err != nil {
			ar.Error(w, err, http.StatusBadRequest, fmt.Sprintf("Cannot read file: %s", err.Error()))
			return
		}
		defer formFile.Close()

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, formFile); err != nil {
			ar.Error(w, err, http.StatusInternalServerError, fmt.Sprintf("Cannot read file as bytes: %s", err.Error()))
			return
		}

		if err = ar.staticFilesStorage.UploadAppleFile(model.AppleFilenames.DeveloperDomainAssociation, buf.Bytes()); err != nil {
			ar.Error(w, err, http.StatusInternalServerError, fmt.Sprintf("Cannot upload file: %s", err.Error()))
			return
		}
		ar.ServeJSON(w, http.StatusOK, nil)
	}
}

// UploadAASAFile is for uploading Apple App Site Association File.
// It is being uploaded as a string, not a file, because it may require manual editing.
func (ar *Router) UploadAASAFile() http.HandlerFunc {
	type file struct {
		Contents string `json:"contents"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		assaFile := new(file)
		if ar.mustParseJSON(w, r, assaFile) != nil {
			return
		}
		if err := ar.staticFilesStorage.UploadAppleFile(model.AppleFilenames.AppSiteAssociation, []byte(assaFile.Contents)); err != nil {
			ar.Error(w, err, http.StatusInternalServerError, fmt.Sprintf("Cannot upload file: %s", err.Error()))
			return
		}
		ar.ServeJSON(w, http.StatusOK, nil)
	}
}
