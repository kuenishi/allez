package main
import (
	"encoding/json"
	"fmt"
	"flag"
	"os/exec"
	"os"
	"log"
	"io"
	"io/ioutil"
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
	nlist  := flag.String("nodes", "nodes", "filename list of machines")
	gname  := flag.String("goal", "", "")
	file   := flag.String("file", "", "set your target file")
	target := flag.String("target", "", "set your target directory")
	// dry-run
	flag.Parse()

	nodes := get_node_list(*nlist)
	fmt.Printf("%v on %v\n", *subcmd, nodes)

	goal,_ := NewGoal(*gname)
	fmt.Printf("%v", goal)

	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	switch *subcmd {
	case "ls":
		for i := range nodes {
			remote := "root@" + nodes[i]
			out,_ := exec.Command("ssh", remote, "ls").Output()
			fmt.Printf(string(out))
		}

	case "cp": // allez -cmd cp [filenames]
		if *file == "" {
			fmt.Printf("no input file\n")
			return
		}
		for i := range nodes {
			fmt.Printf("running on %s, %s => %s\n", nodes[i], *file, *target)
			remote_target := "root@" + nodes[i] + ":~/" + *target
			out,err := exec.Command("scp", *file, remote_target).Output()
			if err != nil {
				fmt.Printf("failed to copy to %s\n", remote_target)
			}
			fmt.Printf("ok: %v %s\n", out, remote_target)
		}

	case "build": // allez -cmd build

	case "start":
	case "stop":
	case "help": usage()
	default:     usage()
	}
}

func get_node_list(nodes string) []string {
	list,_ := readLines(nodes)
	return list
}

// http://stackoverflow.com/questions/5884154/golang-read-text-file-into-string-array-and-write
// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	content,err := ioutil.ReadFile(path)
	if err != nil {
		return nil,err
	}
	lines := strings.Split(string(content),"\n")
	return lines,err
}

type Goal struct {
	Name,File,Num string
	Version float32
	Cp,Build,Start,Stop []string
}


// "riak" => loads "riak.goal" and parses it
func NewGoal(name string) (*Goal, error) {
	filename := name + ".goal"
	fmt.Printf("loading %s\n", filename)
	ifp,err := os.Open(filename)
	defer ifp.Close()

	if err != nil {
		return nil,err
	}
	dec := json.NewDecoder(ifp)
	var g Goal
	for {
		if err = dec.Decode(&g); err == io.EOF {
			return nil,err
		} else if err != nil {
			return nil,err
		} else {
			break
		}
	}
	return &g,err
}
func (g *Goal) DoCopy() error { return nil }
func (g *Goal) DoBuild() error { return nil }
func (g *Goal) DoStart() error { return nil }
func (g *Goal) DoStop() error { return nil }

