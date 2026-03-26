<p align="center">
  <img src="https://raw.githubusercontent.com/tiksync/.github/main/profile/logo-96.png" width="60" alt="TikSync" />
</p>

<h1 align="center">TikSync Go SDK</h1>

<p align="center">
  <strong>TikTok Live SDK for Go</strong> — Real-time chat, gifts, likes, follows & viewer events.<br>
  <a href="https://tik-sync.com">Website</a> · <a href="https://tik-sync.com/docs">Documentation</a> · <a href="https://tik-sync.com/pricing">Pricing</a>
</p>

---

## Installation

```bash
go get github.com/tiksync/tiksync-go
```

## Quick Start

```go
package main

import (
    "fmt"
    tiksync "github.com/tiksync/tiksync-go"
)

func main() {
    client := tiksync.New("username", "your_api_key")

    client.On("chat", func(data map[string]interface{}) {
        fmt.Printf("[%s] %s\n", data["uniqueId"], data["comment"])
    })

    client.On("gift", func(data map[string]interface{}) {
        fmt.Printf("%s sent %s x%v\n", data["uniqueId"], data["giftName"], data["repeatCount"])
    })

    client.On("follow", func(data map[string]interface{}) {
        fmt.Printf("%s followed!\n", data["uniqueId"])
    })

    client.Connect()
}
```

## Configuration

```go
client := tiksync.NewWithConfig("username", tiksync.Config{
    APIKey:               "your_api_key",
    APIURL:               "https://api.tik-sync.com",
    MaxReconnectAttempts: 5,
    ReconnectDelay:       5 * time.Second,
})
```

## Events

| Event | Description |
|-------|-------------|
| `connected` | Connected to stream |
| `chat` | Chat message received |
| `gift` | Gift received (with diamond count, streak info) |
| `like` | Likes received |
| `follow` | New follower |
| `share` | Stream shared |
| `member` | User joined the stream |
| `roomUser` | Viewer count update |
| `streamEnd` | Stream ended |
| `disconnected` | Disconnected |
| `error` | Connection error |

## Get Your API Key

1. Sign up at [tik-sync.com](https://tik-sync.com)
2. Go to Dashboard → API Keys
3. Create a new key

Free tier: 1,000 requests/day, 10 WebSocket connections.

## Why TikSync?

- **< 1ms** signature latency
- **$0** CAPTCHA cost
- **Rust-powered** native backend
- **6 official SDKs** — JS, Python, Go, Java, C#, Rust
- **Direct connection** — no proxy required

## License

MIT — Built by [TikSync](https://tik-sync.com)
