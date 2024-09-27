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
	"binFileDesc":     "-s -n <ROM set name> ![-d | -e] [-o </path/to/output.zip>]\nSupplied with -s when neither -d nor -e are supplied to specify bin file input for splitting into ROM files. Default output file is ./<ROM set name>.zip\n",
	"concatModeDesc":  "-r </path/to/ROM.zip> -n <ROM set name> [-o </path/to/output.bin>]\nConcatenates executable regions of ROM into one file. Default output file is ./<ROM set name>.bin\n",
	"diffModeDesc":    "-r </path/to/clean/ROM.zip> -n <ROM set name> [-o </path/to/output/file.mra>]\nDiffs two ROMs and produces .mra <patch> tags for the differences. Default output file is ./<ROM set name>.mra\n",
	"decryptModeDesc": "-r </path/to/ROM.zip> -n <ROM set name> [-s] [-o <output filename>]\nDecrypt mode. Decrypts and concatenates the executable regions of ROM into a single .bin file, unless the -s flag is set. Default output file is ./<ROM set name>.bin, unless the -s flag is set, in which case it will be ./<ROM set name>.zip\n",
	"encryptModeDesc": "-r </path/to/ROM.zip> -n <ROM set name> [-o <output filename>]\nEncrypt mode. Encrypts and splits the executable regions of ROM back into their MAME format ROM files\n",
	"outputFileDesc":  "Optional flag for specifying output file for operations that output a file\n",
	"patchModeDesc":   "-r </path/to/clean/ROM.zip> -n <ROM set name> [-d] [-o </path/to/output.zip or .bin>]\nPatch mode. Value supplied is the path to a .mra file. Patches a ROM with a .mra patch set. Default output file is <./<ROM set name>_modified.zip>\n",
	"romZipDesc":      "-n <ROM set name>\nRequired when using -d, -e, or -x. Specifies a ROM .zip file to open\n",
	"romSetNameDesc":  "Required. Specifies the ROM set (usually the ZIP name)\n",
	"splitModeDesc":   "-n <ROM set name> [-b </path/to/file.bin> & ![-d | -e]] [-d | -e] [-o </path/to/output.zip>]\nSplit mode. Splits a concatenated binary back into its original MAME files. This flag is usable with -d or -e, but not if -b is set\n",
}

var errorStrings = map[string]string{
	"diffSize":     "binaries differ in size",
	"noBinFile":    ".bin file is required for this operation",
	"noMraFile":    ".mra file is required for this operation",
	"noRomFile":    "ROM file is required for this operation",
	"noRomName":    "ROM set name is required",
	"romParseErr":  "Something went wrong parsing the ROMs",
	"bothEncrypts": "-d and -e are mutually exclusive, in addition to being incompatible with -p and -m",
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
