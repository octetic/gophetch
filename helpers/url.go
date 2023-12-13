package helpers

import (
	"net/url"
)

// Params is a list of tracking parameters to remove from URLs.
// Original list from: https://github.com/mpchadwick/tracking-query-params-registry
var params = []string{"fbclid", "gclid", "gclsrc", "utm_content", "utm_term", "utm_campaign", "utm_medium",
	"utm_source", "utm_id", "_ga", "mc_cid", "mc_eid", "_bta_tid", "_bta_c", "trk_contact", "trk_msg", "trk_module",
	"trk_sid", "gdfms", "gdftrk", "gdffi", "_ke", "redirect_log_mongo_id", "redirect_mongo_id", "sb_referer_host",
	"mkwid", "pcrid", "ef_id", "s_kwcid", "msclkid", "dm_i", "epik", "pk_campaign", "pk_kwd", "pk_keyword",
	"piwik_campaign", "piwik_kwd", "piwik_keyword", "mtm_campaign", "mtm_keyword", "mtm_source", "mtm_medium",
	"mtm_content", "mtm_cid", "mtm_group", "mtm_placement", "matomo_campaign", "matomo_keyword", "matomo_source",
	"matomo_medium", "matomo_content", "matomo_cid", "matomo_group", "matomo_placement", "hsa_cam", "hsa_grp",
	"hsa_mt", "hsa_src", "hsa_ad", "hsa_acc", "hsa_net", "hsa_kw", "hsa_tgt", "hsa_ver", "_branch_match_id", "mkevt",
	"mkcid", "mkrid", "campid", "toolid", "customid", "igshid", "si"}

// CleanURL removes tracking parameters from the given URL.
func CleanURL(u string) string {
	nu, err := url.Parse(u)
	if err != nil {
		return u
	}

	for _, param := range params {
		nu = deleteParam(nu, param)
	}

	return nu.String()
}

// IsURLValid checks if the given URL is valid.
func IsURLValid(u string) bool {
	_, err := url.ParseRequestURI(u)

	// Ensure the url starts with a valid scheme
	if err == nil {
		for _, scheme := range []string{"http", "https"} {
			if u[:len(scheme)+1] == scheme+":" {
				err = nil
				break
			}
		}
	}

	return err == nil
}

func deleteParam(u *url.URL, key string) *url.URL {
	nu := *u
	values := nu.Query()

	values.Del(key)

	nu.RawQuery = values.Encode()
	return &nu
}
