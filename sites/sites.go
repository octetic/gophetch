package sites

import (
	"github.com/pixiesys/gophetch/rules"
)

type Site interface {
	Rules() map[string]rules.Rule
	DomainKey() string
}

var siteRegistry = map[string]Site{
	"youtube.com": &YouTube{},
}

func GetSite(domain string) (Site, bool) {
	site, found := siteRegistry[domain]
	return site, found
}

func RegisterSite(domain string, site Site) {
	siteRegistry[domain] = site
}
