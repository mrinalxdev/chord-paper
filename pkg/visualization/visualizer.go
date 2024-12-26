package visualization

import (
    "encoding/json"
    "net/http"
    "sync"
    
    "github.com/gorilla/websocket"
)

type Visualizer struct {
    clients    map[*websocket.Conn]bool
    broadcast  chan []byte
    register   chan *websocket.Conn
    unregister chan *websocket.Conn
    mutex      sync.Mutex
}

type NodeInfo struct {
    ID        string   `json:"id"`
    Address   string   `json:"address"`
    Successor string   `json:"successor"`
    Fingers   []string `json:"fingers"`
}

func NewVisualizer() *Visualizer {
    return &Visualizer{
        clients:    make(map[*websocket.Conn]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *websocket.Conn),
        unregister: make(chan *websocket.Conn),
    }
}

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func (v *Visualizer) HandleConnections(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    defer ws.Close()

    v.register <- ws

    for {
        _, _, err := ws.ReadMessage()
        if err != nil {
            v.unregister <- ws
            break
        }
    }
}

func (v *Visualizer) Run() {
    for {
        select {
        case client := <-v.register:
            v.mutex.Lock()
            v.clients[client] = true
            v.mutex.Unlock()
        case client := <-v.unregister:
            v.mutex.Lock()
            delete(v.clients, client)
            v.mutex.Unlock()
        case message := <-v.broadcast:
            v.mutex.Lock()
            for client := range v.clients {
                err := client.WriteMessage(websocket.TextMessage, message)
                if err != nil {
                    client.Close()
                    delete(v.clients, client)
                }
            }
            v.mutex.Unlock()
        }
    }
}

func (v *Visualizer) UpdateNetwork(nodes []NodeInfo) {
    data, err := json.Marshal(nodes)
    if err != nil {
        return
    }
    v.broadcast <- data
}