package assets

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
)

var (
	iconBorder = color.RGBA{R: 0x82, G: 0xab, B: 0xf8, A: 0xff}
	iconFill   = color.RGBA{R: 0x20, G: 0x21, B: 0x26, A: 0xff}
	iconDot    = color.RGBA{R: 0x9e, G: 0xce, B: 0x6a, A: 0xff}
)

func drawIcon(size int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	b := size / 16
	if b < 1 {
		b = 1
	}
	dotMin, dotMax := size*2/5, size*3/5
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			c := iconFill
			switch {
			case x < b || y < b || x >= size-b || y >= size-b:
				c = iconBorder
			case x >= dotMin && x <= dotMax && y >= dotMin && y <= dotMax:
				c = iconDot
			}
			img.SetRGBA(x, y, c)
		}
	}
	return img
}
func IconImage(size int) *image.RGBA {
	return drawIcon(size)
}
func IconPNG(size int) []byte {
	var buf bytes.Buffer
	if err := png.Encode(&buf, drawIcon(size)); err != nil {
		return nil
	}
	return buf.Bytes()
}
func TrayIcon() []byte {
	return encodeICO([]int{16})
}
func IconICO() []byte {
	return encodeICO([]int{16, 24, 32, 48, 64, 128, 256})
}
func encodeICO(sizes []int) []byte {
	imgs := make([][]byte, 0, len(sizes))
	offset := 6 + 16*len(sizes)
	for _, size := range sizes {
		data := IconPNG(size)
		if data == nil {
			return nil
		}
		imgs = append(imgs, data)
		offset += len(data)
	}
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, uint16(0))
	binary.Write(&buf, binary.LittleEndian, uint16(1))
	binary.Write(&buf, binary.LittleEndian, uint16(len(sizes)))
	cur := 6 + 16*len(sizes)
	for i, size := range sizes {
		b := size
		if b >= 256 {
			b = 0
		}
		buf.WriteByte(byte(b))
		buf.WriteByte(byte(b))
		buf.WriteByte(0)
		buf.WriteByte(0)
		binary.Write(&buf, binary.LittleEndian, uint16(1))
		binary.Write(&buf, binary.LittleEndian, uint16(32))
		binary.Write(&buf, binary.LittleEndian, uint32(len(imgs[i])))
		binary.Write(&buf, binary.LittleEndian, uint32(cur))
		cur += len(imgs[i])
	}
	for _, data := range imgs {
		buf.Write(data)
	}
	return buf.Bytes()
}
