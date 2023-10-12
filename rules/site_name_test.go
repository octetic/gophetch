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

func TestSiteNameRuleSelectors(t *testing.T) {
	const targetURLStr = "https://www.example.com"

	testCases := []struct {
		desc     string
		mockHTML string
		expected string
		error    error
	}{
		{
			desc:     "Test with og:site_name in meta",
			mockHTML: `<meta property="og:site_name" content="Example"/>`,
			expected: "Example",
		},
		{
			desc:     "Test with name og:site_name in meta",
			mockHTML: `<meta name="og:site_name" content="Example"/>`,
			expected: "Example",
		},
		{
			desc:     "Test with twitter:site_name in meta",
			mockHTML: `<meta property="twitter:site_name" content="Example"/>`,
			expected: "Example",
		},
		{
			desc:     "Test with name twitter:site_name in meta",
			mockHTML: `<meta name="twitter:site_name" content="Example"/>`,
			expected: "Example",
		},
		{
			desc:     "Test with itemprop name in meta",
			mockHTML: `<meta itemprop="name" content="Example"/>`,
			expected: "Example",
		},
		{
			desc:     "Test with name application-name in meta",
			mockHTML: `<meta name="application-name" content="Example"/>`,
			expected: "Example",
		},
		{
			desc:     "Test with multiple selectors, prioritizing og:site_name in meta",
			mockHTML: `<meta property="og:site_name" content="Example"/><meta name="twitter:site_name" content="Example2"/>`,
			expected: "Example",
		},
		{
			desc: "Test with multiple selectors, prioritizing og:site_name in meta",
			mockHTML: `
		        <html>
		        <head>	
		        <meta property="og:site_name" content="Example"/>
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
			expected: "Example",
		},
		{
			desc: "Test no value found",
			mockHTML: `
				<span property="foo:bar">John Foo</span>
				<span property="schema:foo">John Schema</span>
			`,
			expected: "example.com",
			error:    rules.ErrValueNotFound,
		},
	}

	snr := rules.NewSiteNameRule()
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
			result, err := snr.Extract(mockNode, targetURL)
			if err != nil {
				assert.Equal(t, tC.error, err, fmt.Sprintf("Want error %v, got %v", tC.error, err))
			} else {
				assert.Equal(t, tC.expected, result.Value(), fmt.Sprintf("Want %s, got %s", tC.expected, result.Value()))
			}
		})
	}
}
