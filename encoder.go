package xml2map

import (
	"encoding/xml"
	"io"
	"strings"
	"github.com/pkg/errors"
)

var (
	InvalidDocument = errors.New("Invalid document")
)

type Encoder struct {
	r io.Reader
}

func NewEncoder(reader io.Reader) *Encoder {
	return &Encoder{r: reader}
}

func (e *Encoder) Encode() (map[string]interface{}, error) {
	decoder := xml.NewDecoder(e.r)
	tree := &Node{}
	n := 0

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		switch el := t.(type) {
		case xml.StartElement:
			tree = &Node{
				Parent: tree,
				Label:  el.Name.Local,
			}

			if len(el.Attr) > 0 {
				for _, v := range el.Attr {
					tree.AddChild(v.Name.Local, &Node{Label: v.Name.Local, Data: v.Value, Parent: tree})
				}
			}
			n++

		case xml.CharData:
			tree.Data = strings.TrimSpace(string(el))

		case xml.EndElement:
			tree.Parent.AddChild(el.Name.Local, tree)
			if tree.Parent != nil {
				tree = tree.Parent
			}
			n--
		}
	}

	if n > 0 {
		return nil, InvalidDocument
	}

	return tree.ToMap()
}