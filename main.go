package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Difrex/gosway/ipc"
)

func findFocused(nodes []Node) (Node, bool) {
	for _, node := range nodes {
		if node.Focused {
			return node, true
		}
		n, ok := findFocused(node.Nodes)
		if ok {
			return n, ok
		}
	}
	return Node{}, false
}

type Node struct {
	Nodes         []Node `json:"nodes,omitempty"`
	FloatingNodes []Node `json:"floating_nodes,omitempty"`
	Pid           int    `json:"pid,omitempty"`
	Focused       bool   `json:"focused"`
	Focus         []int  `json:"focus"`
}

type Tree struct {
	Name  string `json:"name"`
	Nodes []Node `json:"nodes"`
}

func getTree() (*Tree, error) {
	sc, err := ipc.NewSwayConnection()
	if err != nil {
		return nil, err
	}

	tree := &Tree{}

	data, err := sc.SendCommand(ipc.IPC_GET_TREE, "get_tree")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, tree)
	return tree, err
}

func findFocusedPID() (int, error) {
	tree, err := getTree()
	if err != nil {
		log.Fatal(err)
	}
	node, ok := findFocused(tree.Nodes)
	if ok {
		return node.Pid, nil
	}
	return 0, errors.New("pid not found")
}

func main() {
	pid, err := findFocusedPID()
	if err != nil {
		fmt.Println(os.Getenv("HOME"))
	}
	path, err := os.Readlink(fmt.Sprintf("/proc/%d/cwd", pid))
	if err != nil {
		fmt.Println(os.Getenv("HOME"))
	}
	if path == "/" {
		fmt.Println(os.Getenv("HOME"))
	}
	fmt.Println(path)
}
