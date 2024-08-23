package main
/* tree prints the files, directories in the shape of a tree. This requires
  implementing a data structure. However, it is NOT guaranteed that the 
  output tree follows alphabetical order because directories are accessed in
  parallel. 
  Here's an example output:
.
|-wget.go
|-.gitignore
|-cat.go
|-cp.go
|-hexdump.go
|-ls.go
|-Makefile
|-nslookup.go
|-parrot.go
|-pwd.go
|-README.md
|-request.go
|-server.go
|-touch.go
|-tree.go
|-urls.txt
|-treetest
| |-a.txt
| |-foo
| | |-foo0
| | 
| |-bar
| | |-bar0
| | 
|
	*/

import (
	"os"
	"fmt"
	"strings"
)

type TreeNode struct {
	name string;           // name of current file or directory.
	children []*TreeNode;  // pointer to child nodes.
};

type Tree struct {
	root *TreeNode;
}

func extract_file(url string) string {
	idx := strings.LastIndex(url, "/");
	if idx == -1 {
		return url;
	}
	return url[idx + 1 :];
}

func walker(start string, c chan *TreeNode) {
	/** Since we want to create multiple threads to walk different 
		directories simultaneously, we'll need a mechanism to "join" 
		all walking threads. */
	sp := make([]TreeNode, 1);
	var root *TreeNode = &sp[0];
	root.name = extract_file(start);

	// open directory(or file)
	dobj, err := os.Open(start);
	if err != nil {
		// oops, fail
		c <- root;
		return;
	}

	// check dobj is file or directory.
	info, err := dobj.Stat();
	if err != nil || (!info.IsDir()) {
		// terminate.
		c <- root;
		return;
	}

	// read directory, probably create other walkers.
	entries, err := dobj.Readdir(0);
	if err != nil {
		// stop.
		c <- root;
		return;
	}

	if len(entries) == 0 {
		c <- root;
		return;
	}

	// initialize root node.
	root.children = make([]*TreeNode, len(entries));

	// create several walkers
	channels := make(chan *TreeNode, len(entries));
	for i := 0; i < len(entries); i++ {
		// walk deeper into the file system
		go walker(start + "/" + entries[i].Name(), channels);
	}

	// receive results from subprocess.
	for k := 0; k < len(entries); k++ {
		root.children[k] = <- channels;
	}

	// finish.
	c <- root;
}

func ptreeRecur(root *TreeNode, depth int) {
	if root == nil {
		return;
	}
	for i := 0; i < depth - 1; i++ {
		fmt.Print("| ");
	}
	if depth > 0 {
		fmt.Print("|-");
	}
	fmt.Println(root.name);

	for _, node := range root.children {
		ptreeRecur(node, depth + 1);
	}

	if len(root.children) > 0 {
		for i := 0; i < depth; i++ {
			fmt.Print("| ");
		}
		fmt.Print("\n");
	}
}

func ptree(root *TreeNode) {
	if root == nil {
		return;
	}

	ptreeRecur(root, 0);
}

func main() {
	start := ".";	
	if len(os.Args) == 2 {
		start = os.Args[1];
	}

	rc := make(chan *TreeNode);
	go walker(start, rc);
	ptree(<- rc);
	return;
}
