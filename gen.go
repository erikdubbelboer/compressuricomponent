// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"sort"
	"strings"
)

type dictEntry struct {
	str   string
	score int
}

type dict []dictEntry

func main() {
	all, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(all), "\n")
	finaldict := make(dict, 0)
	originalLengths := make([]int, len(lines))

	for i, _ := range lines {
		originalLengths[i] = len(lines[i])
		lines[i] = url.QueryEscape(lines[i])
		lines[i] = strings.Replace(lines[i], "~", "%7E", -1)
	}

	dict := make(dict, 0)

	for len(finaldict) < 200 {
		already := make(map[string]struct{})
		dict = dict[:0]

		for _, line := range lines {
			fmt.Fprintf(os.Stderr, ".")
			for start := 0; start < len(line)-5; start++ {
				for end := start + 6; end < len(line); end++ {
					if end-start > 28 {
						continue
					}

					s := line[start:end]

					if strings.Contains(s, "~") {
						break
					}

					if _, ok := already[s]; ok {
						continue
					}
					already[s] = struct{}{}

					occurances := 0
					for _, ll := range lines {
						occurances += strings.Count(ll, s)
					}

					if occurances < 10 {
						continue
					}

					score := occurances * (len(s) - 6)

					if score <= 10 {
						continue
					}

					dict = append(dict, dictEntry{
						str:   s,
						score: score,
					})
				}
			}
		}

		if len(dict) == 0 {
			break
		}

		sort.Slice(dict, func(i, j int) bool {
			return dict[i].score > dict[j].score
		})

		s := dict[0].str
		score := dict[0].score
		r := fmt.Sprintf("~%d~", len(finaldict))
		finaldict = append(finaldict, dict[0])

		for i, _ := range lines {
			lines[i] = strings.Replace(lines[i], s, r, -1)
		}

		fmt.Fprintf(os.Stderr, "\n%6d %q\n", score, s)
	}

	fmt.Printf("package %s\n\nvar Dict = []string{\n", os.Args[1])
	for _, d := range finaldict {
		fmt.Printf("\t%q,\n", d.str)
	}
	fmt.Printf("}\n")
}
