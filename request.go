package bricklinkapi

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// RequestHandler defines the request interface
type RequestHandler interface {
	Request(method, uri string) (body []byte, err error)
}

type request struct {
	consumerKey    string
	consumerSecret string
	token          string
	tokenSecret    string
}

// request() handles the request process. It builds of the oauth header,
// sets the request parameters and issues the request.
// The response body is returned as a []byte slice.
func (r request) Request(method, uri string) (body []byte, err error) {
	// new client
	client := http.Client{
		Timeout: time.Second * 30, // Maximum of 5 secs
	}

	// build new request
	req, err := http.NewRequest(method, brickLinkAPIBaseURL+uri, nil)
	if err != nil {
		return body, fmt.Errorf("could not build new request: %v", err)
	}

	// construct timestamp and nonce used in the oauth
	timeUnix := time.Now().UnixNano()
	timestamp := strconv.FormatInt(timeUnix, 10)
	nonce := strconv.FormatInt(rand.New(rand.NewSource(timeUnix)).Int63(), 10)

	// construct values for oauth params
	var oauthParams []string
	oauthParams = append(oauthParams, "oauth_consumer_key="+r.consumerKey)
	oauthParams = append(oauthParams, "oauth_token="+r.token)
	oauthParams = append(oauthParams, "oauth_signature_method="+oauthSignatureMethod)
	oauthParams = append(oauthParams, "oauth_timestamp="+timestamp)
	oauthParams = append(oauthParams, "oauth_nonce="+nonce)
	oauthParams = append(oauthParams, "oauth_version="+oauthVersion)

	// extract uri params from URI and add to oauth params map
	uriSplit := strings.Split(req.URL.String(), "?")
	if len(uriSplit) > 1 {
		uriParamString := strings.Split(uriSplit[1], "&")
		for _, s := range uriParamString {
			oauthParams = append(oauthParams, s)
		}
	}

	// generate signature
	baseURL := generateBaseURL(req, oauthParams)
	signature := generateSignature(baseURL, r.consumerSecret, r.tokenSecret)

	// set header
	req.Header.Set("User-Agent", "bricklinkapi-test")

	// build authorization string for the header
	authorization := "OAuth "
	authorization += "oauth_consumer_key=\"" + r.consumerKey + "\","
	authorization += "oauth_token=\"" + r.token + "\","
	authorization += "oauth_signature_method=\"" + oauthSignatureMethod + "\","
	authorization += "oauth_signature=\"" + signature + "\","
	authorization += "oauth_timestamp=\"" + timestamp + "\","
	authorization += "oauth_nonce=\"" + nonce + "\","
	authorization += "oauth_version=\"" + oauthVersion + "\""

	req.Header.Set("Authorization", authorization)

	// start request
	resp, err := client.Do(req)
	if err != nil {
		return body, err
	}

	// read response body
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	return body, nil
}

// generateBaseURL generates the base URL for the signature
func generateBaseURL(req *http.Request, params []string) string {
	base := req.Method + "&"
	base += encode(strings.Split(req.URL.String(), "?")[0])

	// sort params
	sort.Strings(params)

	paramString := strings.Join(params, "&")
	encodedParamString := encode(paramString)

	if len(encodedParamString) != 0 {
		base += "&" + encodedParamString
	}

	return base
}

// generateSignature generates the OAuth signature for the request.
// It receives the base string, the consumer secret and the token secret
func generateSignature(base, cs, ts string) string {
	// construct the key
	key := encode(cs) + "&" + encode(ts)

	// encrypt with HMAC-SHA1
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(base))
	rawSignature := h.Sum(nil)

	// base64 encode
	base64Signature := make([]byte, base64.StdEncoding.EncodedLen(len(rawSignature)))
	base64.StdEncoding.Encode(base64Signature, rawSignature)

	// percent encode and return
	signature := encode(string(base64Signature))

	return signature
}

// Implements percent encoding. The Golang std library implementation of
// url.QueryEscape is not valid for the oauth spec. Mainly spaces getting
// encoded as "+" instead of "%20"
func encode(s string) string {
	e := []byte(nil)
	for i := 0; i < len(s); i++ {
		b := s[i]
		if encodable(b) {
			e = append(e, '%')
			e = append(e, "0123456789ABCDEF"[b>>4])
			e = append(e, "0123456789ABCDEF"[b&15])
		} else {
			e = append(e, b)
		}
	}
	return string(e)
}

func encodable(b byte) bool {
	return !('A' <= b && b <= 'Z' || 'a' <= b && b <= 'z' ||
		'0' <= b && b <= '9' || b == '-' || b == '.' || b == '_' || b == '~')
}
