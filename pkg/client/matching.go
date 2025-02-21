package client

//go:generate go run golang.org/x/tools/cmd/stringer -type=MatchingAlgorithm -output=matching_string.go
type MatchingAlgorithm int

const (
	// None.
	MatchNone MatchingAlgorithm = iota

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
