package sites

import (
	"github.com/minsoft-io/gophetch/rules"
)

// Site is the interface for a site. Sites are used to define the rules for a specific site. This can be used to
// more accurately extract information from a page when the URL is known and the structure of the page is known.
type Site interface {
	// Rules returns the rules for the site. These rules are used to extract information from a page and should have
	// the same keys as the keys in the gophetch.Extractor.
	Rules() map[string]rules.Rule

	// DomainKey returns the domain key for the site. This is used to match the site to a URL. In most cases this is the
	// domain of the site without the protocol. For example, "www.example.com" or "example.com". However, this can be
	// anything that can be used to match a URL to a site.
	DomainKey() string
}
