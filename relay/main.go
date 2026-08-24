package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

//go:embed partyweb
var partyWebFS embed.FS

const roomIdleTimeout = 30 * time.Minute

// roomRegistry holds one partyServer per room code, created lazily on first
// join and swept once idle so memory doesn't grow unbounded on a long-lived
// free-tier instance.
type roomRegistry struct {
	mu    sync.Mutex
	rooms map[string]*partyServer
}

func newRoomRegistry() *roomRegistry {
	return &roomRegistry{rooms: make(map[string]*partyServer)}
}

func (r *roomRegistry) get(code string) *partyServer {
	r.mu.Lock()
	defer r.mu.Unlock()
	ps, ok := r.rooms[code]
	if !ok {
		ps = newPartyServer()
		r.rooms[code] = ps
	}
	return ps
}

func (r *roomRegistry) sweep() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for code, ps := range r.rooms {
		if ps.playerCount() == 0 && time.Since(ps.idleSince()) > roomIdleTimeout {
			delete(r.rooms, code)
		}
	}
}

func main() {
	registry := newRoomRegistry()

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			registry.sweep()
		}
	}()

	sub, err := fs.Sub(partyWebFS, "partyweb")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(sub)))

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		room := r.URL.Query().Get("room")
		if room == "" {
			http.Error(w, "missing room", http.StatusBadRequest)
			return
		}
		conn, err := partyUpgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		registry.get(room).handleWS(conn)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("merge-kingdom-relay listening on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
