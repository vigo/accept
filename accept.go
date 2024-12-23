package accept

import (
	"slices"
	"strconv"
	"strings"
)

const (
	fallbackMediaType = "text/plain"
)

// MediaType represents a media type and its quality value.
type MediaType struct {
	Type string
	Q    float64
}

// MediaTypes is a collection of MediaType.
type MediaTypes []MediaType

// Sort sorts the MediaTypes in descending order of Q value.
func (mt *MediaTypes) Sort() {
	slices.SortStableFunc(*mt, func(a, b MediaType) int {
		if a.Q > b.Q {
			return -1
		} else if a.Q < b.Q {
			return 1
		}

		return 0
	})
}

// ContentNegotiation holds required arguments.
type ContentNegotiation struct {
	defaultMediaType    string
	supportedMediaTypes []string
}

func (*ContentNegotiation) parseAcceptHeader(header string) MediaTypes {
	mediaTypes := MediaTypes{}
	parts := strings.Split(header, ",")

	for _, part := range parts {
		mediaAndQ := strings.Split(strings.TrimSpace(part), ";")
		mediaType := strings.TrimSpace(mediaAndQ[0])
		qValue := 1.0

		if len(mediaAndQ) > 1 {
			for _, param := range mediaAndQ[1:] {
				if strings.HasPrefix(param, "q=") {
					q, err := strconv.ParseFloat(strings.TrimPrefix(param, "q="), 64)
					if err == nil {
						qValue = q
					}
				}
			}
		}

		mediaTypes = append(mediaTypes, MediaType{Type: mediaType, Q: qValue})
	}

	mediaTypes.Sort()

	return mediaTypes
}

// Negotiate determines the best content type based on the Accept header.
// If no match is found, it returns the default media type.
func (cn *ContentNegotiation) Negotiate(header string) string {
	acceptedTypes := cn.parseAcceptHeader(header)
	for _, acceptedType := range acceptedTypes {
		if acceptedType.Type == "*/*" {
			if len(cn.supportedMediaTypes) > 0 {
				return cn.supportedMediaTypes[0]
			}
		}
		for _, supportedType := range cn.supportedMediaTypes {
			if acceptedType.Type == supportedType {
				return supportedType
			}
		}
	}

	return cn.defaultMediaType
}

// Option represents functional option type.
type Option func(*ContentNegotiation)

// WithSupportedMediaTypes sets supported media types.
func WithSupportedMediaTypes(mediaTypes ...string) Option {
	return func(cn *ContentNegotiation) {
		cn.supportedMediaTypes = mediaTypes
	}
}

// WithDefaultMediaType sets the default fallback media type.
func WithDefaultMediaType(mediaType string) Option {
	return func(cn *ContentNegotiation) {
		cn.defaultMediaType = mediaType
	}
}

// New instantiates http accept header content negotiation functionality.
func New(options ...Option) *ContentNegotiation {
	opts := &ContentNegotiation{
		defaultMediaType: fallbackMediaType,
	}

	for _, opt := range options {
		opt(opts)
	}

	return opts
}
