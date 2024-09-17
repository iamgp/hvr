package main

import (
	"archive/zip"
	"os"
	"path/filepath"
)

func main() {
	// Create a buffer to write our archive to.
	outFile, err := os.Create("testdata/test-lib.zip")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add files to the archive.
	addFiles(w, "testdata/libraries/lib-a", "")

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		panic(err)
	}
}

func addFiles(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := os.ReadDir(basePath)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			dat, err := os.ReadFile(filepath.Join(basePath, file.Name()))
			if err != nil {
				panic(err)
			}

			// Add some files to the archive.
			f, err := w.Create(filepath.Join(baseInZip, file.Name()))
			if err != nil {
				panic(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				panic(err)
			}
		} else if file.IsDir() {
			// Recurse
			newBase := filepath.Join(basePath, file.Name())
			addFiles(w, newBase, filepath.Join(baseInZip, file.Name()))
		}
	}
}
