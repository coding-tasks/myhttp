package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// ErrUnexpectedResponseCode is returned for response code other than 200 OK.
var ErrUnexpectedResponseCode = errors.New("unexpected response code")

const (
	defaultParallelRequests = 10
	defaultTimeout          = 10 * time.Second
)

// Downloader is a http downloader.
type Downloader struct {
	client   *http.Client
	hosts    []string
	parallel int
}

// DownloadOption defines functional options for the downloader.
type DownloadOption func(*Downloader)

// NewDownloader constructs a downloader object.
func NewDownloader(hosts []string, opts ...DownloadOption) *Downloader {
	dl := Downloader{
		client: &http.Client{
			Timeout: defaultTimeout,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: defaultTimeout,
				}).DialContext,
				TLSHandshakeTimeout: defaultTimeout,
			},
		},
		hosts:    hosts,
		parallel: len(hosts),
	}

	if dl.parallel > defaultParallelRequests {
		dl.parallel = defaultParallelRequests
	}

	for _, opt := range opts {
		opt(&dl)
	}

	return &dl
}

// WithParallelRequests sets number of parallel downloads.
func WithParallelRequests(n int) DownloadOption {
	return func(dl *Downloader) {
		dl.parallel = n
	}
}

// WithHTTPClient sets http client for the downloader.
func WithHTTPClient(c *http.Client) DownloadOption {
	return func(dl *Downloader) {
		dl.client = c
	}
}

type result struct {
	url string
	sum string
	err error
}

// Download fetches the content from host, calculates and assign checksum.
func (dl *Downloader) Download() chan result {
	host := make(chan string, dl.parallel)

	var wg sync.WaitGroup
	wg.Add(dl.parallel)

	out := make(chan result, dl.parallel)
	for i := 0; i < dl.parallel; i++ {
		go func() {
			defer wg.Done()
			dl.worker(host, out)
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	go func() {
		for _, url := range dl.hosts {
			host <- sanitizeURL(url)
		}
		close(host)
	}()

	return out
}

func (dl *Downloader) worker(host <-chan string, r chan<- result) {
	for url := range host {
		b, err := dl.fetch(url)
		r <- result{url, checksum(b), err}
	}
}

func (dl *Downloader) fetch(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Adjust-bot/1.0")

	resp, err := dl.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrUnexpectedResponseCode
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func checksum(b []byte) string {
	h := md5.New()
	_, _ = h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

func sanitizeURL(url string) string {
	// We don't care about other schemes like ftp://, file:// etc.
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "https://" + url
	}
	return url
}
