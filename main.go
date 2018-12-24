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

	"github.com/zserge/lorca"

	assetfs "github.com/elazarl/go-bindata-assetfs"
)

var ui lorca.UI
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
		ui.Eval("window.eventEmitter.emit('endSound','" + file + "')")
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

func main() {
	var err error
	ui, err = lorca.New("", "", 480, 320)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go func() {
		// load in bindata
		http.Handle("/", http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "/frontend/build"}))
		log.Fatal(http.Serve(ln, nil))
	}()
	ui.Load(fmt.Sprintf("http://%s", ln.Addr()))

	log.Println("DOM bind")

	ui.Bind("play", controller.Play)
	ui.Bind("pause", controller.Pause)
	ui.Bind("cancel", controller.Cancel)
	ui.Bind("resume", controller.Resume)
	ui.Bind("getAllPanels", controller.GetFromDisk)
	ui.Eval(`
			window.panelController = {
				play,
				pause,
				cancel,
				resume,
				getAllPanels
			}
		`)
	log.Println(err)

	// Wait for the browser window to be closed
	<-ui.Done()
}
