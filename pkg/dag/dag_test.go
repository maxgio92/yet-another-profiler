package dag_test

import (
	"github.com/maxgio92/yap/pkg/dag"
	"testing"
)

type StackTrace struct {
	Syms    []string
	Samples int
}

var (
	trace1 = StackTrace{[]string{"main", "foo", "qux", "grault"}, 4}
	trace5 = StackTrace{[]string{"main", "foo", "quux"}, 2}
	trace2 = StackTrace{[]string{"main", "bar", "quux"}, 2}
	trace3 = StackTrace{[]string{"main", "foo"}, 3}
	trace4 = StackTrace{[]string{"main", "foo", "corge", "garply"}, 5}
	traces = []StackTrace{trace1, trace2, trace3, trace4, trace5}
)

func TestUpsertChildren(t *testing.T) {
	g := dag.NewDAG()
	fillDAG(g, traces)

	testCases := []struct {
		name string
		want int
	}{
		{name: "main", want: 2},
		{name: "foo", want: 3},
		{name: "bar", want: 1},
		{name: "qux", want: 1},
		{name: "quux", want: 0},
	}

	for _, tt := range testCases {
		children := g.Node(tt.name).Children
		if len(children) != tt.want {
			t.Errorf("for %v got %v, want %v", tt.name, len(children), tt.want)
		}
	}
}

func TestUpsertParents(t *testing.T) {
	g := dag.NewDAG()
	fillDAG(g, traces)

	testCases := []struct {
		name string
		want int
	}{
		{name: "main", want: 1},
		{name: "foo", want: 1},
		{name: "bar", want: 1},
		{name: "qux", want: 1},
		{name: "quux", want: 2},
	}

	for _, tt := range testCases {
		parents := g.Node(tt.name).Parents
		if len(parents) != tt.want {
			t.Errorf("for %v got %v, want %v", tt.name, len(parents), tt.want)
		}
	}
}

func fillDAG(g *dag.DAG, traces []StackTrace) {
	var sampleCountTotal int
	for _, v := range traces {
		sampleCountTotal += v.Samples
	}

	for kt, _ := range traces {
		for ks, sym := range traces[kt].Syms {
			var parent string
			if ks > 0 {
				parent = traces[kt].Syms[ks-1]
			}

			// If it's the traced function, that is, the last symbol/IP in the stack trace,
			// update also its weight.
			if ks == len(traces[kt].Syms)-1 {
				g.UpsertNode(sym, parent, float32(traces[kt].Samples)/float32(sampleCountTotal))
			} else {
				g.UpsertNode(sym, parent)
			}
		}
	}
}
