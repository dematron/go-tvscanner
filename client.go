package tvscanner

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/dematron/go-tvscanner/version"
)

type client struct {
	httpClient  *http.Client
	httpTimeout time.Duration
	debug       bool
}

// NewClient return a new Scanner HTTP client
func NewClient() (c *client) {
	return &client{&http.Client{}, 30 * time.Second, false}
}

// NewClientWithCustomHttpConfig returns a new Scanner HTTP client using the predefined http client
func NewClientWithCustomHttpConfig(httpClient *http.Client) (c *client) {
	timeout := httpClient.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &client{httpClient, timeout, false}
}

func (c client) dumpRequest(r *http.Request) {
	if r == nil {
		ContextLogger.Println("dumpReq ok: <nil>")
		return
	}
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		ContextLogger.Println("dumpReq err:", err)
	} else {
		ContextLogger.Println("dumpReq ok:", string(dump))
	}
}

func (c client) dumpResponse(r *http.Response) {
	if r == nil {
		ContextLogger.Println("dumpResponse ok: <nil>")
		return
	}
	dump, err := httputil.DumpResponse(r, true)
	if err != nil {
		ContextLogger.Errorf("dumpResponse err: %v", err)
	} else {
		ContextLogger.Println("dumpResponse ok:", string(dump))
	}
}

// doTimeoutRequest do a HTTP request with timeout
func (c *client) doTimeoutRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	// Do the request in the background so we can check the timeout
	type result struct {
		resp *http.Response
		err  error
	}
	done := make(chan result, 1)
	go func() {
		if c.debug {
			c.dumpRequest(req)
		}
		resp, err := c.httpClient.Do(req)
		if c.debug {
			c.dumpResponse(resp)
		}
		done <- result{resp, err}
	}()
	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		return nil, errors.New("timeout on reading data from TradingView Scanner API")
	}
}

// do prepare and process HTTP request to TradingView Scanner API
func (c *client) do(method string, payload string, authNeeded bool) (response []byte, err error) {
	connectTimer := time.NewTimer(c.httpTimeout)

	rawurl := fmt.Sprintf("%s%s/%s", API_URL, DEFAULT_SCREENER, API_POSTFIX)

	if c.debug {
		fmt.Println("url: ", rawurl)
	}
	req, err := http.NewRequest(method, rawurl, strings.NewReader(payload))
	if err != nil {
		return
	}
	if method == "POST" || method == "PUT" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	}
	req.Header.Add("DNT", "1")
	req.Header.Add("User-Agent", "go-tvscanner/"+version.Version)

	resp, err := c.doTimeoutRequest(connectTimer, req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)
	//fmt.Println(fmt.Sprintf("reponse %s", response), err)
	if err != nil {
		return response, err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		err = errors.New(resp.Status)
	}
	return response, err
}
