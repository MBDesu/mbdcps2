package cps2rom

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/MBDesu/mbdcps2/Resources"
	utils "github.com/MBDesu/mbdcps2/utils"
)

type RomPatch struct {
	Filename     string
	RegionOffset int
	FileOffset   int
	Data         []uint8
}

type IpsPatchFile struct {
	// header: 0x50 0x41 0x54 0x43 0x48 (PATCH)
	Filename string
	Records  []IpsRecord
	// footer: 0x45 0x4f 0x46 (EOF)
}

type IpsRecord struct {
	Offset int // 24 bits
	Size   uint16
	Data   []byte
}

func newIpsPatchFile(filename string) *IpsPatchFile {
	return &IpsPatchFile{
		Filename: filename,
		Records:  make([]IpsRecord, 0, 1024),
	}
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

func mapOffsetToFile(baseOffset int64, offset int64, romRegion RomRegion) (string, int) {
	for _, operation := range romRegion.Operations {
		actualOffset := offset - baseOffset
		if actualOffset >= int64(operation.Offset) && actualOffset < int64(operation.Offset+operation.Length) {
			return operation.Filename, int(actualOffset - int64(operation.Offset))
		}
	}
	return "", -1
}

func ParseMra(mraFile []byte) (*MraXml, error) {
	var mraXml MraXml
	err := xml.Unmarshal(mraFile, &mraXml)
	return &mraXml, err
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
				operationFilename, absoluteOffset := mapOffsetToFile(int64(baseOffset), offset-0x40, romRegion)
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
			l16 := utils.CreateUint16ArrayFromUint8Array(lb)
			r16 := utils.CreateUint16ArrayFromUint8Array(rb)
			bytesChanged := 0
			for i := 0; i < operation.Length/2; i++ {
				data := make([]uint16, 0, 0x1000)
				for ; l16[i] != r16[i]; i++ {
					data = append(data, r16[i])
				}
				if len(data) > 0 {
					data8 := utils.CreateUint8ArrayFromUint16Array(data)
					fileOffset := (i * 2) - len(data8)
					regionOffset := baseOffset + operation.Offset + fileOffset
					romPatches = append(romPatches, RomPatch{operation.Filename, regionOffset, fileOffset, data8})
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

func convertIpsPatchFileToBinary(ipsPatchData []IpsPatchFile) (ipsPatches map[string][]byte) {
	ipsPatches = map[string][]byte{}
	for _, patchData := range ipsPatchData {
		recordsBinary := []byte{0x50, 0x41, 0x54, 0x43, 0x48}
		for _, patchRecord := range patchData.Records {
			recordsBinary = append(recordsBinary, utils.ConvertUintToByteSlice(uint32(patchRecord.Offset), 3)...)
			recordsBinary = append(recordsBinary, utils.ConvertUintToByteSlice(uint32(patchRecord.Size), 2)...)
			for _, patchByte := range patchRecord.Data {
				recordsBinary = append(recordsBinary, patchByte)
			}
		}
		recordsBinary = append(recordsBinary, 0x45, 0x4f, 0x46)
		ipsPatches[patchData.Filename] = recordsBinary
	}
	return
}

func GenerateIpsPatches(patches *[]RomPatch) (ipsPatchFiles map[string][]byte) {
	ipsPatchData := make([]IpsPatchFile, 0, 20)
	filePatchMap := map[string][]RomPatch{}
	for _, patch := range *patches {
		filePatchMap[patch.Filename] = append(filePatchMap[patch.Filename], patch)
	}
	for patchFilename, patches := range filePatchMap {
		newPatchFile := newIpsPatchFile(patchFilename)
		for _, patch := range patches {
			newPatchFile.Records = append(newPatchFile.Records, IpsRecord{
				Offset: patch.FileOffset,
				Size:   uint16(len(patch.Data) & 0xffff),
				Data:   patch.Data,
			})
		}
		ipsPatchData = append(ipsPatchData, *newPatchFile)
	}
	return convertIpsPatchFileToBinary(ipsPatchData)
}

func GenerateMraPatches(patches *[]RomPatch) (patchStrings []string) {
	patchStrings = make([]string, len(*patches), len(*patches)+10)
	var currentFile = ""
	for _, patch := range *patches {
		if currentFile != patch.Filename {
			currentFile = patch.Filename
			patchStrings = append(patchStrings, fmt.Sprintf("<!-- %s -->\n", currentFile))
		}
		patchString := fmt.Sprintf("<patch offset=\"0x%08x\">", patch.RegionOffset+0x40)
		for i, b := range patch.Data {
			if i == len(patch.Data)-1 {
				patchString += fmt.Sprintf("%02x</patch>\n", b)
			} else {
				patchString += fmt.Sprintf("%02x ", b)
			}
		}
		patchStrings = append(patchStrings, patchString)
	}
	return
}
