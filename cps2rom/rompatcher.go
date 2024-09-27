package cps2rom

import (
	"errors"
	"fmt"
	"os"

	"github.com/MBDesu/mbdcps2/Resources"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type RomPatch struct {
	Filename string
	Offset   int
	Data     []uint8
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
