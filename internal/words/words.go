// Package words provides word lists and random test generation.
package words

import (
	"math/rand"
	"strings"
	"time"
)

// Common English words for typing tests. Kept small for simplicity;
// expand as desired.
var common = []string{
	"the", "be", "of", "and", "a", "to", "in", "he", "have", "it",
	"that", "for", "they", "with", "as", "not", "on", "she", "at", "by",
	"this", "we", "you", "do", "but", "from", "or", "which", "one", "would",
	"all", "will", "there", "say", "who", "make", "when", "can", "more", "if",
	"time", "up", "out", "who", "into", "way", "could", "my", "than", "first",
	"water", "been", "call", "oil", "its", "now", "find", "long", "down", "day",
	"did", "get", "come", "made", "may", "part", "over", "new", "sound", "take",
	"only", "little", "work", "know", "place", "year", "live", "me", "back", "give",
	"most", "very", "after", "thing", "our", "just", "name", "good", "sentence", "man",
	"think", "help", "low", "line", "differ", "turn", "cause", "much", "mean", "before",
	"move", "right", "boy", "old", "too", "same", "tell", "does", "set", "three",
	"want", "air", "well", "also", "play", "small", "end", "put", "home", "read",
	"hand", "port", "large", "spell", "add", "even", "land", "here", "must", "big",
	"high", "such", "follow", "act", "why", "ask", "men", "change", "went", "light",
	"kind", "off", "need", "house", "picture", "try", "again", "animal", "point", "mother",
	"world", "near", "build", "self", "earth", "father", "head", "stand", "own", "page",
	"should", "country", "found", "answer", "school", "grow", "study", "still", "learn", "plant",
	"cover", "food", "sun", "four", "between", "state", "keep", "eye", "never", "last",
	"let", "thought", "city", "tree", "cross", "farm", "hard", "start", "might", "story",
	"code", "keyboard", "terminal", "fast", "type", "monkey", "speed", "focus", "rhythm", "flow",
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
