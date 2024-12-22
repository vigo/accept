package accept_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/vigo/accept"
)

func TestContentNegotiation(t *testing.T) {
	tests := []struct {
		name                string
		acceptHeader        string
		supportedMediaTypes []string
		expected            string
	}{
		{
			name:                "Match first supported type",
			acceptHeader:        "application/json,text/html",
			supportedMediaTypes: []string{"application/json", "text/html"},
			expected:            "application/json",
		},
		{
			name:                "Match second supported type with q value",
			acceptHeader:        "text/html;q=0.9,application/json;q=0.8",
			supportedMediaTypes: []string{"application/json", "text/html"},
			expected:            "text/html",
		},
		{
			name:                "Fallback to default media type",
			acceptHeader:        "application/xml",
			supportedMediaTypes: []string{"application/json", "text/html"},
			expected:            "text/plain",
		},
		{
			name:                "Empty Accept header",
			acceptHeader:        "",
			supportedMediaTypes: []string{"application/json", "text/html"},
			expected:            "text/plain",
		},
		{
			name:                "Wildcard match",
			acceptHeader:        "*/*",
			supportedMediaTypes: []string{"application/json", "text/html"},
			expected:            "application/json",
		},
		{
			name:                "Mixed q values with wildcard",
			acceptHeader:        "application/json;q=0.7,*/*;q=0.5",
			supportedMediaTypes: []string{"application/json", "text/html"},
			expected:            "application/json",
		},
		{
			name:                "No supported types, default fallback",
			acceptHeader:        "image/png",
			supportedMediaTypes: []string{"application/json", "text/html"},
			expected:            "text/plain",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cn := accept.New(
				accept.WithSupportedMediaTypes(test.supportedMediaTypes...),
			)

			result := cn.Negotiate(test.acceptHeader)

			if result != test.expected {
				t.Errorf("expected %s, got %s", test.expected, result)
			}
		})
	}
}

func TestContentNegotiation_WithDefaultMediaType(t *testing.T) {
	tests := []struct {
		name                string
		acceptHeader        string
		supportedMediaTypes []string
		defaultMediaType    string
		expected            string
	}{
		{
			name:                "Match first supported type",
			acceptHeader:        "application/xml",
			supportedMediaTypes: []string{"application/json"},
			defaultMediaType:    "text/html",
			expected:            "text/html",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cn := accept.New(
				accept.WithSupportedMediaTypes(test.supportedMediaTypes...),
				accept.WithDefaultMediaType(test.defaultMediaType),
			)

			result := cn.Negotiate(test.acceptHeader)

			if result != test.expected {
				t.Errorf("expected %s, got %s", test.expected, result)
			}
		})
	}
}

func TestMediaTypesSorting(t *testing.T) {
	mediaTypes := accept.MediaTypes{
		{Type: "application/json", Q: 0.8},
		{Type: "text/html", Q: 0.9},
		{Type: "application/xml", Q: 0.7},
	}

	sort.Sort(mediaTypes)

	expectedOrder := []string{"text/html", "application/json", "application/xml"}
	for i, mt := range mediaTypes {
		if mt.Type != expectedOrder[i] {
			t.Errorf("expected %s at position %d, got %s", expectedOrder[i], i, mt.Type)
		}
	}
}

func ExampleContentNegotiation() {
	cn := accept.New(
		accept.WithSupportedMediaTypes("application/json", "text/html"),
	)

	// Simulate an Accept header from a client
	// r.Header.Get("Accept")
	acceptHeader := "application/json,text/html;q=0.9"

	contentType := cn.Negotiate(acceptHeader)

	fmt.Println(contentType)
	// Output: application/json
}

func ExampleContentNegotiation_fallback() {
	cn := accept.New(
		accept.WithSupportedMediaTypes("application/json", "text/html"),
		accept.WithDefaultMediaType("application/xml"),
	)

	// Simulate an Accept header that doesn't match any supported types
	// r.Header.Get("Accept")
	acceptHeader := "image/png"

	contentType := cn.Negotiate(acceptHeader)

	fmt.Println(contentType)
	// Output: application/xml
}

func ExampleContentNegotiation_wildcard() {
	cn := accept.New(
		accept.WithSupportedMediaTypes("application/json", "text/html"),
	)

	// Simulate an Accept header with a wildcard
	// r.Header.Get("Accept")
	acceptHeader := "*/*"

	contentType := cn.Negotiate(acceptHeader)

	fmt.Println(contentType)
	// Output: application/json
}
