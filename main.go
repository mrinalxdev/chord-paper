package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	// "time"
	// "os"

	"github.com/mrinalxdev/chord/pkg/chord"
	"github.com/mrinalxdev/chord/pkg/storage"
	"github.com/mrinalxdev/chord/pkg/visualization"
)

func main() {
    addr := flag.String("addr", ":8080", "HTTP service address")
    dataDir := flag.String("data", "data", "Data directory for BadgerDB")
    introducer := flag.String("introducer", "", "Address of an existing node")
    flag.Parse()

    store, err := storage.NewBadgerStore(filepath.Join(*dataDir, "badger"))
    if err != nil {
        log.Fatal(err)
    }
    defer store.Close()
    node := chord.NewNode(*addr, store)
    
    vis := visualization.NewVisualizer()
    go vis.Run()

    var introducerNode *chord.Node
    if *introducer != "" {
        introducerNode = &chord.Node{Address: *introducer}
    }
    
    if err := node.Join(introducerNode); err != nil {
        log.Fatal(err)
    }

    
    go func() {
        for {
            node.Stabilize()
            node.FixFingers()
            vis.UpdateNetwork([]visualization.NodeInfo{
                {
                    ID:        node.ID.String(),
                    Address:   node.Address,
                    Successor: node.GetSuccessor().Address,
                    Fingers:   make([]string, 0), // Add finger table info
                },
            })
        }
    }()

    http.HandleFunc("/ws", vis.HandleConnections)
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })
    http.Handle("/", http.FileServer(http.Dir("web")))
    
    log.Printf("Starting Chord node at %s\n", *addr)
    log.Fatal(http.ListenAndServe(*addr, nil))
}