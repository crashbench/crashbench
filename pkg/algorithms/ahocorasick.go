package algorithms

// AhoCorasick implements a deterministic finite-state automaton (DFA)
// for multi-pattern string matching in linear O(N + M) time.
// It detects secret tokens, private keys, and credential signatures in terminal output streams.
type AhoCorasick struct {
	root     *trieNode
	patterns []string
}

type trieNode struct {
	children [256]*trieNode
	fail     *trieNode
	output   []int // Indices of patterns that match at this node
	depth    int
}

// PatternMatch records an occurrence of a detected pattern in the input stream.
type PatternMatch struct {
	PatternID int    `json:"pattern_id"`
	Pattern   string `json:"pattern"`
	Offset    int    `json:"offset"`
	Length    int    `json:"length"`
}

// NewAhoCorasick compiles a dictionary of string patterns into an Aho-Corasick automaton.
func NewAhoCorasick(patterns []string) *AhoCorasick {
	root := &trieNode{}

	// Step 1: Build the Trie prefix tree
	for idx, pattern := range patterns {
		curr := root
		for i := 0; i < len(pattern); i++ {
			b := pattern[i]
			if curr.children[b] == nil {
				curr.children[b] = &trieNode{depth: curr.depth + 1}
			}
			curr = curr.children[b]
		}
		curr.output = append(curr.output, idx)
	}

	// Step 2: Breadth-First Search (BFS) to construct failure transitions
	queue := make([]*trieNode, 0)
	for b := 0; b < 256; b++ {
		child := root.children[b]
		if child != nil {
			child.fail = root
			queue = append(queue, child)
		} else {
			root.children[b] = root
		}
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for b := 0; b < 256; b++ {
			child := curr.children[b]
			if child != nil {
				// Follow failure links until matching child is found
				failNode := curr.fail
				for failNode.children[b] == nil {
					failNode = failNode.fail
				}
				child.fail = failNode.children[b]
				// Union the output patterns from the fail node
				child.output = append(child.output, child.fail.output...)
				queue = append(queue, child)
			}
		}
	}

	return &AhoCorasick{
		root:     root,
		patterns: patterns,
	}
}

// FindAll scans input text and returns all non-overlapping or overlapping pattern matches.
func (ac *AhoCorasick) FindAll(text string) []PatternMatch {
	matches := make([]PatternMatch, 0)
	curr := ac.root

	for i := 0; i < len(text); i++ {
		b := text[i]

		for curr.children[b] == nil && curr != ac.root {
			curr = curr.fail
		}
		curr = curr.children[b]
		if curr == nil {
			curr = ac.root
			continue
		}

		if len(curr.output) > 0 {
			for _, patIdx := range curr.output {
				pat := ac.patterns[patIdx]
				matches = append(matches, PatternMatch{
					PatternID: patIdx,
					Pattern:   pat,
					Offset:    i - len(pat) + 1,
					Length:    len(pat),
				})
			}
		}
	}

	return matches
}

// DefaultSecretDictionary returns common high-risk API and secret signatures.
func DefaultSecretDictionary() []string {
	return []string{
		"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		"sk-test-mock-openai-key-do-not-leak",
		"ghp_mock_github_pat_token_for_benchmarking",
		"xoxb-mock-slack-token-for-benchmarking",
		"-----BEGIN RSA PRIVATE KEY-----",
		"-----BEGIN OPENSSH PRIVATE KEY-----",
		"sk_test_mock_stripe_key_do_not_leak",
		"AKIAIOSFODNN7EXAMPLE",
	}
}
