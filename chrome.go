package chromesvc

import (
	"encoding/json"
	"log"
	"time"

	"github.com/wirepair/gcd"
	"github.com/wirepair/gcd/gcdapi"
)

const (
	path = "/usr/bin/chromium-browser"
	dir  = "/tmp/"
	port = "9222"
)

type PageRenderTask struct {
	URL   string
	Proxy string
}

type PageResult struct {
	HTML string
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func startDebugger() *gcd.Gcd {
	defer timeTrack(time.Now(), "start debugger")
	debugger := gcd.NewChromeDebugger()
	debugger.StartProcess(path, dir, port)
	return debugger
}

func loadURL(tab *gcd.ChromeTarget, u string) {
	defer timeTrack(time.Now(), "load URL")
	tab.Page.Enable()
	tab.Network.EnableWithParams(&gcdapi.NetworkEnableParams{})
	wc := make(chan struct{})

	tab.Subscribe("Network.requestWillBeSent", func(target *gcd.ChromeTarget, event []byte) {
		ev := &gcdapi.NetworkRequestWillBeSentEvent{}
		json.Unmarshal(event, ev)
		// log.Print(ev.Params.Request.Url)
	})

	tab.Subscribe("Page.domContentEventFired", func(target *gcd.ChromeTarget, event []byte) {
		log.Print("DOM content event")
	})

	tab.Subscribe("Page.loadEventFired", func(target *gcd.ChromeTarget, event []byte) {
		log.Print("Page load event")
		wc <- struct{}{}
	})

	tab.Page.Navigate(u, "", "")
	<-wc
	close(wc)
}

func grabHTML(debugger *gcd.Gcd, target *gcd.ChromeTarget) (string, error) {
	defer timeTrack(time.Now(), "grab HTML")
	dom := target.DOM
	doc, err := dom.GetDocument(-1, true)
	if err != nil {
		return "", err
	}
	return dom.GetOuterHTMLWithParams(&gcdapi.DOMGetOuterHTMLParams{NodeId: doc.NodeId})
}

// RenderPage renders the page HTML
func RenderPage(t *PageRenderTask) (*PageResult, error) {
	defer timeTrack(time.Now(), "render page total")
	debugger := startDebugger()
	defer debugger.ExitProcess()

	tab, err := debugger.NewTab()
	if err != nil {
		return nil, err
	}
	loadURL(tab, t.URL)
	html, err := grabHTML(debugger, tab)
	if err != nil {
		return nil, err
	}

	return &PageResult{HTML: html}, nil
}
