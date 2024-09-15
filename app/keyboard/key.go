package keyboard

import (
	"github.com/kindlyfire/go-keylogger"
)

var Last string = ""

// работает почти как надо ( не handle  control alt win  и подобные также home end backspace F1-f12 )  детектит буквы  стрелки
func Init() {
	go func() {
		kl := keylogger.NewKeylogger()
		for {
			key := kl.GetKey()

			if !key.Empty {
				Last = string(key.Rune)
			}
		}

	}()
}
func Lastpressendkey() string {
	return Last
}
