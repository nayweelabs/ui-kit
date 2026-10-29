// Command ui-kit-preview serves the kit's gallery on localhost: every part in every state, light and dark.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/ElderlyGate/ui-kit/preview"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:7790", "address to listen on (localhost only)")
	flag.Parse()
	preview.Version = version()
	srv := &http.Server{Addr: *addr, Handler: preview.Handler(), ReadHeaderTimeout: 5 * time.Second}
	fmt.Printf("ui-kit preview on http://%s\n", *addr)
	log.Fatal(srv.ListenAndServe())
}

func version() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, s := range info.Settings {
			if s.Key == "vcs.revision" && len(s.Value) >= 7 {
				return "dev-" + s.Value[:7]
			}
		}
	}
	return "dev"
}
