package assetmanager

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"

	"github.com/DaniZGit/api.stick.it/environment"
	"github.com/DaniZGit/api.stick.it/internal/app"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
)

// Joins path parameters to assets base url that is defined in .env file
func GetPublicAssetsFileUrl(filename string, defaultValue string) string {
	if len(filename) <= 0 {
		return defaultValue;
	}

	url, err := url.JoinPath(environment.AssetsUrl(), filename)

	if err != nil {
		return defaultValue
	}

	return url
}

func GetAssetsFileUrl(paths ...string) string {
	// prepend base assets path
	paths = append([]string{"assets", "public"}, paths...);

	// generate url
	url := filepath.Join(paths...)

	return url
}

func CreateFileWithUUID(f *multipart.FileHeader, ctx *app.ApiContext, folder string, uuidFilename uuid.UUID) (database.File, error) {
	filename := uuidFilename.String()
	
	// append extension to filename if it doesn't have one
	ext := filepath.Ext(f.Filename)
	if len(ext) > 0 {
		filename += ext
	}

	// create folder
	dirErr := os.Mkdir(GetAssetsFileUrl(folder), 0777)
	if dirErr != nil {
		fmt.Printf("error while creating folder '%s' with permission '%s': %s", folder, os.ModeDir, dirErr.Error())
	}

	// upload the file to assets storage
	fileInfo, err := UploadFile(f, GetAssetsFileUrl(folder, filename))
	if err != nil {
		return database.File{}, err
	}
	
	// create file in db
	file, err := ctx.Queries.CreateFile(ctx.Request().Context(), database.CreateFileParams{
		ID: uuidFilename,
		Name: fileInfo.Name(),
		Path: filepath.Join(folder, filename),
	})
	if err != nil {
		return database.File{}, err
	}

	return file, nil
}

func UploadFile(file *multipart.FileHeader, localPath string) (fs.FileInfo, error) {
	src, err := file.Open()
	if err != nil {
		fmt.Println("Error while creating folder", err.Error())
	}

	defer src.Close()

	// Destination
	dst, err := os.Create(localPath)
	if err != nil {
		if serr, ok := err.(*os.PathError); ok {
			fmt.Printf("path: '%s'\nop: '%s'\nerror: '%s'\n", serr.Path, serr.Op, serr.Err.Error())
			return nil, fmt.Errorf("error while creating file '%s': %s", localPath, err.Error())
		}
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return nil, errors.New("could not copy the file contents " + err.Error())
	}

	fileInfo, err := dst.Stat()
	if err != nil {
		return nil, errors.New("could not read file info " + err.Error())
	}

	return fileInfo, nil
}
