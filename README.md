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

## Usage

| Flag | Usage                                                                                        | Description                                                                                                                          |
| :--: | :------------------------------------------------------------------------------------------- | :----------------------------------------------------------------------------------------------------------------------------------- |
| `-d` | `-r </path/to/ROM.zip> -n <ROM set name> [-o <output filename>]`                             | Decrypt mode. Decrypts and concatenates the executable regions of ROM into a single file.                                            |
| `-e` | `-r </path/to/ROM.zip> -n <ROM set name> [-o <output filename>]`                             | Encrypt mode. Encrypts and splits the executable regions of ROM back into their MAME format ROM files.                               |
| `-s` | `-n <ROM set name> [-b </path/to/file.bin> & ![-d \| -e]] [-d \| -e] [-o <output filename>]` | Split mode. Splits a concatenated binary back into its original MAME files. This flag is usable with -d or -e, but not if -b is set. |
| `-b` | `-s -n <ROM set name> ![-d \| -e] [-o <output filename>]`                                    | Supplied with `-s` when neither `-d` nor `-e` are supplied to specify bin file input for splitting into ROM files.                   |
| `-o` | `</path/to/output.file>`                                                                     | Optional flag for specifying output file for operations that output a file.                                                          |
| `-r` | `</path/to/ROM.zip> -n <ROM set name> [-d \| -e]`                                            | Required when using `-d` or `-e`. Specifies the ROM .zip file to open.                                                               |
| `-n` | `<ROM set name>`                                                                             | Required. Specifies the ROM set (usually the ZIP name) to be worked with.                                                            |


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