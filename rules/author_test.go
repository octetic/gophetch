package rules_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"

	"github.com/pixiesys/gophetch/rules"
)

func TestAuthorRuleSelectors(t *testing.T) {
	testCases := []struct {
		desc     string
		mockHTML string
		expected string
		error    error
	}{
		{
			desc:     "Test with valid author.name in JSON-LD",
			mockHTML: `<script type="application/ld+json">{"author": {"name": "John JSON-LD"}}</script>`,
			expected: "John JSON-LD",
		},
		{
			desc:     "Test with brand.name in JSON-LD",
			mockHTML: `<script type="application/ld+json">{"brand": {"name": "Brand JSON-LD"}}</script>`,
			expected: "Brand JSON-LD",
		},
		{
			desc:     "Test with creator.name in JSON-LD",
			mockHTML: `<script type="application/ld+json">{"creator": {"name": "Creator JSON-LD"}}</script>`,
			expected: "Creator JSON-LD",
		},
		{
			desc: "Test with multiple selectors, prioritizing JSON-LD",
			mockHTML: `
				<script type="application/ld+json">{"author": {"name": "Priority JSON-LD"}}</script>
				<meta property="schema:author">John Schema</span>
			`,
			expected: "Priority JSON-LD",
		},
		{
			desc:     "Test with schema:author in span",
			mockHTML: `<span property="schema:author">John Schema</span>`,
			expected: "John Schema",
		},
		{
			desc:     "Test with schema:author in meta",
			mockHTML: `<meta property="schema:author" content="Meta Schema"/>`,
			expected: "Meta Schema",
		},
		{
			desc:     "Test with typeof schema:Person",
			mockHTML: `<div typeof="schema:Person"><span property="schema:name">Person Schema</span></div>`,
			expected: "Person Schema",
		},
		{
			desc:     "Test with dc:creator in span",
			mockHTML: `<span property="dc:creator">John DC</span>`,
			expected: "John DC",
		},
		{
			desc: "Test with multiple selectors, prioritizing schema:author in meta",
			mockHTML: `
		        <html>
		        <head>	
		        <meta property="schema:author" content="Priority Author"/>	
				</head>
				<body>
				<span property="dc:creator">John DC</span>
				<span property="schema:author">Span Author</span>
				<span itemprop="author">John Microdata</span>
				<span class="author">Common Author</span>
				</body>
				</html>
			`,
			expected: "Priority Author",
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

	ar := rules.NewAuthorRule()

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// Parse the mock HTML
			mockNode, err := html.Parse(strings.NewReader(tC.mockHTML))
			if err != nil {
				t.Fatal(err)
			}

			// Call the AuthorRule's Extract method
			result, err := ar.Extract(mockNode)
			if err != nil {
				assert.Equal(t, tC.error, err, fmt.Sprintf("Want error %v, got %v", tC.error, err))
			} else {
				assert.Equal(t, tC.expected, result, fmt.Sprintf("Want %s, got %s", tC.expected, result))
			}
		})
	}
}
