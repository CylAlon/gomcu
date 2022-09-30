package mcu

import (
	"machine"
)

// GpioInit 批量初始化普通GPIO(弥补目前的额tinygo无法批量初始化GPIO的问题)
// 用法：
// led := machine.PA1
// led2 := machine.PA2
// Utils.GpioInit([]machine.Pin{led, led2}, machine.PinOutput)
// led.Low()
// led2.High()
func GpioInit(pin []machine.Pin, mode machine.PinMode) {
	for _, p := range pin {
		p.Configure(machine.PinConfig{Mode: mode})
	}
}
