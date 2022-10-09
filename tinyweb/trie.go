package tinyweb

import (
	"log"
	"strings"
)

type node struct {
	subpath    string // 当前
	children   map[string]*node
	isEndpoint bool // true表示当前节点是路径最后一个端点
}

type trie struct {
	root *node
}

func newTrie() *trie {
	return &trie{
		root: &node{
			children:   make(map[string]*node),
			isEndpoint: false,
		},
	}
}

func dfs(root *node, pattern []string, i int, path *[]string) bool {
	if root == nil {
		return false
	}

	if i >= len(pattern) {
		return root.isEndpoint
	}

	n, ok := root.children[pattern[i]]
	if ok {
		*path = append(*path, n.subpath)
		return dfs(n, pattern, i+1, path)
	} else {
		n, ok := root.children[":"]
		if ok {
			*path = append(*path, n.subpath)
			return dfs(n, pattern, i+1, path)
		} else {
			return false
		}
	}
}

func (t *trie) find(url string) (string, map[string]string, bool) {
	path, ok := splitURL(url)
	if !ok {
		return "", nil, false
	}

	arr := make([]string, 0)
	ok = dfs(t.root, path, 0, &arr)
	if ok {
		params := make(map[string]string)
		for i, x := range path {
			if arr[i][0] == ':' {
				params[arr[i][1:len(arr[i])-1]] = x[0 : len(x)-1]
			}
		}
		handlerKey := strings.Join(arr, "")

		return handlerKey, params, ok
	} else {
		return "", nil, ok
	}
}

func (t *trie) insert(url string) (string, bool) {
	pattern, ok := splitURL(url)
	if !ok {
		log.Fatalf("illegal url: %s\n", url)
		return "", false
	}

	currentNode := t.root
	for _, p := range pattern {
		if p[0] == ':' {
			_, ok = currentNode.children[":"]
			if !ok {
				currentNode.children[":"] = &node{children: make(map[string]*node), subpath: p, isEndpoint: false}
			}
			currentNode = currentNode.children[":"]
		} else {
			_, ok = currentNode.children[p]
			if !ok {
				currentNode.children[p] = &node{subpath: p, children: make(map[string]*node), isEndpoint: false}
			}
			currentNode = currentNode.children[p]
		}
	}
	currentNode.isEndpoint = true
	return strings.Join(pattern, ""), true
}

// fixURL 给缺少`/`的URL添加`/`
func fixURL(url string) string {
	return url + "/"
}

// splitURL
//	(1) /a/b/c/d  fix to   /a/b/c/d/  split to  /, a/, b/, c/, d/
//	(2) / split to /
func splitURL(url string) ([]string, bool) {
	if url == "" || len(url) == 0 || url[0] != '/' {
		return nil, false
	} else {
		if url[len(url)-1] != '/' {
			url = fixURL(url)
		}

		pattern := strings.Split(url, "/")
		pattern = pattern[0 : len(pattern)-1]
		for i := range pattern {
			pattern[i] = pattern[i] + "/"
		}

		return pattern, true
	}
}
