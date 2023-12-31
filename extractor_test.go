package gophetch_test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"

	"github.com/octetic/gophetch"
	"github.com/octetic/gophetch/metadata"
)

// TODO(patrick): Add more tests

func TestExtractorIntegration(t *testing.T) {
	testCases := []struct {
		desc     string
		mockHTML string
		expected metadata.Metadata
	}{
		{
			desc: "Full HTML Page",
			mockHTML: `
				<!DOCTYPE html>
				<html>
				<head>
					<meta property="og:title" content="OG Title"/>
					<meta property="og:description" content="OG Description"/>
					<meta property="article:published_time" content="2022-10-11T15:04:05Z"/>
					<span property="schema:author">John Schema</span>
				</head>
				<body>
					<!-- ... -->
				</body>
				</html>
			`,
			expected: metadata.Metadata{
				Author:      "John Schema",
				Title:       "OG Title",
				Description: "OG Description",
				Date:        "2022-10-11T15:04:05Z",
			},
		},
		// Add more test cases as needed
	}

	ext := gophetch.NewExtractor()
	targetURL, err := url.Parse("https://example.com")
	if err != nil {
		t.Fatal(err)
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// Parse the mock HTML
			mockNode, err := html.Parse(strings.NewReader(tC.mockHTML))
			if err != nil {
				t.Fatal(err)
			}

			// Call the Extractor's method to get all metadata
			result, err := ext.ExtractMetadata(mockNode, targetURL)
			if err != nil {
				t.Fatal(err)
			}

			// Assertion
			assert.Equal(t, tC.expected.Author, result.Author, "They should be equal")
			assert.Equal(t, tC.expected.Title, result.Title, "They should be equal")
			assert.Equal(t, tC.expected.Description, result.Description, "They should be equal")
			assert.Equal(t, tC.expected.Date, result.Date, "They should be equal")
		})
	}
}

func TestExtractorBoundary(t *testing.T) {
	testCases := []struct {
		desc     string
		mockHTML string
		expected metadata.Metadata
	}{
		{
			desc: "Missing All Tags",
			mockHTML: `
				<!DOCTYPE html>
				<html>
				<head>
				</head>
				<body>
				</body>
				</html>
			`,
			expected: metadata.Metadata{},
		},
		{
			desc: "Empty Strings",
			mockHTML: `
				<!DOCTYPE html>
				<html>
				<head>
					<meta property="og:title" content=""/>
					<meta property="og:description" content=""/>
					<meta property="article:published_time" content=""/>
					<span property="schema:author"></span>
				</head>
				<body>
				</body>
				</html>
			`,
			expected: metadata.Metadata{},
		},
		{
			desc: "Null Values",
			mockHTML: `
				<!DOCTYPE html>
				<html>
				<head>
					<meta property="og:title"/>
					<meta property="og:description"/>
					<meta property="article:published_time"/>
					<span property="schema:author"></span>
				</head>
				<body>
				</body>
				</html>
			`,
			expected: metadata.Metadata{},
		},
		{
			desc: "Malformed HTML",
			mockHTML: `
				<!DOCTYPE html>
				<html>
				<head>
					<meta property="og:title" content="OG Title">
				</head
				<body>
				</body>
			`,
			expected: metadata.Metadata{
				Title: "OG Title",
			},
		},
		// Add more test cases as needed
	}

	ext := gophetch.NewExtractor()
	targetURL, err := url.Parse("https://example.com")
	if err != nil {
		t.Fatal(err)
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// Parse the mock HTML
			mockNode, err := html.Parse(strings.NewReader(tC.mockHTML))
			if err != nil {
				t.Fatal(err)
			}

			// Call the Extractor's method to get all metadata
			result, err := ext.ExtractMetadata(mockNode, targetURL)
			if err != nil {
				t.Fatal(err)
			}

			// Assertion
			assert.Equal(t, tC.expected.Title, result.Title, "They should be equal")
		})
	}
}
