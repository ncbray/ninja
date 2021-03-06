package ninja

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleRule(t *testing.T) {
	var b bytes.Buffer
	w := MakeNinjaWriter(&b)
	w.Variable("baz", "1")
	w.Rule("foo", "bar $baz", RuleOptions{})
	assert.Equal(t, "baz = 1\nrule foo\n  command = bar $baz\n", b.String())
}

func TestSimpleBuild(t *testing.T) {
	var b bytes.Buffer
	w := MakeNinjaWriter(&b)
	w.Build([]string{"out1", "out2"}, "a_rule", BuildOptions{})
	assert.Equal(t, "build out1 out2: a_rule\n", b.String())
}

func TestSimpleEscape(t *testing.T) {
	var b bytes.Buffer
	w := MakeNinjaWriter(&b)
	w.Build([]string{"o$u t:"}, "a_rule", BuildOptions{})
	assert.Equal(t, "build o$$u$ t$:: a_rule\n", b.String())
}

func TestIOBuild(t *testing.T) {
	var b bytes.Buffer
	w := MakeNinjaWriter(&b)
	w.Build([]string{"out1", "out2"}, "a_rule", BuildOptions{
		Inputs:          []string{"in1", "in2"},
		ImplicitInputs:  []string{"inImp1", "inImp2"},
		OrderOnlyInputs: []string{"order1", "order2"},
		ImplicitOutputs: []string{"outImp1", "outImp2"},
		Variables:       []Variable{{"foo", "bar"}},
	})
	assert.Equal(t, "build out1 out2 | outImp1 outImp2: a_rule in1 in2 | inImp1 inImp2 || order1 order2\n  foo = bar\n", b.String())
}
