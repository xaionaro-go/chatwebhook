package kickcom

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/xaionaro-go/chatwebhook/pkg/cache"
)

const (
	URLPubKey                       = "https://api.kick.com/public/v1/public-key"
	PubKeyDuration                  = time.Hour * 24
	HTTPHeaderEventMessageID        = "Kick-Event-Message-Id"
	HTTPHeaderEventMessageTimestamp = "Kick-Event-Message-Timestamp"
	HTTPHeaderEventSignature        = "Kick-Event-Signature"
)

type PubKeyVerifier struct {
	CurrentPubKey *rsa.PublicKey
}

func getPubKeySerialized(
	ctx context.Context,
	cache cache.Cache,
) ([]byte, error) {
	var pubKeySerializedStale []byte
	if cache != nil {
		if v, t := cache.Get(ctx, URLPubKey); v != nil {
			if time.Since(t) < PubKeyDuration {
				return v, nil
			}
			pubKeySerializedStale = v
		}
	}

	resp, err := http.Get(URLPubKey)
	if err != nil {
		return pubKeySerializedStale, fmt.Errorf("unable to fetch public key from %s: %w", URLPubKey, err)
	}
	defer resp.Body.Close()

	pubKeyJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		return pubKeySerializedStale, fmt.Errorf("unable to read public key from %s: %w", URLPubKey, err)
	}

	type pubKeyResponse struct {
		Data struct {
			PublicKey string `json:"public_key"`
		} `json:"data"`
		Message string `json:"message"`
	}

	var pubKeyWithMetadata pubKeyResponse
	err = json.Unmarshal(pubKeyJSON, &pubKeyWithMetadata)
	if err != nil {
		return pubKeySerializedStale, fmt.Errorf("unable to unmarshal public key from %s: '%s': %w", URLPubKey, pubKeyJSON, err)
	}
	pubKeySerialized := []byte(pubKeyWithMetadata.Data.PublicKey)

	if cache != nil {
		cache.Set(ctx, URLPubKey, pubKeySerialized)
	}

	return pubKeySerialized, nil
}

func NewPubKeyVerifier(
	ctx context.Context,
	cache cache.Cache,
) (*PubKeyVerifier, error) {
	pubKeySerialized, err := getPubKeySerialized(ctx, cache)
	if err != nil {
		if pubKeySerialized == nil {
			return nil, err
		}
		logger.Errorf(ctx, "using the stale public key due to error: %v", err)
	}

	// Parse the PEM-encoded public key
	pubKey, err := ParsePublicKey(pubKeySerialized)
	if err != nil {
		return nil, fmt.Errorf("unable to parse public key from %s: '%s': %w", URLPubKey, pubKeySerialized, err)
	}

	pubKeyRSA, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key from %s is not an RSA public key", URLPubKey)
	}

	return &PubKeyVerifier{
		CurrentPubKey: pubKeyRSA,
	}, nil
}

func ParsePublicKey(pubKeySerialized []byte) (crypto.PublicKey, error) {
	block, rest := pem.Decode(pubKeySerialized)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}
	if len(rest) > 0 {
		return nil, fmt.Errorf("extra data found after PEM public key block")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse DER encoded public key: %w", err)
	}

	return pubKey, nil
}

// see: https://docs.kick.com/events/webhook-security
func (v *PubKeyVerifier) Verify(msgHash []byte, msgSignature []byte) error {
	return rsa.VerifyPKCS1v15(v.CurrentPubKey, crypto.SHA256, msgHash, msgSignature)
}

// see: https://docs.kick.com/events/webhook-security
func (v *PubKeyVerifier) VerifyRequest(r *http.Request) error {
	msgID := r.Header.Get(HTTPHeaderEventMessageID)
	if msgID == "" {
		return fmt.Errorf("missing %q header", HTTPHeaderEventMessageID)
	}

	msgTimestamp := r.Header.Get(HTTPHeaderEventMessageTimestamp)
	if msgTimestamp == "" {
		return fmt.Errorf("missing %q header", HTTPHeaderEventMessageTimestamp)
	}

	msgSignatureB64 := r.Header.Get(HTTPHeaderEventSignature)
	if msgSignatureB64 == "" {
		return fmt.Errorf("missing %q header", HTTPHeaderEventSignature)
	}

	msgSignature := make([]byte, base64.StdEncoding.DecodedLen(len(msgSignatureB64)))
	n, err := base64.StdEncoding.Decode(msgSignature, []byte(msgSignatureB64))
	if err != nil {
		return fmt.Errorf("unable to decode base64 signature %q: %w", msgSignatureB64, err)
	}
	msgSignature = msgSignature[:n]

	bodyReader, err := r.GetBody()
	if err != nil {
		return fmt.Errorf("unable to read request body: %w", err)
	}
	defer bodyReader.Close()

	body, err := io.ReadAll(bodyReader)
	if err != nil {
		return fmt.Errorf("unable to read request body: %w", err)
	}
	msgHash := sha256.Sum256([]byte(fmt.Sprintf("%s.%s.%s", msgID, msgTimestamp, body)))

	if err := v.Verify(msgHash[:], msgSignature); err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}

	return nil
}
