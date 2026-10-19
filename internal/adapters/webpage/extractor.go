package webpage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/getshelf/backend/internal/modules/item"
	"golang.org/x/net/html/charset"
)

const maxHTMLBytes int64 = 2 << 20 // 2 MiB

type Extractor struct {
	client *http.Client
}

var _ item.PageExtractor = (*Extractor)(nil)

func NewExtractor() *Extractor {
	return &Extractor{
		client: newHTTPClient(),
	}
}

func (extractor *Extractor) Extract(
	ctx context.Context,
	rawURL string,
) (item.PageMetadata, error) {
	pageURL, err := url.Parse(strings.TrimSpace(rawURL))

	if err != nil {
		return item.PageMetadata{}, fmt.Errorf("parse page url: %w", err)
	}

	if err := validateURL(pageURL); err != nil {
		return item.PageMetadata{}, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		pageURL.String(),
		nil,
	)

	if err != nil {
		return item.PageMetadata{}, fmt.Errorf("create page request: %w", err)
	}

	req.Header.Set("User-Agent", "ShelfMetadata/1.0")
	req.Header.Set("Accept", "text/html, application/xhtml+xml")

	res, err := extractor.client.Do(req)
	if err != nil {
		return item.PageMetadata{}, fmt.Errorf("fetch page: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return item.PageMetadata{}, fmt.Errorf(
			"fetch page: HTTP %d",
			res.StatusCode,
		)
	}

	body, err := readBounded(res.Body)
	if err != nil {
		return item.PageMetadata{}, fmt.Errorf("read page: %w", err)
	}

	contentType := res.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(body)
	}

	mediaType, _, err := mime.ParseMediaType(contentType)

	if err != nil {
		return item.PageMetadata{}, fmt.Errorf("parse content type: %w", err)
	}

	if mediaType != "text/html" && mediaType != "application/xhtml+xml" {
		return item.PageMetadata{}, fmt.Errorf(
			"unsupported content type: %s",
			mediaType,
		)
	}

	reader, err := charset.NewReader(bytes.NewReader(body), contentType)
	if err != nil {
		return item.PageMetadata{}, fmt.Errorf("decode page charset: %w", err)
	}

	// Apply the limit again because conversion to UTF-8 can expand the data.
	decoded, err := readBounded(reader)
	if err != nil {
		return item.PageMetadata{}, fmt.Errorf("read decoded page: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(decoded))
	if err != nil {
		return item.PageMetadata{}, fmt.Errorf("parse page HTML: %w", err)
	}

	return parseMetadata(doc, res.Request.URL), nil
}

func readBounded(reader io.Reader) ([]byte, error) {
	body, err := io.ReadAll(io.LimitReader(reader, maxHTMLBytes+1))

	if err != nil {
		return nil, err
	}

	if int64(len(body)) > maxHTMLBytes {
		return nil, fmt.Errorf("HTML exceeds %d bytes", maxHTMLBytes)
	}

	return body, nil
}

func parseMetadata(
	doc *goquery.Document,
	finalURL *url.URL,
) item.PageMetadata {
	meta := make(map[string]string)

	doc.Find("head meta").Each(func (_ int, element *goquery.Selection) {
		value := strings.TrimSpace(element.AttrOr("content", ""))
		if value == "" {
			return
		}

		for _, attribute := range []string{"property", "name"} {
			key := strings.ToLower(
				strings.TrimSpace(element.AttrOr(attribute, "")),
			)

			if key != "" && meta[key] == ""{
				meta[key] = value
			}
		}
	})

	baseURL := finalURL
	baseHref := doc.Find("head base[href]").First().AttrOr("href", "")

	if resolved := resolveURL(finalURL, baseHref); resolved != "" {
		baseURL, _ = url.Parse(resolved)
	}

	var canonicalURL string
	doc.Find("head link[href]").EachWithBreak(
		func(_ int, element *goquery.Selection) bool {
			relations := strings.Fields(
				strings.ToLower(element.AttrOr("rel", "")),
			)

			for _, relation := range relations {
				if relation == "canonical" {
					canonicalURL = resolveURL(
						baseURL,
						element.AttrOr("href", ""),
					)
					break
				}
			}

			return canonicalURL == ""
		},
	)

	return item.PageMetadata{
		Title: normalizeText(firstNonEmpty(
			meta["og:title"],
			meta["twitter:title"],
			doc.Find("head title").First().Text(),
		)),
		Description: normalizeText(firstNonEmpty(
			meta["og:description"],
			meta["twitter:description"],
			meta["description"],
		)),
		ImageURL: firstNonEmpty(
			resolveURL(baseURL, meta["og:image"]),
			resolveURL(baseURL, meta["twitter:image"]),
		),
		SiteName: normalizeText(firstNonEmpty(
			meta["og:site_name"],
			finalURL.Hostname(),
		)),
		CanonicalURL: firstNonEmpty(canonicalURL, finalURL.String()),
		FinalURL:     finalURL.String(),
	}
}

func resolveURL(baseURL *url.URL, rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ""
	}

	reference, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	resolved := baseURL.ResolveReference(reference)
	if err := validateURL(resolved); err != nil {
		return ""
	}

	resolved.Fragment = ""
	return resolved.String()
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}

	return ""
}

func normalizeText(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
