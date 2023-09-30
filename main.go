package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func zipFolder(src, target string) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		fmt.Printf("Given folder does not exist: %s\n", src)
		return err
	}
	fmt.Printf("Zipping folder: %s\n", src)

	file, err := os.Create(target)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := zip.NewWriter(file)
	defer writer.Close()

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Method = zip.Deflate

		header.Name, err = filepath.Rel(filepath.Dir(src), path)
		if err != nil {
			return err
		}

		if info.IsDir() {
			header.Name += "/"
		}

		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(headerWriter, file)
		return err
	})
}

func main() {
	start := time.Now()
	for _, path := range os.Args[1:] {
		err := zipFolder(path, path+".zip")
		if err == nil {
			fmt.Printf("Folder %s zipped.\n\n", path)
		}
	}

	fmt.Println("Script took:", time.Since(start))
}
