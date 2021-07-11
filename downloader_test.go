package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type testServer struct {
	server  *httptest.Server
	baseURL string
}

func NewTestServer() *testServer {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		switch url {
		case "/ankit.pl":
			fallthrough
		case "/github.com":
			fallthrough
		case "/google.com":
			w.WriteHeader(200)
			_, _ = w.Write([]byte(url))
		case "/invalid.com":
			w.WriteHeader(404)
		default:
			w.WriteHeader(500)
		}
	}))

	return &testServer{
		server:  s,
		baseURL: s.URL,
	}
}

func (t *testServer) URL(u string) string {
	return fmt.Sprintf("%s/%s", t.baseURL, u)
}

func TestDownloader_Download(t *testing.T) {
	ts := NewTestServer()
	defer ts.server.Close()

	cases := []struct {
		scenario string
		hosts    []string
		expected map[string]string
	}{
		{
			scenario: "all valid urls",
			hosts: []string{
				ts.URL("ankit.pl"),
				ts.URL("github.com"),
				ts.URL("google.com"),
			},
			expected: map[string]string{
				ts.URL("ankit.pl"):   "480735132954bbad8813ea6e5a4c9a73",
				ts.URL("github.com"): "06bcae6c0f349f5d9c4505a4cbea0cae",
				ts.URL("google.com"): "69ab4493fb4c2dfb0eea78711e6dd410",
			},
		},
		{
			scenario: "with some urls with unexpected response code",
			hosts: []string{
				ts.URL("ankit.pl"),
				ts.URL("invalid.com"),
				ts.URL("error.com"),
			},
			expected: map[string]string{
				ts.URL("ankit.pl"):    "480735132954bbad8813ea6e5a4c9a73",
				ts.URL("invalid.com"): "unexpected response code",
				ts.URL("error.com"):   "unexpected response code",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.scenario, func(t *testing.T) {
			out := NewDownloader(
				tc.hosts,
				WithParallelRequests(1),
				WithHTTPClient(ts.server.Client()),
			).Download()

			actual := make(map[string]string)
			for h := range out {
				if h.err != nil {
					actual[h.url] = h.err.Error()
				} else {
					actual[h.url] = h.sum
				}
			}

			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("expected %+v, got %+v", tc.expected, actual)
			}
		})
	}
}
