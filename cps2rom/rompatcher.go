package cps2rom

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/MBDesu/mbdcps2/Resources"
)

type RomPatch struct {
	Filename string
	Offset   int
	Data     []uint8
}

type MraXml struct {
	XMLName xml.Name `xml:"misterromdescription"`
	Text    string   `xml:",chardata"`
	About   struct {
		Text    string `xml:",chardata"`
		Author  string `xml:"author,attr"`
		Webpage string `xml:"webpage,attr"`
		Source  string `xml:"source,attr"`
		Twitter string `xml:"twitter,attr"`
	} `xml:"about"`
	Name         string `xml:"name"`
	Setname      string `xml:"setname"`
	Rbf          string `xml:"rbf"`
	Mameversion  string `xml:"mameversion"`
	Year         string `xml:"year"`
	Manufacturer string `xml:"manufacturer"`
	Players      string `xml:"players"`
	Joystick     string `xml:"joystick"`
	Rotation     string `xml:"rotation"`
	Region       string `xml:"region"`
	Platform     string `xml:"platform"`
	Category     string `xml:"category"`
	Catver       string `xml:"catver"`
	Mraauthor    string `xml:"mraauthor"`
	Rom          []struct {
		Text    string `xml:",chardata"`
		Index   string `xml:"index,attr"`
		Zip     string `xml:"zip,attr"`
		Type    string `xml:"type,attr"`
		Md5     string `xml:"md5,attr"`
		Address string `xml:"address,attr"`
		Part    []struct {
			Text   string `xml:",chardata"`
			Name   string `xml:"name,attr"`
			Crc    string `xml:"crc,attr"`
			Length string `xml:"length,attr"`
		} `xml:"part"`
		Patch []struct {
			Data   string `xml:",chardata"`
			Offset string `xml:"offset,attr"`
		} `xml:"patch"`
		Interleave []struct {
			Text   string `xml:",chardata"`
			Output string `xml:"output,attr"`
			Part   []struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
				Crc  string `xml:"crc,attr"`
				Map  string `xml:"map,attr"`
			} `xml:"part"`
		} `xml:"interleave"`
	} `xml:"rom"`
	Nvram struct {
		Text  string `xml:",chardata"`
		Index string `xml:"index,attr"`
		Size  string `xml:"size,attr"`
	} `xml:"nvram"`
	Buttons struct {
		Text    string `xml:",chardata"`
		Names   string `xml:"names,attr"`
		Default string `xml:"default,attr"`
		Count   string `xml:"count,attr"`
	} `xml:"buttons"`
}

func createUint8ArrayFromUint16Array(arr []uint16) []uint8 {
	newArr := make([]uint8, len(arr)*2)
	for i := 0; i < len(arr); i++ {
		val := uint8((arr[i] & 0xff00) >> 8)
		newArr[i*2] = val
		val = uint8(arr[i] & 0xff)
		newArr[i*2+1] = val
	}

	return newArr
}

func ParseMra(mraFile []byte) (*MraXml, error) {
	var mraXml MraXml
	err := xml.Unmarshal(mraFile, &mraXml)
	return &mraXml, err
}

func mapOffsetToFile(baseOffset int64, offset int64, romRegion RomRegion) (string, int) {
	for _, operation := range romRegion.Operations {
		actualOffset := offset - 0x40 - baseOffset
		if actualOffset >= int64(operation.Offset) && actualOffset < int64(operation.Offset+operation.Length) {
			return operation.Filename, int(actualOffset - int64(operation.Offset))
		}
	}
	return "", -1
}

func PatchRomRegionWithMra(romZip *zip.ReadCloser, mra MraXml, romRegion RomRegion, fileContentMap map[string][]byte, baseOffset int, outputFilepath string) error {
	for _, rom := range mra.Rom {
		lastOperationFilename := ""
		for _, patch := range rom.Patch {
			offset, err := strconv.ParseInt(patch.Offset, 0, 32)
			if err != nil {
				return err
			}
			data := make([]uint8, 0, len(patch.Data)*3)

			for _, byteString := range strings.Split(patch.Data, " ") {
				byte16, err := strconv.ParseInt(byteString, 16, 16)
				if err != nil {
					return err
				}
				byte8 := uint8(byte16 & 0xff)
				data = append(data, byte8)
			}
			if len(data) > 0 {
				operationFilename, absoluteOffset := mapOffsetToFile(int64(baseOffset), offset, romRegion)
				if operationFilename == "" || absoluteOffset == -1 {
					continue
				} else if operationFilename != lastOperationFilename {
					Resources.Logger.Info(fmt.Sprintf("  Patching %s", operationFilename))
					lastOperationFilename = operationFilename
				}
				for i, dataByte := range data {
					fileContentMap[operationFilename][absoluteOffset+i] = dataByte
				}
			}
		}
	}
	return nil
}

func createUint16ArrayFromUint8Array(arr []uint8) []uint16 {
	length := len(arr)
	newArr := make([]uint16, length/2)
	i := 0
	j := 0
	for i < length {
		val := uint16(arr[i+1]) << 8
		val |= uint16(arr[i])
		newArr[j] = val
		i += 2
		j++
	}
	return newArr
}

func DiffRomRegion(baseOffset int, region RomRegion, first *zip.ReadCloser, second *zip.ReadCloser) (*[]RomPatch, error) {
	var romPatches []RomPatch
	for _, operation := range region.Operations {
		if operation.Filename != "" {
			l, err := first.Open(operation.Filename)
			if err != nil {
				return nil, err
			}
			r, err := second.Open(operation.Filename)
			if err != nil {
				return nil, err
			}
			lb, err := io.ReadAll(l)
			if err != nil {
				return nil, err
			}
			rb, err := io.ReadAll(r)
			if err != nil {
				return nil, err
			}
			l16 := createUint16ArrayFromUint8Array(lb)
			r16 := createUint16ArrayFromUint8Array(rb)
			bytesChanged := 0
			for i := 0; i < operation.Length/2; i++ {
				data := make([]uint16, 0, 0x1000)
				for ; l16[i] != r16[i]; i++ {
					data = append(data, r16[i])
				}
				if len(data) > 0 {
					data8 := createUint8ArrayFromUint16Array(data)
					romPatches = append(romPatches, RomPatch{operation.Filename, (baseOffset + operation.Offset + (i * 2)) - len(data8), data8})
					bytesChanged += len(data8)
				}
			}
			logStr := fmt.Sprintf("  %s", operation.Filename)
			if bytesChanged > 0 {
				logStr += fmt.Sprintf(": %d bytes changed", bytesChanged)
				Resources.Logger.Error(logStr)
			} else {
				Resources.Logger.Info(logStr)
			}
		}
	}
	return &romPatches, nil
}

func GenerateMraPatches(patches *[]RomPatch) []string {
	patchStrings := make([]string, len(*patches), len(*patches)+10)
	var currentFile = ""
	for _, patch := range *patches {
		if currentFile != patch.Filename {
			currentFile = patch.Filename
			patchStrings = append(patchStrings, fmt.Sprintf("<!-- %s -->\n", currentFile))
		}
		patchString := fmt.Sprintf("<patch offset=\"0x%08x\">", patch.Offset+0x40)
		for i, b := range patch.Data {
			if i == len(patch.Data)-1 {
				patchString += fmt.Sprintf("%02x</patch>\n", b)
			} else {
				patchString += fmt.Sprintf("%02x ", b)
			}
		}
		patchStrings = append(patchStrings, patchString)
	}
	return patchStrings
}
