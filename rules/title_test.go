package rules_test

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"

	"github.com/minsoft-io/gophetch/rules"
)

func TestTitleRuleSelectors(t *testing.T) {
	testCases := []struct {
		desc     string
		mockHTML string
		expected string
		error    error
	}{
		{
			desc:     "Test with og:title in meta",
			mockHTML: `<meta property="og:title" content="OG Title"/>`,
			expected: "OG Title",
		},
		{
			desc:     "Test with title tag",
			mockHTML: `<title>HTML Title</title>`,
			expected: "HTML Title",
		},
		{
			desc:     "Test with headline in JSON-LD",
			mockHTML: `<script type="application/ld+json">{"headline": "JSON-LD Title"}</script>`,
			expected: "JSON-LD Title",
		},
		{
			desc:     "Test with .post-title class",
			mockHTML: `<div class="post-title">Post Title</div>`,
			expected: "Post Title",
		},
		{
			desc: "Test with multiple selectors, prioritizing og:title",
			mockHTML: `
		        <html>
		        <head>	
				<meta property="og:title" content="Priority OG Title"/>
				<title>HTML Title</title>
				</head>
				<body>
				<span class="post-title">Post Title</span>
				</body>
				</html>
			`,
			expected: "Priority OG Title",
		},
		{
			desc: "Test no value found",
			mockHTML: `
				<meta property="foo:bar" content="2022-01-01T15:04:05Z"/>
			`,
			expected: "",
			error:    rules.ErrValueNotFound,
		},
	}

	tr := rules.NewTitleRule()
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

			// Call the TitleRule's Extract method
			result, err := tr.Extract(mockNode, targetURL)
			if err != nil {
				assert.Equal(t, tC.error, err, fmt.Sprintf("Want error %v, got %v", tC.error, err))
			} else {
				assert.Equal(t, tC.expected, result.Value(), fmt.Sprintf("Want %s, got %s", tC.expected, result.Value()))
			}

		})
	}
}
