package xml2map

import "fmt"

type Node struct {
	Parent   *Node
	Children map[string]Nodes
	Data     string
	Label    string
}

type Nodes []*Node

func (n *Node) AddChild(s string, c *Node) {
	if n.Children == nil {
		n.Children = map[string]Nodes{}
	}
	n.Children[s] = append(n.Children[s], c)
}

func (n *Node) ToMap() (map[string]interface{}, error) {
	result := map[string]interface{}{}

	for _, child := range n.Children {
		for _, node := range child {
			if len(node.Data) > 0 {
				if len(node.Children) > 0 {
					return nil, fmt.Errorf("tag '%s' contains invalid value: '%s'", node.Label, node.Data)
				}

				if len(node.Parent.Children[node.Label]) > 1 {
					if val, ok := result[node.Label]; ok {
						val = append(val.([]string), node.Data)
						result[node.Label] = val
					} else {
						result[node.Label] = []string{node.Data}
					}
				} else {
					result[node.Label] = node.Data
				}
			} else {
				m, err := node.ToMap()
				if err != nil {
					return nil, err
				}
				if len(node.Parent.Children[node.Label]) > 1 {
					if val, ok := result[node.Label]; ok {
						val = append(val.([]map[string]interface{}), m)
						result[node.Label] = val
					} else {

						result[node.Label] = []map[string]interface{}{m}
					}
				} else {
					result[node.Label] = m
				}
			}
		}
	}

	return result, nil
}


