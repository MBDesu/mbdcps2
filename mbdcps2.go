package main

import (
	"archive/zip"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/MBDesu/mbdcps2/Resources"
	"github.com/MBDesu/mbdcps2/cps2crypt"
	"github.com/MBDesu/mbdcps2/cps2rom"
	file_utils "github.com/MBDesu/mbdcps2/utils"
)

type Flags struct {
	isConcatMode  bool
	isDecryptMode bool
	isEncryptMode bool
	isSplitMode   bool
	outputFile    string
}

var binFile []byte
var flags Flags
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
	concatModePtr := flag.Bool("c", false, Resources.Strings.Flag["concatModeDesc"])
	flag.Func("b", Resources.Strings.Flag["binFileDesc"], func(binFilepath string) error {
		if !(*splitModePtr || *concatModePtr) {
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
	flag.Func("n", Resources.Strings.Flag["romSetNameDesc"], func(romname string) error {
		if romname == "" {
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
	outputFilePtr := flag.String("o", "out", Resources.Strings.Flag["outputFileDesc"])
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
	flags = Flags{*concatModePtr, *decryptModePtr, *encryptModePtr, *splitModePtr, *outputFilePtr}
	if *decryptModePtr && *encryptModePtr {
		throw(Resources.Strings.Error["bothEncrypts"])
	}
}

func checkErr(err error) {
	if err != nil {
		throw(err.Error())
	}
}

func handleEncryptionOperation(decryptMode cps2crypt.Direction) {
	hasFileExtension := len(strings.Split(flags.outputFile, ".")) > 1
	executableRegionBinary, err := cps2rom.ProcessRegion(romZip, romDef.Maincpu)
	checkErr(err)
	res, err := cps2crypt.Crypt(decryptMode, romDef, romZip, executableRegionBinary)
	checkErr(err)
	if flags.isSplitMode {
		if !hasFileExtension {
			flags.outputFile += ".zip"
		}
		file_utils.SplitRegionToFiles(romDef.Maincpu, res, flags.outputFile)
	} else {
		if !hasFileExtension {
			flags.outputFile += ".bin"
		}
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

func main() {
	parseFlags()
	if !hasRomFile && (flags.isDecryptMode || flags.isEncryptMode) {
		flag.Usage()
		throw(Resources.Strings.Error["noRomFile"])
	} else if !hasRomFile && flags.isSplitMode {
		flag.Usage()
		throw(Resources.Strings.Error["noBinFile"])
	}
	err := cps2rom.ValidateRomZip(romDef, romZip)
	if err != nil {
		throw(err.Error())
	}
	if flags.isDecryptMode || flags.isEncryptMode {
		handleEncryptionOperation(cps2crypt.Direction(flags.isDecryptMode))
	} else if flags.isSplitMode {
		err = file_utils.SplitRegionToFiles(romDef.Maincpu, binFile, flags.outputFile)
		checkErr(err)
		Resources.Logger.Done(fmt.Sprintf("%s maincpu files written to %s", romName, flags.outputFile))
	} else if flags.isConcatMode {
	}
	defer romZip.Close()
	os.Exit(0)
}
