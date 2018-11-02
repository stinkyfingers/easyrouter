package easyrouter

// package prefixtree implements a simple prefix tree for string searches in urls with '*' wildcards

type Node struct {
	children  map[rune]*Node
	leaf      bool
	wildcard  bool
	leafValue int
}

// '*' char is wild--matches any chars until the next '/'

func NewTree() *Node {
	return &Node{
		children: make(map[rune]*Node),
	}
}

// Insert adds a value to a Tree
func (n *Node) Insert(value string, leafValue int) {
	current := n
	runes := stringToRunes(value)
	var inWildMode bool
	for i := 0; i < len(runes); i++ {
		if inWildMode {
			if runes[i] == '/' {
				inWildMode = false
			} else {
				continue
			}
		}
		next := current.findNode(runes[i])

		// node found; set current to next; continue
		if next != nil {
			current.children[runes[i]] = next
			current = next
			continue
		}

		// create node
		next = &Node{
			children: make(map[rune]*Node),
		}
		if i == len(runes)-1 {
			next.leaf = true
			next.leafValue = leafValue
		}
		current.children[runes[i]] = next
		current = next

		if runes[i] == '*' {
			inWildMode = true
		}
	}
}

// Find returns true and the tree value if a Tree contains value
func (n *Node) Find(value string) (bool, int) {
	current := n
	runes := stringToRunes(value)
	var inWildMode bool
	for i := 0; i < len(runes); i++ {
		if inWildMode {
			if runes[i] == '/' {
				inWildMode = false
			} else {
				continue
			}
		}

		var newNode *Node
		newNode = current.findNode(runes[i])
		if newNode == nil {
			// check wildcard
			newNode = current.findNode('*')
			if newNode == nil {
				return false, -1
			} else {
				inWildMode = true
			}
		}
		if i == len(runes)-1 && !newNode.leaf {
			return false, -1
		}
		current = newNode
	}
	return true, current.leafValue
}

// findNode returns a pointer to a node if r is a child of n
func (n *Node) findNode(r rune) *Node {
	if node, ok := n.children[r]; ok {
		return node
	}
	return nil
}

// stringToRunes returns a []rune for a string
func stringToRunes(value string) []rune {
	var runes []rune
	for _, ch := range value {
		runes = append(runes, ch)
	}
	return runes
}
