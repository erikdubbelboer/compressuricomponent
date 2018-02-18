package useragents

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/erikdubbelboer/compressuricomponent"
)

func TestDict(t *testing.T) {
	content, err := ioutil.ReadFile("ua.txt")
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(string(content), "\n")
	lines = lines[:1]

	encoder := compressuricomponent.NewEncoder(Dict)
	decoder := compressuricomponent.NewDecoder(Dict)

	for _, line := range lines {
		t.Run(line, func(t *testing.T) {
			encoded := encoder.Encode(line)
			t.Log(encoded)
			if decoded, err := decoder.Decode(encoded); err != nil {
				t.Fatal(err)
			} else if decoded != line {
				t.Fatalf("%q", decoded)
			}
		})
	}
}

func BenchmarkDict(b *testing.B) {
	content, err := ioutil.ReadFile("ua.txt")
	if err != nil {
		b.Fatal(err)
	}

	lines := strings.Split(string(content), "\n")

	encoder := compressuricomponent.NewEncoder(Dict)
	decoder := compressuricomponent.NewDecoder(Dict)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		line := lines[i%len(lines)]
		encoded := encoder.Encode(line)
		if decoded, err := decoder.Decode(encoded); err != nil {
			b.Fatal(err)
		} else if decoded != line {
			b.Fatalf("%q", decoded)
		}
	}
}
