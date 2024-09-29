package main

import (
	"archive/zip"
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/MBDesu/mbdcps2/Resources"
	"github.com/MBDesu/mbdcps2/cps2crypt"
	"github.com/MBDesu/mbdcps2/cps2rom"
	file_utils "github.com/MBDesu/mbdcps2/utils"
)

// Ideal Workflow
//
// Decrypt `.zip` to `.bin`, edit hex, encrypt to `.zip`, diff with clean `.zip` to generate `.mra` patches
//
// Then encrypt should take a `.bin` and a `.zip` as input, adding the missing files back to the `.zip`
//
// | Mode             | Flag | Priority | Input File Format | Output File Format | ROM set name |
// | ---------------- | :--: | :------: | :---------------: | :----------------: | :----------: |
// | Convert          |  c   |    4     |       .zip        |        .zip        |     N/A      |
// | Decrypt          |  d   |    1     |       .zip        |        .bin        |   Required   |
// | Encrypt          |  e   |    2     |     .bin+.zip     |        .zip        |   Required   |
// | Generate .mra    |  m   |    5     |       .zip        |        .mra        |   Required   |
// | Patch            |  p   |    3     |       .zip        |        .zip        |   Required   |
//
// | Argument        | Flag | Required With |
// | :-------------- | :--: | :-----------: |
// | Output filepath |  o   |      N/A      |
// | Input zip       |  z   |  c, d, m, p   |
// | Input bin       |  b   |       e       |
// | ROM set name    |  n   |  d, e, m, p   |
// | Input diff zip  |  x   |       m       |

type Flags struct {
	isDecryptMode   bool
	isEncryptMode   bool
	isPatchMode     bool
	isMraMode       bool
	romSetName      string
	binFilepath     string
	outputFilepath  string
	zipFilepath     string
	diffZipFilepath string
	mraFilepath     string
}

var flags Flags

func parseFlags() {
	decryptMode := flag.Bool("d", false, Resources.Strings.Flag["decryptModeDesc"])
	encryptMode := flag.Bool("e", false, Resources.Strings.Flag["encryptModeDesc"])
	patchMode := flag.Bool("p", false, Resources.Strings.Flag["patchModeDesc"])
	diffMode := flag.Bool("m", false, Resources.Strings.Flag["diffModeDesc"])
	romName := flag.String("n", "", Resources.Strings.Flag["romNameDesc"])
	binFile := flag.String("b", "", Resources.Strings.Flag["binFileDesc"])
	outputFile := flag.String("o", "", Resources.Strings.Flag["outputFileDesc"])
	zipFile := flag.String("z", "", Resources.Strings.Flag["zipFileDesc"])
	diffZipFile := flag.String("x", "", Resources.Strings.Flag["diffZipDesc"])
	mraFile := flag.String("r", "", Resources.Strings.Flag["mraFileDesc"])

	flag.Parse()
	flags = Flags{*decryptMode, *encryptMode, *patchMode, *diffMode, *romName, *binFile, *outputFile, *zipFile, *diffZipFile, *mraFile}
	validateFlags()
}

func validateFlags() {
	zipFileRequired := flags.isDecryptMode || flags.isEncryptMode || flags.isMraMode || flags.isPatchMode
	if zipFileRequired && flags.zipFilepath == "" {
		flag.Usage()
		throw(Resources.Strings.Error["noRomFile"])
	}
	binFileRequired := flags.isEncryptMode
	if binFileRequired && flags.binFilepath == "" {
		flag.Usage()
		throw(Resources.Strings.Error["noBinFile"])
	}
	romSetNameRequired := flags.isDecryptMode || flags.isEncryptMode || flags.isMraMode || flags.isPatchMode
	if romSetNameRequired && flags.romSetName == "" {
		flag.Usage()
		throw(Resources.Strings.Error["noRomSetName"])
	}
	diffZipFileRequired := flags.isMraMode
	if diffZipFileRequired && flags.diffZipFilepath == "" {
		flag.Usage()
		throw(Resources.Strings.Error["noDiffRomFile"])
	}
	mraFileRequired := flags.isPatchMode
	if mraFileRequired && flags.mraFilepath == "" {
		flag.Usage()
		throw(Resources.Strings.Error["noMraFile"])

	}
}

func throw(errorString string) {
	fmt.Println(Resources.LogText.Red(Resources.LogText.Bold("[!]")) + " " + errorString)
	os.Exit(1)
}

func check(err error) {
	if err != nil {
		throw(err.Error())
	}
}

func decrypt() {
	if flags.outputFilepath == "" {
		flags.outputFilepath = flags.romSetName + ".bin"
	}
	romZipFile, romDef, err := cps2rom.ParseRomZip(flags.zipFilepath, flags.romSetName)
	check(err)
	defer romZipFile.Close()
	romBinary, err := cps2rom.ProcessRegionFromZip(romZipFile, romDef.Maincpu)
	check(err)
	decryptedRomBinary, err := cps2crypt.Crypt(cps2crypt.Decrypt, *romDef, romZipFile, romBinary)
	check(err)
	err = file_utils.WriteBytesToFile(flags.outputFilepath, decryptedRomBinary)
	check(err)
	Resources.Logger.Done(fmt.Sprintf("Decrypted ROM written to %s!", flags.outputFilepath))
}

func encrypt() {
	if flags.outputFilepath == "" {
		flags.outputFilepath = flags.romSetName + ".zip"
	}
	romZipFile, romDef, err := cps2rom.ParseRomZip(flags.zipFilepath, flags.romSetName)
	check(err)
	decryptedRomBinary, err := file_utils.GetFileContents(flags.binFilepath)
	check(err)
	defer romZipFile.Close()
	encryptedRegion, err := cps2crypt.Crypt(cps2crypt.Encrypt, *romDef, romZipFile, decryptedRomBinary)
	check(err)
	err = cps2rom.SplitRegionToFiles(romDef.Maincpu, encryptedRegion, flags.outputFilepath+"_enc")
	check(err)
	encryptedRegionZip, err := zip.OpenReader(flags.outputFilepath + "_enc")
	check(err)
	defer encryptedRegionZip.Close()
	err = cps2rom.WriteModifiedRegionToZip(flags.outputFilepath, romZipFile, encryptedRegionZip, romDef.Maincpu)
	check(err)
	err = file_utils.DeleteFile(flags.outputFilepath + "_enc")
	check(err)
	Resources.Logger.Done(fmt.Sprintf("Encrypted ROM written to %s!", flags.outputFilepath))
}

func patch() {
	if flags.outputFilepath == "" {
		flags.outputFilepath = flags.romSetName + ".zip"
	}
	romZipFile, romDef, err := cps2rom.ParseRomZip(flags.zipFilepath, flags.romSetName)
	check(err)
	mraFile, err := file_utils.GetFileContents(flags.mraFilepath)
	// TODO: make it so you can patch more than just maincpu
	err = cps2rom.PatchRomRegionWithMra(romZipFile, mraFile, romDef.Maincpu, flags.outputFilepath)
	check(err)
	Resources.Logger.Done(fmt.Sprintf("Patched ROM written to %s!", flags.outputFilepath))
}

func diff() {
	if flags.outputFilepath == "" {
		flags.outputFilepath = flags.romSetName + ".mra"
	}
	var lBytes []uint8
	var rBytes []uint8
	romZipFile, romDef, err := cps2rom.ParseRomZip(flags.zipFilepath, flags.romSetName)
	check(err)
	modifiedRomZip, _, err := cps2rom.ParseRomZip(flags.diffZipFilepath, flags.romSetName)
	// TODO: make it so you can diff more than just maincpu for patchering
	lBytes, err = cps2rom.ProcessRegionFromZip(romZipFile, romDef.Maincpu)
	check(err)
	rBytes, err = cps2rom.ProcessRegionFromZip(modifiedRomZip, romDef.Maincpu)
	check(err)
	patches, err := cps2rom.DiffTwoBins(flags.romSetName, lBytes, rBytes, romDef.Maincpu, false)
	check(err)
	patchStrings := cps2rom.GenerateMraPatches(patches)
	patchFile, err := file_utils.CreateFile(flags.outputFilepath)
	check(err)
	_, err = patchFile.WriteString(Resources.Strings.Info["mraHeader"])
	check(err)
	for _, patch := range patchStrings {
		_, err = patchFile.WriteString(patch)
		check(err)
	}
	defer patchFile.Close()
	Resources.Logger.Done(fmt.Sprintf("Patches written to %s!", flags.outputFilepath))
}

func main() {
	err := cps2rom.ParseRoms()
	check(err)
	parseFlags()
	if flags.isDecryptMode {
		decrypt()
	} else if flags.isEncryptMode {
		encrypt()
	} else if flags.isPatchMode {
		patch()
	} else if flags.isMraMode {
		diff()
	}
	os.Exit(0)
}
