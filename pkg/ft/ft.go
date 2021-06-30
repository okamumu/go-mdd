package ft

import (
	"fmt"
	"github.com/okamumu/go-mdd/pkg/bdd"
)

type minPath struct {
	len int
	set [][]string
}

func Ftmcs(b *bdd.BDD, f *bdd.Node) [][]string {
	vars := make(map[string]*bdd.Node)
	v := b.GetVarLabels()
	x := b.GetVars()
	for i, n := range v {
		vars[n] = x[i]
	}

	result := [][]string{}
	r := f
	for r != b.Zero() {
		s, minf := ftmcs(b, vars, r)
		for _, x := range s {
			result = append(result, x)
		}
		r = b.Not(b.Imp(r, minf))
	}
	return result
}

func ftmcs(b *bdd.BDD, vars map[string]*bdd.Node, f *bdd.Node) ([][]string, *bdd.Node) {
	path := make([]string, 0, len(vars))
	s := &minPath{
		len: len(vars),
		set: make([][]string, 0),
	}
	bddres := b.Zero()
	getMinPath(b, f, path, s)
	for _, ms := range s.set {
		tmp := b.One()
		for _, x := range ms {
			tmp = b.And(tmp, vars[x])
		}
		bddres = b.Or(bddres, tmp)
	}
	return s.set, bddres
}

func copyPath(path []string) []string {
	p := make([]string, len(path))
	for i, x := range path {
		p[i] = x
	}
	return p
}

func getMinPath(b *bdd.BDD, f *bdd.Node, path []string, s *minPath) {
	if s.len < len(path) {
		return
	}
	switch {
	case f == b.One():
		switch {
		case s.len > len(path):
			s.len = len(path)
			s.set = [][]string{copyPath(path)}
		case s.len == len(path):
			s.set = append(s.set, copyPath(path))
		}
	case f != b.One() && f != b.Zero():
		getMinPath(b, f.Zero(), path, s)
		path = append(path, f.Label())
		getMinPath(b, f.One(), path, s)
	}
}

////

func Ftmcs0(b *bdd.BDD0, f *bdd.Node01) [][]string {
	vars := make(map[string]*bdd.Node01)
	v := b.GetVarLabels()
	x := b.GetVars()
	for i, n := range v {
		vars[n] = x[i]
	}

	result := [][]string{}
	r := f
	for !(r.Node0 == b.Zero && r.Neg == false) {
		s, minf := ftmcs0(b, vars, r)
		fmt.Println("minf:")
		fmt.Println(minf)
		b.ToDot(minf)
		for _, x := range s {
			result = append(result, x)
		}
		r = b.Not(b.Imp(r, minf))
	}
	return result
}

func ftmcs0(b *bdd.BDD0, vars map[string]*bdd.Node01, f *bdd.Node01) ([][]string, *bdd.Node01) {
	path := make([]string, 0, len(vars))
	s := &minPath{
		len: len(vars),
		set: make([][]string, 0),
	}
	bddres := &bdd.Node01{b.Zero, false}
	getMinPath0(b, f.Node0, f.Neg, path, s)
	for _, ms := range s.set {
		tmp := &bdd.Node01{b.Zero, true}
		for _, x := range ms {
			tmp = b.And(tmp, vars[x])
		}
		bddres = b.Or(bddres, tmp)
	}
	return s.set, bddres
}

func getMinPath0(b *bdd.BDD0, f *bdd.Node0, fneg bool, path []string, s *minPath) {
	if s.len < len(path) {
		return
	}
	switch {
	case f == b.Zero && fneg == true:
		switch {
		case s.len > len(path):
			s.len = len(path)
			s.set = [][]string{copyPath(path)}
		case s.len == len(path):
			s.set = append(s.set, copyPath(path))
		}
	case f != b.Zero:
		n0, b0 := f.Zero(fneg)
		getMinPath0(b, n0, b0, path, s)

		path = append(path, f.Label())
		n1, b1 := f.One(fneg)
		getMinPath0(b, n1, b1, path, s)
	}
}
