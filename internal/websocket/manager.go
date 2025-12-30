package websocket

import (
	"log"
	"sync"
)

type Manager struct {
	clients    map[*Client]bool
	symbols    map[string]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
	mu         sync.RWMutex
}

type Message struct {
	Symbol string      `json:"symbol"`
	Type   string      `json:"type"` // "price", "anomaly", "error"
	Data   interface{} `json:"data"`
}

func NewManager() *Manager {
	return &Manager{
		clients:    make(map[*Client]bool),
		symbols:    make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
	}
}

func (m *Manager) Run() {
	for {
		select {
		case client := <-m.register:
			m.mu.Lock()
			m.clients[client] = true
			m.mu.Unlock()
			log.Println("New client registered")

		case client := <-m.unregister:
			m.mu.Lock()
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				// Clean up subscriptions
				for symbol := range m.symbols {
					delete(m.symbols[symbol], client)
				}
				close(client.send)
			}
			m.mu.Unlock()
			log.Println("Client unregistered")

		case message := <-m.broadcast:
			m.mu.RLock()
			if subscribers, ok := m.symbols[message.Symbol]; ok {
				for client := range subscribers {
					select {
					case client.send <- message:
					default:
						// Handle slow client
						log.Printf("Slow client detected, dropping message for %s", message.Symbol)
					}
				}
			}
			m.mu.RUnlock()
		}
	}
}

func (m *Manager) Subscribe(client *Client, symbol string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.symbols[symbol]; !ok {
		m.symbols[symbol] = make(map[*Client]bool)
	}
	m.symbols[symbol][client] = true
	log.Printf("Client subscribed to %s", symbol)
}

func (m *Manager) Unsubscribe(client *Client, symbol string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if subs, ok := m.symbols[symbol]; ok {
		delete(subs, client)
		if len(subs) == 0 {
			delete(m.symbols, symbol)
		}
	}
	log.Printf("Client unsubscribed from %s", symbol)
}

func (m *Manager) Broadcast(msg Message) {
	m.broadcast <- msg
}

func (m *Manager) GetSubscribedSymbols() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	symbols := make([]string, 0, len(m.symbols))
	for s := range m.symbols {
		symbols = append(symbols, s)
	}
	return symbols
}
