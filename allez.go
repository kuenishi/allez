package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: allez -cmd cp|build|start|stop -nodes nodes [ARGV]\n")
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {

	flag.Usage = usage
	subcmd := flag.String("cmd", "help", "allez command")
	nlist := flag.String("nodes", "nodes", "filename list of machines")
	gname := flag.String("goal", "", "")
	// dry-run
	flag.Parse()

	nodes := get_node_list(*nlist)
	fmt.Printf("%v on %v\n", *subcmd, nodes)

	goal, err := NewGoal(*gname)
	if err != nil {
		fmt.Printf("%s: %v\n", *gname, err)
		return
	} else {
		fmt.Printf("%v\n", goal)
	}

	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	switch *subcmd {
	case "ls":
		for i := range nodes {
			remote := "root@" + nodes[i]
			out, _ := exec.Command("ssh", remote, "ls").Output()
			fmt.Printf(string(out))
		}

	case "cp": // allez -cmd cp [filenames]
		goal.DoCopy(nodes)

	case "build": // allez -cmd build
		goal.DoBuild(nodes)
	case "start":
		goal.DoStart(nodes)
	case "stop":
		goal.DoStop(nodes)
	case "help":
		usage()
	default:
		usage()
	}
}

func get_node_list(nodes string) []string {
	list, _ := readLines(nodes)
	return list
}

// http://stackoverflow.com/questions/5884154/golang-read-text-file-into-string-array-and-write
// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	return lines, err
}

type NotFoundError struct {
}

func (e *NotFoundError) Error() string {
	return "not_found"
}

type Goal struct {
	Name, File, Num    string
	Version            float32
	Cp                 string
	Build, Start, Stop []string
}

// "riak" => loads "riak.goal" and parses it
func NewGoal(name string) (*Goal, error) {
	filename := name + ".goal"
	fmt.Printf("loading %s\n", filename)
	ifp, err := os.Open(filename)
	defer ifp.Close()

	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(ifp)
	var g Goal
	for {
		if err = dec.Decode(&g); err == io.EOF {
			return nil, err
		} else if err != nil {
			return nil, err
		} else {
			break
		}
	}
	return &g, err
}

func (g *Goal) DoCopy(nodes []string) error {
	file := g.File
	if g.File == "" {
		fmt.Printf("no input file\n")
		return &NotFoundError{}
	}
	for i := range nodes {
		fmt.Printf("working on %s, %s => %s\n", nodes[i], file, g.Cp)
		target := "root@" + nodes[i] + ":~/" + g.Cp
		out, err := exec.Command("scp", file, target).Output()
		if err != nil {
			fmt.Printf("failed to copy %s to %s\n", file, target)
			return err
		}
		fmt.Printf("ok: %v %s\n", out, target)
	}

	return nil
}
func (g *Goal) DoBuild(nodes []string) error { return nil }
func (g *Goal) DoStart(nodes []string) error { return nil }
func (g *Goal) DoStop(nodes []string) error  { return nil }
