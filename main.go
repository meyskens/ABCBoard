package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/zserge/webview"
)

var w webview.WebView
var controller = PanelController{Panels: []Panel{}, playCancelers: map[string]context.CancelFunc{}, playPausers: map[string]*sync.Mutex{}}

// Panel is one panel to be shown
type Panel struct {
	Name     string `json:"name"`
	Shortcut string `json:"shortcut"`
	File     string `json:"file"`
}

// PanelController helps with controling the pannels
type PanelController struct {
	Panels        []Panel `json:"panels"`
	panelsForFile map[string]*Panel
	playCancelers map[string]context.CancelFunc
	playPausers   map[string]*sync.Mutex
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

	p.panelsForFile = map[string]*Panel{}
	for i := range p.Panels {
		p.panelsForFile[p.Panels[i].File] = &p.Panels[i]
	}
	return p.Panels
}

// Play starts playing a specific file
func (p *PanelController) Play(file string) {
	fmt.Println(file)
	ctx, cancel := context.WithCancel(context.Background())
	p.playCancelers[file] = cancel
	pause := sync.Mutex{}
	p.playPausers[file] = &pause
	go func() {
		playMP3(ctx, &pause, file)
		cancel()
		p.playCancelers[file] = nil
		w.Dispatch(func() {
			w.Eval("window.eventEmitter.emit('endSound','" + file + "')")
		})
	}()
}

// Cancel stops playing a specific file
func (p *PanelController) Cancel(file string) {
	if cancel := p.playCancelers[file]; cancel != nil {
		cancel()
	}
}

// Pause pauses playing a specific file
func (p *PanelController) Pause(file string) {
	if mutex := p.playPausers[file]; mutex != nil {
		mutex.Lock()
	}
}

// Resume pauses playing a specific file
func (p *PanelController) Resume(file string) {
	if mutex := p.playPausers[file]; mutex != nil {
		mutex.Unlock()
	}
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
	w = webview.New(webview.Settings{Debug: true, Title: "ABCBoard", Width: 800, Height: 600, Resizable: true, URL: "http://" + ln.Addr().String()})
	w.Dispatch(func() {
		w.Bind("panelController", &controller)
		w.Eval("window.panelController = panelController")
	})
	w.Run()
}
