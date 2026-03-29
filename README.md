<p align="center">
  <img src="https://raw.githubusercontent.com/tiksync/.github/main/profile/logo-96.png" width="60" alt="TikSync" />
</p>

<h1 align="center">TikSync Go SDK</h1>

<p align="center">
  <strong>TikTok Live SDK for Go</strong> - Real-time chat, gifts, likes, follows & viewer events.<br>
  <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT"></a>
  <a href="https://discord.gg/2RkymdBNa7"><img src="https://img.shields.io/discord/1487514051427700886?label=Discord&logo=discord&logoColor=white" alt="Discord"></a><br>
  <a href="https://tik-sync.com">Website</a> - <a href="https://tik-sync.com/docs">Documentation</a> - <a href="https://tik-sync.com/pricing">Pricing</a>
</p>

---

## Why TikSync?

- **No Puppeteer/Chromium** - Pure Rust signing engine, no browser dependency
- **Fast** - Sub-millisecond signature generation
- **Production-ready** - Used by 50+ TikTok Live streamers daily
- **Built-in reliability** - Auto-reconnection and error handling
- **6 SDKs, 1 API** - Same design across JS, Python, Go, Rust, Java, C#

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

## Get Started

1. Sign up at [tik-sync.com](https://tik-sync.com)
2. Create a free API key in your dashboard
3. Install the SDK and start building

Free tier available. See [pricing](https://tik-sync.com/pricing) for details.

## All SDKs

| Language | Package | Install |
|----------|---------|---------|
| JavaScript | [npm](https://www.npmjs.com/package/tiksync) | `npm install tiksync` |
| Python | [PyPI](https://pypi.org/project/tiksync/) | `pip install tiksync` |
| Go | [go.dev](https://pkg.go.dev/github.com/tiksync/tiksync-go) | `go get github.com/tiksync/tiksync-go` |
| Rust | [crates.io](https://crates.io/crates/tiksync) | `cargo add tiksync` |
| Java | [Maven Central](https://central.sonatype.com/artifact/io.github.0xwolfsync/tiksync) | See docs |
| C# | [NuGet](https://www.nuget.org/packages/TikSync) | `dotnet add package TikSync` |

## Community

- [Discord](https://discord.gg/2RkymdBNa7) - Get help and chat with other developers
- [Documentation](https://tik-sync.com/docs) - Full API reference
- [Blog](https://tik-sync.com/blog) - Technical deep-dives
- [Status](https://tik-sync.com/status) - Service uptime

## License

MIT - Built by [TikSync](https://tik-sync.com)
