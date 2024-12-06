package types_test

import (
	"github.com/aarioai/airis/pkg/types"
	"reflect"
	"testing"
)

type FakeInt24 int32
type FakeImgSrc struct {
	Provider      int    `json:"provider"`       // 图片处理ID，如阿里云图片处理、网易云图片处理等
	CropPattern   string `json:"crop_pattern"`   // e.g.  https://xxx/img.jpg?width=${WIDTH}&height=${HEIGHT}
	ResizePattern string `json:"resize_pattern"` // e.g. https://xxx/img.jpg?maxwidth=${MAXWIDTH}
	Origin        string `json:"origin"`         // 不一定是真实的
	Path          string `json:"path"`           // path 可能是 filename，也可能是 带文件夹的文件名
	/*
	   不要独立出来 filename，一方面太多内容了；另一方面增加业务侧复杂度
	*/
	//Filename  string `json:"filename"`  // basename + extension  直接交path给服务端处理
	Filetype FakeInt24 `json:"filetype"`
	Size     int       `json:"size"`
	Width    int       `json:"width"`
	Height   int       `json:"height"`
	Allowed  [][2]int  `json:"allowed"` // 允许的width,height
	Jsonkey  string    `json:"jsonkey"` // 特殊约定字段
}

func f0(i any) reflect.Kind {
	return types.PrimitiveType(&i)
}

// 获取原始类型
func TestPrimitiveType(t *testing.T) {
	type safeInt FakeInt24
	type b safeInt
	type c b
	x := c(100)
	t.Log(types.PrimitiveType(&x), types.PType(&x), types.PType(x))
	type g1 struct {
		A int64 `json:"a"`
	}
	type g2 g1
	g := g2{A: 200}
	gg := &g
	gg2 := &gg

	t.Log(types.PrimitiveType(&g), types.PType(&g), types.PType(g))
	t.Log(types.PrimitiveType(&gg), types.PType(&gg), types.PType(gg))
	t.Log(types.PrimitiveType(&gg2), f0(gg2), types.PType(&gg2), types.PType(gg2))

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
	t.Log(types.PType(a.Y.Img))
}
