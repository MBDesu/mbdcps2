package utils

func ConvertUintToByteSlice(data uint32, size int) []byte {
	switch size {
	case 2:
		return ConvertUint16ToByteSlice(uint16(data & 0xffff))
	case 3:
		return ConvertUint24ToByteSlice(data)
	case 4:
		return ConvertUint32ToByteSlice(data)
	}
	return nil
}

func ConvertUint24ToByteSlice(x uint32) (bytes []byte) {
	bytes = []byte{
		byte(x & 0xff0000 >> 16),
		byte(x & 0x00ff00 >> 8),
		byte(x & 0xff),
	}
	return
}

func ConvertUint32ToByteSlice(x uint32) (bytes []byte) {
	bytes = make([]byte, 4)
	for i := 3; i >= 0; i-- {
		bytes[3-i] = byte((x >> (i * 8)) & 0xff)
	}
	return
}

func ConvertUint16ToByteSlice(x uint16) (bytes []byte) {
	bytes = []byte{
		byte(x & 0xff00 >> 8),
		byte(x & 0xff),
	}
	return
}

func CreateUint8ArrayFromUint16Array(arr []uint16) []uint8 {
	newArr := make([]uint8, len(arr)*2)
	for i := 0; i < len(arr); i++ {
		val := uint8((arr[i] & 0xff00) >> 8)
		newArr[i*2] = val
		val = uint8(arr[i] & 0xff)
		newArr[i*2+1] = val
	}

	return newArr
}

func CreateUint16ArrayFromUint8Array(arr []uint8) []uint16 {
	length := len(arr)
	newArr := make([]uint16, length/2)
	i := 0
	j := 0
	for i < length {
		val := uint16(arr[i+1])
		val |= uint16(arr[i]) << 8
		newArr[j] = val
		i += 2
		j++
	}
	return newArr
}
