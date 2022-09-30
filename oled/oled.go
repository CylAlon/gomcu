package oled

import (
	"machine"
	"time"
)

func init() {

	CsPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	DcPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	RstPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	SckPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	SdaPin.Configure(machine.PinConfig{Mode: machine.PinOutput})

	RstPin.Low()
	time.Sleep(100 * time.Millisecond)
	RstPin.High()
	RegisterWR(0xAE, WriteCmd) // --turn off oled panel
	RegisterWR(0x00, WriteCmd) // ---set low column address
	RegisterWR(0x10, WriteCmd) // ---set high column address
	RegisterWR(0x40, WriteCmd) // --set start line address  Set Mapping RAM Display Start Line (0x00~0x3F)
	RegisterWR(0x81, WriteCmd) // --set contrast control register
	RegisterWR(0xCF, WriteCmd) // Set SEG Output Current Brightness
	RegisterWR(0xA1, WriteCmd) // --Set SEG/Column Mapping     0xa0   ҷ    0xa1
	RegisterWR(0xC8, WriteCmd) // Set COM/Row Scan Direction   0xc0   ·    0xc8
	RegisterWR(0xA6, WriteCmd) // --set normal display
	RegisterWR(0xA8, WriteCmd) // --set multiplex ratio(1 to 64)
	RegisterWR(0x3f, WriteCmd) // --1/64 duty
	RegisterWR(0xD3, WriteCmd) // -set display offset	Shift Mapping RAM Counter (0x00~0x3F)
	RegisterWR(0x00, WriteCmd) // -not offset
	RegisterWR(0xd5, WriteCmd) // --set display clock divide ratio/oscillator frequency
	RegisterWR(0x80, WriteCmd) // --set divide ratio, Set Clock as 100 Frames/Sec
	RegisterWR(0xD9, WriteCmd) // --set pre-charge period
	RegisterWR(0xF1, WriteCmd) // Set Pre-Charge as 15 Clocks & Discharge as 1 Clock
	RegisterWR(0xDA, WriteCmd) // --set com pins hardware configuration
	RegisterWR(0x12, WriteCmd)
	RegisterWR(0xDB, WriteCmd) // --set vcomh
	RegisterWR(0x40, WriteCmd) // Set VCOM Deselect Level
	RegisterWR(0x20, WriteCmd) // -Set Page Addressing Mode (0x00/0x01/0x02)
	RegisterWR(0x02, WriteCmd) //
	RegisterWR(0x8D, WriteCmd) // --set Charge Pump enable/disable
	RegisterWR(0x14, WriteCmd) // --set(0x10) disable
	RegisterWR(0xA4, WriteCmd) // Disable Entire Display On (0xa4/0xa5)
	RegisterWR(0xA6, WriteCmd) // Disable Inverse Display On (0xa6/a7)
	RegisterWR(0xAF, WriteCmd) // --turn on oled panel
	//
	RegisterWR(0xAF, WriteCmd) /*display ON*/
	Clear()
	SetPos(0, 0)
	Clear()
}

func RegisterWR(dat byte, cmd uint8) {
	// 硬件spi
	// OLED_CS_L;
	// cmd ? OLED_DC_H : OLED_DC_L; // H是数据模式，L是命令模式
	// // HAL_SPI_Transmit(&hspi1, &dat, 1, 10);
	// SPI_Write(&dat,1);
	// OLED_CS_H;
	// 	软件spi

	if cmd != 0 {
		DcPin.High()
	} else {
		DcPin.Low()
	}
	CsPin.Low()
	for i := 0; i < 8; i++ {
		SckPin.Low()
		if dat&0x80 != 0 {
			SdaPin.High()
		} else {
			SdaPin.Low()
		}
		dat <<= 1
		SckPin.High()
	}
	CsPin.High()
	DcPin.High()
}

func SetPos(x, y uint8) {
	RegisterWR(0xb0+y, WriteCmd)
	RegisterWR(((x&0xf0)>>4)|0x10, WriteCmd)
	RegisterWR((x&0x0f)|0x01, WriteCmd)
}
func DisplayOn() {
	RegisterWR(0x8D, WriteCmd) // SET DCDC命令
	RegisterWR(0x14, WriteCmd) // DCDC ON
	RegisterWR(0xAF, WriteCmd) // DISPLAY ON
}
func DisplayOff() {
	RegisterWR(0x8D, WriteCmd) // SET DCDC命令
	RegisterWR(0x10, WriteCmd) // DCDC OFF
	RegisterWR(0xAE, WriteCmd) // DISPLAY OFF
}
func Clear() {
	var i, n uint8
	for i = 0; i < 8; i++ {
		RegisterWR(0xb0+i, WriteCmd) // 设置页地址（0~7）
		RegisterWR(0x00, WriteCmd)   // 设置显示位置—列低地址
		RegisterWR(0x10, WriteCmd)   // 设置显示位置—列高地址
		for n = 0; n < 128; n++ {
			RegisterWR(0, WriteData)
		}
	} // 更新显示
}

// 在指定位置显示一个字符,包括部分字符
// x:0~127
// y:0~63
// size:选择字体 16/12

func ShowChar(x, y uint8, chr byte, size bool) {
	c := chr - ' '
	if x > MaxColumn-1 {
		x = 0
		y = y + 2
	}
	SetPos(x, y)
	if size {
		for i := 0; i < 8; i++ {
			RegisterWR(CHAR8x16[c][i], WriteData)
		}
		SetPos(x, y+1)
		for i := 0; i < 8; i++ {
			RegisterWR(CHAR8x16[c][8+i], WriteData)
		}
	} else {
		for i := 0; i < 6; i++ {
			RegisterWR(CHAR6x8[c][i], WriteData)
		}
	}
}
func ShowChineseOnce(x, y uint8, str string) {
	lens := len(Che16)
	for i := 0; i < lens; i++ {
		if Che16[i].Word == str {
			SetPos(x, y)
			for j := 0; j < 16; j++ {
				RegisterWR(Che16[i].Buff[j], WriteData)
			}
			SetPos(x, y+1)
			for j := 0; j < 16; j++ {
				RegisterWR(Che16[i].Buff[j+16], WriteData)
			}
			break
		}
	}
}

func PrintOLed(x, y uint8, str string) {
	ix, iy := x, y
	for i := 0; i < len(str); i++ {
		if str[i] >= 0x80 {
			ShowChineseOnce(ix, iy, str[i:i+3])
			ix = ix + 16
			i += 2
		} else {
			ShowChar(ix, iy, str[i], true)
			ix = ix + 8
		}
		if ix >= XWidth {
			ix = 0
			iy = iy + 2
		}
	}
}

// func DrawBMP(x0, y0, x1, y1 uint8, bmp []byte) {
// 	var j, b, x, y uint8
// 	if x1%8 == 0 {
// 		b = x1 / 8
// 	} else {
// 		b = x1/8 + 1
// 	}
// 	for y = y0; y < y1; y++ {
// 		SetPos(x0, y)
// 		for x = x0; x < b; x++ {
// 			RegisterWR(bmp[j], WriteData)
// 			j++
// 		}
// 	}
// }
