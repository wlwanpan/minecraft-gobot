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

	"github.com/wlwanpan/minecraft-gobot/awsclient"
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
	return fmt.Sprintf("backup_%s", b.id)
}

func (b *backer) start() error {
	// zip the world directory
	zipFilepath, err := b.zipworld()
	if err != nil {
		log.Printf("Error zipping '/world': error='%s'", err)
		b.state = BACKER_STATE_FAILED
		return err
	}
	log.Println("Successfully zipped!")

	// upload zip file to S3.
	resp, err := b.upload(zipFilepath)
	if err != nil {
		log.Printf("Error uploading zipped file: error='%s'", err)
		b.state = BACKER_STATE_FAILED
		return err
	}

	log.Printf("Backup successsful: url='%s'", resp.S3URL)
	b.state = BACKER_STATE_DONE
	b.lastUrl = resp.S3URL
	return nil
}

func (b *backer) zipworld() (string, error) {
	b.state = BACKER_STATE_ZIPPING

	wd, _ := os.Getwd()
	sourceDir := filepath.Join(wd, "world")
	zipPath := filepath.Join(wd, "backups", b.filename())
	log.Printf("Zipping world dir: path='%s'", sourceDir)

	if err := zipfiles(sourceDir, zipPath); err != nil {
		return "", err
	}
	return zipPath, nil
}

func (b *backer) upload(zipPath string) (*awsclient.S3StoreFileResp, error) {
	b.state = BACKER_STATE_UPLOADING

	client, err := awsclient.New()
	if err != nil {
		return nil, err
	}
	return client.StoreFile(zipPath, b.filename())
}

func zipfiles(s string, dest string) error {
	zipFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	w := zip.NewWriter(zipFile)
	defer w.Close()

	if err := addFile(w, s, ""); err != nil {
		return err
	}

	return err
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
