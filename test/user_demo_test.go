package test

import (
	"fmt"
	"testing"
)

func TestValues(t *testing.T) {
	//en := json.NewEncoder(os.Stdout)
	//
	//s := make(map[string]string)
	//s["sds"] = "dsds"
	//en.Encode(s)

	s := "/hello/:demo"
	b := "/hello/:dsd/dsds"
	c := "/xxx/:vvv/ggg"
	n := node{}
	n.insert(s, []string{"hello", ":demo"}, 0)
	n.insert(b, []string{"hello", ":dsd", "dsds"}, 0)
	n.insert(c, []string{"xxx", ":vvv", "ggg"}, 0)
	for _, v := range n.children {
		fmt.Println(v)
	}
}


func demo(args ...interface{}) {
	fmt.Println(args...)
}

type node struct {
	patter   string
	part     string
	children []*node
	isWild   bool
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(patter string, parts []string, height int) {
	if len(parts) == height {
		n.patter = patter
		return
	}
	part := parts[height]
	child := n.matchChild(part)

	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}

		n.children = append(n.children, child)
	}
	child.insert(patter, parts, height+1)
}
