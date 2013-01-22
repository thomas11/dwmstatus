package dwmstatus

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"log"
	"time"
)

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
