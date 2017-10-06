package fileupload

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mobingilabs/mobingi-sdk-go/pkg/debug"
	"github.com/pkg/errors"
)

func ProcessFileUpload(r *http.Request) error {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		return errors.Wrap(err, "formfile failed")
	}

	defer file.Close()
	path := r.FormValue("uploadpath")
	if path == "" {
		path = os.TempDir()
	}

	_, fstr := filepath.Split(handler.Filename)
	debug.Info("path:", path)
	debug.Info("file:", handler.Filename)
	f, err := os.Create(path + "/" + fstr)
	if err != nil {
		return errors.Wrap(err, "create file failed")
	}

	defer f.Close()
	_, err = io.Copy(f, file)
	return err
}
