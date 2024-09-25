package cps2rom

import (
	"archive/zip"
	"fmt"
	"io"
	"strings"

	"github.com/MBDesu/mbdcps2/Resources"
)

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
		_, ok := hasFiles[name]
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
		return fmt.Errorf("missing %d files: %s", numMissingFiles, Resources.LogText.Bold(strings.Join(missingFiles, ", ")))
	}

	return nil
}

func ProcessRegion(romZip *zip.ReadCloser, region RomRegion) ([]uint8, error) {
	regionBinary := make([]uint8, region.Size)
	for i := range len(region.Operations) {
		operation := region.Operations[i]
		var bufPtr = operation.Offset
		var operationFile *zip.File
		for _, file := range romZip.File {
			if operation.Filename == file.Name {
				operationFile = file
				break
			}
		}
		if operation.Type != strings.ToLower("load") {
			continue
		}
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
		Resources.Logger.Info(fmt.Sprintf("Processing %s, starting at offset +0x%06X; bytes in region: 0x%06X", operationFile.Name, bufPtr, bytesLeft))

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
				for k := 0; k < operation.GroupSize && bytesLeft > 0; k++ {
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
	return regionBinary, nil
}
