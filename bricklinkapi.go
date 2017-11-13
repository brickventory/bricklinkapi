package bricklinkapi

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	brickLinkAPIBaseURL  = "https://api.bricklink.com/api/store/v1"
	oauthVersion         = "1.0"
	oauthSignatureMethod = "HMAC-SHA1"
)

var (
	itemTypes = []string{"MINIFIG", "PART", "SET", "BOOK", "GEAR", "CATALOG", "INSTRUCTION", "UNSORTED_LOT", "ORIGINAL_BOX"}
)

// Bricklink is the main handler for the Bricklink API requests
type Bricklink struct {
	ConsumerKey    string
	ConsumerSecret string
	Token          string
	TokenSecret    string
}

// New returns a Bricklink handler ready to use
func New(consumerKey, consumerSecret, token, tokenSecret string) *Bricklink {
	return &Bricklink{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		Token:          token,
		TokenSecret:    tokenSecret,
	}
}

// GetItem issues a GET request to the Bricklink API
func (bl *Bricklink) GetItem(itemType, itemNumber string) (response string, err error) {
	// set default itemtype to "part"
	if itemType == "" || !stringInSlice(itemType, itemTypes) {
		return response, errors.New("itemType is not specified or valid")
	}

	// validate itemNumber
	if itemNumber == "" {
		return response, errors.New("itemNumber is not specified")
	}
	// build uri
	uri := "/items/" + itemType + "/" + itemNumber

	body, err := bl.request("GET", uri)
	if err != nil {
		return response, err
	}

	return string(body), nil
}

// request() handles the request process. It builds of the oauth header,
// sets the request parameters and issues the request.
// The response body is returned as a []byte slice.
func (bl Bricklink) request(method, uri string) (body []byte, err error) {
	// new client
	client := http.Client{
		Timeout: time.Second * 5, // Maximum of 5 secs
	}

	// build new request
	req, err := http.NewRequest(method, brickLinkAPIBaseURL+uri, nil)
	if err != nil {
		return body, fmt.Errorf("could not build new request: %v", err)
	}

	// construct timestamp and nonce used in the oauth
	timeUnix := time.Now().Unix()
	timestamp := strconv.FormatInt(timeUnix, 10)
	nonce := strconv.FormatInt(rand.New(rand.NewSource(timeUnix)).Int63(), 10)

	// construct values for request
	params := url.Values{}
	params.Add("oauth_consumer_key", bl.ConsumerKey)
	params.Add("oauth_token", bl.Token)
	params.Add("oauth_signature_method", oauthSignatureMethod)
	params.Add("oauth_timestamp", timestamp)
	params.Add("oauth_nonce", nonce)
	params.Add("oauth_version", oauthVersion)

	// generate signature
	baseURL := generateBaseURL(req, params)
	signature := generateSignature(baseURL, bl.ConsumerSecret, bl.TokenSecret)
	params.Add("oauth_signature", signature)

	// set header
	req.Header.Set("User-Agent", "bricklinkapi-test")

	// build authorization string
	authorization := "OAuth "
	authorization += "oauth_consumer_key=\"" + bl.ConsumerKey + "\","
	authorization += "oauth_token=\"" + bl.Token + "\","
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
func generateBaseURL(req *http.Request, params url.Values) string {
	base := req.Method + "&"
	base += encode(req.URL.String()) + "&"
	base += encode(params.Encode())

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

// helper function to check if a string is in a slice
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
