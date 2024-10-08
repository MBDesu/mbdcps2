package cps2rom

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/MBDesu/mbdcps2/Resources"
	file_utils "github.com/MBDesu/mbdcps2/utils"
)

func SplitRegionToFiles(romRegion RomRegion, binary []byte, zipPath string) error {
	f, err := file_utils.CreateFile(zipPath)
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

func ValidateRomZip(romDefinition RomDefinition, zip *zip.ReadCloser) error {
	var numFiles = len(romDefinition.Maincpu.Operations) + len(romDefinition.Audiocpu.Operations) + len(romDefinition.Gfx.Operations) + len(romDefinition.Qsound.Operations) + len(romDefinition.Key.Operations)
	regions := []RomRegion{romDefinition.Audiocpu, romDefinition.Gfx, romDefinition.Maincpu, romDefinition.Qsound, romDefinition.Key}
	requiredFiles := make([]string, 0, numFiles)

	for _, region := range regions {
		if len(region.Operations) > 0 {
			for _, operation := range region.Operations {
				if operation.Filename != "" {
					requiredFiles = append(requiredFiles, operation.Filename)
				}
			}
		}
	}

	hasFiles := make(map[string]bool)
	for _, filename := range requiredFiles {
		hasFiles[filename] = false
	}
	for _, file := range zip.File {
		var name = file.Name
		_, ok := hasFiles[name] // using extance of the key
		if ok {
			hasFiles[name] = true
		}
	}

	numMissingFiles := 0
	missingFiles := make([]string, 0, numFiles)
	for filename, hasFile := range hasFiles {
		if !hasFile {
			numMissingFiles = numMissingFiles + 1
			missingFiles = append(missingFiles, filename)
		}
	}
	if numMissingFiles > 0 {
		logString := fmt.Sprintf("missing %d files:\n", numMissingFiles)
		for _, missingFile := range missingFiles {
			logString += "    " + Resources.LogText.Bold(missingFile) + "\n"
		}
		return fmt.Errorf(logString)
	}

	return nil
}

func ProcessRegionFromZip(romZip *zip.ReadCloser, region RomRegion) ([]uint8, error) {
	Resources.Logger.Warn("Processing binary...")
	regionBinary := make([]uint8, region.Size)

	for _, operation := range region.Operations {
		var bufPtr = operation.Offset

		// TODO: implement other operations
		switch strings.ToLower(operation.Type) {
		case "fill":
			for range operation.Length {
				regionBinary[bufPtr] = uint8(operation.FillValue & 0xff)
			}
			continue
		case "continue":
			bufPtr += operation.Length
			continue
		default:
			break
		}
		if operation.Type != strings.ToLower("load") {
			continue
		}
		var operationFile *zip.File
		for _, file := range romZip.File {
			if operation.Filename == filepath.Base(file.Name) {
				operationFile = file
				break
			}
		}
		if operationFile != nil {
			r, err := operationFile.Open()
			if err != nil {
				return nil, err
			}
			p, err := io.ReadAll(r)
			if err != nil {
				return nil, err
			}
			defer r.Close()
			bytesLeft := operation.Length
			skip := operation.Skip + operation.GroupSize
			Resources.Logger.Info(fmt.Sprintf("Processing %s, starting at offset +0x%06X", operation.Filename, bufPtr))

			if (operation.GroupSize == 1 || !operation.Reverse) && operation.Skip == 0 {
				for j := range bytesLeft {
					regionBinary[bufPtr] = p[j]
					bufPtr++
				}
			} else if operation.GroupSize == 1 {
				for j := range bytesLeft {
					regionBinary[bufPtr] = p[j]
					bufPtr++
				}
			} else if !operation.Reverse {
				for j := 0; j < operation.Length && bytesLeft > 0; j++ {
					// fmt.Printf("bufPtr = %06x\n", bufPtr)
					for k := 0; k < operation.GroupSize && bytesLeft > 0; k++ {
						// fmt.Printf("%06x\n", k+bufPtr)
						regionBinary[k+bufPtr] = p[j+k]
						bytesLeft--
					}
					bufPtr += skip
				}
			} else {
				bytesWritten := 0
				for bytesWritten < operation.Length {
					for j := operation.GroupSize - 1; j >= 0; j-- {
						regionBinary[bufPtr+j] = p[bytesWritten]
						bytesWritten++
					}
					bufPtr += skip
				}
			}
		}
	}
	Resources.Logger.Done("Done processing binary!")
	return regionBinary, nil
}

func ParseRomZip(file_path string, romSetName string) (*zip.ReadCloser, *RomDefinition, error) {
	Resources.Logger.Warn(fmt.Sprintf("Parsing %s...", filepath.Clean(file_path)))
	romZipFile, err := file_utils.GetZipFileReader(file_path)
	if err != nil {
		return nil, nil, err
	}
	romDef, ok := (*RomDefinitions)[romSetName]
	if !ok {
		return nil, nil, errors.New(fmt.Sprintf("ROM set %s is invalid or unsupported", romSetName))
	}
	err = ValidateRomZip(romDef, romZipFile)
	if err == nil {
		Resources.Logger.Done("ROM OK")
	}
	return romZipFile, &romDef, err
}

func copyZippedFileToNewZip(file *zip.File, newZip *zip.Writer) error {
	romFile, err := file.Open()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(file.FileInfo())
	if err != nil {
		return err
	}
	header.Name = file.Name
	newZippedFile, err := newZip.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(newZippedFile, romFile)
	if err != nil {
		return err
	}
	return err
}

func WriteModifiedRegionToZip(outputFilepath string, romZip *zip.ReadCloser, modifiedRegionZip *zip.ReadCloser, region RomRegion) error {
	f, err := file_utils.CreateFile(outputFilepath)
	if err != nil {
		return err
	}
	excludedRegionFilenames := make([]string, len(region.Operations))
	for i, operation := range region.Operations {
		excludedRegionFilenames[i] = operation.Filename
	}
	newZip := zip.NewWriter(f)
	for _, file := range romZip.File {
		isInExcludedRegion := false
		for _, filename := range excludedRegionFilenames {
			if file.Name == filename {
				isInExcludedRegion = true
				break
			}
		}
		if isInExcludedRegion {
			continue
		}
		err = copyZippedFileToNewZip(file, newZip)
		if err != nil {
			return err
		}
	}
	for _, file := range modifiedRegionZip.File {
		err = copyZippedFileToNewZip(file, newZip)
		if err != nil {
			return err
		}
	}
	newZip.Close()
	return err
}
