package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"gonum.org/v1/plot/vg"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

const (
	max, min = 0.0, 199.0
)

func main() {

	var class [][]float64

	//重み
	w := []float64{18.6, 11.1, 81.1}

	//初期の重み
	firstW := []float64{18.6, 11.1, 81.1}

	// 図の生成
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	//任意の点
	dots := make(plotter.XYs, 2)

	//クラス1
	x1, y1 := 8.0, 2.0
	dots[0].X = x1
	dots[0].Y = y1

	//クラス2
	x2, y2 := 3.0, 6.0
	dots[1].X = x2
	dots[1].Y = y2

	//各クラスのサンプル
	n := 100
	class1, plotdata1 := randomPoints(n, x1, y1)
	class2, plotdata2 := randomPoints(n, x2, y2)
	class = append(class1, class2...)

	//教師データ作成
	b := make([]float64, n*2)
	makeTrainData(b, n)

	var errGraph []float64
	var beforeError float64
	var afterError float64
	afterError = 10000
	var count int

	for {
		count++

		beforeError = afterError

		if count < 100 {
			errGraph = append(errGraph, beforeError)
		}

		//誤差関数の微分
		differentialErrorFunc := differentialErrorFunc(class, w, b)

		//重みの更新
		w = weightCalc(w, differentialErrorFunc, 0.0001)

		//誤差を求める
		afterError = errorFunc(class, w, b)

		//前回と今回の誤差の二乗が閾値以下だったら終了
		if (beforeError-afterError)*(beforeError-afterError) < 0.0000001 {
			break
		}
	}

	fmt.Println(w)
	p2, err := plot.New()
	if err != nil {
		panic(err)
	}

	// Make a line plotter and set its style.
	l, err := plotter.NewLine(lineGraph(errGraph))
	if err != nil {
		panic(err)
	}
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(0)}
	l.LineStyle.Color = color.RGBA{B: 245, A: 125}

	p2.Add(l)
	p2.Title.Text = "Plotutil example"
	p2.X.Label.Text = "X"
	p2.Y.Label.Text = "Y"

	p2.Legend.Add("line", l)
	// Save the plot to a PNG file.
	if err := p2.Save(4*vg.Inch, 4*vg.Inch, "points.png"); err != nil {
		panic(err)
	}

	//初期境界線のplot--------------------------------------------------------
	border := plotter.NewFunction(func(x float64) float64 {
		//x2 = -(w1 / w2)*x1 - w0 / w2
		return -(firstW[1]/firstW[2])*x - (firstW[0] / firstW[2])
	})
	border.Color = color.RGBA{B: 155, A: 5}
	//----------------------------------------------------------------------

	//最終境界線のplot--------------------------------------------------------
	lastBorder := plotter.NewFunction(func(x float64) float64 {
		//x2 = -(w1 / w2)*x1 - w0 / w2
		return -(w[1]/w[2])*x - (w[0] / w[2])
	})
	lastBorder.Color = color.RGBA{B: 255, A: 255}
	//----------------------------------------------------------------------

	//label
	p.Title.Text = "Points Example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// Draw a grid behind the data
	p.Add(plotter.NewGrid())

	// Make a scatter plotter and set its style.
	s, err := plotter.NewScatter(plotdata1)
	if err != nil {
		panic(err)
	}

	y, err := plotter.NewScatter(plotdata2)
	if err != nil {
		panic(err)
	}

	r, err := plotter.NewScatter(dots)
	if err != nil {
		panic(err)
	}

	s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 55}
	y.GlyphStyle.Color = color.RGBA{R: 155, B: 128, A: 255}
	r.GlyphStyle.Color = color.RGBA{R: 128, B: 0, A: 0}
	p.Add(s)
	p.Add(y)
	p.Add(r)
	p.Add(lastBorder)
	p.Add(border)
	p.Legend.Add("class1", s)
	p.Legend.Add("class2", y)

	// Axis ranges
	p.X.Min = 0
	p.X.Max = 10
	p.Y.Min = 0
	p.Y.Max = 10

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, 6*vg.Inch, "report.png"); err != nil {
		panic(err)
	}
}

func randomCount(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float64()*(max-min) + min
}

//ガウス分布
func random(axis float64) float64 {
	//分散
	dispersion := 1.0
	rand.Seed(time.Now().UnixNano())
	return rand.NormFloat64()*dispersion + axis
}

//学習データの生成
func randomPoints(n int, x, y float64) ([][]float64, plotter.XYs) {
	matrix := make([][]float64, n)
	pts := make(plotter.XYs, n)
	for i := range matrix {
		l := random(x)
		m := random(y)
		matrix[i] = []float64{1.0, l, m}
		pts[i].X = l
		pts[i].Y = m
	}
	return matrix, pts
}

//学習
//前半に-1,後半に1を格納
func makeTrainData(b []float64, n int) {
	for i := 0; i < n*2; i++ {
		if i >= n {
			b[i] = 1
		} else {
			b[i] = -1
		}
	}
}

//重みの更新
func weightCalc(w, differentialErrorFunc []float64, p float64) []float64 {
	//fmt.Println(innerProduct)
	//fmt.Println("内積")
	w[0] = w[0] - p*differentialErrorFunc[0]
	w[1] = w[1] - p*differentialErrorFunc[1]
	w[2] = w[2] - p*differentialErrorFunc[2]

	return w
}

//誤差関数の微分
func differentialErrorFunc(class [][]float64, w, b []float64) (result []float64) {
	result = []float64{0.0, 0.0, 0.0}
	for i, x := range class {
		innerProduct := innerProduct(w, x)
		result[0] += (innerProduct - b[i]) * x[0]
		result[1] += (innerProduct - b[i]) * x[1]
		result[2] += (innerProduct - b[i]) * x[2]

	}
	return
}

//内積計算
func innerProduct(w, x []float64) (f float64) {
	if len(w) != len(x) {
		panic("エラーですよ")
	}

	for i := range w {
		f += w[i] * x[i]
	}

	return
}

func errorFunc(class [][]float64, w, b []float64) (result float64) {
	result = 0.0
	for i, x := range class {
		calc := innerProduct(w, x) - b[i]
		calc = calc * calc
		result += calc
	}
	return
}

// 誤差関数の出力
func lineGraph(n []float64) plotter.XYs {
	pts := make(plotter.XYs, len(n))
	for i, m := range n {
		//fmt.Println(m)
		pts[i].X = float64(i)
		pts[i].Y = m
	}
	return pts
}
