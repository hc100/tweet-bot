package xclient

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/hc100/tweet-bot/internal/config"
)

func buildAuthorizationHeader(method string, rawURL string, query url.Values, credentials config.Credentials) string {
	oauthParams := map[string]string{
		"oauth_consumer_key":     credentials.APIKey,
		"oauth_nonce":            randomNonce(),
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        fmt.Sprintf("%d", time.Now().Unix()),
		"oauth_token":            credentials.AccessToken,
		"oauth_version":          "1.0",
	}

	signatureParams := make(url.Values, len(oauthParams))
	for key, value := range oauthParams {
		signatureParams.Set(key, value)
	}
	for key, values := range query {
		for _, value := range values {
			signatureParams.Add(key, value)
		}
	}

	oauthParams["oauth_signature"] = sign(method, rawURL, signatureParams, credentials)

	keys := make([]string, 0, len(oauthParams))
	for key := range oauthParams {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf(`%s="%s"`, percentEncode(key), percentEncode(oauthParams[key])))
	}

	return "OAuth " + strings.Join(parts, ", ")
}

func sign(method string, rawURL string, params url.Values, credentials config.Credentials) string {
	baseURL := normalizeURL(rawURL)
	paramString := normalizeParams(params)
	baseString := strings.Join([]string{
		strings.ToUpper(method),
		percentEncode(baseURL),
		percentEncode(paramString),
	}, "&")

	key := percentEncode(credentials.APISecret) + "&" + percentEncode(credentials.AccessTokenSecret)
	mac := hmac.New(sha1.New, []byte(key))
	_, _ = mac.Write([]byte(baseString))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func normalizeURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	scheme := strings.ToLower(parsed.Scheme)
	host := strings.ToLower(parsed.Hostname())
	port := parsed.Port()

	if (scheme == "https" && port == "443") || (scheme == "http" && port == "80") {
		port = ""
	}
	if port != "" {
		host = host + ":" + port
	}

	path := parsed.EscapedPath()
	if path == "" {
		path = "/"
	}

	return scheme + "://" + host + path
}

func normalizeParams(params url.Values) string {
	type pair struct {
		key   string
		value string
	}

	pairs := make([]pair, 0)
	for key, values := range params {
		encodedKey := percentEncode(key)
		for _, value := range values {
			pairs = append(pairs, pair{
				key:   encodedKey,
				value: percentEncode(value),
			})
		}
	}

	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].key == pairs[j].key {
			return pairs[i].value < pairs[j].value
		}
		return pairs[i].key < pairs[j].key
	})

	parts := make([]string, 0, len(pairs))
	for _, p := range pairs {
		parts = append(parts, p.key+"="+p.value)
	}

	return strings.Join(parts, "&")
}

func percentEncode(v string) string {
	encoded := url.QueryEscape(v)
	encoded = strings.ReplaceAll(encoded, "+", "%20")
	encoded = strings.ReplaceAll(encoded, "*", "%2A")
	return strings.ReplaceAll(encoded, "%7E", "~")
}

func randomNonce() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}

	return base64.RawURLEncoding.EncodeToString(buf)
}
