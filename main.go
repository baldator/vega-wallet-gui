package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	runtimedebug "runtime/debug"
	"strconv"
	"time"

	"code.vegaprotocol.io/go-wallet/fsutil"
	"code.vegaprotocol.io/go-wallet/wallet"
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	bootstrap "github.com/asticode/go-astilectron-bootstrap"
)

// Constants
const version = "0.1.0"
const licence = `Copyright 2021 Marco Torello<br>
<br>
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:<br>
<br>
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.<br>
<br>
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.`

var rootPath string
var walletPath string
var l *log.Logger

// Vars injected via ldflags by bundler
var (
	AppName            string
	BuiltAt            string
	VersionAstilectron string
	VersionElectron    string
)

// Application Vars
var (
	fs    = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	debug = fs.Bool("d", true, "enables the debug mode")
	w     *astilectron.Window
)

func main() {
	rootPath = fsutil.DefaultVegaDir()
	walletPath = rootPath + "/wallets/"

	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// Create logger
	l = log.New(log.Writer(), log.Prefix(), log.Flags())
	l.SetOutput(file)

	// Parse flags
	fs.Parse(os.Args[1:])

	// Run bootstrap
	l.Printf("Running app built at %s\n", BuiltAt)
	if err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDarwinPath:  "resources/icon.icns",
			AppIconDefaultPath: "resources/icon.png",
			SingleInstance:     true,
			VersionAstilectron: VersionAstilectron,
			VersionElectron:    VersionElectron,
		},
		Debug:  *debug,
		Logger: l,
		MenuOptions: []*astilectron.MenuItemOptions{{
			Label: astikit.StrPtr("File"),
			SubMenu: []*astilectron.MenuItemOptions{
				{
					Label: astikit.StrPtr("About"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						htmlAbout := `Vega Wallet GUI version ` + version + `<br>
							Vega wallet version ` + getVersion() + `<br><br>
							Copyright ` + strconv.Itoa(time.Now().UTC().Year()) + ` Marco Torello under MIT licence.`
						if err := bootstrap.SendMessage(w, "about", htmlAbout, func(m *bootstrap.MessageIn) {
							// Unmarshal payload
							var s string
							if err := json.Unmarshal(m.Payload, &s); err != nil {
								l.Println(fmt.Errorf("unmarshaling payload failed: %w", err))
								return
							}
							l.Printf("About modal has been displayed and payload is %s!\n", s)
						}); err != nil {
							l.Println(fmt.Errorf("sending about event failed: %w", err))
						}
						return
					},
				},
				{
					Label: astikit.StrPtr("Licence"),
					OnClick: func(e astilectron.Event) (deleteListener bool) {
						if err := bootstrap.SendMessage(w, "licence", licence, func(m *bootstrap.MessageIn) {
							// Unmarshal payload
							var s string
							if err := json.Unmarshal(m.Payload, &s); err != nil {
								l.Println(fmt.Errorf("unmarshaling payload failed: %w", err))
								return
							}
							l.Printf("Licence modal has been displayed and payload is %s!\n", s)
						}); err != nil {
							l.Println(fmt.Errorf("sending licence event failed: %w", err))
						}
						return
					},
				},
				{Role: astilectron.MenuItemRoleClose},
			},
		}},
		OnWait: func(_ *astilectron.Astilectron, ws []*astilectron.Window, _ *astilectron.Menu, _ *astilectron.Tray, _ *astilectron.Menu) error {
			w = ws[0]
			go func() {
				time.Sleep(5 * time.Second)
				//if err := bootstrap.SendMessage(w, "check.out.menu", "Don't forget to check out the menu!"); err != nil {
				//	l.Println(fmt.Errorf("sending check.out.menu event failed: %w", err))
				//}
			}()
			return nil
		},
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astikit.StrPtr("#333"),
				Center:          astikit.BoolPtr(true),
				Height:          astikit.IntPtr(700),
				Width:           astikit.IntPtr(700),
			},
		}},
	}); err != nil {
		l.Fatal(fmt.Errorf("running bootstrap failed: %w", err))
	}
}

func getVersion() string {
	var wallet wallet.Wallet
	wallet.Owner = ""
	bi, ok := runtimedebug.ReadBuildInfo()
	if !ok {
		log.Printf("Failed to read build info")
		return ""
	}

	for _, dep := range bi.Deps {
		l.Printf("dependency: %+v\n", dep)
		if dep.Path == "code.vegaprotocol.io/go-wallet" {
			return dep.Version
		}

	}
	return ""
}
