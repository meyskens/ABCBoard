package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/zserge/webview"
)

var controller = PanelController{Panels: []Panel{}}

// Panel is one panel to be shown
type Panel struct {
	Name     string `json:"name"`
	Shortcut string `json:"shortcut"`
	File     string `json:"file"`
}

// PanelController helps with controling the pannels
type PanelController struct {
	Panels []Panel `json:"panels"`
}

// GetFromDisk loads the panels from the config file
func (p *PanelController) GetFromDisk() []Panel {
	f, err := os.Open("./panels.json")
	jsonParser := json.NewDecoder(f)
	jsonParser.Decode(&p.Panels)
	if err != nil {
		fmt.Println(err)
		return p.Panels
	}
	return p.Panels
}

// Play starts playing a specific file
func (p *PanelController) Play(file string) {
	fmt.Println(file)
	go playMP3(file)
}

func handleAPIPanels(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	out, _ := json.Marshal(controller.GetFromDisk())
	w.Write(out)
}

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go func() {
		// load in bindata
		http.Handle("/api/panels", http.HandlerFunc(handleAPIPanels))
		http.Handle("/", http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "/frontend/build"}))
		log.Fatal(http.Serve(ln, nil))
	}()
	w := webview.New(webview.Settings{Debug: true, Title: "ABCBoard", Width: 800, Height: 600, Resizable: true, URL: "http://" + ln.Addr().String()})
	w.Dispatch(func() {
		w.Bind("panelController", &controller)
		w.Eval("window.panelController = panelController")
	})
	w.Run()
}
