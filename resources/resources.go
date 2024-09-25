package Resources

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type StringResources struct {
	Flag  map[string]string
	Error map[string]string
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
	"splitModeDesc":   "-s -n <ROM set name> [-b </path/to/file.bin> & ![-d | -e]] [-d | -e] [-o <output filename>]\nSplit mode. Splits a concatenated binary back into its original MAME files",
	"binFileDesc":     "-s -n <ROM set name> ![-d | -e] [-o <output filename>]\nSupplied with -s when neither -d nor -e are supplied to specify bin file input for splitting into ROM files",
	"concatModeDesc":  "-b </path/to/file.bin> -n <ROM set name> [-o <output filename>]\nConcatenates executable regions of ROM into one file",
	"decryptModeDesc": "-r </path/to/ROM.zip> -n <ROM set name> [-o <output filename>]\nDecrypt mode. Decrypts and concatenates the executable regions of ROM into a single file",
	"encryptModeDesc": "-r </path/to/ROM.zip> -n <ROM set name> [-o <output filename>]\nEncrypt mode. Encrypts and splits the executable regions of ROM back into their MAME format ROM files",
	"outputFileDesc":  "-o </path/to/output.file>\nOptional flag for specifying output file for operations that output a file.",
	"romZipDesc":      "</path/to/ROM.zip> -n <ROM set name>\nRequired when using -d, -e, or -m. Specifies the ROM .zip file to open",
	"romSetNameDesc":  "<ROM set name>\nRequired. Specifies the ROM set (usually the ZIP name)",
}

var errorStrings = map[string]string{
	"noBinFile":    ".bin file is required for this operation",
	"noRomFile":    "ROM file is required for this operation",
	"romParseErr":  "Something went wrong parsing the ROMs",
	"bothEncrypts": "-d and -e are mutually exclusive, in addition to being incompatible with -p and -m",
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

var Strings = StringResources{flagStrings, errorStrings}
var LogText = LogAliases{blue, bold, green, red, yellow}
var Logger = Log{info, warn, error, done}
