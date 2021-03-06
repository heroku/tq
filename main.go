package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pelletier/go-toml"
	"github.com/pelletier/go-toml/query"
)

const (
	version = "unknown"
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
		os.Exit(2)
	}
	return args[0]
}

func printKey(v interface{}) {
	switch t := v.(type) {
	case string:
		fmt.Printf("%s\n", t)
	case int, int32, int64, uint, uint8, uint16, uint32, uint64:
		fmt.Printf("%d\n", t)
	case float32, float64:
		fmt.Printf("%f\n", t)
	case bool:
		fmt.Printf("%t\n", t)
	case []interface{}:
		for i := range t {
			printKey(t[i])
		}
	default:
		fmt.Printf("%s\n", t)
	}
}

func printTree(t *toml.Tree) {
	for _, k := range t.Keys() {
		printKey(t.Get(k))
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
	for _, v := range result.Values() {
		switch t := v.(type) {
		case *toml.Tree:
			printTree(t)
		default:
			printKey(v)
		}
	}
}
