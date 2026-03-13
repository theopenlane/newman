package compose

import (
	"encoding/json"
	"maps"
	"net/url"
	"time"
)

// BuildTemplateURLs builds the URL map for template data from a Config, with optional per-key overrides
// Keys are derived from the URLConfig json struct tags to match template variable names under .URLS
func BuildTemplateURLs(config Config, overrides map[string]any) (map[string]any, error) {
	urls, err := structToMap(config.URLS)
	if err != nil {
		return nil, err
	}

	maps.Copy(urls, overrides)

	return urls, nil
}

// BuildTemplateData builds a template data map with the standard email fields from a Config and Recipient
// Keys are derived from json struct tags to match template variable names
func BuildTemplateData(config Config, recipient Recipient, overrides map[string]any) (map[string]any, error) {
	base, err := structToMap(config)
	if err != nil {
		return nil, err
	}

	base["Year"] = resolveYear(config.Year)

	recipientMap, err := structToMap(recipient)
	if err != nil {
		return nil, err
	}

	base["Recipient"] = recipientMap

	maps.Copy(base, overrides)

	return base, nil
}

// AddTokenToURL appends a token query parameter to a base URL
func AddTokenToURL(baseURL, token string) (string, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return "", ErrInvalidURL
	}

	values := parsed.Query()
	values.Set("token", token)
	parsed.RawQuery = values.Encode()

	return parsed.String(), nil
}

// resolveYear returns the configured year or defaults to the current year
func resolveYear(year int) int {
	if year != 0 {
		return year
	}

	return time.Now().Year()
}

// structToMap converts a struct to map[string]any using json field tags
func structToMap(v any) (map[string]any, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, ErrStructMarshal
	}

	var m map[string]any

	if err := json.Unmarshal(data, &m); err != nil {
		return nil, ErrStructMarshal
	}

	return m, nil
}
