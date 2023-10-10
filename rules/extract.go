package rules

import (
	"strings"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

import (
	"encoding/json"
	// import other required packages
)

func ExtractJSONLD(node *html.Node, selectors []string) (string, bool) {
	jsonLdNodes := cascadia.QueryAll(node, cascadia.MustCompile(`script[type="application/ld+json"]`))
	for _, jsonLdNode := range jsonLdNodes {
		if jsonLdNode.FirstChild != nil {
			var obj map[string]interface{}
			err := json.Unmarshal([]byte(jsonLdNode.FirstChild.Data), &obj)
			if err != nil {
				continue
			}

			for _, selector := range selectors {
				keys := strings.Split(selector, ".")
				val := obj
				for _, key := range keys {
					if nextVal, ok := val[key].(map[string]interface{}); ok {
						val = nextVal
					} else if finalVal, ok := val[key].(string); ok {
						return finalVal, true
					} else {
						break
					}
				}
			}
		}
	}
	return "", false
}

func ExtractMeta(node *html.Node, selectors []string) (string, bool) {
	for _, selector := range selectors {
		metaNode := cascadia.Query(node, cascadia.MustCompile(selector))
		if metaNode != nil {
			for _, attr := range metaNode.Attr {
				if attr.Key == "content" {
					return attr.Val, true
				}
			}
		}
	}
	return "", false
}

func ExtractCSS(node *html.Node, selectors []string) (string, bool) {
	for _, selector := range selectors {
		cssNode := cascadia.Query(node, cascadia.MustCompile(selector))
		if cssNode != nil && cssNode.FirstChild != nil {
			return strings.TrimSpace(cssNode.FirstChild.Data), true
		}
	}
	return "", false
}

func ExtractTime(node *html.Node, selectors []string) (string, bool) {
	for _, selector := range selectors {
		cssNode := cascadia.Query(node, cascadia.MustCompile(selector))
		if cssNode != nil {
			for _, attr := range cssNode.Attr {
				if attr.Key == "datetime" {
					return attr.Val, true
				}
			}
		}
	}
	return "", false
}
