package theme

import (
	"sort"
	"strings"
)

// Match is a theme that matched a fuzzy query.
type Match struct {
	Theme   Theme
	Score   int
	Indices []int // byte positions in Theme.Name that matched the query
}

// Filter fuzzy-matches query against theme names and returns matches,
// best first. An empty query returns all themes in name order.
//
// Matching is subsequence based (fzf-style): every query rune must appear in
// order. Consecutive runs and matches at word starts score higher.
func Filter(themes []Theme, query string) []Match {
	q := normalize(query)
	out := make([]Match, 0, len(themes))
	for _, t := range themes {
		if score, idx, ok := fuzzyScore(t.Name, q); ok {
			out = append(out, Match{Theme: t, Score: score, Indices: idx})
		}
	}
	if q != "" {
		sort.SliceStable(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	}
	return out
}

func fuzzyScore(name, q string) (int, []int, bool) {
	if q == "" {
		return 0, nil, true
	}
	// Fast path: substring matches beat any scattered match.
	if i := strings.Index(name, q); i >= 0 {
		idx := make([]int, len(q))
		for k := range idx {
			idx[k] = i + k
		}
		score := 1000 - i - len(name)
		if i == 0 || name[i-1] == '_' {
			score += 100
		}
		return score, idx, true
	}

	idx := make([]int, 0, len(q))
	score, qi, prev := 0, 0, -2
	for ni := 0; ni < len(name) && qi < len(q); ni++ {
		if name[ni] != q[qi] {
			continue
		}
		switch {
		case ni == prev+1:
			score += 10 // consecutive
		case ni == 0 || name[ni-1] == '_':
			score += 8 // word start
		default:
			score -= ni - prev // gap penalty
		}
		idx = append(idx, ni)
		prev = ni
		qi++
	}
	if qi < len(q) {
		return 0, nil, false
	}
	return score - len(name), idx, true
}
