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

func TestDateRuleSelectors(t *testing.T) {
	testCases := []struct {
		desc     string
		mockHTML string
		expected string
		error    error
	}{
		{
			desc:     "Test with datePublished in JSON-LD",
			mockHTML: `<script type="application/ld+json">{"datePublished": "2022-10-11"}</script>`,
			expected: "2022-10-11",
		},
		{
			desc:     "Test with article:published_time in meta",
			mockHTML: `<meta property="article:published_time" content="2022-10-11T15:04:05Z"/>`,
			expected: "2022-10-11T15:04:05Z",
		},
		{
			desc:     "Test with time element",
			mockHTML: `<time datetime="2022-10-11T15:04:05Z"></time>`,
			expected: "2022-10-11T15:04:05Z",
		},
		{
			desc:     "Test with .post-date class",
			mockHTML: `<div class="post-date">2022-10-11</div>`,
			expected: "2022-10-11",
		},
		{
			desc: "Test with multiple selectors, prioritizing datePublished in JSON-LD",
			mockHTML: `
				<script type="application/ld+json">{"datePublished": "2022-10-11"}</script>
				<meta property="article:published_time" content="2022-01-01T15:04:05Z"/>
			`,
			expected: "2022-10-11",
		},
		{
			desc: "Test with multiple selectors, prioritizing article:published_time in meta",
			mockHTML: `
				<meta property="article:published_time" content="2022-01-01T15:04:05Z"/>
				<time datetime="2022-10-11T15:04:05Z"></time>	
			`,
			expected: "2022-01-01T15:04:05Z",
		},
		{
			desc: "Test no value found",
			mockHTML: `
				<meta property="foo:bar" content="2022-01-01T15:04:05Z"/>
				<time foobar="2022-10-11T15:04:05Z"></time>
				<div class="foo">2022-10-11</div>
			`,
			expected: "",
			error:    rules.ErrValueNotFound,
		},
	}

	dr := rules.NewDateRule( /* initialize with appropriate params */ )
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

			// Call the DateRule's Extract method
			result, err := dr.Extract(mockNode, targetURL)
			if err != nil {
				assert.Equal(t, tC.error, err, fmt.Sprintf("Want error %v, got %v", tC.error, err))
			} else {
				assert.Equal(t, tC.expected, result.Value(), fmt.Sprintf("Want %s, got %s", tC.expected, result.Value()))
			}
		})
	}
}
