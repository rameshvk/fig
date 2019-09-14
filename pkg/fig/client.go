// Package fig implements the Golang client for fig
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

// Client implements the main fig client API.
type Client struct {
	*http.Client
	URL string
}

// New creates a new client based on the provided URL prefix.
func New(url string) *Client {
	return &Client{&http.Client{}, url}
}

func (c *Client) GetSince(version int) (int, map[string]string) {
	u := mustParse(c.URL)
	q, err := url.ParseQuery(u.RawQuery)
	check(err)

	q.Add("version", strconv.Itoa(version))
	u.RawQuery = q.Encode()
	u = u.ResolveReference(mustParse("items"))
	resp, err := c.Client.Get(u.String())
	if err == nil && resp.StatusCode != 200 {
		body, err2 := ioutil.ReadAll(resp.Body)
		check(err2)
		err = errors.New("http.request failed " + resp.Status + "\n" + string(body))
	}
	check(err)
	defer resp.Body.Close()
	got := struct {
		Version int
		Config  map[string]string
	}{}
	check(json.NewDecoder(resp.Body).Decode(&got))
	return got.Version, got.Config
}

func (c *Client) Set(key, val string) {
	u := mustParse(c.URL)
	u = u.ResolveReference(mustParse("items/" + url.PathEscape(key)))
	resp, err := c.Client.Post(u.String(), "application/json", strings.NewReader(val))
	if err == nil && resp.StatusCode != 200 {
		body, err2 := ioutil.ReadAll(resp.Body)
		check(err2)
		err = errors.New("http.request failed " + resp.Status + "\n" + string(body))
	}
	check(err)
	resp.Body.Close()
}

func (c *Client) History(key, epoch string) (string, []string) {
	u := mustParse(c.URL)
	q, err := url.ParseQuery(u.RawQuery)
	check(err)

	q.Add("epoch", epoch)
	u.RawQuery = q.Encode()
	u = u.ResolveReference(mustParse("items/" + url.PathEscape(key)))
	resp, err := c.Client.Get(u.String())
	if err == nil && resp.StatusCode != 200 {
		body, err2 := ioutil.ReadAll(resp.Body)
		check(err2)
		err = errors.New("http.request failed " + resp.Status + "\n" + string(body))
	}
	check(err)
	defer resp.Body.Close()
	got := struct {
		Epoch   string
		History []string
	}{}
	check(json.NewDecoder(resp.Body).Decode(&got))
	return got.Epoch, got.History
}

func mustParse(s string) *url.URL {
	u, err := url.Parse(s)
	check(err)
	return u
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
