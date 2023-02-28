package client

//go:generate stringer -type=MatchingAlgorithm -output=matching_string.go
type MatchingAlgorithm int

const (
	MatchUnspecified MatchingAlgorithm = iota

	// Any word.
	MatchAny

	// All words.
	MatchAll

	// Exact match.
	MatchLiteral

	// Regular expression.
	MatchRegex

	// Fuzzy word.
	MatchFuzzy

	// Automatic using a document classification model.
	MatchAuto
)
