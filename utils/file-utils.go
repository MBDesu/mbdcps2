package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func CreateFile(file_path string) (*os.File, error) {
	_, err := os.Stat(filepath.Dir(file_path))
	if err != nil {
		return nil, err
	}
	if os.IsNotExist(err) {
		return nil, err
	}
	f, err := os.Create(filepath.Clean(file_path))
	if err != nil {
		return nil, err
	}
	return f, err
}

func DeleteFile(file_path string) error {
	_, err := os.Stat(filepath.Dir(file_path))
	if err != nil {
		return err
	}
	if os.IsNotExist(err) {
		return err
	}
	err = os.Remove(filepath.Clean(file_path))
	return err
}

func GetFileContents(file_path string) ([]byte, error) {
	return os.ReadFile(filepath.Clean(file_path))
}

func GetZipFileReader(file_path string) (*zip.ReadCloser, error) {
	return zip.OpenReader(filepath.Clean(file_path))
}

func WriteBytesToFile(file_path string, bytes []byte) error {
	_, err := os.Stat(filepath.Dir(file_path))
	if err != nil {
		return err
	}
	if os.IsNotExist(err) {
		return err
	}
	err = os.WriteFile(filepath.Base(file_path), bytes, 0644)
	return err
}

func SwapFileBytes(file_path string) ([]byte, error) {
	fileContents, err := GetFileContents(file_path)
	if err != nil {
		return nil, err
	}
	i := 0
	for i < len(fileContents) {
		tmp := fileContents[i]
		fileContents[i] = fileContents[i+1]
		fileContents[i+1] = tmp
		i += 2
	}
	return fileContents, err
}

func UnzipFilesToFilenameContentMap(zipFile *zip.ReadCloser) (map[string][]byte, error) {
	var bytes = make(map[string][]byte, len(zipFile.File))
	for _, file := range zipFile.File {
		r, err := file.Open()
		if err != nil {
			return nil, err
		}
		fileContents, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		bytes[file.Name] = fileContents
		defer r.Close()
	}
	return bytes, nil
}
