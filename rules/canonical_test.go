package rules_test

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"

	"github.com/octetic/gophetch/rules"
)

func TestCanonicalRuleSelectors(t *testing.T) {
	const targetURLStr = "https://example.com"

	testCases := []struct {
		desc     string
		mockHTML string
		expected string
		error    error
	}{
		{
			desc:     "Test with og:url in meta",
			mockHTML: `<meta property="og:url" content="https://example.com"/>`,
			expected: "https://example.com",
		},
		{
			desc:     "Test with name twitter:url in meta",
			mockHTML: `<meta name="twitter:url" content="https://example.com"/>`,
			expected: "https://example.com",
		},
		{
			desc:     "Test with property twitter:url in meta",
			mockHTML: `<meta property="twitter:url" content="https://example.com"/>`,
			expected: "https://example.com",
		},
		{
			desc:     "Test with rel canonical in link",
			mockHTML: `<link rel="canonical" href="https://example.com"/>`,
			expected: "https://example.com",
		},
		{
			desc:     "Test with rel alternate hreflang x-default in link",
			mockHTML: `<link rel="alternate" hreflang="x-default" href="https://example.com"/>`,
			expected: "https://example.com",
		},
		{
			desc:     "Test with multiple selectors, prioritizing og:url in meta",
			mockHTML: `<meta property="og:url" content="https://example.com"/><link rel="canonical" href="https://example.net"/>`,
			expected: "https://example.com",
		},
		{
			desc: "Test with multiple selectors, prioritizing og:url in meta",
			mockHTML: `
		        <html>
		        <head>	
		        <meta property="og:url" content="https://example.com"/>
		        <link rel="canonical" href="https://example.net"/>
				</head>
				<body>
				<span property="dc:creator">John DC</span>
				<span property="schema:author">Span Author</span>
				<span itemprop="author">John Microdata</span>
				<span class="author">Common Author</span>
				</body>
				</html>
			`,
			expected: "https://example.com",
		},
		{
			desc: "Test no value found",
			mockHTML: `
				<span property="foo:bar">John Foo</span>
				<span property="schema:foo">John Schema</span>
			`,
			expected: targetURLStr,
			error:    rules.ErrValueNotFound,
		},
	}

	cr := rules.NewCanonicalRule()
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
			result, err := cr.Extract(mockNode, targetURL)
			if err != nil {
				assert.Equal(t, tC.error, err, fmt.Sprintf("Want error %v, got %v", tC.error, err))
			} else {
				assert.Equal(t, tC.expected, result.Value(), fmt.Sprintf("Want %s, got %s", tC.expected, result.Value()))
			}
		})
	}
}
