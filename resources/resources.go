package Resources

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type StringResources struct {
	Flag  map[string]string
	Error map[string]string
	Info  map[string]string
}

type LogAliases struct {
	Blue   func(a ...interface{}) string
	Bold   func(a ...interface{}) string
	Green  func(a ...interface{}) string
	Red    func(a ...interface{}) string
	Yellow func(a ...interface{}) string
}

type Log struct {
	Info  func(msg string)
	Warn  func(msg string)
	Error func(msg string)
	Done  func(msg string)
}

var flagStrings = map[string]string{
	"concatModeDesc":  "-z </path/to/ROM.zip> -n <ROM set name> [-o </path/to/output/file.bin>]\nConcatenation mode. Concatenates the maincpu region into a single binary file\n",
	"decryptModeDesc": "-z </path/to/ROM.zip> -n <ROM set name> [-o </path/to/output/file.bin>]\nDecrypt mode. Decrypts a ROM's opcodes. Output is a concatenation of the decrypted binary .bin\n",
	"encryptModeDesc": "-b </path/to/decrypted.bin> -z </path/to/ROM.zip> -n <ROM set name> [-o </path/to/output/file.zip>]\nEncrypt mode. Encrypts a ROM's opcodes. Output is a full ROM .zip\n",
	"patchModeDesc":   "-z </path/to/ROM.zip> -n <ROM set name> [-o </path/to/output/file.zip]\nPatch mode. Patches a ROM .zip with a .mra file's <patch>es. Output is a full ROM .zip\n",
	"diffModeDesc":    "-z </path/to/ROM.zip> -x </path/to/modified/ROM.zip> -n <ROM set name> [-o </path/to/output/file.zip>]\nDiff mode. Diffs two ROMs of the same ROM set and produces a file with .mra style patches in it. Output is said .mra file\n",
	"romSetNameDesc":  "Specifies the ROM set name for the ROM set you are working with. Usually the .zip filename. Required with the c, d, e, m, p flags\n",
	"binFileDesc":     "Specifies an input .bin file. Required with the e flag\n",
	"outputFileDesc":  "Specifies an output file path. Optional\n",
	"zipFileDesc":     "Specifies an input ROM .zip. Required with c, d, m, p flags\n",
	"diffZipDesc":     "Specifies an input ROM .zip to diff against the z flag for generating .mra patches. Required with the m flag\n",
	"mraFileDesc":     "Specifies an input .mra to patch the z flag input with. Required with the p flag\n",
}

var errorStrings = map[string]string{
	"diffSize":      "binaries differ in size",
	"noBinFile":     "-b input .bin file is required for this operation",
	"noMraFile":     "-r input .mra is required for this operation",
	"noRomFile":     "-z input ROM .zip is required for this operation",
	"noDiffRomFile": "-x input modified ROM .zip is required for this operation",
	"noRomSetName":  "-n ROM set name is required for this operation",
	"romParseErr":   "Something went wrong parsing the ROMs",
}

var infoStrings = map[string]string{
	"mraHeader": "<!--\n  these patches are for use with .mra files and are not the actual offsets; to get\n  the actual offsets, subtract 0x40 from these\n-->\n",
}

var blue = color.New(color.FgBlue).SprintFunc()
var bold = color.New(color.Bold).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()

func info(msg string) {
	log(blue(bold("[-]")), msg)
}
func warn(msg string) {
	log(yellow(bold("[+]")), msg)
}
func error(msg string) {
	log(red(bold("[!]")), msg)
}
func done(msg string) {
	log(green(bold("[+]")), msg)
}
func log(glyph string, msg string) {
	if !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	fmt.Printf("%s %s", glyph, msg)
}

var Strings = StringResources{flagStrings, errorStrings, infoStrings}
var LogText = LogAliases{blue, bold, green, red, yellow}
var Logger = Log{info, warn, error, done}
