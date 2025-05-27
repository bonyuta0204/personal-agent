package model

import (
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type DocumentId string

// represent a document in the knowledge base
type Document struct {
	ID        DocumentId
	StoreId   StoreId
	Path      string
	Content   string
	Embedding []float64
	Tags      []string
	SHA       string

	ModifiedAt time.Time // The time when the document was last modified. This is used to detect changes in the document.
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// represent a document entry in the knowledge base
type DocumentEntry struct {
	Path       string
	ModifiedAt time.Time
}

// set document tag from its content (Obsidian markdown style, including YAML frontmatter)
func (d *Document) SetTagsFromContent() {
	tagSet := make(map[string]struct{})

	// 1. Extract tags from YAML frontmatter if present
	content := d.Content
	if strings.HasPrefix(content, "---") {
		end := strings.Index(content[3:], "---")
		if end != -1 {
			yamlBlock := content[3 : 3+end]
			var fm map[string]interface{}
			if err := yaml.Unmarshal([]byte(yamlBlock), &fm); err == nil {
				if tags, ok := fm["tags"]; ok {
					switch v := tags.(type) {
					case []interface{}:
						for _, t := range v {
							if tagStr, ok := t.(string); ok {
								tagSet[tagStr] = struct{}{}
							}
						}
					case string:
						tagSet[v] = struct{}{}
					}
				}
			}
		}
	}

	// 2. Extract #tag style tags from the whole content
	re := regexp.MustCompile(`#([\p{L}\p{N}_\-/]+)`)
	matches := re.FindAllStringSubmatch(content, -1)
	for _, m := range matches {
		if len(m) > 1 {
			tag := m[1]
			tagSet[tag] = struct{}{}
		}
	}

	// 3. Set to d.Tags (deduplicated)
	d.Tags = make([]string, 0, len(tagSet))
	for tag := range tagSet {
		d.Tags = append(d.Tags, tag)
	}
}
