package rules_test

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"

	"github.com/pixiesys/gophetch/rules"
)

func TestFaviconRuleSelectors(t *testing.T) {
	const targetURLStr = "https://example.com"

	testCases := []struct {
		desc     string
		mockHTML string
		expected string
		error    error
	}{
		{
			desc:     "Test with apple-touch-icon in link",
			mockHTML: `<link rel="apple-touch-icon" href="https://example.com/apple-touch-icon.png"/>`,
			expected: "https://example.com/apple-touch-icon.png",
		},
		{
			desc:     "Test with apple-touch-icon-precomposed in link",
			mockHTML: `<link rel="apple-touch-icon-precomposed" href="https://example.com/apple-touch-icon-precomposed.png"/>`,
			expected: "https://example.com/apple-touch-icon-precomposed.png",
		},
		{
			desc:     "Test with icon in link",
			mockHTML: `<link rel="icon" href="https://example.com/icon.png"/>`,
			expected: "https://example.com/icon.png",
		},
		{
			desc:     "Test with shortcut icon in link",
			mockHTML: `<link rel="shortcut icon" href="https://example.com/shortcut-icon.png"/>`,
			expected: "https://example.com/shortcut-icon.png",
		},
		{
			desc:     "Test with mask-icon in link",
			mockHTML: `<link rel="mask-icon" href="https://example.com/mask-icon.png"/>`,
			expected: "https://example.com/mask-icon.png",
		},
		{
			desc:     "Test with multiple selectors, prioritizing icon in link",
			mockHTML: `<link rel="apple-touch-icon" href="https://example.com/apple-touch-icon.png"/><link rel="icon" href="https://example.com/icon.png"/>`,
			expected: "https://example.com/icon.png",
		},
		{
			desc: "Test with multiple selectors, prioritizing icon in meta",
			mockHTML: `
		        <html>
		        <head>	
		        <link rel="canonical" href="https://example.net"/>
			    <link rel="apple-touch-icon" href="https://example.com/apple-touch-icon.png"/>
		        <link rel="icon" href="https://example.com/icon.png"/>
				<link rel="shortcut icon" href="https://example.com/shortcut-icon.png"/>
				</head>
				<body>
				<span property="dc:creator">John DC</span>
				<span property="schema:author">Span Author</span>
				<span itemprop="author">John Microdata</span>
				<span class="author">Common Author</span>
				</body>
				</html>
			`,
			expected: "https://example.com/icon.png",
		},
		{
			desc: "Test no value found",
			mockHTML: `
				<span property="foo:bar">John Foo</span>
				<span property="schema:foo">John Schema</span>
			`,
			expected: "",
			error:    rules.ErrValueNotFound,
		},
	}

	fr := rules.NewFaviconRule()
	targetURL, err := url.Parse(targetURLStr)
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

			// Call the AuthorRule's Extract method
			result, err := fr.Extract(mockNode, targetURL)
			if err != nil {
				assert.Equal(t, tC.error, err, fmt.Sprintf("Want error %v, got %v", tC.error, err))
			} else {
				assert.Equal(t, tC.expected, result.Value(), fmt.Sprintf("Want %s, got %s", tC.expected, result.Value()))
			}
		})
	}
}
