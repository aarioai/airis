package types_test

import (
	"github.com/aarioai/airis/pkg/types"
	"reflect"
	"testing"
)

type FakeInt24 int32
type FakeImgSrc struct {
	Provider      int    `json:"provider"`       // Image-processing ID, such as Alibaba Cloud or NetEase Cloud image processing.
	CropPattern   string `json:"crop_pattern"`   // e.g.  https://xxx/img.jpg?width={width:int}&height={height:int}
	ResizePattern string `json:"resize_pattern"` // e.g. https://xxx/img.jpg?maxwidth={max_width:int}
	Origin        string `json:"origin"`         // Not necessarily the real origin.
	Path          string `json:"path"`           // path may be a filename, or a filename that includes a directory
	/*
	   Do not split filename into its own field. There is already too much content, and it increases complexity on the business side.
	*/
	//Filename  string `json:"filename"`  // basename + extension; hand path to the server instead
	Filetype FakeInt24 `json:"filetype"`
	Size     int       `json:"size"`
	Width    int       `json:"width"`
	Height   int       `json:"height"`
	Allowed  [][2]int  `json:"allowed"` // Allowed width and height.
	Jsonkey  string    `json:"jsonkey"` // Field with a special convention.
}

func f0(i any) reflect.Kind {
	return types.PrimitiveType(&i)
}

// Get the primitive type.
func TestPrimitiveType(t *testing.T) {
	type safeInt FakeInt24
	type b safeInt
	type c b
	x := c(100)
	if types.PType(&x) != types.PType(x) || types.PType(x) != reflect.Int32 {
		t.Errorf("types ptype invalid: x:= c(100) %v", types.PType(&x))
	}

	type g1 struct {
		A int64 `json:"a"`
	}
	type g2 g1
	g := g2{A: 200}
	gg := &g
	gg2 := &gg

	t.Log(types.PrimitiveType(&g), types.PType(&g), types.PType(g))

	if types.PType(&g) != types.PType(g) || types.PType(g) != reflect.Struct {
		t.Errorf("types invalid: g %v", types.PType(&g))
	}
	if types.PType(&gg) != types.PType(gg) || types.PType(gg) != reflect.Struct {
		t.Errorf("types invalid: gg %v", types.PType(&gg))
	}
	if types.PType(&gg2) != types.PType(gg2) || types.PType(gg2) != reflect.Struct {
		t.Errorf("types ptypeinvalid: gg2 %v", types.PType(&gg2))
	}

	type y struct {
		A string `json:"a"`
		B int64  `json:"b"`
		C int    `json:"c"`
	}
	a := y{A: "LOVE", B: 100, C: 300}
	t.Log(f0(&a), types.PType(&a), types.PType(a))
}
func TestPrimitiveType2(t *testing.T) {
	type y struct {
		Tmp FakeImgSrc  `json:"-"`
		t   FakeImgSrc  `json:"images"`
		Y   *int        `json:"y"`
		Img *FakeImgSrc `json:"img"`
	}
	type x struct {
		A string `json:"a"`
		B int64  `json:"b"`
		C int    `json:"c"`
		Y *y     `json:"y"`
	}
	yy := 10000
	y0 := y{Y: &yy}
	a := x{A: "LOVE", B: 100, C: 300, Y: &y0}
	if types.PType(&a.Y.Img) != reflect.Struct {
		t.Errorf("types ptypeinvalid: a.Y.Img %v", types.PType(&a.Y.Img))
	}
}
