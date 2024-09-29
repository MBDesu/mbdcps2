package cps2crypt

import (
	"archive/zip"
	"fmt"
	"io"

	"github.com/MBDesu/mbdcps2/Resources"
	"github.com/MBDesu/mbdcps2/cps2rom"
)

type Direction bool

const (
	Encrypt = false
	Decrypt = true
)

var fn1_groupA = []int{10, 4, 6, 7, 2, 13, 15, 14}
var fn1_groupB = []int{0, 1, 3, 5, 8, 9, 11, 12}

var fn2_groupA = []int{6, 0, 2, 13, 1, 4, 14, 7}
var fn2_groupB = []int{3, 5, 9, 10, 8, 15, 12, 11}

func bit8(x uint16, n int) uint8 {
	return uint8(((x >> n) & 1) & 0xff)
}

func bit16(x uint16, n int) uint16 {
	return (x >> n) & 1
}

func bit8_16(x uint8, n int) uint16 {
	return (uint16(x) >> n) & 1
}

func bit32(x uint32, n int) uint32 {
	return (x >> n) & 1
}

type sbox struct {
	table   [64]uint8
	inputs  [6]int
	outputs [2]int
}

func (sbox) extract_inputs(inputs [6]int, val uint32) int {
	var res int = 0
	for i := range 6 {
		if inputs[i] >= 0 {
			res |= int(bit32(val, inputs[i]) << i)
		}
	}
	return res
}

type optimized_sbox struct {
	input_lookup [256]uint8
	output       [64]uint8
}

func (optimized_sbox) optimize(optimized *optimized_sbox, in sbox) {
	for i := range 256 {
		optimized.input_lookup[i] = uint8(in.extract_inputs(in.inputs, uint32(i)))
	}
	for i := range 64 {
		o := in.table[i]
		optimized.output[i] = 0
		if o&1 > 0 {
			optimized.output[i] |= 1 << in.outputs[0]
		}
		if o&2 > 0 {
			optimized.output[i] |= 1 << in.outputs[1]
		}
	}
}

func (optimized_sbox) fn(optimized *optimized_sbox, in uint8, key uint32) uint8 {
	return optimized.output[optimized.input_lookup[in]^uint8(key&0x3f)]
}

var fn1_r1_boxes = []sbox{
	{
		[64]uint8{
			0, 2, 2, 0, 1, 0, 1, 1, 3, 2, 0, 3, 0, 3, 1, 2, 1, 1, 1, 2, 1, 3, 2, 2, 2, 3, 3, 2, 1, 1, 1, 2,
			2, 2, 0, 0, 3, 1, 3, 1, 1, 1, 3, 0, 0, 1, 0, 0, 1, 2, 2, 1, 2, 3, 2, 2, 2, 3, 1, 3, 2, 0, 1, 3,
		},
		[6]int{3, 4, 5, 6, -1, -1},
		[2]int{3, 6},
	},
	{
		[64]uint8{
			3, 0, 2, 2, 2, 1, 1, 1, 1, 2, 1, 0, 0, 0, 2, 3, 2, 3, 1, 3, 0, 0, 0, 2, 1, 2, 2, 3, 0, 3, 3, 3,
			0, 1, 3, 2, 3, 3, 3, 1, 1, 1, 1, 2, 0, 1, 2, 1, 3, 2, 3, 1, 1, 3, 2, 2, 2, 3, 1, 3, 2, 3, 0, 0,
		},
		[6]int{0, 1, 2, 4, 7, -1},
		[2]int{2, 7},
	},
	{
		[64]uint8{
			3, 0, 3, 1, 1, 0, 2, 2, 3, 1, 2, 0, 3, 3, 2, 3, 0, 1, 0, 1, 2, 3, 0, 2, 0, 2, 0, 1, 0, 0, 1, 0,
			2, 3, 1, 2, 1, 0, 2, 0, 2, 1, 0, 1, 0, 2, 1, 0, 3, 1, 2, 3, 1, 3, 1, 1, 1, 2, 0, 2, 2, 0, 0, 0,
		},
		[6]int{0, 1, 2, 3, 6, 7},
		[2]int{0, 1},
	},
	{
		[64]uint8{
			3, 2, 0, 3, 0, 2, 2, 1, 1, 2, 3, 2, 1, 3, 2, 1, 2, 2, 1, 3, 3, 2, 1, 0, 1, 0, 1, 3, 0, 0, 0, 2,
			2, 1, 0, 1, 0, 1, 0, 1, 3, 1, 1, 2, 2, 3, 2, 0, 3, 3, 2, 0, 2, 1, 3, 3, 0, 0, 3, 0, 1, 1, 3, 3,
		},
		[6]int{0, 1, 3, 5, 6, 7},
		[2]int{4, 5},
	},
}

var fn1_r2_boxes = []sbox{
	{
		[64]uint8{
			3, 3, 2, 0, 3, 0, 3, 1, 0, 3, 0, 1, 0, 2, 1, 3, 1, 3, 0, 3, 3, 1, 3, 3, 3, 2, 3, 2, 2, 3, 1, 2,
			0, 2, 2, 1, 0, 1, 2, 0, 3, 3, 0, 1, 3, 2, 1, 2, 3, 0, 1, 3, 0, 1, 2, 2, 1, 2, 1, 2, 0, 1, 3, 0,
		},
		[6]int{0, 1, 2, 3, 6, -1},
		[2]int{1, 6},
	},
	{
		[64]uint8{
			1, 2, 3, 2, 1, 3, 0, 1, 1, 0, 2, 0, 0, 2, 3, 2, 3, 3, 0, 1, 2, 2, 1, 0, 1, 0, 1, 2, 3, 2, 1, 3,
			2, 2, 2, 0, 1, 0, 2, 3, 2, 1, 2, 1, 2, 1, 0, 3, 0, 1, 2, 3, 1, 2, 1, 3, 2, 0, 3, 2, 3, 0, 2, 0,
		},
		[6]int{2, 4, 5, 6, 7, -1},
		[2]int{5, 7},
	},
	{
		[64]uint8{
			0, 1, 0, 2, 1, 1, 0, 1, 0, 2, 2, 2, 1, 3, 0, 0, 1, 1, 3, 1, 2, 2, 2, 3, 1, 0, 3, 3, 3, 2, 2, 2,
			1, 1, 3, 0, 3, 1, 3, 0, 1, 3, 3, 2, 1, 1, 0, 0, 1, 2, 2, 2, 1, 1, 1, 2, 2, 0, 0, 3, 2, 3, 1, 3,
		},
		[6]int{1, 2, 3, 4, 5, 7},
		[2]int{0, 3},
	},
	{
		[64]uint8{
			2, 1, 0, 3, 3, 3, 2, 0, 1, 2, 1, 1, 1, 0, 3, 1, 1, 3, 3, 0, 1, 2, 1, 0, 0, 0, 3, 0, 3, 0, 3, 0,
			1, 3, 3, 3, 0, 3, 2, 0, 2, 1, 2, 2, 2, 1, 1, 3, 0, 1, 0, 1, 0, 1, 1, 1, 1, 3, 1, 0, 1, 2, 3, 3,
		},
		[6]int{0, 1, 3, 4, 6, 7},
		[2]int{2, 4},
	},
}

var fn1_r3_boxes = []sbox{
	{
		[64]uint8{
			0, 0, 0, 3, 3, 1, 1, 0, 2, 0, 2, 0, 0, 0, 3, 2, 0, 1, 2, 3, 2, 2, 1, 0, 3, 0, 0, 0, 0, 0, 2, 3,
			3, 0, 0, 1, 1, 2, 3, 3, 0, 1, 3, 2, 0, 1, 3, 3, 2, 0, 0, 1, 0, 2, 0, 0, 0, 3, 1, 3, 3, 3, 3, 3,
		},
		[6]int{0, 1, 5, 6, 7, -1},
		[2]int{0, 5},
	},
	{
		[64]uint8{
			2, 3, 2, 3, 0, 2, 3, 0, 2, 2, 3, 0, 3, 2, 0, 2, 1, 0, 2, 3, 1, 1, 1, 0, 0, 1, 0, 2, 1, 2, 2, 1,
			3, 0, 2, 1, 2, 3, 3, 0, 3, 2, 3, 1, 0, 2, 1, 0, 1, 2, 2, 3, 0, 2, 1, 3, 1, 3, 0, 2, 1, 1, 1, 3,
		},
		[6]int{2, 3, 4, 6, 7, -1},
		[2]int{6, 7},
	},
	{
		[64]uint8{
			3, 0, 2, 1, 1, 3, 1, 2, 2, 1, 2, 2, 2, 0, 0, 1, 2, 3, 1, 0, 2, 0, 0, 2, 3, 1, 2, 0, 0, 0, 3, 0,
			2, 1, 1, 2, 0, 0, 1, 2, 3, 1, 1, 2, 0, 1, 3, 0, 3, 1, 1, 0, 0, 2, 3, 0, 0, 0, 0, 3, 2, 0, 0, 0,
		},
		[6]int{0, 2, 3, 4, 5, 6},
		[2]int{1, 4},
	},
	{
		[64]uint8{
			0, 1, 0, 0, 2, 1, 3, 2, 3, 3, 2, 1, 0, 1, 1, 1, 1, 1, 0, 3, 3, 1, 1, 0, 0, 2, 2, 1, 0, 3, 3, 2,
			1, 3, 3, 0, 3, 0, 2, 1, 1, 2, 3, 2, 2, 2, 1, 0, 0, 3, 3, 3, 2, 2, 3, 1, 0, 2, 3, 0, 3, 1, 1, 0,
		},
		[6]int{0, 1, 2, 3, 5, 7},
		[2]int{2, 3},
	},
}

var fn1_r4_boxes = []sbox{
	{
		[64]uint8{
			1, 1, 1, 1, 1, 0, 1, 3, 3, 2, 3, 0, 1, 2, 0, 2, 3, 3, 0, 1, 2, 1, 2, 3, 0, 3, 2, 3, 2, 0, 1, 2,
			0, 1, 0, 3, 2, 1, 3, 2, 3, 1, 2, 3, 2, 0, 1, 2, 2, 0, 0, 0, 2, 1, 3, 0, 3, 1, 3, 0, 1, 3, 3, 0,
		},
		[6]int{1, 2, 3, 4, 5, 7},
		[2]int{0, 4},
	},
	{
		[64]uint8{
			3, 0, 0, 0, 0, 1, 0, 2, 3, 3, 1, 3, 0, 3, 1, 2, 2, 2, 3, 1, 0, 0, 2, 0, 1, 0, 2, 2, 3, 3, 0, 0,
			1, 1, 3, 0, 2, 3, 0, 3, 0, 3, 0, 2, 0, 2, 0, 1, 0, 3, 0, 1, 3, 1, 1, 0, 0, 1, 3, 3, 2, 2, 1, 0,
		},
		[6]int{0, 1, 2, 3, 5, 6},
		[2]int{1, 3},
	},
	{
		[64]uint8{
			0, 1, 1, 2, 0, 1, 3, 1, 2, 0, 3, 2, 0, 0, 3, 0, 3, 0, 1, 2, 2, 3, 3, 2, 3, 2, 0, 1, 0, 0, 1, 0,
			3, 0, 2, 3, 0, 2, 2, 2, 1, 1, 0, 2, 2, 0, 0, 1, 2, 1, 1, 1, 2, 3, 0, 3, 1, 2, 3, 3, 1, 1, 3, 0,
		},
		[6]int{0, 2, 4, 5, 6, 7},
		[2]int{2, 6},
	},
	{
		[64]uint8{
			0, 1, 2, 2, 0, 1, 0, 3, 2, 2, 1, 1, 3, 2, 0, 2, 0, 1, 3, 3, 0, 2, 2, 3, 3, 2, 0, 0, 2, 1, 3, 3,
			1, 1, 1, 3, 1, 2, 1, 1, 0, 3, 3, 2, 3, 2, 3, 0, 3, 1, 0, 0, 3, 0, 0, 0, 2, 2, 2, 1, 2, 3, 0, 0,
		},
		[6]int{0, 1, 3, 4, 6, 7},
		[2]int{5, 7},
	},
}

var fn2_r1_boxes = []sbox{
	{
		[64]uint8{
			2, 0, 2, 0, 3, 0, 0, 3, 1, 1, 0, 1, 3, 2, 0, 1, 2, 0, 1, 2, 0, 2, 0, 2, 2, 2, 3, 0, 2, 1, 3, 0,
			0, 1, 0, 1, 2, 2, 3, 3, 0, 3, 0, 2, 3, 0, 1, 2, 1, 1, 0, 2, 0, 3, 1, 1, 2, 2, 1, 3, 1, 1, 3, 1,
		},
		[6]int{0, 3, 4, 5, 7, -1},
		[2]int{6, 7},
	},
	{
		[64]uint8{
			1, 1, 0, 3, 0, 2, 0, 1, 3, 0, 2, 0, 1, 1, 0, 0, 1, 3, 2, 2, 0, 2, 2, 2, 2, 0, 1, 3, 3, 3, 1, 1,
			1, 3, 1, 3, 2, 2, 2, 2, 2, 2, 0, 1, 0, 1, 1, 2, 3, 1, 1, 2, 0, 3, 3, 3, 2, 2, 3, 1, 1, 1, 3, 0,
		},
		[6]int{1, 2, 3, 4, 6, -1},
		[2]int{3, 5},
	},
	{
		[64]uint8{
			1, 0, 2, 2, 3, 3, 3, 3, 1, 2, 2, 1, 0, 1, 2, 1, 1, 2, 3, 1, 2, 0, 0, 1, 2, 3, 1, 2, 0, 0, 0, 2,
			2, 0, 1, 1, 0, 0, 2, 0, 0, 0, 2, 3, 2, 3, 0, 1, 3, 0, 0, 0, 2, 3, 2, 0, 1, 3, 2, 1, 3, 1, 1, 3,
		},
		[6]int{1, 2, 4, 5, 6, 7},
		[2]int{1, 4},
	},
	{
		[64]uint8{
			1, 3, 3, 0, 3, 2, 3, 1, 3, 2, 1, 1, 3, 3, 2, 1, 2, 3, 0, 3, 1, 0, 0, 2, 3, 0, 0, 0, 3, 3, 0, 1,
			2, 3, 0, 0, 0, 1, 2, 1, 3, 0, 0, 1, 0, 2, 2, 2, 3, 3, 1, 2, 1, 3, 0, 0, 0, 3, 0, 1, 3, 2, 2, 0,
		},
		[6]int{0, 2, 3, 5, 6, 7},
		[2]int{0, 2},
	},
}

var fn2_r2_boxes = []sbox{
	{
		[64]uint8{
			3, 1, 3, 0, 3, 0, 3, 1, 3, 0, 0, 1, 1, 3, 0, 3, 1, 1, 0, 1, 2, 3, 2, 3, 3, 1, 2, 2, 2, 0, 2, 3,
			2, 2, 2, 1, 1, 3, 3, 0, 3, 1, 2, 1, 1, 1, 0, 2, 0, 3, 3, 0, 0, 2, 0, 0, 1, 1, 2, 1, 2, 1, 1, 0,
		},
		[6]int{0, 2, 4, 6, -1, -1},
		[2]int{4, 6},
	},
	{
		[64]uint8{
			0, 3, 0, 3, 3, 2, 1, 2, 3, 1, 1, 1, 2, 0, 2, 3, 0, 3, 1, 2, 2, 1, 3, 3, 3, 2, 1, 2, 2, 0, 1, 0,
			2, 3, 0, 1, 2, 0, 1, 1, 2, 0, 2, 1, 2, 0, 2, 3, 3, 1, 0, 2, 3, 3, 0, 3, 1, 1, 3, 0, 0, 1, 2, 0,
		},
		[6]int{1, 3, 4, 5, 6, 7},
		[2]int{0, 3},
	},
	{
		[64]uint8{
			0, 0, 2, 1, 3, 2, 1, 0, 1, 2, 2, 2, 1, 1, 0, 3, 1, 2, 2, 3, 2, 1, 1, 0, 3, 0, 0, 1, 1, 2, 3, 1,
			3, 3, 2, 2, 1, 0, 1, 1, 1, 2, 0, 1, 2, 3, 0, 3, 3, 0, 3, 2, 2, 0, 2, 2, 1, 2, 3, 2, 1, 0, 2, 1,
		},
		[6]int{0, 1, 3, 4, 5, 7},
		[2]int{1, 7},
	},
	{
		[64]uint8{
			0, 2, 1, 2, 0, 2, 2, 0, 1, 3, 2, 0, 3, 2, 3, 0, 3, 3, 2, 3, 1, 2, 3, 1, 2, 2, 0, 0, 2, 2, 1, 2,
			2, 3, 3, 3, 1, 1, 0, 0, 0, 3, 2, 0, 3, 2, 3, 1, 1, 1, 1, 0, 1, 0, 1, 3, 0, 0, 1, 2, 2, 3, 2, 0,
		},
		[6]int{1, 2, 3, 5, 6, 7},
		[2]int{2, 5},
	},
}

var fn2_r3_boxes = []sbox{
	{
		[64]uint8{
			2, 1, 2, 1, 2, 3, 1, 3, 2, 2, 1, 3, 3, 0, 0, 1, 0, 2, 0, 3, 3, 1, 0, 0, 1, 1, 0, 2, 3, 2, 1, 2,
			1, 1, 2, 1, 1, 3, 2, 2, 0, 2, 2, 3, 3, 3, 2, 0, 0, 0, 0, 0, 3, 3, 3, 0, 1, 2, 1, 0, 2, 3, 3, 1,
		},
		[6]int{2, 3, 4, 6, -1, -1},
		[2]int{3, 5},
	},
	{
		[64]uint8{
			3, 2, 3, 3, 1, 0, 3, 0, 2, 0, 1, 1, 1, 0, 3, 0, 3, 1, 3, 1, 0, 1, 2, 3, 2, 2, 3, 2, 0, 1, 1, 2,
			3, 0, 0, 2, 1, 0, 0, 2, 2, 0, 1, 0, 0, 2, 0, 0, 1, 3, 1, 3, 2, 0, 3, 3, 1, 0, 2, 2, 2, 3, 0, 0,
		},
		[6]int{0, 1, 3, 5, 7, -1},
		[2]int{0, 2},
	},
	{
		[64]uint8{
			2, 2, 1, 0, 2, 3, 3, 0, 0, 0, 1, 3, 1, 2, 3, 2, 2, 3, 1, 3, 0, 3, 0, 3, 3, 2, 2, 1, 0, 0, 0, 2,
			1, 2, 2, 2, 0, 0, 1, 2, 0, 1, 3, 0, 2, 3, 2, 1, 3, 2, 2, 2, 3, 1, 3, 0, 2, 0, 2, 1, 0, 3, 3, 1,
		},
		[6]int{0, 1, 2, 3, 5, 7},
		[2]int{1, 6},
	},
	{
		[64]uint8{
			1, 2, 3, 2, 0, 2, 1, 3, 3, 1, 0, 1, 1, 2, 2, 0, 0, 1, 1, 1, 2, 1, 1, 2, 0, 1, 3, 3, 1, 1, 1, 2,
			3, 3, 1, 0, 2, 1, 1, 1, 2, 1, 0, 0, 2, 2, 3, 2, 3, 2, 2, 0, 2, 2, 3, 3, 0, 2, 3, 0, 2, 2, 1, 1,
		},
		[6]int{0, 2, 4, 5, 6, 7},
		[2]int{4, 7},
	},
}

var fn2_r4_boxes = []sbox{
	{
		[64]uint8{
			2, 0, 1, 1, 2, 1, 3, 3, 1, 1, 1, 2, 0, 1, 0, 2, 0, 1, 2, 0, 2, 3, 0, 2, 3, 3, 2, 2, 3, 2, 0, 1,
			3, 0, 2, 0, 2, 3, 1, 3, 2, 0, 0, 1, 1, 2, 3, 1, 1, 1, 0, 1, 2, 0, 3, 3, 1, 1, 1, 3, 3, 1, 1, 0,
		},
		[6]int{0, 1, 3, 6, 7, -1},
		[2]int{0, 3},
	},
	{
		[64]uint8{
			1, 2, 2, 1, 0, 3, 3, 1, 0, 2, 2, 2, 1, 0, 1, 0, 1, 1, 0, 1, 0, 2, 1, 0, 2, 1, 0, 2, 3, 2, 3, 3,
			2, 2, 1, 2, 2, 3, 1, 3, 3, 3, 0, 1, 0, 1, 3, 0, 0, 0, 1, 2, 0, 3, 3, 2, 3, 2, 1, 3, 2, 1, 0, 2,
		},
		[6]int{0, 1, 2, 4, 5, 6},
		[2]int{4, 7},
	},
	{
		[64]uint8{
			2, 3, 2, 1, 3, 2, 3, 0, 0, 2, 1, 1, 0, 0, 3, 2, 3, 1, 0, 1, 2, 2, 2, 1, 3, 2, 2, 1, 0, 2, 1, 2,
			0, 3, 1, 0, 0, 3, 1, 1, 3, 3, 2, 0, 1, 0, 1, 3, 0, 0, 1, 2, 1, 2, 3, 2, 1, 0, 0, 3, 2, 1, 1, 3,
		},
		[6]int{0, 2, 3, 4, 5, 7},
		[2]int{1, 2},
	},
	{
		[64]uint8{
			2, 0, 0, 3, 2, 2, 2, 1, 3, 3, 1, 1, 2, 0, 0, 3, 1, 0, 3, 2, 1, 0, 2, 0, 3, 2, 2, 3, 2, 0, 3, 0,
			1, 3, 0, 2, 2, 1, 3, 3, 0, 1, 0, 3, 1, 1, 3, 2, 0, 3, 0, 2, 3, 2, 1, 3, 2, 3, 0, 0, 1, 3, 2, 1,
		},
		[6]int{2, 3, 4, 5, 6, 7},
		[2]int{5, 6},
	},
}

func bitswap8(val uint16, b7 int, b6 int, b5 int, b4 int, b3 int, b2 int, b1 int, b0 int) uint8 {
	return ((bit8(val, b7) << 7) |
		(bit8(val, b6) << 6) |
		(bit8(val, b5) << 5) |
		(bit8(val, b4) << 4) |
		(bit8(val, b3) << 3) |
		(bit8(val, b2) << 2) |
		(bit8(val, b1) << 1) |
		(bit8(val, b0)))
}

func expandKey(keyNum int, destKey *[]uint32, srcKey []uint32) {
	bitSets := make([][]uint8, 2)
	bitSets[0] = []uint8{
		33, 58, 49, 36, 0, 31,
		22, 30, 3, 16, 5, 53,
		10, 41, 23, 19, 27, 39,
		43, 6, 34, 12, 61, 21,
		48, 13, 32, 35, 6, 42,
		43, 14, 21, 41, 52, 25,
		18, 47, 46, 37, 57, 53,
		20, 8, 55, 54, 59, 60,
		27, 33, 35, 18, 8, 15,
		63, 1, 50, 44, 16, 46,
		5, 4, 45, 51, 38, 25,
		13, 11, 62, 29, 48, 2,
		59, 61, 62, 56, 51, 57,
		54, 9, 24, 63, 22, 7,
		26, 42, 45, 40, 23, 14,
		2, 31, 52, 28, 44, 17,
	}
	bitSets[1] = []uint8{
		34, 9, 32, 24, 44, 54,
		38, 61, 47, 13, 28, 7,
		29, 58, 18, 1, 20, 60,
		15, 6, 11, 43, 39, 19,
		63, 23, 16, 62, 54, 40,
		31, 3, 56, 61, 17, 25,
		47, 38, 55, 57, 5, 4,
		15, 42, 22, 7, 2, 19,
		46, 37, 29, 39, 12, 30,
		49, 57, 31, 41, 26, 27,
		24, 36, 11, 63, 33, 16,
		56, 62, 48, 60, 59, 32,
		12, 30, 53, 48, 10, 0,
		50, 35, 3, 59, 14, 49,
		51, 45, 44, 2, 21, 33,
		55, 52, 23, 28, 8, 26,
	}
	bits := bitSets[keyNum]
	(*destKey)[0] = 0
	(*destKey)[1] = 0
	(*destKey)[2] = 0
	(*destKey)[3] = 0
	for i := range 96 {
		(*destKey)[i/24] |= bit32(srcKey[bits[i]/32], int(bits[i]%32)) << (i % 24)
	}
}

func expandSubkey(subkey *[]uint32, seed uint16) {
	seed &= 0xffff
	bits := []int{
		5, 10, 14, 9, 4, 0, 15, 6, 1, 8, 3, 2, 12, 7, 13, 11,
		5, 12, 7, 2, 13, 11, 9, 14, 4, 1, 6, 10, 8, 0, 15, 3,
		4, 10, 2, 0, 6, 9, 12, 1, 11, 7, 15, 8, 13, 5, 14, 3,
		14, 11, 12, 7, 4, 5, 2, 10, 1, 15, 0, 9, 8, 6, 13, 3,
	}
	(*subkey)[0] = 0
	(*subkey)[1] = 0

	for i := range 64 {
		(*subkey)[i/32] |= uint32(bit16(seed, bits[i])) << (i % 32)
	}
}

func fn(input uint8, sboxes []optimized_sbox, key uint32) uint8 {
	return sboxes[0].fn(&sboxes[0], input, key>>0) |
		sboxes[1].fn(&sboxes[1], input, key>>6) |
		sboxes[2].fn(&sboxes[2], input, key>>12) |
		sboxes[3].fn(&sboxes[3], input, key>>18)
}

func feistel(val uint16,
	bitsA []int, bitsB []int,
	boxes1 []optimized_sbox, boxes2 []optimized_sbox,
	boxes3 []optimized_sbox, boxes4 []optimized_sbox,
	key1 uint32, key2 uint32,
	key3 uint32, key4 uint32) uint16 {
	var l = bitswap8(val, bitsB[7], bitsB[6], bitsB[5], bitsB[4], bitsB[3], bitsB[2], bitsB[1], bitsB[0]) & 0xff
	var r = bitswap8(val, bitsA[7], bitsA[6], bitsA[5], bitsA[4], bitsA[3], bitsA[2], bitsA[1], bitsA[0]) & 0xff

	l ^= fn(r, boxes1, key1)
	r ^= fn(l, boxes2, key2)
	l ^= fn(r, boxes3, key3)
	r ^= fn(l, boxes4, key4)

	return (bit8_16(l, 0) << bitsA[0]) |
		(bit8_16(l, 1) << bitsA[1]) |
		(bit8_16(l, 2) << bitsA[2]) |
		(bit8_16(l, 3) << bitsA[3]) |
		(bit8_16(l, 4) << bitsA[4]) |
		(bit8_16(l, 5) << bitsA[5]) |
		(bit8_16(l, 6) << bitsA[6]) |
		(bit8_16(l, 7) << bitsA[7]) |
		(bit8_16(r, 0) << bitsB[0]) |
		(bit8_16(r, 1) << bitsB[1]) |
		(bit8_16(r, 2) << bitsB[2]) |
		(bit8_16(r, 3) << bitsB[3]) |
		(bit8_16(r, 4) << bitsB[4]) |
		(bit8_16(r, 5) << bitsB[5]) |
		(bit8_16(r, 6) << bitsB[6]) |
		(bit8_16(r, 7) << bitsB[7])

}

func optimizeSBoxes(out []optimized_sbox, input []sbox) {
	for box := range 4 {
		out[box].optimize(&out[box], input[box])
	}
}

func initializeOptimizedSBoxes(sboxes []optimized_sbox) {
	for i := range 4 {
		sboxes[i] = optimized_sbox{}
	}
}

func decodeKey() {
	decoded := [10]uint16{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for b := range 160 {
		bit := (317 - b) % 160
		if (key[bit/8] >> ((bit ^ 7) % 8) & 1) > 0 {
			decoded[b/16] |= (0x8000 >> (b % 16))
		}
	}
	masterKey_1 = (uint32(decoded[0]) << 16) | uint32(decoded[1])
	masterKey_2 = (uint32(decoded[2]) << 16) | uint32(decoded[3])
	upperLimit = 0xffffff
	lowerLimit = 0xff0000
	if decoded[9] != 0xffff {
		upperLimit = ((((int64(^decoded[9])) & 0x3ff) << 14) | 0x3fff) + 1
		lowerLimit = 0
		Resources.Logger.Info(fmt.Sprintf("Master key 1 = 0x%08x", masterKey_1))
		Resources.Logger.Info(fmt.Sprintf("Master key 2 = 0x%08x", masterKey_2))
		Resources.Logger.Info(fmt.Sprintf("Lower limit = 0x%06x", lowerLimit))
		Resources.Logger.Info(fmt.Sprintf("Upper limit = 0x%06x", upperLimit))
	}
	upperLimit /= 2
	lowerLimit /= 2
}

func parseKey(romDef cps2rom.RomDefinition, romZip *zip.ReadCloser) error {
	keyFilename := romDef.Key.Operations[0].Filename
	var keyFile zip.File
	for _, file := range romZip.File {
		if file.Name == keyFilename {
			keyFile = *file
		}
	}
	if keyFile.Name != "" {
		Resources.Logger.Info(fmt.Sprintf("key found: %s (0x%01x bytes)\n", keyFile.Name, keyFile.UncompressedSize64))
	}
	r, err := keyFile.Open()
	if err != nil {
		return err
	}
	defer r.Close()
	p, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	keyBytes := make([]uint8, keyFile.UncompressedSize64, 0x14)
	for b := range len(p) {
		keyBytes[b] = uint8(p[b])
	}
	key = keyBytes
	decodeKey()
	return nil
}

func createUint8ArrayFromUint16Array(arr []uint16) []uint8 {
	newArr := make([]uint8, len(arr)*2)
	for i := 0; i < len(arr); i++ {
		if encDir == Decrypt {
			val := uint8((arr[i] & 0xff00) >> 8)
			newArr[i*2] = val
			val = uint8(arr[i] & 0xff)
			newArr[i*2+1] = val
		} else {
			val := uint8(arr[i] & 0xff)
			newArr[i*2] = val
			val = uint8((arr[i] & 0xff00) >> 8)
			newArr[i*2+1] = val
		}
	}

	return newArr
}

func createUint16ArrayFromUint8Array(arr []uint8) []uint16 {
	length := len(arr)
	newArr := make([]uint16, length/2)
	i := 0
	j := 0
	for i < length {
		val := uint16(arr[i]) << 8
		val |= uint16(arr[i+1])
		newArr[j] = val
		i += 2
		j++
	}
	return newArr
}

var key []uint8
var lowerLimit int64 = 0xffffff
var masterKey_1 uint32 = 0x0
var masterKey_2 uint32 = 0x0
var rom []uint16
var upperLimit int64 = 0xff0000
var encDir Direction

func Crypt(direction Direction, romDef cps2rom.RomDefinition, romZip *zip.ReadCloser, romBinary []uint8) ([]uint8, error) {
	encDir = direction
	err := parseKey(romDef, romZip)
	if err != nil {
		return nil, err
	}
	rom = createUint16ArrayFromUint8Array(romBinary)

	key1 := make([]uint32, 4)
	dec := make([]uint16, len(rom))
	sboxes10 := make([]optimized_sbox, 4)
	sboxes11 := make([]optimized_sbox, 4)
	sboxes12 := make([]optimized_sbox, 4)
	sboxes13 := make([]optimized_sbox, 4)
	sboxes20 := make([]optimized_sbox, 4)
	sboxes21 := make([]optimized_sbox, 4)
	sboxes22 := make([]optimized_sbox, 4)
	sboxes23 := make([]optimized_sbox, 4)
	sboxes := [][]optimized_sbox{
		sboxes10, sboxes11, sboxes12, sboxes13,
		sboxes20, sboxes21, sboxes22, sboxes23,
	}
	for _, sbox_set := range sboxes {
		initializeOptimizedSBoxes(sbox_set)
	}
	length := len(rom)
	optimizeSBoxes(sboxes10, fn1_r1_boxes)
	optimizeSBoxes(sboxes11, fn1_r2_boxes)
	optimizeSBoxes(sboxes12, fn1_r3_boxes)
	optimizeSBoxes(sboxes13, fn1_r4_boxes)
	optimizeSBoxes(sboxes20, fn2_r1_boxes)
	optimizeSBoxes(sboxes21, fn2_r2_boxes)
	optimizeSBoxes(sboxes22, fn2_r3_boxes)
	optimizeSBoxes(sboxes23, fn2_r4_boxes)
	masterKey := []uint32{masterKey_1, masterKey_2}
	expandKey(0, &key1, masterKey)

	key1[0] ^= bit32(key1[0], 1) << 4
	key1[0] ^= bit32(key1[0], 2) << 5
	key1[0] ^= bit32(key1[0], 8) << 11
	key1[1] ^= bit32(key1[1], 0) << 5
	key1[1] ^= bit32(key1[1], 8) << 11
	key1[2] ^= bit32(key1[2], 1) << 5
	key1[2] ^= bit32(key1[2], 8) << 11

	for i := 0; i <= 0xffff; i++ {
		subkey := make([]uint32, 2)
		key2 := make([]uint32, 4)

		seed := feistel(uint16(i&0xffff), fn1_groupA, fn1_groupB, sboxes10, sboxes11, sboxes12, sboxes13, key1[0], key1[1], key1[2], key1[3])
		// fmt.Printf("seed = %04x\n", seed)
		expandSubkey(&subkey, seed)

		subkey[0] ^= masterKey_1
		subkey[1] ^= masterKey_2

		expandKey(1, &key2, subkey)

		key2[0] ^= bit32(key2[0], 0) << 5
		key2[0] ^= bit32(key2[0], 6) << 11
		key2[1] ^= bit32(key2[1], 0) << 5
		key2[1] ^= bit32(key2[1], 1) << 4
		key2[2] ^= bit32(key2[2], 2) << 5
		key2[2] ^= bit32(key2[2], 3) << 4
		key2[2] ^= bit32(key2[2], 7) << 11
		key2[3] ^= bit32(key2[3], 1) << 5

		for a := i; a < length; a += 0x10000 {
			if int64(a) >= lowerLimit && int64(a) <= upperLimit {
				if direction { // decrypt
					dec[a] = feistel(rom[a], fn2_groupA, fn2_groupB, sboxes20, sboxes21, sboxes22, sboxes23, key2[0], key2[1], key2[2], key2[3])
				} else { // encrypt
					dec[a] = feistel(rom[a], fn2_groupA, fn2_groupB, sboxes23, sboxes22, sboxes21, sboxes20, key2[3], key2[2], key2[1], key2[0])
				}
			} else {
				dec[a] = rom[a]
			}
		}

	}
	return createUint8ArrayFromUint16Array(dec), nil
}
