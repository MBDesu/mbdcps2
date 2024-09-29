package cps2rom

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/MBDesu/mbdcps2/Resources"
	file_utils "github.com/MBDesu/mbdcps2/utils"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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

func parseMra(mraFile []byte) (*MraXml, error) {
	var mraXml MraXml
	err := xml.Unmarshal(mraFile, &mraXml)
	return &mraXml, err
}

func mapOffsetToFile(offset int64, romRegion RomRegion) (string, int) {
	for _, operation := range romRegion.Operations {
		actualOffset := offset - 0x40
		if actualOffset >= int64(operation.Offset) && offset < int64(operation.Offset+operation.Length) {
			return operation.Filename, int(actualOffset - int64(operation.Offset))
		}
	}
	return "", -1
}

func PatchRomRegionWithMra(romZip *zip.ReadCloser, mraFile []byte, romRegion RomRegion, outputFilepath string) error {
	Resources.Logger.Warn("Patching ROM...")
	fileContentMap, err := file_utils.UnzipFilesToFilenameContentMap(romZip)
	if err != nil {
		return err
	}
	mra, err := parseMra(mraFile)
	if err != nil {
		return err
	}
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
			operationFilename, absoluteOffset := mapOffsetToFile(offset, romRegion)
			if operationFilename != lastOperationFilename {
				Resources.Logger.Info(fmt.Sprintf("Patching %s", operationFilename))
				lastOperationFilename = operationFilename
			}
			for i, dataByte := range data {
				fileContentMap[operationFilename][absoluteOffset+i] = dataByte
			}
		}
	}
	Resources.Logger.Done("Done patching ROM!")
	Resources.Logger.Warn("Writing files to .zip...")
	f, err := file_utils.CreateFile(outputFilepath)
	if err != nil {
		return err
	}
	w := *zip.NewWriter(f)
	for file, content := range fileContentMap {
		x, err := w.Create(file)
		if err != nil {
			return err
		}
		x.Write(content)
	}
	w.Close()
	return err
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

func DiffTwoBins(romName string, one []uint8, two []uint8, region RomRegion, showTable bool) (*[]RomPatch, error) {
	Resources.Logger.Warn("Diffing ROMs...")
	romPatches := make([]RomPatch, 0, 0x1000)
	one_16 := createUint16ArrayFromUint8Array(one)
	two_16 := createUint16ArrayFromUint8Array(two)
	if len(one_16) != len(two_16) {
		return nil, errors.New(Resources.Strings.Error["diffSize"])
	}
	for _, operation := range region.Operations {
		Resources.Logger.Info(fmt.Sprintf("Diffing %s, starting at offset +0x%06x", operation.Filename, operation.Offset))
		for i := (operation.Offset / 2); i < (operation.Offset/2)+(operation.Length/2); i++ {
			data := make([]uint16, 0, 0x10000)
			for ; one_16[i] != two_16[i]; i++ {
				data = append(data, two_16[i])
			}
			if len(data) > 0 {
				data8 := createUint8ArrayFromUint16Array(data)
				romPatches = append(romPatches, RomPatch{operation.Filename, (i * 2) - len(data8), data8})
			}
		}
	}
	one_8 := createUint8ArrayFromUint16Array(one_16)
	if len(romPatches) > 0 && showTable {
		for _, patch := range romPatches {
			patch.Offset *= 2
		}
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"File", "Offset", "1", "2", "Num bytes"})
		t.SetStyle(table.StyleColoredBlackOnBlueWhite)
		t.Style().Color.RowAlternate = text.Colors{text.BgBlue, text.FgYellow}
		t.Style().Color.Row = text.Colors{text.BgBlue, text.FgYellow}
		fmt.Println()
		for _, patch := range romPatches {
			left := ""
			right := ""
			for i := range patch.Data {
				left += fmt.Sprintf("%02x ", one_8[patch.Offset+i])
				right += fmt.Sprintf("%02x ", patch.Data[i])
				if i > 0 && i%0xf == 0 {
					left += "\n"
					right += "\n"
				}
			}
			t.AppendRow(table.Row{patch.Filename, fmt.Sprintf("%06x", patch.Offset), left, right, fmt.Sprintf("0x%02x", len(patch.Data))})
		}
		t.Render()
		fmt.Println()
	}
	Resources.Logger.Done("Done diffing ROMs!")
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
