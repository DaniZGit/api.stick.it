package assetmanager

import (
	"errors"
	"io"
	"io/fs"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"

	"github.com/DaniZGit/api.stick.it/environment"
	"github.com/DaniZGit/api.stick.it/internal/app"
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Joins path parameters to assets base url that is defined in .env file
func GetPublicAssetsFileUrl(filename string, returnEmpty bool) string {
	if returnEmpty && len(filename) <= 0 {
		return "";
	}

	url, err := url.JoinPath(environment.AssetsUrl(), filename)

	if err != nil {
		return ""
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

func CreateFile(f *multipart.FileHeader, ctx *app.ApiContext, folder, filename string) (database.File, error) {
	// append extension to filename if it doesn't have one
	ext := filepath.Ext(f.Filename)
	if len(ext) > 0 {
		filename += ext
	}

	// upload the file to assets storage
	fileInfo, err := UploadFile(f, GetAssetsFileUrl(folder, filename))
	if err != nil {
		return database.File{}, err
	}
	
	// create file in db
	file, err := ctx.Queries.CreateFile(ctx.Request().Context(), database.CreateFileParams{
		ID: pgtype.UUID{Bytes: uuid.New(), Valid: true},
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
		return nil, errors.New("could not read the file")
	}

	defer src.Close()

	// Destination
	dst, err := os.Create(localPath)
	if err != nil {
		return nil, errors.New("could not create the file")
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return nil, errors.New("could not copy the file contents")
	}

	fileInfo, err := dst.Stat()
	if err != nil {
		return nil, errors.New("could not read file info")
	}

	return fileInfo, nil
}
