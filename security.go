package tiksync

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SecurityCore struct {
	apiKey      string
	deviceID    string
	fingerprint string
	primaryKey  []byte
	secondaryKey []byte
}

func newSecurityCore(apiKey string) *SecurityCore {
	deviceID := "tsd_" + strings.ReplaceAll(uuid.New().String(), "-", "")
	fingerprint := computeFingerprint(deviceID)

	sc := &SecurityCore{
		apiKey:      apiKey,
		deviceID:    deviceID,
		fingerprint: fingerprint,
	}
	sc.primaryKey = deriveKey(apiKey, deviceID, []byte("tiksync-primary-v1"))
	sc.secondaryKey = deriveKey(apiKey, deviceID, []byte("tiksync-secondary-v1"))
	return sc
}

func deriveKey(apiKey, deviceID string, salt []byte) []byte {
	h := sha256.New()
	h.Write([]byte(apiKey))
	h.Write([]byte("|"))
	h.Write([]byte(deviceID))
	h.Write([]byte("|"))
	h.Write(salt)
	return h.Sum(nil)
}

func computeFingerprint(deviceID string) string {
	h := sha256.New()
	h.Write([]byte(deviceID))
	h.Write([]byte(runtime.GOOS))
	h.Write([]byte(runtime.GOARCH))
	h.Write([]byte("1.0.0"))

	hostname, _ := os.Hostname()
	h.Write([]byte(hostname))

	username := os.Getenv("USERNAME")
	if username == "" {
		username = os.Getenv("USER")
	}
	uh := sha256.Sum256([]byte(username))
	h.Write([]byte(hex.EncodeToString(uh[:])))

	return hex.EncodeToString(h.Sum(nil))
}

func (sc *SecurityCore) sign(method, path, body string) (signature, timestamp, nonce string) {
	ts := fmt.Sprintf("%d", time.Now().UnixMilli())
	n := generateNonce()

	payload := fmt.Sprintf("%s|%s|%s|%s|%s", method, path, ts, n, body)

	primarySig := hmacSHA256(sc.primaryKey, []byte(payload))
	secondaryInput := fmt.Sprintf("%s|%s", primarySig, n)
	secondarySig := hmacSHA512(sc.secondaryKey, []byte(secondaryInput))

	combined := fmt.Sprintf("ts1.%s.%s", primarySig[:16], secondarySig[:24])
	return combined, ts, n
}

func (sc *SecurityCore) generateToken() string {
	epoch := time.Now().Unix() / 180
	nonce := make([]byte, 8)
	rand.Read(nonce)
	nonceHex := hex.EncodeToString(nonce)

	payload := fmt.Sprintf("%d|%s|%s", epoch, nonceHex, sc.deviceID)

	masterKey := deriveKey(sc.apiKey, sc.deviceID, []byte("tiksync-token-master-v1"))
	mac := hmac.New(sha256.New, masterKey)
	mac.Write([]byte(payload))
	sig := hex.EncodeToString(mac.Sum(nil))

	return fmt.Sprintf("tst1.%d.%s.%s", epoch, nonceHex, sig[:32])
}

func (sc *SecurityCore) getHeaders() map[string][]string {
	sig, ts, nonce := sc.sign("GET", "/v1/connect", "")
	token := sc.generateToken()

	return map[string][]string{
		"x-ts-device":      {sc.deviceID},
		"x-ts-signature":   {sig},
		"x-ts-timestamp":   {ts},
		"x-ts-nonce":       {nonce},
		"x-ts-token":       {token},
		"x-ts-fingerprint": {sc.fingerprint},
		"x-ts-version":     {"1.0.0"},
	}
}

func generateNonce() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func hmacSHA256(key, data []byte) string {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil))
}

func hmacSHA512(key, data []byte) string {
	mac := hmac.New(sha512.New, key)
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil))
}
