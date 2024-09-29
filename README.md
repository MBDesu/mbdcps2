# mbdcps2

`mbdcps2` is a CLI utility for modifying and reverse engineering CPS2 ROMs.



## Features

### Implemented

- [x] Decryption (MAME -> `.bin`, MAME -> MAME, `.bin` -> `.bin`) ✅ 2024-09-25
- [x] Encryption (`.bin` -> MAME, MAME -> MAME, `.bin` -> `.bin`) ✅ 2024-09-25
- [x] Diff between clean and modified ROMs to produce `.mra` patches automagically ✅ 2024-09-26
- [x] `.mra` patching ✅ 2024-09-27


### TODO

- [ ] Splitting (`.bin` -> MAME\[encryption spits out a full ROM .zip from a `.bin` file, but that's the only splitting exposed to the user right now])
- [ ] Concatenating (MAME -> `.bin` \[decryption spits out a concatenated `maincpu` region, but that's the only concatenation exposed to the user right now])
- [ ] MAME <-> Darksoft conversion
- [ ] Unshuffling graphics
- [ ] Patching of graphics/audio/etc. regions


## Installation

Navigate to [Releases](https://github.com/MBDesu/mbdcps2/releases) and find the .zip for your OS and architecture. Unzip it into a location of your choosing, but it is recommended that you place it somewhere on your PATH environment variable so you may run it from anywhere.


## Usage

```
  -b string
        Specifies an input .bin file. Required with the e flag
    
  -d    -z </path/to/ROM.zip> -n <ROM set name> [-o </path/to/output/file.bin>]
        Decrypt mode. Decrypts a ROM's opcodes. Output is a concatenation of the decrypted binary .bin
    
  -e    -b </path/to/decrypted.bin> -z </path/to/ROM.zip> -n <ROM set name> [-o </path/to/output/file.zip>]
        Encrypt mode. Encrypts a ROM's opcodes. Output is a full ROM .zip
    
  -m    -z </path/to/ROM.zip> -x </path/to/modified/ROM.zip> -n <ROM set name> [-o </path/to/output/file.zip>]
        Diff mode. Diffs two ROMs of the same ROM set and produces a file with .mra style patches in it. Output is said .mra file
    
  -n string
    
  -o string
        Specifies an output file path. Optional
    
  -p    -z </path/to/ROM.zip> -n <ROM set name> [-o </path/to/output/file.zip]
        Patch mode. Patches a ROM .zip with a .mra file's <patch>es. Output is a full ROM .zip
    
  -r string
        Specifies an input .mra to patch the z flag input with. Required with the p flag
    
  -x string
        Specifies an input ROM .zip to diff against the z flag for generating .mra patches. Required with the m flag
    
  -z string
        Specifies an input ROM .zip. Required with d, m, p flags

```

You can find an example workflow [here](https://gist.github.com/MBDesu/c332f919a653044f7ba2f20316e88f07).


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
