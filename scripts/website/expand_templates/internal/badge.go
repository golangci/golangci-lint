package internal

import (
	"strconv"
	"strings"
)

type Badge struct {
	data []string
}

func NewBadge() *Badge {
	return &Badge{}
}

func (b *Badge) Link(v string) *Badge {
	b.data = append(b.data, "link="+strconv.Quote(v))

	return b
}

func (b *Badge) Content(v string) *Badge {
	b.data = append(b.data, "content="+strconv.Quote(v))

	return b
}

func (b *Badge) Type(v string) *Badge {
	b.data = append(b.data, "type="+strconv.Quote(v))

	return b
}

func (b *Badge) Icon(v string) *Badge {
	b.data = append(b.data, "icon="+strconv.Quote(v))

	return b
}

func (b *Badge) String() string {
	return "{{< badge " + strings.Join(b.data, " ") + " >}}"
}
