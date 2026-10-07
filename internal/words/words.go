// Package words provides word lists and random test generation.
package words

import (
	"math/rand"
	"strings"
	"time"
)

// Official Monkeytype default English wordlist (english.json, 200 words).
// Source: https://github.com/monkeytypegame/monkeytype/blob/master/frontend/static/languages/english.json
var common = []string{
	"the", "be", "of", "and", "a", "to", "in", "he", "have", "it",
	"that", "for", "they", "I", "with", "as", "not", "on", "she", "at",
	"by", "this", "we", "you", "do", "but", "from", "or", "which", "one",
	"would", "all", "will", "there", "say", "who", "make", "when", "can", "more",
	"if", "no", "man", "out", "other", "so", "what", "time", "up", "go",
	"about", "than", "into", "could", "state", "only", "new", "year", "some", "take",
	"come", "these", "know", "see", "use", "get", "like", "then", "first", "any",
	"work", "now", "may", "such", "give", "over", "think", "most", "even", "find",
	"day", "also", "after", "way", "many", "must", "look", "before", "great", "back",
	"through", "long", "where", "much", "should", "well", "people", "down", "own", "just",
	"because", "good", "each", "those", "feel", "seem", "how", "high", "too", "place",
	"little", "world", "very", "still", "nation", "hand", "old", "life", "tell", "write",
	"become", "here", "show", "house", "both", "between", "need", "mean", "call", "develop",
	"under", "last", "right", "move", "thing", "general", "school", "never", "same", "another",
	"begin", "while", "number", "part", "turn", "real", "leave", "might", "want", "point",
	"form", "off", "child", "few", "small", "since", "against", "ask", "late", "home",
	"interest", "large", "person", "end", "open", "public", "follow", "during", "present", "without",
	"again", "hold", "govern", "around", "possible", "head", "consider", "word", "program", "problem",
	"however", "lead", "system", "set", "order", "eye", "plan", "run", "keep", "face",
	"fact", "group", "play", "stand", "increase", "early", "course", "change", "help", "line",
}

// Generate returns n random words joined for a typing test.
func Generate(n int) string {
	return strings.Join(GenerateSlice(n), " ")
}

// GenerateSlice returns n random words.
func GenerateSlice(n int) []string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	out := make([]string, n)
	for i := range out {
		out[i] = common[r.Intn(len(common))]
	}
	return out
}
