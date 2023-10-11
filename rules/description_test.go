package rules_test

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
	// import your own packages here

	"github.com/pixiesys/gophetch/rules"
)

func TestDescriptionRuleSelectors(t *testing.T) {
	testCases := []struct {
		desc     string
		mockHTML string
		expected []string
		error    error
	}{
		{
			desc:     "Test with og:description in meta",
			mockHTML: `<meta property="og:description" content="OG Description"/>`,
			expected: []string{"OG Description"},
		},
		{
			desc:     "Test with description in JSON-LD",
			mockHTML: `<script type="application/ld+json">{"description": "JSON-LD Description"}</script>`,
			expected: []string{"JSON-LD Description"},
		},
		{
			desc:     "Test with .post-description class",
			mockHTML: `<div class="post-description">Post Description</div>`,
			expected: []string{"Post Description"},
		},
		{
			desc: "Test with multiple selectors, prioritizing og:description",
			mockHTML: `
				<meta property="og:description" content="Priority OG Description"/>
				<div class="post-description">Post Description</div>
			`,
			expected: []string{"Priority OG Description"},
		},
		{
			desc: "Test no description found",
			mockHTML: `
				<div class="post-foo">Post Description</div>
			`,
			expected: []string{},
			error:    rules.ErrValueNotFound,
		},
	}

	dr := rules.NewDescriptionRule( /* initialize with appropriate params */ )
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

			// Call the DescriptionRule's Extract method
			result, err := dr.Extract(mockNode, targetURL)
			if err != nil {
				assert.Equal(t, tC.error, err, fmt.Sprintf("Want error %v, got %v", tC.error, err))
			} else {
				assert.Equal(t, tC.expected, result, fmt.Sprintf("Want %s, got %s", tC.expected, result))
			}
		})
	}
}
