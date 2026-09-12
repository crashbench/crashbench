package algorithms

import (
	"errors"
	"fmt"
)

// ProcessNode represents a single OS process in the process execution tree.
type ProcessNode struct {
	PID      int            `json:"pid"`
	PPID     int            `json:"ppid"`
	Command  string         `json:"command"`
	Children []*ProcessNode `json:"children"`
	IsZombie bool           `json:"is_zombie"`
}

// ProcessDAG tracks hierarchy, detect cycles, and plans topological process reaping.
type ProcessDAG struct {
	nodes map[int]*ProcessNode
	roots []*ProcessNode
}

// NewProcessDAG constructs an empty process execution graph.
func NewProcessDAG() *ProcessDAG {
	return &ProcessDAG{
		nodes: make(map[int]*ProcessNode),
		roots: make([]*ProcessNode, 0),
	}
}

// AddProcess adds or updates a process node in the graph.
func (dag *ProcessDAG) AddProcess(pid, ppid int, command string, isZombie bool) {
	node, exists := dag.nodes[pid]
	if !exists {
		node = &ProcessNode{
			PID:      pid,
			PPID:     ppid,
			Command:  command,
			Children: make([]*ProcessNode, 0),
			IsZombie: isZombie,
		}
		dag.nodes[pid] = node
	} else {
		node.PPID = ppid
		node.Command = command
		node.IsZombie = isZombie
	}
}

// BuildHierarchy resolves all parent-child pointers and identifies root nodes.
func (dag *ProcessDAG) BuildHierarchy() {
	dag.roots = make([]*ProcessNode, 0)
	for _, node := range dag.nodes {
		node.Children = make([]*ProcessNode, 0)
	}

	for _, node := range dag.nodes {
		parent, hasParent := dag.nodes[node.PPID]
		if hasParent && parent.PID != node.PID {
			parent.Children = append(parent.Children, node)
		} else {
			dag.roots = append(dag.roots, node)
		}
	}
}

// Color constants for Tarjan cycle detection
const (
	colorWhite = 0 // Unvisited
	colorGray  = 1 // In progress (in recursion stack)
	colorBlack = 2 // Fully visited
)

// DetectCycles verifies that the process graph is a strict DAG without cyclic fork bombs.
func (dag *ProcessDAG) DetectCycles() (bool, error) {
	colors := make(map[int]int)
	for pid := range dag.nodes {
		colors[pid] = colorWhite
	}

	var dfs func(node *ProcessNode) bool
	dfs = func(node *ProcessNode) bool {
		colors[node.PID] = colorGray
		for _, child := range node.Children {
			if colors[child.PID] == colorGray {
				return true // Cycle detected!
			}
			if colors[child.PID] == colorWhite {
				if dfs(child) {
					return true
				}
			}
		}
		colors[node.PID] = colorBlack
		return false
	}

	for _, root := range dag.roots {
		if colors[root.PID] == colorWhite {
			if dfs(root) {
				return true, errors.New("process DAG cycle detected: fork-bomb condition")
			}
		}
	}

	return false, nil
}

// TopologicalReapOrder returns PIDs in post-order (leaves first, roots last)
// to guarantee child processes are terminated before parent pipes break.
func (dag *ProcessDAG) TopologicalReapOrder() []int {
	visited := make(map[int]bool)
	reapOrder := make([]int, 0, len(dag.nodes))

	var postOrder func(node *ProcessNode)
	postOrder = func(node *ProcessNode) {
		visited[node.PID] = true
		for _, child := range node.Children {
			if !visited[child.PID] {
				postOrder(child)
			}
		}
		reapOrder = append(reapOrder, node.PID)
	}

	for _, root := range dag.roots {
		if !visited[root.PID] {
			postOrder(root)
		}
	}

	return reapOrder
}

// CountOrphanZombies returns the count of detached child processes that survived.
func (dag *ProcessDAG) CountOrphanZombies() int {
	zombies := 0
	for _, node := range dag.nodes {
		if node.IsZombie {
			zombies++
		}
	}
	return zombies
}

// String returns a visual ASCII tree of the process hierarchy.
func (dag *ProcessDAG) String() string {
	res := ""
	var printNode func(node *ProcessNode, indent string)
	printNode = func(node *ProcessNode, indent string) {
		res += fmt.Sprintf("%s└─ PID %d (%s)%s\n", indent, node.PID, node.Command, map[bool]string{true: " [ZOMBIE]", false: ""}[node.IsZombie])
		for _, child := range node.Children {
			printNode(child, indent+"  ")
		}
	}

	for _, root := range dag.roots {
		printNode(root, "")
	}
	return res
}
