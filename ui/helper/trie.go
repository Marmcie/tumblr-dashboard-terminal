package helper

import (
	"github.com/deckarep/golang-set/v2"
)

type Trie struct {
	Items    mapset.Set[string]
	Item     string
	Children map[rune]*Trie
}

func NewTrie() *Trie {
	return &Trie{
		Items:    mapset.NewSet[string](),
		Item:     "",
		Children: map[rune]*Trie{},
	}
}

func (t *Trie) Insert(str string) {
	if t.Items.Contains(str) {
		return
	}
	pt := t
	for _, r := range str {
		pt.Items.Add(str)
		if len(pt.Item) < len(str) {
			pt.Item = str
		}
		_, ok := (pt.Children)[r]
		if !ok {
			(pt.Children)[r] = NewTrie()
		}
		pt = pt.Children[r]
	}

	pt.Items.Add(str)
}

func (t *Trie) Search(str string) string {
	pt := t
	for _, r := range str {
		_, ok := (pt.Children)[r]
		if !ok {
			return ""
		}
		pt = pt.Children[r]
	}
	return pt.Item
}
