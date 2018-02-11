package fitting
/*
#cgo LDFLAGS: -L. -lopencv_core
#include "cfitting.h"
*/
import "C"
import (
	"unsafe"
	"math"
	"fmt"
//	"image"
//	"image/color"
//	"image/png"
//	"os"
//	"time"
//	"path/filepath"
)
var (
	Template string  = "template"
)
func ShowTest(lis []int64){
	var tmp,k int64
	for _,li := range lis {
		str:=""
		for i:= 63;i>=0;i--{
			tmp = li >> uint(i)
			if (tmp | 1)  == tmp {
				k =1
			}else{
				k =0
			}
			str = fmt.Sprintf("%d%s",k,str)
		}
		fmt.Println(str)
	}
	var cmd string
	fmt.Scanf("%s",cmd)
	cmd = ""
}

func GetCurveFittingWeight(X []float64,Y []float64,MaxLen int, W []float64) bool {
	_x := (*C.double)(unsafe.Pointer(&X[0]))
	_y := (*C.double)(unsafe.Pointer(&Y[0]))
	Len := C.int(len(X))
	Max := C.int(MaxLen)
//	fmt.Println(W)
	_w := (*C.double)(unsafe.Pointer(&W[0]))
	Out := C.GetCurveWeight(_x,_y,Len,Max,_w)
	if Out != 0 {
		return true
	}
	return false
}

func Rounding(val float64) float64 {
	x,y:= math.Modf(val)
	if y>0.4 {
		x++
	}
	return x
}

func MatrixToIntArray(Matrix [][64]bool) (ts []int64){

	Dx := len(Matrix)
	Dy := 64
	ts = make([]int64,Dx)
	for i:=0;i<Dx;i++ {
		for j:=0;j<Dy;j++{
			ts[i] = ts[i] << 1
			if Matrix[i][j] {
				ts[i] ++
			}
		}
	}
	return ts

}
func GetCoverageRateTest(X []float64,Y []float64,ts []int64) float64 {

	dx := len(ts)
	dy := 64
	var x int
	var y int
	var ty,_ty int64
	var Out float64 = 0
	var testTy []int64 = make([]int64,64)
	for i,_x := range X {
		x = int(Rounding(_x*float64(dx)))-1
		if x <0 {
			x = 0
		}
		if x >= dx {
			Out+=1
			continue
		}
	//	fmt.Println("x",x)
		_ty = ts[x]
		y = int(Rounding(Y[i]*float64(dy)))-1
		if y <0 {
			y = 0
		}
		ty = 1
		ty = ty << (63-uint(y))
		testTy[x] = ty
		if (ty | _ty ) != _ty {
			Out+=1
		}
	}
	ShowTest(ts)
	var cmd string
	fmt.Scanf("%s\r",cmd)
	ShowTest(testTy)
	fmt.Println(Out/float64(len(X)))
	fmt.Scanf("%s\r",cmd)
	return Out/float64(len(X))

}
func GetCoverageRate(X []float64,Y []float64,ts []int64) float64 {

	dx := len(ts)
	dy := 64
	var x int
	var y int
	var ty,_ty int64
	var Out float64 = 0
//	var testTy []int64 = make([]int64,64)
	for i,_x := range X {
		x = int(Rounding(_x*float64(dx))) - 1
		if x <0 {
			x = 0
		}
		if x >= dx {
			Out+=1
			continue
		}
	//	fmt.Println("x",x)
		_ty = ts[x]
		y = int(Rounding(Y[i]*float64(dy))) - 1
		if y < 0 {
		//	Out+=1
		//	continue
			y = 0
		}
		ty = 1
		ty = ty << (63-uint(y))
//		testTy[x] = ty
		if (ty | _ty ) != _ty {
			Out+=1
		}
	}
/**
	ShowTest(ts)
	var cmd string
	fmt.Scanf("%s\r",cmd)
	ShowTest(testTy)
	fmt.Println(Out/float64(len(X)))
	fmt.Scanf("%s\r",cmd)
**/
	return Out/float64(len(X))

}

func MappingMatrix(Matrix [][64]bool,W []float64) bool {

	dx := len(Matrix)
	dy := 64
	var tmpY float64
	var rect [][2]int
	for i:=1 ; i <= dx ; i++ {
		_x := float64(i)/float64(dx)
		tmpY = 0
		for j,_w := range W {
			tmpY += math.Pow(_x,float64(j))*_w
		}
		if tmpY>1.2 || tmpY < -0.2 {
			return false
		}
		rect = append(rect,[2]int{int(Rounding(_x*float64(dx)))-1,int(Rounding(tmpY*float64(dy)))-1})
	}
	for _,r := range rect {
		if r[1] >= dy {
			r[1] = dy -1
		}else if r[1] < 0 {
			r[1] = 0
		}
	//	fmt.Println(r)
		Matrix[r[0]][r[1]] = true
	//	img.Set(r[0],r[1],color.RGBA{10,55,55,55})
	}
	return true

}
func FillMatrix(Matrix [][64]bool) {
	Dx := len(Matrix)
	Dy := 64
	var b,e int
	for i:=0;i<Dx;i++ {
		for j:=0;j<Dy;j++{
		//	_,_,_,a := img.At(i,j).RGBA()
			if Matrix[i][j] {
				b = j+1
				break
			}
		}
		for j:=Dy-1;j>=0;j--{
		//	_,_,_,a := img.At(i,j).RGBA()
		//	if a >0 {
			if Matrix[i][j] {
				e = j
				break
			}
		}
		for j:=b;j<e;j++ {
			Matrix[i][j] = true
	//		img.Set(i,j,color.RGBA{10,55,55,55})
		}
	}
}
func DrawMatrixOEM(X []float64,Y []float64,MaxX int) ([]int64,float64,[]float64,error) {

	Matrix:=make([][64]bool,MaxX)
	lastW := GetBastCols(X,Y)
	var W []float64 = nil
	for _,w := range lastW {
		if MappingMatrix(Matrix,w) {
			W = w
			break
		}
	}
	if W == nil {
		return nil,0,W,fmt.Errorf("Get Best Cols is err")
	}
	var YMin,YMax []float64
	var XMin,XMax []float64

	var tmpY float64
	for i,_x := range X {
		tmpY = 0
		for j,_w := range W {
			tmpY += math.Pow(_x,float64(j))*_w
		}
		if Y[i] >= tmpY {
			YMin = append(YMin,Y[i])
			XMin = append(XMin,_x)
		}
		if Y[i] <= tmpY {
			YMax = append(YMax,Y[i])
			XMax = append(XMax,_x)
		}
	}

	MinW := GetBastCols(XMin,YMin)
	for _,_w := range MinW {
		if MappingMatrix(Matrix,_w) {
			break
		}
	}

	MaxW := GetBastCols(XMax,YMax)
	for _,_w := range MaxW {
		if MappingMatrix(Matrix,_w) {
			break
		}
	}

	FillMatrix(Matrix)
	ts := MatrixToIntArray(Matrix)
	rate:= GetCoverageRate(X,Y,ts)
	return ts,rate,W,nil
}
func DrawMatrix(X []float64,Y []float64,Dx int,Wlen int) ([]int64,error) {

	W := make([]float64,Wlen)
	if (!GetCurveFittingWeight(X,Y,Wlen,W)){
		return nil,fmt.Errorf("err GetCurveFitting")
	}
	Matrix := make([][64]bool,Dx)
	if MappingMatrix(Matrix,W) {
		return MatrixToIntArray(Matrix),nil
	}

//	lastW := GetBastCols(X,Y)
//	for _,w := range lastW {
////		fmt.Println(_i,"--------",len(w))
//		if MappingMatrix(Matrix,w) {
//			return MatrixToIntArray(Matrix),nil
//		}
//	}
	return nil ,fmt.Errorf("W == nil")

}

func CheckCurveFitting(X []float64,Y []float64,W []float64) (valErr float64) {
	var tmpY float64
	for i,_x := range X {
		tmpY = 0
		for j,_w := range W {
			tmpx := math.Pow(_x,float64(j))
		//	fmt.Println(_x,j,tmpx,_w)
			tmpY += _w * tmpx
		}
		valErr += math.Pow((Y[i] - tmpY),2)
		//fmt.Println(i,Y[i] , tmpY)
	}
	return 1/float64(len(X)) * valErr
}
func GetBastCols(X []float64,Y []float64) ([][]float64) {
	var valErr float64 = -1
	var OutW [][]float64
	for i:=2;i<20;i++ {
		W := make([]float64,i)
		if (!GetCurveFittingWeight(X,Y,i,W)){
			break
		}
		tmp := CheckCurveFitting(X,Y,W)
//		fmt.Println(i,tmp)
		if valErr < 0 {
			valErr = tmp
		}else{
			if tmp >= valErr {
				break
			}
			valErr = tmp
		}
		OutW = append(OutW,W)
	}

//	for j:=len(OutW)-1;j>=0;j-- {
//		Ws = append(Ws,OutW[j])
//	}
	for from, to := 0,len(OutW)-1;from<to;from,to = from+1,to-1{
		OutW[from],OutW[to] = OutW[to],OutW[from]
	}
	return OutW

}
