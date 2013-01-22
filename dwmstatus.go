// Package dwmstatus is a Go implementation of
// http://dwm.suckless.org/dwmstatus/, a utility to set the title of the X root
// window, which dwm uses to get the content of the status bar.
//
// This package does not provide any content for the status bar. You can pass
// your own GenTitleFunc to generate any content you like. dwmstatus/main.go
// contains an example.
//
// Thomas Kappler <tkappler@gmail.com>
package dwmstatus

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"log"
	"time"
)

// Pass a function with this signature to Run. It will be called repeatedly and
// whatever you write to Buffer b will be the new value of the status bar. The
// current time is provided as a convenience. b will be empty every time.
type GenTitleFunc func(now time.Time, b *bytes.Buffer)

func setWindowTitle(title []byte, X *xgb.Conn, window xproto.Window) {
	xproto.ChangeProperty(X, xproto.PropModeReplace, window, xproto.AtomWmName,
		xproto.AtomString, byte(8), uint32(len(title)), title)
}

func setStatus(status []byte, X *xgb.Conn) {
	screen := xproto.Setup(X).DefaultScreen(X)
	fmt.Println(string(status))
	setWindowTitle(status, X, screen.Root)
}

// Start the process of updating the status bar. genTitle will be called
// repeatedly in the given interval.
func Run(interval time.Duration, genTitle GenTitleFunc) {
	X, err := xgb.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	var status bytes.Buffer

	c := time.Tick(interval)
	for now := range c {
		genTitle(now, &status)
		setStatus(status.Bytes(), X)
		status.Reset()
	}
}
