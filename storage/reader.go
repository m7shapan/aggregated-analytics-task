package storage

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"io"
	"os"
)

const (
	EventsFile  CSVFile = "data/events.csv"
	ReposFile   CSVFile = "data/repos.csv"
	ActorsFile  CSVFile = "data/actors.csv"
	CommitsFile CSVFile = "data/commits.csv"
)

type CSVFile string

type Reader interface {
	Read() ([]string, error)
}

type fileReader struct {
	filePath string
	csvFile  CSVFile
	reader   *csv.Reader
}

func NewFileReader(filePath string, loadedFile CSVFile) (*fileReader, error) {
	f := fileReader{
		filePath: filePath,
		csvFile:  loadedFile,
	}

	return &f, f.open()
}

func (fr *fileReader) open() error {
	f, err := os.Open(fr.filePath)
	if err != nil {
		return err
	}

	defer f.Close()

	gzf, err := gzip.NewReader(f)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(gzf)

	for {
		header, err := tarReader.Next()

		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if header.Name != string(fr.csvFile) {
			continue
		}

		buf := bytes.NewBuffer(make([]byte, 0, 16))
		_, err = io.Copy(buf, tarReader)
		fr.reader = csv.NewReader(buf)
	}

	return nil

}

func (fr *fileReader) Read() (record []string, err error) {
	return fr.reader.Read()
}
