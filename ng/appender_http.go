package ng

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//******************* HTTP APPENDER ********************

type Client struct {
	httpClient            *http.Client
	DialTimeout           time.Duration
	DialKeepAliveTimeout  time.Duration
	MaxIdleConnections    int
	IdleConnTimeout       time.Duration
	TlsHandshakeTimeout   time.Duration
	ResponseHeaderTimeout time.Duration
	//ExpectContinueTimeout    time.Duration	// will disable HTTP2 if used
	HttpClientRequestTimeout time.Duration
	disableVerifyCert        bool
}
type Auth struct {
	Username []byte
	Password []byte
}

type HTTPAppender struct {
	*OutAppender
	url    string
	path   string
	method string
	client *Client
	auth   *Auth
}

func NewHTTPAppender(filter, url, username, password string) (*HTTPAppender, error) {
	if len(url) == 0 {
		return nil, fmt.Errorf("url required")
	}
	oa := newOutAppender(filter, "")
	t := new(HTTPAppender)
	t.OutAppender = oa
	t.url = url
	t.method = "POST"
	t.disableColor = true

	tc := new(Client)
	tc.DialTimeout = 30 * time.Second
	tc.DialKeepAliveTimeout = 30 * time.Second
	tc.MaxIdleConnections = 100
	tc.IdleConnTimeout = 90 * time.Second
	tc.TlsHandshakeTimeout = 10 * time.Second
	tc.ResponseHeaderTimeout = 10 * time.Second
	//t.ExpectContinueTimeout = 1 * time.Second
	tc.HttpClientRequestTimeout = 30 * time.Second

	t.client = tc

	if len(username) > 0 || len(password) > 0 {
		a := new(Auth)
		a.Username = []byte(username)
		a.Password = []byte(password)
		t.auth = a
	}
	return t, nil
}
func (f *HTTPAppender) Name() string {
	if len(f.name) > 0 {
		return f.name
	}
	return fmt.Sprintf("%T", f)
}
func (f *HTTPAppender) DisableColor() bool {
	return f.disableColor
}
func (f *HTTPAppender) Applicable(msg string) bool {
	if f.filter == "*" {
		return true
	}
	if strings.Index(msg, f.filter) > -1 {
		return true
	}
	return false
}

func (f *HTTPAppender) Process(msg []byte) {
	// Send via HTTP
	//fmt.Printf("sent: |%v|\n", strings.TrimSpace(string(msg)))
	resp, err := f.client.Fetch(f.method, f.url, f.auth, nil, bytes.NewBuffer([]byte(strings.TrimSpace(string(msg)))))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		if resp != nil {
			fmt.Println("on " + f.url + strconv.Itoa(resp.StatusCode))
		}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("on " + f.url + " 500 unable to read response")
	}
	if len(body) > 0 {
		fmt.Printf("response: |%v|", string(body))
	}
}

func (c *Client) Fetch(method, url string, auth *Auth, header map[string]string, body io.Reader) (*http.Response, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.disableVerifyCert,
	}
	var netTransport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   c.DialTimeout, // time spent establishing a TCP connection
			KeepAlive: c.DialKeepAliveTimeout,
			//DualStack: true,		// now set by default and deprecated
		}).DialContext,
		MaxIdleConns:        c.MaxIdleConnections,
		IdleConnTimeout:     c.IdleConnTimeout,
		TLSHandshakeTimeout: c.TlsHandshakeTimeout, // time spent performing the TLS handshake
		//ExpectContinueTimeout: c.ExpectContinueTimeout, //time client will wait between sending request headers and receiving the go-ahead to send the body
		ResponseHeaderTimeout: c.ResponseHeaderTimeout, //time spent reading the headers of the response
		TLSClientConfig:       tlsConfig,
	}
	if c.httpClient == nil {
		c.httpClient = &http.Client{
			Timeout:   c.HttpClientRequestTimeout, //entire exchange, from Dial to reading the body
			Transport: netTransport,
		}
	}

	req, _ := http.NewRequest(method, url, body)
	//req = req.WithContext(ctx)
	req.Close = true
	if auth != nil {
		req.SetBasicAuth(string(auth.Username), string(auth.Password))
	}
	// Add any required headers.
	for headerKey, value := range header {
		req.Header.Add(headerKey, value)

		if headerKey == "Content-Length" {
			v, _ := strconv.ParseInt(value, 10, 64)
			req.ContentLength = v
		}
		if headerKey == "Content-Type" {
			req.Header.Set(headerKey, value)
		}
	}

	// Perform said network call.
	res, err := c.httpClient.Do(req)
	if err != nil {
		//glog.Error(err.Error()) // use glog it's amazing
		return nil, err
	}

	// If response from network call is not 200, return error too.
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted {
		return res, errors.New(res.Status)
	}
	return res, nil
}
