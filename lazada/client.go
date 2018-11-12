// Package lazada provides the client for interacting with the Lazada open platform
package lazada

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
)

type Client struct {
	BaseURL *url.URL

	client *http.Client

	common service

	secret string
	appKey string

	accessToken string

	// The product service used for making API calls related to products
	Products *ProductService

	// The auth service used for making API calls related to authorization or OAuth
	Auth *AuthService
}

type service struct {
	client *Client
}

// ListOptions used for various API calls
type ListOptions struct {
	// Used to paginate through results
	Offset int `url:"offset"`

	// Limit how many items are returned
	Limit int `url:"limit"`
}

type Payload struct {
	Payload string `url:"payload"`
}

// Default list options
var DefaultListOptions = ListOptions{
	Limit:  100,
	Offset: 0,
}

type LazadaResponse struct {
	Code      string          `json:"code"`
	Data      json.RawMessage `json:"data"`
	RequestID string          `json:"request_id"`

	Type    string          `json:"type"`
	Message string          `json:"message"`
	Detail  []*ErrorDetails `json:"detail"`
}

// NewClient takes in the application key, secret, and Lazada region and returns a client.
func NewClient(appKey, secret string, region Region) *Client {
	baseURL, _ := url.Parse(endpoints[region])

	c := &Client{
		client:  http.DefaultClient,
		appKey:  appKey,
		secret:  secret,
		BaseURL: baseURL,
	}

	initServices(c)
	return c
}

func initServices(c *Client) {
	c.common.client = c
	c.Products = (*ProductService)(&c.common)
	c.Auth = (*AuthService)(&c.common)
}

// NewTokenClient takes a client access token and returns a copy of the client with the token set.
func (c *Client) NewTokenClient(token string) *Client {
	newC := *c
	newC.accessToken = token
	initServices(&newC)
	return &newC
}

// SetRegion changes the region on the client
func (c *Client) SetRegion(region Region) {
	baseURL, _ := url.Parse(endpoints[region])
	c.BaseURL = baseURL
}

// addOptions sets the query string using the query encoding library
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// NewRequest returns an http request conforming to the open platform
// Any body supplied will be encoded to XML
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasPrefix(urlStr, "https") {
		urlStr = fmt.Sprintf("rest%s", urlStr)
	}

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "cant parse url")
	}

	var req *http.Request

	// If we are sending a body url encode everything (even the xml for some reason)
	if body != nil {
		buf := new(bytes.Buffer)
		buf.Write([]byte(xml.Header))

		enc := xml.NewEncoder(buf)
		err := enc.Encode(body)
		if err != nil {
			return nil, errors.Wrap(err, "cant encode body")
		}

		reqParams := url.Values{}
		reqParams.Set("payload", buf.String())
		reqParams.Set("sign_method", "sha256")
		reqParams.Set("timestamp", fmt.Sprintf("%d", time.Now().Unix()*1000))
		reqParams.Set("app_key", c.appKey)

		if c.accessToken != "" {
			reqParams.Set("access_token", c.accessToken)
		}
		sig := c.Signature(strings.TrimPrefix(u.Path, "/rest"), reqParams, strings.NewReader(reqParams.Encode()))
		reqParams.Set("sign", sig)

		req, err = http.NewRequest(method, u.String(), strings.NewReader(reqParams.Encode()))
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
		return req, nil
	}

	req, err = http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if body != nil {
	}

	return req, nil
}

// Do runs a http.Request adding in the various required query parameters if they weren't set by the body already.
// It will marshal the data returned into the provided interface.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*LazadaResponse, error) {
	var q url.Values

	if req.Body == nil {
		q = req.URL.Query()
		q.Set("sign_method", "sha256")
		q.Set("timestamp", fmt.Sprintf("%d", time.Now().Unix()*1000))
		q.Set("app_key", c.appKey)

		if c.accessToken != "" {
			q.Set("access_token", c.accessToken)
		}

		var body io.ReadCloser
		var err error
		if req.Body != nil {
			body, err = req.GetBody()
			if err != nil {
				return nil, err
			}
		}

		sig := c.Signature(strings.TrimPrefix(req.URL.Path, "/rest"), q, body)
		q.Set("sign", sig)
	}

	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	lazResp, err := CheckResponse(resp)
	if err != nil {
		return nil, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			decErr := json.Unmarshal(lazResp.Data, v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = errors.Wrap(decErr, "unable to unmarshal into struct")
			}
		}
	}

	return lazResp, err
}

// CheckResponse makes sure we didn't receive an error from the platform and if we did it returns the error properly.
func CheckResponse(r *http.Response) (*LazadaResponse, error) {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		lazResp := &LazadaResponse{}
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, errors.Wrap(err, "unable to read body")
		}
		decErr := json.Unmarshal(data, lazResp)
		if decErr != nil {
			return nil, errors.Wrap(decErr, "unable to decode response")
		}

		switch lazResp.Code {
		case "0":
			//Reset the body
			r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
			return lazResp, nil
		default:
			return nil, &ErrorResponse{Response: r, Code: lazResp.Code, Message: lazResp.Message,
				Type: lazResp.Type, RequestID: lazResp.RequestID, Detail: lazResp.Detail}
		}
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}

	return nil, errorResponse
}

// Signature calculates the signature for the query parameter needed with every request
func (c *Client) Signature(api string, val url.Values, body io.Reader) string {
	var buf bytes.Buffer

	buf.WriteString(api)
	keys := make([]string, 0, len(val))
	for k := range val {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		vs := val[k]
		keyEscaped := url.QueryEscape(k)

		for _, v := range vs {
			buf.WriteString(keyEscaped)
			buf.WriteString(v)
		}
	}

	// So even though the documentation says u need to add the body to the signature you dont
	//if body != nil {
	//	b, _ := ioutil.ReadAll(body)
	//	buf.Write(b)
	//}

	signer := hmac.New(sha256.New, []byte(c.secret))
	signer.Write(buf.Bytes())
	sig := signer.Sum(nil)

	return strings.ToUpper(hex.EncodeToString(sig))
}

// SliceString takes in a slice of strings and returns a string used in query parameters for the open platform
func SliceString(in []string) string {
	var builder strings.Builder

	builder.WriteString("[")
	for i, st := range in {
		builder.WriteString("\"")
		builder.WriteString(st)
		builder.WriteString("\"")

		if len(in)-1 != i {
			builder.WriteString(",")
		}
	}

	builder.WriteString("]")

	return builder.String()
}
