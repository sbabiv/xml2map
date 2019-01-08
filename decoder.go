package xml2map

import (
	"encoding/xml"
	"errors"
	"io"
	"strings"
)

const (
	defaultAttrPrefix = "@"
	defaultTextPrefix = "#text"
)

var (
	//ErrInvalidDocument invalid document err
	ErrInvalidDocument = errors.New("invalid document")

	//ErrInvalidRoot data at the root level is invalid err
	ErrInvalidRoot = errors.New("data at the root level is invalid")
)

type node struct {
	Parent  *node
	Value   map[string]interface{}
	Attrs   []xml.Attr
	Label   string
	Text    string
	HasMany bool
}

// Decoder instance
type Decoder struct {
	r 			io.Reader
	attrPrefix	*string
	textPrefix	*string
}

// NewDecoder create new decoder instance
func NewDecoder(reader io.Reader) *Decoder {
	return &Decoder{r: reader}
}

func NewDecoderWithCustomPrefixes(reader io.Reader, attrPrefix, textPrefix *string) *Decoder {
	return &Decoder{r: reader, attrPrefix: attrPrefix, textPrefix: textPrefix}
}

//Decode xml string to map[string]interface{}
func (d *Decoder) Decode() (map[string]interface{}, error) {
	decoder := xml.NewDecoder(d.r)
	n := &node{}
	stack := make([]*node, 0)

	attrPrefix := defaultAttrPrefix
	textPrefix := defaultTextPrefix

	if d.attrPrefix != nil {
		attrPrefix = *d.attrPrefix
	}

	if d.textPrefix != nil {
		textPrefix = *d.textPrefix
	}

	for {
		token, err := decoder.Token()
		if err != nil && err != io.EOF {
			return nil, err
		}

		if token == nil {
			break
		}

		switch tok := token.(type) {
		case xml.StartElement:
			{
				n = &node{
					Label:  tok.Name.Local,
					Parent: n,
					Value:  map[string]interface{}{tok.Name.Local: map[string]interface{}{}},
					Attrs:  tok.Attr,
				}

				setAttrs(attrPrefix, n, &tok)
				stack = append(stack, n)

				if n.Parent != nil {
					n.Parent.HasMany = true
				}
			}

		case xml.CharData:
			data := strings.TrimSpace(string(tok))
			if len(stack) > 0 {
				stack[len(stack)-1].Text = data
			} else if len(data) > 0 {
				return nil, ErrInvalidRoot
			}

		case xml.EndElement:
			{
				length := len(stack)
				stack, n = stack[:length-1], stack[length-1]

				if !n.HasMany {
					if len(n.Attrs) > 0 {
						m := n.Value[n.Label].(map[string]interface{})
						m[textPrefix] = n.Text
					} else {
						n.Value[n.Label] = n.Text
					}
				}

				if len(stack) == 0 {
					return n.Value, nil
				}

				setNodeValue(n)
				n = n.Parent
			}
		}
	}

	return nil, ErrInvalidDocument
}

func setAttrs(attrPrefix string, n *node, tok *xml.StartElement) {
	if len(tok.Attr) > 0 {
		m := make(map[string]interface{})
		for _, attr := range tok.Attr {
			if len(attr.Name.Space) > 0 {
				m[attrPrefix+attr.Name.Space+":"+attr.Name.Local] = attr.Value
			} else {
				m[attrPrefix+attr.Name.Local] = attr.Value
			}
		}
		n.Value[tok.Name.Local] = m
	}
}

func setNodeValue(n *node) {
	if v, ok := n.Parent.Value[n.Parent.Label]; ok {
		m := v.(map[string]interface{})
		if v, ok = m[n.Label]; ok {
			switch item := v.(type) {
			case string:
				m[n.Label] = []string{item, n.Value[n.Label].(string)}
			case []string:
				m[n.Label] = append(item, n.Value[n.Label].(string))
			case map[string]interface{}:
				vm := getMap(n)
				if vm != nil {
					m[n.Label] = []map[string]interface{}{item, vm}
				}
			case []map[string]interface{}:
				vm := getMap(n)
				if vm != nil {
					m[n.Label] = append(item, vm)
				}
			}
		} else {
			m[n.Label] = n.Value[n.Label]
		}

	} else {
		n.Parent.Value[n.Parent.Label] = n.Value[n.Label]
	}
}

func getMap(node *node) map[string]interface{} {
	if v, ok := node.Value[node.Label]; ok {
		switch v.(type) {
		case string:
			return map[string]interface{}{node.Label: v}
		case map[string]interface{}:
			return node.Value[node.Label].(map[string]interface{})
		}
	}

	return nil
}
