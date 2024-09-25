package main

import (
	"archive/zip"
	_ "embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MBDesu/mbdcps2/Resources"
	"github.com/MBDesu/mbdcps2/cps2crypt"
	"github.com/MBDesu/mbdcps2/cps2rom"
)

type Flags struct {
	isDecryptMode bool
	isEncryptMode bool
	outputFile    string
}

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
	outputFilePtr := flag.String("o", "out.bin", Resources.Strings.Flag["outputFileDesc"])
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
	flags = Flags{*decryptModePtr, *encryptModePtr, *outputFilePtr}
	if *decryptModePtr && *encryptModePtr {
		throw(Resources.Strings.Error["bothEncrpts"])
	}
}

func checkErr(err error) {
	if err != nil {
		throw(err.Error())
	}
}

func main() {
	parseFlags()
	if !hasRomFile {
		flag.Usage()
		throw(Resources.Strings.Error["noRomFile"])
	}
	err := cps2rom.ValidateRomZip(romDef, romZip)
	if err != nil {
		throw(err.Error())
	}
	if flags.isDecryptMode {
		executableRegionBinary, err := cps2rom.ProcessRegion(romZip, romDef.Maincpu)
		checkErr(err)
		dec, err := cps2crypt.Crypt(cps2crypt.Decrypt, romDef, romZip, executableRegionBinary)
		checkErr(err)
		_, err = os.Stat(filepath.Dir(flags.outputFile))
		if os.IsNotExist(err) {
			throw(err.Error())
		}
		err = os.WriteFile(filepath.Base(flags.outputFile), dec, 0644)
		checkErr(err)
		Resources.Logger.Done(fmt.Sprintf("Decrypted binary written to %s", flags.outputFile))
	}
	defer romZip.Close()
}
