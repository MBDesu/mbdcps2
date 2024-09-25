# mbdcps2

`mbdcps2` is a CLI utility for modifying and reverse engineering CPS2 ROMs.



## Features

### Implemented

- [x] Decryption (MAME -> `.bin`, MAME -> MAME, `.bin` -> `.bin`) ✅ 2024-09-25
- [x] Encryption (`.bin` -> MAME, MAME -> MAME, `.bin` -> `.bin`) ✅ 2024-09-25
- [x] Splitting (MAME -> `.bin`) ✅ 2024-09-25


### TODO

- [ ] Concatenating (`.bin` -> MAME)
- [ ] MAME <-> Darksoft conversion
- [ ] `.mra`/`.ips` patching
- [ ] Unshuffling graphics
- [ ] Diff between clean and modified ROMs to produce `.ips` and `.mra` patches automagically


## Installation

Navigate to [Releases](https://github.com/MBDesu/mbdcps2/releases) and find the .zip for your OS and architecture. Unzip it into a location of your choosing, but it is recommended that you place it somewhere on your PATH environment variable so you may run it from anywhere.


## Usage

| Flag | Usage                                                                                        | Description                                                                                                                          |
| :--: | :------------------------------------------------------------------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------- |
| `-d` | `-r </path/to/ROM.zip> -n <ROM set name> [-o <output filename>]`                             | Decrypt mode. Decrypts and concatenates the executable regions of ROM into a single .bin file, unless the -s flag is set.                                           |
| `-e` | `-r </path/to/ROM.zip> -n <ROM set name> [-o <output filename>]`                             | Encrypt mode. Encrypts and splits the executable regions of ROM back into their MAME format ROM files.                               |
| `-s` | `-n <ROM set name> [-b </path/to/file.bin> & ![-d \| -e]] [-d \| -e] [-o <output filename>]` | Split mode. Splits a concatenated binary back into its original MAME files. This flag is usable with -d or -e, but not if -b is set. |
| `-b` | `-s -n <ROM set name> ![-d \| -e] [-o <output filename>]`                                    | Supplied with `-s` when neither `-d` nor `-e` are supplied to specify bin file input for splitting into ROM files.                   |
| `-o` | `</path/to/output.file>`                                                                     | Optional flag for specifying output file for operations that output a file.                                                          |
| `-r` | `</path/to/ROM.zip> -n <ROM set name> [-d \| -e]`                                            | Required when using `-d` or `-e`. Specifies the ROM .zip file to open.                                                               |
| `-n` | `<ROM set name>`                                                                             | Required. Specifies the ROM set (usually the ZIP name) to be worked with.                                                            |


### Supported ROM Sets

| ROM Set   |          |           |           |
| --------- | -------- | --------- | --------- |
| 1944      | ddsom    | spf2th    | vhuntjr2  |
| 1944j     | ddsoma   | spf2tu    | vsav      |
| 1944u     | ddsomar1 | spf2xj    | vsav2     |
| 19xx      | ddsomb   | ssf2      | vsava     |
| 19xxa     | ddsomh   | ssf2a     | vsavb     |
| 19xxar1   | ddsomj   | ssf2ar1   | vsavh     |
| 19xxb     | ddsomjr1 | ssf2h     | vsavj     |
| 19xxh     | ddsomjr2 | ssf2j     | vsavu     |
| 19xxj     | ddsomr1  | ssf2jr1   | xmcota    |
| 19xxjr1   | ddsomr2  | ssf2jr2   | xmcotaa   |
| 19xxjr2   | ddsomr3  | ssf2r1    | xmcotaar1 |
| 19xxu     | ddsomu   | ssf2t     | xmcotaar2 |
| armwar    | ddsomur1 | ssf2ta    | xmcotab   |
| armwara   | ddtod    | ssf2tb    | xmcotah   |
| armwarar1 | ddtoda   | ssf2tba   | xmcotahr1 |
| armwarb   | ddtodar1 | ssf2tbh   | xmcotaj   |
| armwarr1  | ddtodh   | ssf2tbj   | xmcotaj1  |
| armwaru   | ddtodhr1 | ssf2tbj1  | xmcotaj2  |
| armwaru1  | ddtodhr2 | ssf2tbr1  | xmcotaj3  |
| avsp      | ddtodj   | ssf2tbu   | xmcotajr  |
| avspa     | ddtodjr1 | ssf2th    | xmcotar1  |
| avsph     | ddtodjr2 | ssf2tu    | xmcotau   |
| avspj     | ddtodr1  | ssf2tur1  | xmvsf     |
| avspu     | ddtodu   | ssf2u     | xmvsfa    |
| batcir    | ddtodur1 | ssf2us2   | xmvsfar1  |
| batcira   | dimahoo  | ssf2xj    | xmvsfar2  |
| batcirj   | dimahoou | ssf2xjr1  | xmvsfar3  |
| choko     | dstlk    | ssf2xjr1r | xmvsfb    |
| csclub    | dstlka   | uecology  | xmvsfh    |
| csclub1   | dstlkh   | vampj     | xmvsfj    |
| cscluba   | pgear    | vampja    | xmvsfjr1  |
| csclubh   | sgemfa   | vampjr1   | xmvsfjr2  |
| csclubj   | sgemfh   | vhunt2    | xmvsfjr3  |
| csclubjy  | smbomb   | vhunt2r1  | xmvsfr1   |
| cybots    | smbombr1 | vhuntj    | xmvsfu    |
| cybotsj   | spf2t    | vhuntjr1  | xmvsfur1  |
| cybotsu   | spf2ta   | vhuntjr1s | xmvsfur2  |


## Building

Clone this repository and run `go build -ldflags="-w -s" -gcflags=all=-l -o /path/to/output/mbdcps2` to build for your architecture.

To build for other architectures, run `env GOOS=<OS> GOARCH=<arch> go build -ldflags="-w -s" -gcflags=all=-l -o /path/to/output/mbdcps2`.

The following builds are available in each release:

|    Target      | `amd64` | `arm` | `arm64` |
| -------------: | :-----: | :---: | :-----: |
| macOS (darwin) |   ✅    |   ❌   |    ✅   |
| FreeBSD        |   ✅    |   ✅   |    ✅   |
| Linux          |   ✅    |   ✅   |    ✅   |
| NetBSD         |   ✅    |   ✅   |    ✅   |
| OpenBSD        |   ✅    |   ✅   |    ✅   |
| Windows        |   ✅    |   ✅   |    ✅   |

The .zip file for each OS and arch is formatted as `mbdcps2-<OS>-<arch>-<version>.zip`. For example, the `amd64` for macOS build of `mbdcps2` v0.0.1 is named `mbdcps2-darwin-amd64-0.0.1.zip`.


## Contributing

Feel free to make a pull request with any changes you'd like.