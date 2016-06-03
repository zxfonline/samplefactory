package main

import (
	"fmt"

	"github.com/zxfonline/csvconfig"
	"github.com/zxfonline/samplefactory"
)

type OBJ struct {
	Attr1 int
	Attr2 string
	Att3  map[string]string
	Att4  []*OBJ1
}
type OBJ1 struct {
	AP1 string
	AP2 float32
}

func (o *OBJ) String() string {
	return fmt.Sprintf("attr1=%+v,attr2=%+v,att3=%+v,att4=%+v", &o.Attr1, &o.Attr2, *&o.Att3, &o.Att4)
}

func NewObj() *OBJ {
	n := &OBJ{
		Attr1: 1,
		Attr2: "2",
		Att3:  make(map[string]string),
		Att4:  make([]*OBJ1, 0, 4),
	}
	n.Att3["a31"] = "311"
	n.Att3["a32"] = "321"
	n.Att3["a33"] = "331"
	n.Att4 = append(n.Att4, &OBJ1{"41", 41})
	n.Att4 = append(n.Att4, &OBJ1{"42", 42})
	n.Att4 = append(n.Att4, &OBJ1{"43", 43})
	n.Att4 = append(n.Att4, &OBJ1{"44", 44})
	return n
}

type KVAttr struct {
	Sid   int `json:",string"`
	Value int `json:",string"`
	DESC  string
	Obj   *OBJ `json:"-" lua:"-"` // lua:"-"
}

func (s *KVAttr) GetSid() int {
	return s.Sid
}
func (s *KVAttr) Clone() samplefactory.Sample {
	//这里只做了浅层拷贝，从csv文件而来的字段，只需要浅层拷贝足够
	cn := *s
	return &cn
}

func main() {
	//初始化配置文件信息
	csvconfig.Init("", "")
	err1 := csvconfig.Load([]string{"kv"})
	if err1 != nil {
		panic(err1)
	}
	factory, err := samplefactory.CreateSampleFacotry("kv", func() samplefactory.Sample {
		return &KVAttr{
			Obj: NewObj(),
		}
	})
	if err != nil {
		panic(err)
	}
	//	vsample := factory.GetSample(1).(*KVAttr)

	//	fmt.Printf("vsample=%+v,%+v,%+v,%+v,%+v,[%+v],%+v\n", &vsample, vsample, &vsample.Sid, &vsample.Value, &vsample.DESC, vsample.Obj.String(), vsample.Obj.Attr1)
	//	vsample.Value = 1
	//	vsample.Obj.Attr1 = 2
	//	fmt.Println("----")
	//	fmt.Printf("vsample=%+v,%+v,%+v,%+v,%+v,[%+v],%+v\n", &vsample, vsample, &vsample.Sid, &vsample.Value, &vsample.DESC, vsample.Obj.String(), vsample.Obj.Attr1)
	//	fmt.Println("----")

	//	vsample = factory.GetSample(1).(*KVAttr)
	//	fmt.Printf("vsample=%+v,%+v,%+v,%+v,%+v,[%+v],%+v\n", &vsample, vsample, &vsample.Sid, &vsample.Value, &vsample.DESC, vsample.Obj.String(), vsample.Obj.Attr1)
	//	vsample.Value = 2
	//	fmt.Println("----")
	//	vsample = factory.GetSample(1).(*KVAttr)
	//	fmt.Printf("vsample=%+v,%+v,%+v,%+v,%+v,[%+v],%+v\n", &vsample, vsample, &vsample.Sid, &vsample.Value, &vsample.DESC, vsample.Obj.String(), vsample.Obj.Attr1)

	err = factory.Store("./tmp", 4, func() samplefactory.Sample { return new(KVAttr) })
	if err != nil {
		panic(err)
	}
}
