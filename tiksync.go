package tiksync

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const DefaultAPIURL = "https://api.tik-sync.com"

type EventHandler func(data map[string]interface{})

type Config struct {
	APIKey               string
	APIURL               string
	MaxReconnectAttempts int
	ReconnectDelay       time.Duration
}

type TikSync struct {
	uniqueID  string
	config    Config
	conn      *websocket.Conn
	listeners map[string][]EventHandler
	mu        sync.RWMutex
	running   bool
	reconnects int
	security   *SecurityCore
	dynNative  *dynamicNativeCore
}

func New(uniqueID string, apiKey string) *TikSync {
	return NewWithConfig(uniqueID, Config{
		APIKey:               apiKey,
		APIURL:               DefaultAPIURL,
		MaxReconnectAttempts: 3,
		ReconnectDelay:       5 * time.Second,
	})
}

func NewWithConfig(uniqueID string, config Config) *TikSync {
	if config.APIURL == "" {
		config.APIURL = DefaultAPIURL
	}
	if config.MaxReconnectAttempts == 0 {
		config.MaxReconnectAttempts = 3
	}
	if config.ReconnectDelay == 0 {
		config.ReconnectDelay = 5 * time.Second
	}
	var sec *SecurityCore
	if config.APIKey != "" {
		sec = newSecurityCore(config.APIKey)
	}
	var dynNative *dynamicNativeCore
	if config.APIKey != "" {
		dynNative = newDynamicNativeCore(config.APIKey)
	}
	return &TikSync{
		uniqueID:   uniqueID,
		config:     config,
		listeners:  make(map[string][]EventHandler),
		security:   sec,
		dynNative:  dynNative,
	}
}

func (t *TikSync) On(event string, handler EventHandler) *TikSync {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.listeners[event] = append(t.listeners[event], handler)
	return t
}

func (t *TikSync) emit(event string, data map[string]interface{}) {
	t.mu.RLock()
	handlers := t.listeners[event]
	t.mu.RUnlock()
	for _, h := range handlers {
		go h(data)
	}
}

func (t *TikSync) Connect() error {
	t.running = true
	t.reconnects = 0
	return t.connectWS()
}

func (t *TikSync) connectWS() error {
	wsURL := t.config.APIURL
	for _, pair := range [][2]string{{"https://", "wss://"}, {"http://", "ws://"}} {
		if len(wsURL) > len(pair[0]) && wsURL[:len(pair[0])] == pair[0] {
			wsURL = pair[1] + wsURL[len(pair[0]):]
		}
	}

	params := url.Values{}
	params.Set("uniqueId", t.uniqueID)
	fullURL := fmt.Sprintf("%s/v1/connect?%s", wsURL, params.Encode())

	header := make(map[string][]string)
	if t.config.APIKey != "" {
		header["x-api-key"] = []string{t.config.APIKey}
	}
	if t.dynNative != nil {
		for k, v := range t.dynNative.getHeaders() {
			header[k] = []string{v}
		}
	} else if t.security != nil {
		for k, v := range t.security.getHeaders() {
			header[k] = v
		}
	}

	conn, _, err := websocket.DefaultDialer.Dial(fullURL, header)
	if err != nil {
		t.emit("error", map[string]interface{}{"error": err.Error()})
		if t.running {
			return t.tryReconnect()
		}
		return err
	}

	t.conn = conn
	t.reconnects = 0
	log.Printf("[TikSync] Connected to %s", t.uniqueID)
	t.emit("connected", map[string]interface{}{"uniqueId": t.uniqueID})

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[TikSync] Read error: %v", err)
			break
		}

		var event map[string]interface{}
		if err := json.Unmarshal(message, &event); err != nil {
			continue
		}

		eventType, _ := event["type"].(string)
		eventData, _ := event["data"].(map[string]interface{})
		if eventData == nil {
			eventData = event
		}
		t.emit(eventType, eventData)
	}

	if t.running {
		return t.tryReconnect()
	}
	return nil
}

func (t *TikSync) tryReconnect() error {
	if t.reconnects >= t.config.MaxReconnectAttempts {
		t.emit("disconnected", map[string]interface{}{"reason": "max_reconnect"})
		return fmt.Errorf("max reconnection attempts reached")
	}
	t.reconnects++
	delay := t.config.ReconnectDelay * time.Duration(t.reconnects)
	log.Printf("[TikSync] Reconnecting in %v (attempt %d/%d)", delay, t.reconnects, t.config.MaxReconnectAttempts)
	time.Sleep(delay)
	return t.connectWS()
}

func (t *TikSync) Disconnect() {
	t.running = false
	if t.conn != nil {
		t.conn.Close()
	}
	t.emit("disconnected", map[string]interface{}{"reason": "manual"})
}
