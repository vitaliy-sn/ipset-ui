package ipset

import (
	"ipset-ui/internal/utils"
	"strings"
)

// FilterEntries filters entries by CIDR, IP, or comment substring.
// If filter is empty, returns entries as is.
func FilterEntries(entries []EntryWithComment, filter string) []EntryWithComment {
	if filter == "" {
		return entries
	}

	var filteredEntries []EntryWithComment
	if utils.IsValidCIDR(filter) {
		// If the filter is a CIDR, find matching CIDRs and CIDRs that include this CIDR
		for _, entry := range entries {
			if utils.IsValidCIDR(entry.Entry) {
				if utils.IsCIDRInCIDR(entry.Entry, filter) || entry.Entry == filter {
					filteredEntries = append(filteredEntries, entry)
				}
			}
		}
	} else if utils.IsValidIP(filter) {
		// If the filter is an IP address, find matching IP addresses and all CIDRs that include this IP address
		for _, entry := range entries {
			if utils.IsValidIP(entry.Entry) && entry.Entry == filter {
				filteredEntries = append(filteredEntries, entry)
			} else if utils.IsValidCIDR(entry.Entry) && utils.IsIPInCIDR(filter, entry.Entry) {
				filteredEntries = append(filteredEntries, entry)
			}
		}
	} else {
		// Otherwise, filter by substring in comment (case-insensitive)
		for _, entry := range entries {
			if entry.Comment != "" && strings.Contains(strings.ToLower(entry.Comment), strings.ToLower(filter)) {
				filteredEntries = append(filteredEntries, entry)
			}
		}
	}
	return filteredEntries
}
