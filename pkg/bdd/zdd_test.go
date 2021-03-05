package bdd

import (
	"fmt"
	"testing"
	"sort"
)

// utils

func (d *ZDD) todot(x *Node) {
	x.ToDot(func(x *Node) string {
		if x == d.one {
			return "1"
		} else {
			return "0"
		}
	})
}

func NewZDDWithLabels(labels []int) (*ZDD, map[int]*Node) {
	d := NewZDD()
	for _, i := range labels {
		d.NewHeader(fmt.Sprint(i), []DomainInt{0, 1})
	}
	nodes := make(map[int]*Node)
	for _, i := range labels {
		// make create vector
		x := []int{1}
		for _, j := range labels {
			if i != j {
				x = append(x, 0)
			} else {
				x = append(x, 1)
			}
		}
		nodes[i] = d.CreateNode(x)
	}
	return d, nodes
}

func linrange(x0 float64, x1 float64, n int) []float64 {
	ans := make([]float64, n, n)
	x := x0
	d := (x1 - x0) / float64(n-1)
	for i,_ := range ans {
		ans[i] = x
		x += d
	}
	return ans
}

func makeTestData() ([]int, [][]int) {
		data := [][]float64{
		{0.2799439284864673, 0.019039179685146124},
		{0.17006659269016278, 0.26812585079180584},
		{0.37160535186424815, 0.28336464179809084},
		{0.39279646612146735, 0.2789501222723816},
		{0.44911286867346534, 0.2067605663915406},
		{0.505207192733002, 0.07778618522601977},
		{0.381127966318632, 0.36580119057192695},
		{0.14314120834324617, 0.5282542334011777},
		{0.2236688207291726, 0.5027191237033151},
		{0.5865981007905481, 0.05016684053503706},
		{0.2157117712338983, 0.5699545901561343},
		{0.6600618683347792, 0.006513992462842788},
		{0.6964756269215944, 0.031164261499776913},
		{0.5572474263734104, 0.5457354609512821},
		{0.38370575109517757, 0.6870518034076929},
		{0.14047278702240318, 0.8099471630562083},
		{0.6117903795750386, 0.6200985639530681},
		{0.8350140149860443, 0.26002375370524433},
		{0.621745085645081, 0.6249760808944675},
		{0.9223171788742697, 0.040441694559945285},
		{0.40157225733130186, 0.8622123559544623},
		{0.5654235033016655, 0.7840149804945578},
		{0.8605048496383341, 0.48642029259985065},
		{0.5627138851056968, 0.8499394786290626},
		{0.7124617313668333, 0.7347698978106127},
		{0.9656307414336753, 0.3647058735973785},
		{0.9944967296698335, 0.548297306757731},
		{0.5733819926662398, 0.9813641372820436},
		{0.9236020954856745, 0.7540471034450749},
		{0.8910887808888235, 0.8901974734237881},
	}
	labels := []int{13, 19, 2, 14, 29, 23, 3, 26, 25, 7, 9, 27, 12, 30, 17, 24, 8, 4, 18, 5, 20, 21, 28, 1, 16, 10, 15, 6, 11, 22}
	ddx := linrange(0, 1, 30)
	ddy := linrange(0, 1, 30)
	result := [][]int{}
	for _, x := range ddx {
		for _, y := range ddy {
			v := []int{}
			for i, p := range data {
				tmpx := x - p[0]
				tmpy := y - p[1]
				if tmpx*tmpx+tmpy*tmpy < 0.3*0.3 {
					v = append(v, labels[i])
				}
			}
			result = append(result, v)
		}
	}
	return labels, result
}

func (d *ZDD) getmcs(x *Node) [][]string {
	switch {
	case x == d.one:
		return [][]string{[]string{}}
	case x == d.zero:
		return [][]string{}
	default:
		result0 := d.getmcs(x.nodes[0])
		result1 := d.getmcs(x.nodes[1])
		for i,_ := range result1 {
			result1[i] = append(result1[i], x.header.label)
		}
		ans := append(result0, result1...)
		sort.Slice(ans, func(i, j int) bool {
			if len(ans[i]) < len(ans[j]) {
				return true
			} else {
				return false
			}
		})
		s := [][]string{}
		for i := 0; i < len(ans); i++ {
			if isIncludeList(ans[i], s) == false {
				s = append(s, ans[i])
			}
		}
		return s
	}
}

func isIncludeList(x []string, lst [][]string) bool {
	for _, y := range lst {
		if isInclude(y, x) == true {
			return true
		}
	}
	return false
}

func isInclude(x []string, y []string) bool {
	for _, w := range x {
		if hasWord(w, y) == false {
			return false
		}
	}
	return true
}

func hasWord(w string, s []string) bool {
	for _, x := range s {
		if x == w {
			return true
		}
	}
	return false
}

// tests

func TestZDD1(t *testing.T) {
	f := NewZDD()
	one := f.terminals[true]
	zero := f.terminals[false]
	fmt.Println(one)
	fmt.Println(zero)
	h1 := f.NewHeader("x", []DomainInt{0, 1})
	h2 := f.NewHeader("y", []DomainInt{0, 1})
	fmt.Println(h1, h2)
	x := f.Node(h2, f.Node(h1, one, zero), f.Node(h1, one, zero))
	// v1 := f.Node(h1, zero, zero)
	// v2 := f.Node(h1, one, one)
	// v1.ToDot(func(x *Node) string {
	// 	if x == one {
	// 		return "1"
	// 	} else {
	// 		return "0"
	// 	}
	// })
	// v2.ToDot(func(x *Node) string {
	// 	if x == one {
	// 		return "1"
	// 	} else {
	// 		return "0"
	// 	}
	// })
	y := f.Node(h2, f.Node(h1, zero, zero), f.Node(h1, one, one))
	z := f.Intersect(x, y)
	fmt.Println(z)
	f.todot(x)
	f.todot(y)
	f.todot(z)
}

func TestCreateZDD1(t *testing.T) {
	d := NewZDD()
	d.NewHeader("x1", []DomainInt{0, 1})
	d.NewHeader("x2", []DomainInt{0, 1})
	d.NewHeader("x3", []DomainInt{0, 1})
	x := d.CreateNode([]int{1, 0, 1, 0})
	d.todot(x)
}

func TestZDDFT1(t *testing.T) {
	labels, result := makeTestData()
	d, vars := NewZDDWithLabels(labels)
	fts := []*Node{}
	for _, xs := range result {
		ans := vars[xs[0]]
		for i := 1; i < len(xs); i++ {
			ans = d.Product(ans, vars[xs[i]])
		}
		fts = append(fts, ans)
	}
	ft := d.Union(fts...)
	d.todot(ft)
}

func TestZDDFT2(t *testing.T) {
	// fmt.Println(isInclude([]string{"1", "2"}, []string{"1", "2", "3"}))
	labels := []int{1, 2, 3, 4, 5, 6, 7}
	d, vars := NewZDDWithLabels(labels)
	ft := d.Union(d.Product(vars[2]), d.Product(vars[4], vars[6], vars[7]), d.Product(vars[3], vars[5], vars[2]),
		d.Product(vars[1], vars[7], vars[3]), d.Product(vars[5], vars[7], vars[3]), d.Product(vars[1], vars[2], vars[6]))
	d.todot(ft)
	result := d.getmcs(ft)
	fmt.Println(result)
}

func TestZDDFTMCS1(t *testing.T) {
	labels, result := makeTestData()
	d, vars := NewZDDWithLabels(labels)
	fts := []*Node{}
	for _, xs := range result {
		ans := vars[xs[0]]
		for i := 1; i < len(xs); i++ {
			ans = d.Product(ans, vars[xs[i]])
		}
		fts = append(fts, ans)
	}
	ft := d.Union(fts...)
	fmt.Println(ft)
	mcs := d.getmcs(ft)
	fmt.Println(mcs)
}

