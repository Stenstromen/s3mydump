package mygzip

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
)

func GzipFile(filename string) error {
	log.Printf("Gzipping file %s", filename)

	source, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening source file: %v", err)
	}
	defer source.Close()

	target, err := os.Create(filename + ".gz")
	if err != nil {
		return fmt.Errorf("error creating gzip file: %v", err)
	}
	defer target.Close()

	gw, err := gzip.NewWriterLevel(target, gzip.BestCompression)
	if err != nil {
		return fmt.Errorf("error creating gzip writer: %v", err)
	}
	defer gw.Close()

	if _, err := io.Copy(gw, source); err != nil {
		return fmt.Errorf("error writing to gzip file: %v", err)
	}

	if err := os.Remove(filename); err != nil {
		return fmt.Errorf("error removing original file: %v", err)
	}

	return nil
}
