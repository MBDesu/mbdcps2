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
// | Concat           |  c   |    5     |       .zip        |        .bin        |   Required   |
// | Decrypt          |  d   |    1     |       .zip        |        .bin        |   Required   |
// | Encrypt          |  e   |    2     |     .bin+.zip     |        .zip        |   Required   |
// | Generate .mra    |  m   |    4     |       .zip        |        .mra        |   Required   |
// | Patch            |  p   |    3     |       .zip        |        .zip        |   Required   |
// | Decode gfx       |  g   |    6     |       .zip        |        .bin        |   Required   |
//
// | Argument        | Flag |   Required With     |
// | :-------------- | :--: | :-----------------: |
// | Output filepath |  o   |       N/A           |
// | Input zip       |  z   |    c, d, g, m, p    |
// | Input bin       |  b   |       e             |
// | ROM set name    |  n   |  c, d, e, g, m, p   |
// | Input diff zip  |  x   |       m             |

type Flags struct {
	isConcatMode    bool
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
	binFile := flag.String("b", "", Resources.Strings.Flag["binFileDesc"])
	concatMode := flag.Bool("c", false, Resources.Strings.Flag["concatModeDesc"])
	decryptMode := flag.Bool("d", false, Resources.Strings.Flag["decryptModeDesc"])
	encryptMode := flag.Bool("e", false, Resources.Strings.Flag["encryptModeDesc"])
	diffMode := flag.Bool("m", false, Resources.Strings.Flag["diffModeDesc"])
	romName := flag.String("n", "", Resources.Strings.Flag["romNameDesc"])
	outputFile := flag.String("o", "", Resources.Strings.Flag["outputFileDesc"])
	patchMode := flag.Bool("p", false, Resources.Strings.Flag["patchModeDesc"])
	mraFile := flag.String("r", "", Resources.Strings.Flag["mraFileDesc"])
	diffZipFile := flag.String("x", "", Resources.Strings.Flag["diffZipDesc"])
	zipFile := flag.String("z", "", Resources.Strings.Flag["zipFileDesc"])

	flag.Parse()
	flags = Flags{*concatMode, *decryptMode, *encryptMode, *patchMode, *diffMode, *romName, *binFile, *outputFile, *zipFile, *diffZipFile, *mraFile}
	validateFlags()
}

func validateFlags() {
	zipFileRequired := flags.isDecryptMode || flags.isEncryptMode || flags.isMraMode || flags.isPatchMode || flags.isConcatMode
	if zipFileRequired && flags.zipFilepath == "" {
		flag.Usage()
		throw(Resources.Strings.Error["noRomFile"])
	}
	binFileRequired := flags.isEncryptMode
	if binFileRequired && flags.binFilepath == "" {
		flag.Usage()
		throw(Resources.Strings.Error["noBinFile"])
	}
	romSetNameRequired := flags.isDecryptMode || flags.isEncryptMode || flags.isMraMode || flags.isPatchMode || flags.isConcatMode
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

func concat() {
	if flags.outputFilepath == "" {
		flags.outputFilepath = flags.romSetName + ".bin"
	}
	romZipFile, romDef, err := cps2rom.ParseRomZip(flags.zipFilepath, flags.romSetName)
	check(err)
	defer romZipFile.Close()
	romBinary, err := cps2rom.ProcessRegionFromZip(romZipFile, romDef.Maincpu)
	check(err)
	f, err := file_utils.CreateFile(flags.outputFilepath)
	check(err)
	f.Write(romBinary)
	f.Close()
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
	check(err)
	mra, err := cps2rom.ParseMra(mraFile)
	check(err)
	fileContentMap, err := file_utils.UnzipFilesToFilenameContentMap(romZipFile)
	check(err)
	Resources.Logger.Warn("Patching ROM...")
	romRegions := []struct {
		regionName string
		region     cps2rom.RomRegion
	}{{"maincpu", romDef.Maincpu}, {"audiocpu", romDef.Audiocpu}, {"qsound", romDef.Qsound}, {"gfx", romDef.Gfx}}
	baseOffset := 0
	for _, region := range romRegions {
		Resources.Logger.Info(fmt.Sprintf("%s (+0x%06x):", region.regionName, baseOffset))
		err = cps2rom.PatchRomRegionWithMra(romZipFile, *mra, region.region, fileContentMap, baseOffset, flags.outputFilepath)
		if region.regionName == "audiocpu" {
			baseOffset += 0x40000
		} else {
			baseOffset += region.region.Size
		}
	}
	check(err)
	Resources.Logger.Done("Done patching ROM!")
	Resources.Logger.Warn("Writing files to .zip...")
	f, err := file_utils.CreateFile(flags.outputFilepath)
	w := *zip.NewWriter(f)
	for file, content := range fileContentMap {
		x, err := w.Create(file)
		check(err)
		x.Write(content)
	}
	w.Close()
	Resources.Logger.Done(fmt.Sprintf("Patched ROM written to %s!", flags.outputFilepath))
}

func diff() {
	if flags.outputFilepath == "" {
		flags.outputFilepath = flags.romSetName + ".mra"
	}
	var patches []cps2rom.RomPatch
	firstRom, romDef, err := cps2rom.ParseRomZip(flags.zipFilepath, flags.romSetName)
	check(err)
	secondRom, _, err := cps2rom.ParseRomZip(flags.diffZipFilepath, flags.romSetName)
	check(err)
	romRegions := []struct {
		regionName string
		region     cps2rom.RomRegion
	}{{"maincpu", romDef.Maincpu}, {"audiocpu", romDef.Audiocpu}, {"qsound", romDef.Qsound}, {"gfx", romDef.Gfx}}
	baseOffset := 0
	Resources.Logger.Warn("Diffing ROMs...")
	for _, region := range romRegions {
		Resources.Logger.Info(fmt.Sprintf("%s (+0x%06x):", region.regionName, baseOffset))
		regionPatches, err := cps2rom.DiffRomRegion(baseOffset, region.region, firstRom, secondRom)
		check(err)
		patches = append(patches, *regionPatches...)
		if region.regionName == "audiocpu" {
			baseOffset += 0x40000
		} else {
			baseOffset += region.region.Size
		}
	}
	patchStrings := cps2rom.GenerateMraPatches(&patches)
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
	} else if flags.isConcatMode {
		concat()
	}
	os.Exit(0)
}
