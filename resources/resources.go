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
	"decryptModeDesc": "-z </path/to/ROM.zip> -n <ROM set name> [-o </path/to/output/file.bin>]",
	"encryptModeDesc": "-b </path/to/decrypted.bin> -z </path/to/ROM.zip> -n <ROM set name> [-o </path/to/output/file.zip>]",
	"patchModeDesc":   "-z </path/to/ROM.zip> -n <ROM set name> [-o </path/to/output/file.zip]",
	"diffModeDesc":    "-z </path/to/ROM.zip>",
	"romSetNameDesc":  "-z </path/to/ROM.zip> -x </path/to/modified/ROM.zip> -n <ROM set name> [-o </path/to/output/file.zip>]",
	"binFileDesc":     "Specifies an input .bin file. Required with the e flag",
	"outputFileDesc":  "Specifies an output file path. Optional",
	"zipFileDesc":     "Specifies an input ROM .zip. Required with d, m, p flags",
	"diffZipDesc":     "Specifies an input ROM .zip to diff against the z flag for generating .mra patches. Required with the m flag",
	"mraFileDesc":     "Specifies an input .mra to patch the z flag input with. Required with the p flag",
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
