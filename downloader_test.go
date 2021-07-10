package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestDownloader_Download(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	defer server.Close()

	baseURL := server.URL

	cases := []struct {
		scenario string
		hosts    []string
		expected map[string]string
	}{
		{
			scenario: "all valid urls",
			hosts:    []string{baseURL + "/ankit.pl", baseURL + "/github.com", baseURL + "/google.com"},
			expected: map[string]string{
				baseURL + "/ankit.pl":   "480735132954bbad8813ea6e5a4c9a73",
				baseURL + "/github.com": "06bcae6c0f349f5d9c4505a4cbea0cae",
				baseURL + "/google.com": "69ab4493fb4c2dfb0eea78711e6dd410",
			},
		},
		{
			scenario: "with some urls with unexpected response code",
			hosts:    []string{baseURL + "/ankit.pl", baseURL + "/invalid.com", baseURL + "/error.com"},
			expected: map[string]string{
				baseURL + "/ankit.pl":    "480735132954bbad8813ea6e5a4c9a73",
				baseURL + "/invalid.com": "unexpected response code",
				baseURL + "/error.com":   "unexpected response code",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.scenario, func(t *testing.T) {
			actual := NewDownloader(tc.hosts, WithHTTPClient(server.Client())).Download()

			if !reflect.DeepEqual(tc.expected, actual) {
				t.Errorf("expected %+v, got %+v", tc.expected, actual)
			}
		})
	}
}
