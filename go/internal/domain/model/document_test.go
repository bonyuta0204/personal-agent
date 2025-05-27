package model

import (
	"reflect"
	"testing"
)

func TestSetTagsFromContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
		expect  []string
	}{
		{
			name:    "single tag",
			content: "This is a #test document.",
			expect:  []string{"test"},
		},
		{
			name:    "multiple tags",
			content: "#foo #bar Some text #baz",
			expect:  []string{"foo", "bar", "baz"},
		},
		{
			name:    "duplicate tags",
			content: "#dup #dup #unique",
			expect:  []string{"dup", "unique"},
		},
		{
			name:    "nested tags",
			content: "#project/ai #project/ml",
			expect:  []string{"project/ai", "project/ml"},
		},
		{
			name:    "underscore and hyphen",
			content: "#foo_bar #foo-bar",
			expect:  []string{"foo_bar", "foo-bar"},
		},
		{
			name:    "no tags",
			content: "No tags here.",
			expect:  []string{},
		},
		{
			name:    "multibyte tags",
			content: "日本語タグ #タグ #タグ/サブ",
			expect:  []string{"タグ", "タグ/サブ"},
		},
		// --- YAML frontmatter (array) ---
		{
			name:    "frontmatter tags (array)",
			content: "---\ntags:\n  - company\n  - ai\nstatus: Active\n---\n本文 #foo",
			expect:  []string{"company", "ai", "foo"},
		},
		// --- YAML frontmatter (string) ---
		{
			name:    "frontmatter tags (string)",
			content: "---\ntags: solo\nstatus: Rejected\n---\n#solo #extra",
			expect:  []string{"solo", "extra"},
		},
		// --- Both frontmatter and #tag, with overlap ---
		{
			name:    "frontmatter and #tag overlap",
			content: "---\ntags:\n  - overlap\n  - onlyfm\n---\n#overlap #onlytag",
			expect:  []string{"overlap", "onlyfm", "onlytag"},
		},
		// --- Edge: frontmatter present but no tags ---
		{
			name:    "frontmatter present, no tags",
			content: "---\nstatus: Only status\n---\n#foo",
			expect:  []string{"foo"},
		},
		// --- Edge: frontmatter tags empty array ---
		{
			name:    "frontmatter tags empty array",
			content: "---\ntags: []\n---\n#foo",
			expect:  []string{"foo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Document{Content: tt.content}
			d.SetTagsFromContent()
			// Check if all expected tags are present (order doesn't matter)
			if !equalStringSliceIgnoreOrder(d.Tags, tt.expect) {
				t.Errorf("got tags %v, want %v", d.Tags, tt.expect)
			}
		})
	}
}

func equalStringSliceIgnoreOrder(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	ma := make(map[string]int)
	mb := make(map[string]int)
	for _, v := range a {
		ma[v]++
	}
	for _, v := range b {
		mb[v]++
	}
	return reflect.DeepEqual(ma, mb)
}
