package file_utils

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MBDesu/mbdcps2/Resources"
	"github.com/MBDesu/mbdcps2/cps2rom"
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

func SplitRegionToFiles(romRegion cps2rom.RomRegion, binary []byte, zipPath string) error {
	f, err := CreateFile(zipPath)
	if err != nil {
		return err
	}
	w := zip.NewWriter(f)
	for _, operation := range romRegion.Operations {
		Resources.Logger.Info(fmt.Sprintf("Writing %s from 0x%06x to 0x%06x...", operation.Filename, operation.Offset, operation.Offset+operation.Length))
		regionBytes := binary[operation.Offset : operation.Offset+operation.Length]
		fr, err := w.Create(operation.Filename)
		if err != nil {
			return err
		}
		_, err = fr.Write(regionBytes)
		if err != nil {
			return err
		}
	}
	err = w.Close()
	if err != nil {
		return err
	}
	err = f.Close()
	return err
}
