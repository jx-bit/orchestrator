package childenv

import "testing"

func TestGOMAXPROCSAddsCapWhenUnset(t *testing.T) {
	env := []string{"PATH=/bin", "MANIFEST_KEY=abc"}
	got := GOMAXPROCS(env, "4")
	want := []string{"PATH=/bin", "MANIFEST_KEY=abc", "GOMAXPROCS=4"}
	if len(got) != len(want) {
		t.Fatalf("GOMAXPROCS(%v) = %v, want %v", env, got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("GOMAXPROCS(%v) = %v, want %v", env, got, want)
		}
	}
}

func TestGOMAXPROCSRespectsExplicitSetting(t *testing.T) {
	env := []string{"PATH=/bin", "GOMAXPROCS=2"}
	got := GOMAXPROCS(env, "4")
	if len(got) != len(env) {
		t.Fatalf("GOMAXPROCS(%v) = %v, want env unchanged", env, got)
	}
	for i := range env {
		if got[i] != env[i] {
			t.Fatalf("GOMAXPROCS(%v) = %v, want env unchanged", env, got)
		}
	}
}
