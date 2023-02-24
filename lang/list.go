package lang

import (
	"go/doc/comment"
	"strconv"
)

// List defines a list block element in the documentation for a symbol or
// package.
type List struct {
	blankBetween bool
	items        []*Item
}

// NewList initializes a list from the equivalent type from the comment package.
func NewList(cfg *Config, docList *comment.List) *List {
	var l List
	l.items = make([]*Item, len(docList.Items))
	for i, item := range docList.Items {
		l.items[i] = NewItem(cfg.Inc(0), item)
	}

	l.blankBetween = docList.BlankBetween()

	return &l
}

// BlankBetween returns true if there should be a blank line between list items.
func (l *List) BlankBetween() bool {
	return l.blankBetween
}

// Items returns the slice of items in the list.
func (l *List) Items() []*Item {
	return l.items
}

// ItemKind identifies the kind of item
type ItemKind string

const (
	// OrderedItem identifies an ordered (i.e. numbered) item.
	OrderedItem ItemKind = "ordered"

	// UnorderedItem identifies an unordered (i.e. bulletted) item.
	UnorderedItem ItemKind = "unordered"
)

// Item defines a single item in a list in the documentation for a symbol or
// package.
type Item struct {
	blocks []*Block
	kind   ItemKind
	number int
}

// NewItem initializes a list item from the equivalent type from the comment
// package.
func NewItem(cfg *Config, docItem *comment.ListItem) *Item {
	var (
		num  int
		kind ItemKind
	)
	if n, err := strconv.Atoi(docItem.Number); err == nil {
		num = n
		kind = OrderedItem
	} else {
		kind = UnorderedItem
	}

	return &Item{
		blocks: ParseBlocks(cfg, docItem.Content, true),
		kind:   kind,
		number: num,
	}
}

// Blocks returns the blocks of documentation in a list item.
func (i *Item) Blocks() []*Block {
	return i.blocks
}

// Kind returns the kind of the list item.
func (i *Item) Kind() ItemKind {
	return i.kind
}

// Number returns the number of the list item. Only populated if the item is of
// the OrderedItem kind.
func (i *Item) Number() int {
	return i.number
}
