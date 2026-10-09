package theme

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultExists(t *testing.T) {
	if _, ok := Get(Default); !ok {
		t.Fatalf("default theme %q missing", Default)
	}
	if len(All()) < 100 {
		t.Fatalf("expected a large theme list, got %d", len(All()))
	}
}

func TestGetNormalizes(t *testing.T) {
	for _, name := range []string{"serika_dark", "Serika Dark", "  serika dark "} {
		if th, ok := Get(name); !ok || th.Name != "serika_dark" {
			t.Fatalf("Get(%q) = %v, %v", name, th.Name, ok)
		}
	}
	if _, ok := Get("nope"); ok {
		t.Fatal("want miss for unknown theme")
	}
}

func TestFilterRanksPrefixFirst(t *testing.T) {
	ms := Filter(All(), "nord")
	if len(ms) == 0 || ms[0].Theme.Name != "nord" {
		t.Fatalf("want nord first, got %+v", ms)
	}
	ms = Filter(All(), "srkdk")
	if len(ms) == 0 || ms[0].Theme.Name != "serika_dark" {
		t.Fatalf("want serika_dark first for scattered query, got %v", ms[0].Theme.Name)
	}
	for _, m := range ms {
		if len(m.Indices) != len("srkdk") {
			t.Fatalf("indices %v for %s", m.Indices, m.Theme.Name)
		}
	}
	if got := Filter(All(), ""); len(got) != len(All()) {
		t.Fatalf("empty query should return all, got %d", len(got))
	}
	if got := Filter(All(), "zzzzzz"); len(got) != 0 {
		t.Fatalf("want no matches, got %d", len(got))
	}
}

func TestSaveLoad(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	if Load().Name != Default {
		t.Fatal("want default with no config")
	}
	if err := Save("dracula"); err != nil {
		t.Fatal(err)
	}
	if got := Load().Name; got != "dracula" {
		t.Fatalf("got %q, want dracula", got)
	}
	if _, err := os.Stat(filepath.Join(dir, "termitype", "config.json")); err != nil {
		t.Fatal(err)
	}
	if err := Save("not_a_theme"); err == nil {
		t.Fatal("want error saving unknown theme")
	}
}
