package main

import (
	"archive/zip"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/MBDesu/mbdcps2/Resources"
	"github.com/MBDesu/mbdcps2/cps2crypt"
	"github.com/MBDesu/mbdcps2/cps2rom"
	file_utils "github.com/MBDesu/mbdcps2/utils"
)

type Flags struct {
	// isConcatMode  bool
	isDecryptMode bool
	isDiffMode    bool
	isEncryptMode bool
	isSplitMode   bool
	outputFile    string
}

var binFile []byte
var diffMode string
var flags Flags
var modifiedRomBin *os.File
var modifiedRomZip *zip.ReadCloser
var romDef cps2rom.RomDefinition
var romZip *zip.ReadCloser
var romName string = ""
var hasRomFile bool = false

func throw(errorString string) {
	fmt.Println(Resources.LogText.Red(Resources.LogText.Bold("[!]")) + " " + errorString)
	os.Exit(1)
}

func parseFlags() {
	decryptModePtr := flag.Bool("d", false, Resources.Strings.Flag["decryptModeDesc"])
	encryptModePtr := flag.Bool("e", false, Resources.Strings.Flag["encryptModeDesc"])
	splitModePtr := flag.Bool("s", false, Resources.Strings.Flag["splitModeDesc"])
	// concatModePtr := flag.Bool("c", false, Resources.Strings.Flag["concatModeDesc"])
	var isDiffMode bool = false
	flag.Func("b", Resources.Strings.Flag["binFileDesc"], func(binFilepath string) error {
		if !(*splitModePtr /*|| *concatModePtr*/) {
			return nil
		}
		b, err := os.ReadFile(binFilepath)
		if err != nil {
			return err
		}
		binFile = b
		hasRomFile = true
		return err
	})
	flag.Func("x", Resources.Strings.Flag["diffModeDesc"], func(modifiedRomFilepath string) error {
		pathParts := strings.Split(filepath.Base(modifiedRomFilepath), ".")
		diffMode = pathParts[len(pathParts)-1]
		cleanFilepath := filepath.Clean(modifiedRomFilepath)
		if diffMode == "bin" {
			f, err := os.Open(cleanFilepath)
			if err != nil {
				return err
			}
			modifiedRomBin = f
		} else if diffMode == "zip" {
			z, err := zip.OpenReader(cleanFilepath)
			if err != nil {
				return err
			}
			modifiedRomZip = z
		} else {
			return errors.New("modified ROM file must be a .bin or .zip")
		}
		isDiffMode = true
		return nil
	})
	flag.Func("n", Resources.Strings.Flag["romSetNameDesc"], func(romname string) error {
		if romname == "" {
			flag.Usage()
			return fmt.Errorf("%s %s", Resources.LogText.Red(Resources.LogText.Bold("[!]")), "ROM set name is required")
		}
		romName = romname
		tmpRom, ok := (*cps2rom.RomDefinitions)[romName]
		if !ok {
			flag.Usage()
			return fmt.Errorf("%s is not a valid or supported ROM set. If %s is a hack of an existing ROM set but not named such, pass the ROM set name with the -n flag", romName, romName)
		}
		romDef = tmpRom
		return nil
	})
	outputFilePtr := flag.String("o", "", Resources.Strings.Flag["outputFileDesc"])
	flag.Func("r", Resources.Strings.Flag["romZipDesc"], func(zipname string) error {
		r, err := zip.OpenReader(zipname)
		if err != nil {
			return err
		}
		romZip = r
		err = cps2rom.ParseRoms()
		if err != nil {
			return err
		}

		hasRomFile = true
		return err
	})

	flag.Parse()
	flags = Flags{ /**concatModePtr,*/ *decryptModePtr, isDiffMode, *encryptModePtr, *splitModePtr, *outputFilePtr}
	if *decryptModePtr && *encryptModePtr {
		throw(Resources.Strings.Error["bothEncrypts"])
	}
}

func checkErr(err error) {
	if err != nil {
		throw(err.Error())
	}
}

// TODO: refactor flag spaghetti
func handleEncryptionOperation(decryptMode cps2crypt.Direction) {
	if flags.outputFile == "" {
		if flags.isSplitMode {
			flags.outputFile = romName + ".zip"
		} else {
			flags.outputFile = romName + ".bin"
		}
	}
	executableRegionBinary, err := cps2rom.ProcessRegionFromZip(romZip, romDef.Maincpu)
	checkErr(err)
	res, err := cps2crypt.Crypt(decryptMode, romDef, romZip, executableRegionBinary)
	checkErr(err)
	if flags.isSplitMode {
		file_utils.SplitRegionToFiles(romDef.Maincpu, res, flags.outputFile)
	} else {
		file_utils.WriteBytesToFile(flags.outputFile, res)
	}
	operation := "Encrypted"
	if decryptMode {
		operation = "Decrypted"
	}
	object := "binary"
	if flags.isSplitMode {
		object = "zip"
	}
	Resources.Logger.Done(fmt.Sprintf("%s %s written to %s", operation, object, flags.outputFile))
}

func handleDiffOperation() {
	if flags.outputFile == "" {
		flags.outputFile = romName + ".mra"
	}
	var lBytes []uint8
	var rBytes []uint8
	// TODO: make it so you can diff more than just maincpu for patchering
	lBytes, err := cps2rom.ProcessRegionFromZip(romZip, romDef.Maincpu)
	checkErr(err)
	if diffMode == "zip" {
		rBytes, err = cps2rom.ProcessRegionFromZip(modifiedRomZip, romDef.Maincpu)
		checkErr(err)
	} else {
		rBytes, err = io.ReadAll(modifiedRomBin)
		checkErr(err)
	}
	patches, err := cps2rom.DiffTwoBins(romName, lBytes, rBytes, romDef.Maincpu, false)
	checkErr(err)
	patchStrings := cps2rom.GenerateMraPatches(patches)
	patchFile, err := file_utils.CreateFile(flags.outputFile)
	checkErr(err)
	for _, patch := range patchStrings {
		_, err = patchFile.WriteString(patch)
		checkErr(err)
	}
	defer patchFile.Close()
}

func main() {
	parseFlags()
	if romName == "" {
		flag.Usage()
		throw(Resources.Strings.Error["noRomName"])
	}
	if !hasRomFile && (flags.isDecryptMode || flags.isEncryptMode || flags.isDiffMode) {
		flag.Usage()
		throw(Resources.Strings.Error["noRomFile"])
	} else if !hasRomFile && flags.isSplitMode {
		flag.Usage()
		throw(Resources.Strings.Error["noBinFile"])
	}
	if flags.isDecryptMode || flags.isEncryptMode {
		err := cps2rom.ValidateRomZip(romDef, romZip)
		if err != nil {
			throw(err.Error())
		}
		handleEncryptionOperation(cps2crypt.Direction(flags.isDecryptMode))
	} else if flags.isSplitMode {
		err := file_utils.SplitRegionToFiles(romDef.Maincpu, binFile, flags.outputFile)
		checkErr(err)
		Resources.Logger.Done(fmt.Sprintf("%s maincpu files written to %s", romName, flags.outputFile))
	} else if flags.isDiffMode {
		handleDiffOperation()
		Resources.Logger.Done(fmt.Sprintf("%s .mra patches written to %s", romName, flags.outputFile))
	}
	defer romZip.Close()
	os.Exit(0)
}
