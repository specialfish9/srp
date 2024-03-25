package main

import (
	"fmt"
	"os"
)

type RequestSaving struct {
	Method string
	Sender string
	Url    string
	Path   string
}

type Saver interface {
	Save(*RequestSaving) error
}

type FileSaver struct {
	fileName string
	file     *os.File
}

var _ Saver = (*FileSaver)(nil)

func NewFileSaver(fileName string) *FileSaver {
	return &FileSaver{
		fileName: fileName,
	}
}

func (fs *FileSaver) Init() error {
	file, err := os.OpenFile(fs.fileName, os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		return err
	}

	fs.file = file

	return nil
}

func (fs *FileSaver) Close() error {
	if fs.file != nil {
		return fs.file.Close()
	}

	return nil
}

func (fs *FileSaver) Save(req *RequestSaving) error {
	content := fmt.Sprintf("%s;%s;%s;%s\n", req.Url, req.Path, req.Method, req.Sender)
	_, err := fs.file.WriteString(content)

	return err
}
