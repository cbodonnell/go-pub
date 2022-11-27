package media

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type Media struct {
	File     multipart.File
	MimeType string
	Name     string
	UUID     string
	FileExt  string
}

func ParseMedia(r *http.Request, name string) (Media, error) {
	file, header, err := r.FormFile("file")
	if err != nil {
		return Media{}, err
	}
	defer file.Close()

	if header.Size > 15*1024*1024 {
		return Media{}, fmt.Errorf("file too large: %d", header.Size)
	}

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		return Media{}, err
	}

	filetype := http.DetectContentType(buff)
	if filetype != "audio/mpeg" {
		return Media{}, fmt.Errorf("invalid file type: %s", filetype)
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return Media{}, err
	}

	m := Media{
		File:     file,
		MimeType: filetype,
		Name:     header.Filename,
		UUID:     uuid.New().String(),
		FileExt:  filepath.Ext(header.Filename),
	}

	return m, nil
}

func Delete(path string) error {
	log.Printf("deleting: %s\n", path)
	return os.Remove(path)
}

func (m *Media) Save(dir string) error {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf(dir + m.UUID + m.FileExt))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, m.File)
	if err != nil {
		return err
	}

	return nil
}
