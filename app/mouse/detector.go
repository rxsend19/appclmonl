package mouse

import (
	"strconv"
	"strings"

	"github.com/abdfnx/gosh"
)

func Run() (int, int) {
	_, xs, _ := gosh.RunOutput(`Add-Type -Path 'C:\Windows\Microsoft.NET\Framework64\v4.0.30319\System.Windows.Forms.dll';[System.Windows.Forms.Cursor]::Position.x `)
	_, ys, _ := gosh.RunOutput(`Add-Type -Path 'C:\Windows\Microsoft.NET\Framework64\v4.0.30319\System.Windows.Forms.dll';[System.Windows.Forms.Cursor]::Position.y`)
	x, _ := strconv.Atoi(strings.TrimSpace(xs))
	y, _ := strconv.Atoi(strings.TrimSpace(ys))

	return x, y
}
