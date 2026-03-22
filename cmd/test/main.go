package main

import (
	"fmt"
	"sync"
	"time"

	tiksync "github.com/tiksync/tiksync-go"
)

func main() {
	client := tiksync.NewWithConfig("croc_mi", tiksync.Config{
		APIKey:               "sl_int_synclive_electron_prod_2026",
		APIURL:               "https://api.synclive.fr",
		MaxReconnectAttempts: 1,
		ReconnectDelay:       5 * time.Second,
	})

	var wg sync.WaitGroup
	wg.Add(1)
	count := 0

	client.On("connected", func(data map[string]interface{}) {
		fmt.Println("[Go] connected: OK")
	})
	client.On("roomInfo", func(data map[string]interface{}) {
		fmt.Println("[Go] roomInfo: OK")
	})
	client.On("chat", func(data map[string]interface{}) {
		fmt.Printf("[Go] chat: OK - %v\n", data["uniqueId"])
		count++
		if count >= 3 {
			wg.Done()
		}
	})
	client.On("gift", func(data map[string]interface{}) {
		fmt.Printf("[Go] gift: OK - %v\n", data["uniqueId"])
	})
	client.On("member", func(data map[string]interface{}) {
		fmt.Printf("[Go] member: OK - %v\n", data["uniqueId"])
	})
	client.On("roomUser", func(data map[string]interface{}) {
		fmt.Println("[Go] roomUser: OK")
	})
	client.On("error", func(data map[string]interface{}) {
		fmt.Printf("[Go] ERROR: %v\n", data["error"])
		wg.Done()
	})
	client.On("disconnected", func(data map[string]interface{}) {
		fmt.Printf("[Go] disconnected: %v\n", data["reason"])
	})

	fmt.Println("Connecting to croc_mi...")
	go func() {
		err := client.Connect()
		if err != nil {
			fmt.Printf("[Go] Connect error: %v\n", err)
			wg.Done()
		}
	}()

	wg.Wait()
	client.Disconnect()
	fmt.Println("[Go] disconnected: manual")
}
