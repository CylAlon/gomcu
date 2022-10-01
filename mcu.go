package gomcu

import (
	"device/stm32"
	"machine"
	"time"
)

/********************************************外部中断********************************************/
type ExitConfig struct {
	Pin    machine.Pin
	Mode   machine.PinMode
	Change machine.PinChange
}

func NewExit(config ExitConfig, callback func()) *ExitConfig {
	exit := ExitConfig{
		Pin:    config.Pin,
		Mode:   config.Mode,
		Change: config.Change,
	}
	exit.Pin.Configure(machine.PinConfig{Mode: exit.Mode})
	exit.Pin.SetInterrupt(exit.Change, func(pin machine.Pin) {
		exit.Iqr(false) // 关闭中断
		go func() {
			callback()
			exit.Iqr(true) // 开启中断
		}()
	})
	return &exit
}
func (s *ExitConfig) Iqr(flag bool) {
	if flag {
		stm32.EXTI.IMR.Set(1 << s.Pin) // 开启中断
	} else {
		stm32.EXTI.IMR.ClearBits(1 << s.Pin) // 关闭中断
	}
}

// Jitter 消抖
func Jitter(pin machine.Pin, hight bool, callback func()) {
	lambda := func(pin machine.Pin) bool {
		return pin.Get() && hight || pin.Get() && !hight
	}
	if lambda(pin) {
		time.Sleep(time.Millisecond * 150)
		if lambda(pin) {
			for lambda(pin) {
				time.Sleep(time.Millisecond * 10)
			}
			callback()
			time.Sleep(time.Millisecond * 100)
		}
	}
}
