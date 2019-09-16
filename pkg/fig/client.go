// Package fig implements the Golang client for fig
//
// The typical use of fig to fetch configuration is
// something like this;
//
//      import "github.com/rameshvk/fig/pkg/fig"
//      ...
//      var cfg := fig.Config(url, key, secret, time.Second)
//      ...
//      val, err := cfg.Get("my.entry", map[string]string{"user": 42})
//
// The arg provided to Get is typically a map or a struct whose
// fields can be accessed by the setting definition to provide the
// specific cofiguration value.  The setting is a simple expression
// which is evaluated for the matching arguments.
package fig

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Client implements the raw fig client API.
//
// This is typically not used by services. For fetching
// configuration, the Config() function is a lot simpler.
type Client struct {
	*http.Client
	URL         string
	AddAuthInfo func(r *http.Request) *http.Request
}

// New creates a new client based on the provided URL prefix.
//
// This implemention is not cached. Use Config() for a cached
// configuration fetcher
func New(url string) *Client {
	return &Client{&http.Client{}, url, func(r *http.Request) *http.Request { return r }}
}

// WithKey sets up the client to make calls with the provided API key
func (c *Client) WithKey(key, secret string) *Client {
	c.AddAuthInfo = func(r *http.Request) *http.Request {
		r.SetBasicAuth(key, secret)
		return r
	}
	return c
}

func (c *Client) GetSince(version int) (int, map[string]string) {
	u := mustParse(c.URL)
	q, err := url.ParseQuery(u.RawQuery)
	check(err)

	q.Add("version", strconv.Itoa(version))
	u.RawQuery = q.Encode()
	u = u.ResolveReference(mustParse("items"))
	var got struct {
		Version int
		Config  map[string]string
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	check(err)
	r, err := c.Client.Do(c.AddAuthInfo(req))
	checkResponse(r, err, &got)
	return got.Version, got.Config
}

func (c *Client) Set(key, val string) {
	u := mustParse(c.URL)
	u = u.ResolveReference(mustParse("items/" + url.PathEscape(key)))
	req, err := http.NewRequest("POST", u.String(), strings.NewReader(val))
	check(err)
	req.Header.Set("Content-Type", "application/json")
	r, err := c.Client.Do(c.AddAuthInfo(req))
	checkResponse(r, err, nil)
}

func (c *Client) History(key, epoch string) (string, []string) {
	u := mustParse(c.URL)
	q, err := url.ParseQuery(u.RawQuery)
	check(err)

	q.Add("epoch", epoch)
	u.RawQuery = q.Encode()
	u = u.ResolveReference(mustParse("items/" + url.PathEscape(key)))
	var got struct {
		Epoch   string
		History []string
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	check(err)
	r, err := c.Client.Do(c.AddAuthInfo(req))
	checkResponse(r, err, &got)
	return got.Epoch, got.History
}

func mustParse(s string) *url.URL {
	u, err := url.Parse(s)
	check(err)
	return u
}

func checkResponse(resp *http.Response, err error, v interface{}) {
	if err == nil && resp.StatusCode != 200 {
		body, err2 := ioutil.ReadAll(resp.Body)
		check(err2)
		err = errors.New("http.request failed " + resp.Status + "\n" + string(body))
	}
	defer resp.Body.Close()
	check(err)
	if v != nil {
		check(json.NewDecoder(resp.Body).Decode(v))
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
