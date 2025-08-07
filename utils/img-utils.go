package utils

type Palette_Color struct {
	R uint8
	G uint8
	B uint8
}

func (p *Palette_Color) GetColorBytes() (color []byte) {
	color = []byte{p.B, p.G, p.R, 0x0}
	return
}

func NewPaletteColor(colorBytes [3]byte) (paletteColor *Palette_Color) {
	paletteColor = &Palette_Color{R: colorBytes[0], G: colorBytes[1], B: colorBytes[2]}
	return
}

type Palette struct {
	colors [16]Palette_Color
}

func (p *Palette) Data() (data []byte) {
	for _, color := range p.colors {
		data = append(data, color.GetColorBytes()...)
	}
	return
}

func NewPalette(colors []Palette_Color) (p *Palette) {
	paletteColors := [16]Palette_Color{}
	for i := 0; i < 16; i++ {
		paletteColors[i] = colors[i]
	}
	for i := 16 - len(colors); i > 0; i-- {
		paletteColors[i] = *NewPaletteColor([3]byte{0, 0, 0})
	}
	p = &Palette{colors: paletteColors}
	return
}

type Bmp_Header struct {
	FileSize        uint32
	Padding         uint32
	DataOffset      uint32
	WindowsMode     uint32
	Width           uint32
	Height          uint32
	Planes          uint16
	BitsPerPixel    uint16
	Compression     uint32
	DataSize        uint32
	HorizonalPpm    uint32
	VerticalPpm     uint32
	NumColors       uint32
	ImportantColors uint32
}

type Bitmap struct {
	Header      *Bmp_Header
	PaletteData *Palette
	PixelData   []byte
}

func (b *Bmp_Header) Data() (header []byte) {
	header = []byte{
		0x42, 0x4d,
	}
	header = append(header, ConvertUint32ToByteSlice(b.FileSize)...)
	header = append(header, ConvertUint32ToByteSlice(b.Padding)...)
	header = append(header, ConvertUint32ToByteSlice(b.DataOffset)...)
	header = append(header, ConvertUint32ToByteSlice(b.WindowsMode)...)
	header = append(header, ConvertUint32ToByteSlice(b.Width)...)
	header = append(header, ConvertUint32ToByteSlice(b.Height)...)
	header = append(header, ConvertUint16ToByteSlice(b.Planes)...)
	header = append(header, ConvertUint16ToByteSlice(b.BitsPerPixel)...)
	header = append(header, ConvertUint32ToByteSlice(b.Compression)...)
	header = append(header, ConvertUint32ToByteSlice(b.DataSize)...)
	header = append(header, ConvertUint32ToByteSlice(b.HorizonalPpm)...)
	header = append(header, ConvertUint32ToByteSlice(b.VerticalPpm)...)
	header = append(header, ConvertUint32ToByteSlice(b.NumColors)...)
	header = append(header, ConvertUint32ToByteSlice(b.ImportantColors)...)
	return
}

func NewBmpHeader() *Bmp_Header {
	// format: BITMAPINFOHEADER
	return &Bmp_Header{
		0x00000000, // file size
		0x00000000, // reserved
		0x76000000, // offset to pixel data
		0x28000000, // DIB header size
		0x00000000, // width
		0x00000000, // height
		0x0100,     // planes = 1
		0x0400,     // bits per pixel = 4
		0x00000000, // compression = none
		0x00000000, // size of bitmap data
		0x00000000, // print resolution, horizontal (pixels per meter)
		0x00000000, // print resolution, vertical (pixels per meter)
		0x00000000, // number of colors in palette
		0x00000000, // important colors = all
	}
}

func NewBitmap() *Bitmap {
	return &Bitmap{
		NewBmpHeader(),
		NewPalette(
			[]Palette_Color{
				// Default
				// *NewPaletteColor([3]byte{0x00, 0x00, 0x00}),
				// *NewPaletteColor([3]byte{0x12, 0x19, 0x25}),
				// *NewPaletteColor([3]byte{0x78, 0xf0, 0x3b}),
				// *NewPaletteColor([3]byte{0x62, 0x14, 0x0e}),
				// *NewPaletteColor([3]byte{0x7f, 0x00, 0xfa}),
				// *NewPaletteColor([3]byte{0xff, 0x63, 0xff}),
				// *NewPaletteColor([3]byte{0xff, 0x7f, 0x00}),
				// *NewPaletteColor([3]byte{0xff, 0x7f, 0xff}),
				// *NewPaletteColor([3]byte{0x00, 0x80, 0x00}),
				// *NewPaletteColor([3]byte{0x00, 0x80, 0xff}),
				// *NewPaletteColor([3]byte{0x00, 0xff, 0x00}),
				// *NewPaletteColor([3]byte{0x00, 0xff, 0xff}),
				// *NewPaletteColor([3]byte{0xff, 0x80, 0x00}),
				// *NewPaletteColor([3]byte{0xff, 0x80, 0xff}),
				// *NewPaletteColor([3]byte{0xff, 0xff, 0x00}),
				// *NewPaletteColor([3]byte{0xff, 0xff, 0xff}),

				// DE LP
				*NewPaletteColor([3]byte{0x44, 0x44, 0x33}),
				*NewPaletteColor([3]byte{0xff, 0xee, 0xaa}),
				*NewPaletteColor([3]byte{0xff, 0xbb, 0x99}),
				*NewPaletteColor([3]byte{0xee, 0x99, 0x77}),
				*NewPaletteColor([3]byte{0xcc, 0x88, 0x66}),
				*NewPaletteColor([3]byte{0xff, 0xdd, 0x00}),
				*NewPaletteColor([3]byte{0xff, 0x00, 0x00}),
				*NewPaletteColor([3]byte{0x99, 0x55, 0x11}),
				*NewPaletteColor([3]byte{0x55, 0x00, 0x00}),
				*NewPaletteColor([3]byte{0x33, 0x44, 0x55}),
				*NewPaletteColor([3]byte{0x44, 0x66, 0x77}),
				*NewPaletteColor([3]byte{0x66, 0x88, 0x99}),
				*NewPaletteColor([3]byte{0x88, 0xaa, 0xbb}),
				*NewPaletteColor([3]byte{0xbb, 0xcc, 0xdd}),
				*NewPaletteColor([3]byte{0xff, 0xff, 0xff}),
				*NewPaletteColor([3]byte{0x00, 0x00, 0x00}),

				// GA LP
				// *NewPaletteColor([3]byte{0x11, 0x11, 0x11}),
				// *NewPaletteColor([3]byte{0x44, 0x33, 0x55}),
				// *NewPaletteColor([3]byte{0x55, 0x55, 0x77}),
				// *NewPaletteColor([3]byte{0x55, 0x77, 0x99}),
				// *NewPaletteColor([3]byte{0x77, 0x88, 0xbb}),
				// *NewPaletteColor([3]byte{0x99, 0xaa, 0xcc}),
				// *NewPaletteColor([3]byte{0xdd, 0xcc, 0xee}),
				// *NewPaletteColor([3]byte{0xff, 0xee, 0xff}),
				// *NewPaletteColor([3]byte{0x88, 0x33, 0xdd}),
				// *NewPaletteColor([3]byte{0x66, 0x33, 0x99}),
				// *NewPaletteColor([3]byte{0x22, 0x22, 0x77}),
				// *NewPaletteColor([3]byte{0x99, 0x66, 0x00}),
				// *NewPaletteColor([3]byte{0xee, 0xaa, 0x00}),
				// *NewPaletteColor([3]byte{0xff, 0xff, 0x88}),
				// *NewPaletteColor([3]byte{0xff, 0x99, 0x99}),
				// *NewPaletteColor([3]byte{0x00, 0x00, 0x00}),

				// DGA LP
				// *NewPaletteColor([3]byte{17, 17, 17}),
				// *NewPaletteColor([3]byte{68, 34, 34}),
				// *NewPaletteColor([3]byte{85, 17, 34}),
				// *NewPaletteColor([3]byte{102, 34, 68}),
				// *NewPaletteColor([3]byte{119, 51, 68}),
				// *NewPaletteColor([3]byte{119, 119, 102}),
				// *NewPaletteColor([3]byte{153, 136, 102}),
				// *NewPaletteColor([3]byte{187, 170, 136}),
				// *NewPaletteColor([3]byte{51, 68, 68}),
				// *NewPaletteColor([3]byte{34, 51, 51}),
				// *NewPaletteColor([3]byte{17, 34, 34}),
				// *NewPaletteColor([3]byte{68, 0, 0}),
				// *NewPaletteColor([3]byte{136, 0, 0}),
				// *NewPaletteColor([3]byte{170, 0, 0}),
				// *NewPaletteColor([3]byte{153, 0, 85}),
				// *NewPaletteColor([3]byte{0, 0, 0}),
			},
		),
		[]byte{},
	}
}

func (bmp *Bitmap) Data() []byte {
	data := append(bmp.Header.Data(), bmp.PaletteData.Data()...)
	data = append(data, bmp.PixelData...)
	return data
}

func GfxToBitmapPixelData(interleavedGfx []byte, tileSize int) (pixelData []byte) {
	width := 0x1000
	height := (len(interleavedGfx) * 2) / 0x1000

	pixelNibbles := parsePixelNibbles(interleavedGfx)
	pixelData = placePixels(pixelNibbles, width, height, tileSize)
	pixelData = packPixelData(pixelData)
	return
}

func parsePixelNibbles(decodedGfxBin []byte) (pixelNibbles []byte) {
	pixelNibbles = make([]byte, len(decodedGfxBin)*2)
	for i := 0; i < len(decodedGfxBin)*8; i++ {
		j := ((((i >> 3) & 0x1ffff0) >> 1) | (((i >> 3) & 8) << 17) | ((i >> 3) & 0xffe00007)) << 3
		k := (((i ^ 7) & 7) << 3) | ((j & 0x18) >> 3) | ((j & (^0x1f)) << 1)
		pixelNibbles[k/8] |= ((decodedGfxBin[i/8] >> (i & 7)) & 1) << (k & 7)
	}
	return
}

// TODO: needs work; tiles are backwards
func placePixels(pixelNibbles []byte, width int, height int, tileSize int) (pixelData []byte) {
	pixelData = make([]byte, len(pixelNibbles))
	for y := 0; y < height; y += tileSize {
		for x := 0; x < width; x += tileSize {
			offset := (y/tileSize)*(width/tileSize) + (x / tileSize)
			sourceIndex := offset * tileSize * tileSize
			destinationIndex := y*width + x
			for yy := 0; yy < tileSize; yy++ {
				for xx := 0; xx < tileSize; xx++ {
					pixelData[destinationIndex+xx] = pixelNibbles[sourceIndex+xx]
				}
				sourceIndex += tileSize
				destinationIndex += width
			}
		}
	}
	return
}

func packPixelData(pixelNibbles []byte) (pixelData []byte) {
	pixelData = make([]byte, len(pixelNibbles)/2)
	for i := 0; i < len(pixelData); i++ {
		index := len(pixelNibbles) - 2 - (i * 2)
		pixelData[i] = ((pixelNibbles[index+1] << 4) | (pixelNibbles[index] & 0xf))
	}
	return
}
