# mbdcps2

`mbdcps2` is a CLI utility for modifying and reverse engineering CPS2 ROMs.



## Features

### Implemented

- [x] Decryption ✅ 2024-09-25

### TODO

- [ ] Encryption (`.bin` -> MAME)
- [ ] MAME <-> Darksoft conversion
- [ ] `.mra`/`.ips` patching
- [ ] Unshuffling graphics
- [ ] Diff between clean and modified ROMs to produce `.ips` and `.mra` patches automagically

## Usage

| Flag | Usage                                                            | Description                                                                                           |
| :--: | :--------------------------------------------------------------- | :---------------------------------------------------------------------------------------------------- |
| `-d` | `-r </path/to/ROM.zip> -n <ROM set name> [-o <output filename>]` | Decrypt mode. Decrypts and concatenates the executable regions of ROM into a single file              |
| `-e` | `-r </path/to/ROM.zip> -n <ROM set name> [-o <output filename>]` | Encrypt mode. Encrypts and splits the executable regions of ROM back into their MAME format ROM files |

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