package engine

import (
	"sort"
	"testing"
)

func TestBindingRegisterAffected(t *testing.T) {
	bt := NewBindingTracker()

	bt.Register("/name", "text1")
	bt.Register("/name", "text2")
	bt.Register("/email", "text3")

	affected := bt.Affected([]string{"/name"})
	sort.Strings(affected)
	if len(affected) != 2 || affected[0] != "text1" || affected[1] != "text2" {
		t.Errorf("affected = %v, want [text1 text2]", affected)
	}
}

func TestBindingUnregister(t *testing.T) {
	bt := NewBindingTracker()

	bt.Register("/name", "text1")
	bt.Register("/name", "text2")
	bt.Unregister("text1")

	affected := bt.Affected([]string{"/name"})
	if len(affected) != 1 || affected[0] != "text2" {
		t.Errorf("affected = %v, want [text2]", affected)
	}
}

func TestBindingPathOverlap(t *testing.T) {
	bt := NewBindingTracker()

	bt.Register("/user/name", "comp1")

	// Changing parent should affect children bindings
	affected := bt.Affected([]string{"/user"})
	if len(affected) != 1 || affected[0] != "comp1" {
		t.Errorf("parent change: affected = %v, want [comp1]", affected)
	}

	// Changing child should affect parent bindings
	bt2 := NewBindingTracker()
	bt2.Register("/user", "comp2")
	affected2 := bt2.Affected([]string{"/user/name"})
	if len(affected2) != 1 || affected2[0] != "comp2" {
		t.Errorf("child change: affected = %v, want [comp2]", affected2)
	}
}

func TestBindingNoOverlap(t *testing.T) {
	bt := NewBindingTracker()

	bt.Register("/name", "comp1")

	affected := bt.Affected([]string{"/email"})
	if len(affected) != 0 {
		t.Errorf("unrelated path: affected = %v, want []", affected)
	}
}

func TestBindingMultiplePaths(t *testing.T) {
	bt := NewBindingTracker()

	bt.Register("/a", "comp1")
	bt.Register("/b", "comp2")
	bt.Register("/a", "comp3")

	affected := bt.Affected([]string{"/a", "/b"})
	sort.Strings(affected)
	if len(affected) != 3 {
		t.Errorf("affected = %v, want 3 items", affected)
	}
}
