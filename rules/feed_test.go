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

func TestFeedRuleSelectors(t *testing.T) {
	const targetURLStr = "https://example.com"

	testCases := []struct {
		desc     string
		mockHTML string
		expected []string
		error    error
	}{
		{
			desc:     "Test with atom in link",
			mockHTML: `<link rel="alternate" type="application/atom+xml" href="https://example.com/atom.xml"/>`,
			expected: []string{"https://example.com/atom.xml"},
		},
		{
			desc:     "Test with rss in link",
			mockHTML: `<link rel="alternate" type="application/rss+xml" href="https://example.com/rss.xml"/>`,
			expected: []string{"https://example.com/rss.xml"},
		},
		{
			desc:     "Test with multiple selectors, returning all feed URLs",
			mockHTML: `<link rel="alternate" type="application/atom+xml" href="https://example.com/atom.xml"/><link rel="alternate" type="application/rss+xml" href="https://example.com/rss.xml"/>`,
			expected: []string{"https://example.com/rss.xml", "https://example.com/atom.xml"},
		},
		{
			desc: "Test with multiple selectors, prioritizing icon in meta",
			mockHTML: `
		        <html>
		        <head>	
				 <link rel="alternate" type="application/atom+xml" href="https://example.com/atom.xml"/>
				 <link rel="alternate" type="application/rss+xml" href="https://example.com/rss.xml"/>
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
			expected: []string{"https://example.com/rss.xml", "https://example.com/atom.xml"},
		},
		{
			desc: "Test no value found",
			mockHTML: `
				<span property="foo:bar">John Foo</span>
				<span property="schema:foo">John Schema</span>
			`,
			expected: []string{},
			error:    rules.ErrValueNotFound,
		},
	}

	fr := rules.NewFeedRule()
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
