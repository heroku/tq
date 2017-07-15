package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pelletier/go-toml"
	"github.com/pelletier/go-toml/query"
)

const (
	version = "0.1"
)

var (
	keyFilter = flag.String("key", "", "Single key to filter on")
)

func handleArgs() string {
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Queries toml read from stdin.")
		fmt.Println("  ex: < path/to/file.toml tq '$.foo.bar'")
		fmt.Println()
		fmt.Println("Usage: tq [options] <query>")
		fmt.Println("version", version)
		fmt.Println()
		fmt.Println("<query> : see https://godoc.org/github.com/pelletier/go-toml/query for docs")
		fmt.Println()
		fmt.Println("options:")
		flag.PrintDefaults()
		fmt.Println()
		os.Exit(2)
	}
	return args[0]
}

func printKey(k, filter string, v interface{}) {
	if filter == "" {
		fmt.Printf("%s = %q\n", k, v)
		return
	}
	if filter == k {
		fmt.Printf("%s\n", v)
	}
}

func printTree(t *toml.Tree, filter string) {
	for _, k := range t.Keys() {
		printKey(k, filter, t.Get(k))
	}
}

func main() {
	q := handleArgs()

	tree, err := toml.LoadReader(os.Stdin)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	result, err := query.CompileAndExecute(q, tree)
	if err != nil {
		fmt.Printf("Error: %T - %s", err, err)
		os.Exit(1)
	}
	for i, v := range result.Values() {
		switch t := v.(type) {
		case *toml.Tree:
			printTree(t, *keyFilter)
		default:
			fmt.Printf("default: %d, %T, %s", i, t, t)
		}
	}
}
