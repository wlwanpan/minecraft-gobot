package mcs

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type backerState int

const (
	BACKER_STATE_ZIPPING backerState = iota
	BACKER_STATE_UPLOADING
	BACKER_STATE_DONE
	BACKER_STATE_FAILED
)

type backer struct {
	id      string
	state   backerState
	lastUrl string
}

func newBacker() *backer {
	return &backer{
		id: strconv.Itoa(int(time.Now().Unix())),
	}
}

func (b *backer) filename() string {
	return fmt.Sprintf("backup_%s.zip", b.id)
}

func (b *backer) start() error {
	// 1. zip the world directory
	if err := b.zipworld(); err != nil {
		log.Printf("Error zipping /world: error='%s'", err)
		return err
	}
	log.Println("Successfully zipped!")

	return nil
}

func (b *backer) zipworld() error {
	wd, _ := os.Getwd()
	sourceDir := filepath.Join(wd, "world")
	zipFile := filepath.Join(wd, "backups", b.filename())
	log.Printf("Zipping world dir: path='%s'", sourceDir)

	return zipfiles(sourceDir, zipFile)
}

func zipfiles(s string, dest string) error {
	zipFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	w := zip.NewWriter(zipFile)
	defer w.Close()

	return addFile(w, s, "")
}

func addFile(w *zip.Writer, base string, baseInZip string) error {
	log.Printf("Adding files: base='%s'", base)
	files, err := ioutil.ReadDir(base)
	if err != nil {
		return err
	}

	for _, file := range files {
		absPath := filepath.Join(base, file.Name())
		relPath := filepath.Join(baseInZip, file.Name())

		if file.IsDir() {
			addFile(w, absPath, relPath)
			continue
		}

		f, err := ioutil.ReadFile(absPath)
		if err != nil {
			return err
		}

		newZip, err := w.Create(relPath)
		if err != nil {
			return err
		}
		_, err = newZip.Write(f)
		if err != nil {
			return err
		}
	}

	return nil
}
