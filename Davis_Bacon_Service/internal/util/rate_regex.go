package util

import "regexp"

// RatePattern is a compiled regex for parsing wage/fringe lines.
var RatePattern = regexp.MustCompile(
	`^(?P<label>[A-Za-z0-9\s/\-\(\):]+?)\.*\s*\$?(?P<base>\d+\.\d{2}|\*\*)\s+(?P<fringe>\d+\.\d{2}|\*\*)`,
)
