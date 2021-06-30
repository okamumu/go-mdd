package ft

import (
	"fmt"
	"github.com/okamumu/go-mdd/pkg/bdd"
	"testing"
)

func linrange(x0 float64, x1 float64, n int) []float64 {
	ans := make([]float64, n, n)
	x := x0
	d := (x1 - x0) / float64(n-1)
	for i, _ := range ans {
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
	// for _, x := range result {
	// 	for i := 0; i < 20; i++ {
	// 		if i < len(x) {
	// 			fmt.Print(x[i], " ")
	// 		} else {
	// 			fmt.Print(0, " ")
	// 		}
	// 	}
	// 	fmt.Println()
	// }
	return labels, result
}

func TestFT1(t *testing.T) {
	fmt.Println("Run TestFT1")
	labels := []int{0, 1, 2, 3, 4, 5, 6, 7}
	vars := make([]*bdd.Node, len(labels))
	b := bdd.NewBDD()
	for i, x := range labels {
		vars[i] = b.Var(fmt.Sprint(x))
	}
	ft := b.Or(
		b.And(vars[2]),
		b.And(vars[4], vars[0], vars[7]),
		b.And(vars[3], vars[5], vars[2]),
		b.And(vars[1], vars[7], vars[3]),
		b.And(vars[5], vars[7], vars[3]),
		b.And(vars[1], vars[2], vars[6]),
	)
	b.ToDot(ft)
	fmt.Println(Ftmcs(b, ft))
}

func TestFT0(t *testing.T) {
	fmt.Println("Run TestFT0")
	labels := []int{0, 1, 2, 3, 4, 5, 6, 7}
	vars := make([]*bdd.Node01, len(labels))
	b := bdd.NewBDD0()
	for i, x := range labels {
		vars[i] = b.Var(fmt.Sprint(x))
	}
	ft := b.Or(
		b.And(vars[2]),
		b.And(vars[4], vars[0], vars[7]),
		b.And(vars[3], vars[5], vars[2]),
		b.And(vars[1], vars[7], vars[3]),
		b.And(vars[5], vars[7], vars[3]),
		b.And(vars[1], vars[2], vars[6]),
	)
	b.ToDot(ft)
	fmt.Println(Ftmcs0(b, ft))
}

// func BenchmarkBDDFT(b *testing.B) {
// 	labels, result := makeTestData()
// 	d := bdd.NewBDD()
// 	vars := make(map[int]*bdd.Node)
// 	for _, x := range labels {
// 		vars[x] = d.Var(fmt.Sprint(x))
// 	}
// 	b.ResetTimer()
// 	fts := []*bdd.Node{}
// 	for _, xs := range result {
// 		ans := make([]*bdd.Node, 0, len(xs))
// 		for _, x := range xs {
// 			ans = append(ans, vars[x])
// 		}
// 		fts = append(fts, d.And(ans...))
// 	}
// 	ft := d.Or(fts...)
// 	fmt.Println(ft)
// }

func TestBDDFTMCS1(t *testing.T) {
	fmt.Println("Run TestBDDFTMCS1")
	labels, result := makeTestData()
	d := bdd.NewBDD()
	vars := make(map[int]*bdd.Node)
	for _, x := range labels {
		vars[x] = d.Var(fmt.Sprint(x))
	}
	fts := []*bdd.Node{}
	for _, xs := range result {
		ans := make([]*bdd.Node, 0, len(xs))
		for _, x := range xs {
			ans = append(ans, vars[x])
		}
		fts = append(fts, d.And(ans...))
	}
	ft := d.Or(fts...)
	mcs := Ftmcs(d, ft)
	fmt.Println(mcs)
}

func TestBDDFTMCS0(t *testing.T) {
	fmt.Println("Run TestBDDFTMCS0")
	labels, result := makeTestData()
	d := bdd.NewBDD0()
	vars := make(map[int]*bdd.Node01)
	vars2 := make(map[string]*bdd.Node01)
	for _, x := range labels {
		vars[x] = d.Var(fmt.Sprint(x))
		vars2[fmt.Sprint(x)] = vars[x]
	}
	fts := []*bdd.Node01{}
	for _, xs := range result {
		ans := make([]*bdd.Node01, 0, len(xs))
		for _, x := range xs {
			ans = append(ans, vars[x])
		}
		fts = append(fts, d.And(ans...))
	}
	ft := d.Or(fts...)
	fmt.Println(ft)
	// ss, minf := ftmcs0(d, vars2, ft)
	// fmt.Println(ss)
	// r := d.Not(d.Imp(ft, minf))
	// ss, minf = ftmcs0(d, vars2, r)
	// fmt.Println(ss)
	// r = d.Not(d.Imp(r, minf))
	// ss, minf = ftmcs0(d, vars2, r)
	// fmt.Println(ss)
	// r = d.Not(d.Imp(r, minf))
	// ss, minf = ftmcs0(d, vars2, r)
	// fmt.Println(ss)
	// r = d.Not(d.Imp(r, minf))
	// ss, minf = ftmcs0(d, vars2, r)
	// fmt.Println(ss)
	// r = d.Not(d.Imp(r, minf))
	// ss, minf = ftmcs0(d, vars2, r)
	// fmt.Println(ss)
	// r = d.Not(d.Imp(r, minf))
	// ss, minf = ftmcs0(d, vars2, r)
	// fmt.Println(ss)
	// r = d.Not(d.Imp(r, minf))
	// ss, minf = ftmcs0(d, vars2, r)
	// fmt.Println(ss)
	// r = d.Not(d.Imp(r, minf))
	// ss, minf = ftmcs0(d, vars2, r)
	// fmt.Println(ss)
	// r = d.Not(d.Imp(r, minf))
	// ss, minf = ftmcs0(d, vars2, r)
	// fmt.Println(ss)
	// r = d.Not(d.Imp(r, minf))
	// ss, minf = ftmcs0(d, vars2, r)
	// fmt.Println(ss)
	// d.ToDot(ft)
	mcs := Ftmcs0(d, ft)
	fmt.Println(mcs)
}

// func BenchmarkBDDFT1(b *testing.B) {
// 	labels, result := makeTestData()
// 	d := bdd.NewBDD()
// 	vars := make(map[int]*bdd.Node)
// 	vars2 := make(map[string]*bdd.Node)
// 	for _, x := range labels {
// 		vars[x] = d.Var(fmt.Sprint(x))
// 		vars2[fmt.Sprint(x)] = vars[x]
// 	}
// 	b.ResetTimer()
// 	fts := []*bdd.Node{}
// 	for _, xs := range result {
// 		ans := make([]*bdd.Node, 0, len(xs))
// 		for _, x := range xs {
// 			ans = append(ans, vars[x])
// 		}
// 		fts = append(fts, d.And(ans...))
// 	}
// 	ft := d.Or(fts...)
// 	mcs := Ftmcs(d, ft)
// 	fmt.Println(mcs)
// }

func BenchmarkBDDFT0(b *testing.B) {
	fmt.Println("Run BenchBDDFT0")
	labels, result := makeTestData()
	d := bdd.NewBDD0()
	vars := make(map[int]*bdd.Node01)
	vars2 := make(map[string]*bdd.Node01)
	for _, x := range labels {
		vars[x] = d.Var(fmt.Sprint(x))
		vars2[fmt.Sprint(x)] = vars[x]
	}
	b.ResetTimer()
	fts := []*bdd.Node01{}
	for _, xs := range result {
		ans := make([]*bdd.Node01, 0, len(xs))
		for _, x := range xs {
			ans = append(ans, vars[x])
		}
		fts = append(fts, d.And(ans...))
	}
	ft := d.Or(fts...)
	mcs := Ftmcs0(d, ft)
	fmt.Println(mcs)
}
