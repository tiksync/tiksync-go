# TikSync Go SDK

**TikTok Live SDK for Go** — Real-time chat, gifts, likes, follows & viewer events.

```go
package main

import (
    "fmt"
    "github.com/tiksync/tiksync-go"
)

func main() {
    client := tiksync.New("charlidamelio", "your_api_key")

    client.On("chat", func(data map[string]interface{}) {
        fmt.Printf("[%s] %s\n", data["uniqueId"], data["comment"])
    })

    client.On("gift", func(data map[string]interface{}) {
        fmt.Printf("%s sent %s\n", data["uniqueId"], data["giftName"])
    })

    client.Connect()
}
```

## Installation

```bash
go get github.com/tiksync/tiksync-go
```

## Events

| Event | Description |
|-------|-------------|
| `chat` | Chat messages |
| `gift` | Gift events with diamond values |
| `like` | Like events |
| `follow` | New followers |
| `share` | Stream shares |
| `member` | User joins |
| `roomUser` | Viewer count updates |
| `streamEnd` | Stream ended |
| `connected` | Connected |
| `disconnected` | Disconnected |
| `error` | Error |

## License

MIT — built by [SyncLive](https://synclive.fr) | [tik-sync.com](https://tik-sync.com)
