package internal

import (
	"strconv"
	"strings"
)

type Cards struct {
	cards []*Card
	cols  int
}

func NewCards() *Cards {
	return &Cards{}
}

func (c *Cards) Cols(cols int) *Cards {
	c.cols = cols

	return c
}

func (c *Cards) Add(card *Card) *Cards {
	c.cards = append(c.cards, card)

	return c
}

func (c *Cards) String() string {
	v := "{{< cards"
	if c.cols > 0 {
		v += "  cols=" + strconv.Itoa(c.cols)
	}
	v += " >}}\n"

	for _, card := range c.cards {
		v += card.String()
	}

	v += "{{< /cards >}}\n"

	return v
}

type Card struct {
	data []string
}

func NewCard() *Card {
	return &Card{}
}

func (c *Card) Link(v string) *Card {
	c.data = append(c.data, "link="+strconv.Quote(v))

	return c
}

func (c *Card) Title(v string) *Card {
	c.data = append(c.data, "title="+strconv.Quote(v))

	return c
}

func (c *Card) Subtitle(v string) *Card {
	c.data = append(c.data, "subtitle="+strconv.Quote(v))

	return c
}

func (c *Card) Icon(v string) *Card {
	c.data = append(c.data, "icon="+strconv.Quote(v))

	return c
}

func (c *Card) Image(v string) *Card {
	c.data = append(c.data, "image="+strconv.Quote(v))

	return c
}

func (c *Card) ImageStyle(v string) *Card {
	c.data = append(c.data, "imageStyle="+strconv.Quote(v))

	return c
}

func (c *Card) Width(v string) *Card {
	c.data = append(c.data, "width="+strconv.Quote(v))

	return c
}

func (c *Card) Height(v string) *Card {
	c.data = append(c.data, "height="+strconv.Quote(v))
	return c
}

func (c *Card) Tag(v, tType string) *Card {
	c.data = append(c.data, "tag="+strconv.Quote(v))

	if tType != "" {
		c.data = append(c.data, "tagType="+strconv.Quote(tType))
	}

	return c
}

func (c *Card) String() string {
	return "  {{< card " + strings.Join(c.data, " ") + " >}}\n"
}
