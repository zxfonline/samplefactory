package main

import (
	"bytes"
	"fmt"
	"strings"

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

type ItemKv map[int16]*GDPreItem

func (kvm ItemKv) MarshalText() ([]byte, error) {
	nb := bytes.NewBuffer(nil)
	i := 0
	for _, v := range kvm {
		if i > 0 {
			nb.WriteString(",")
		}
		nb.WriteString(fmt.Sprintf("%d,%d", v.ItemId, v.ItemCnt))
		i++
	}
	return nb.Bytes(), nil
}
func (kvm ItemKv) UnmarshalText(text []byte) error {
	ps := strings.Split(string(text), ",")
	for i := 0; i < len(ps); i += 2 {
		k := csvconfig.Sto16(ps[i])
		v := csvconfig.Sto32(ps[i+1])
		gi := &GDPreItem{ItemId: k, ItemCnt: v}
		kvm[k] = gi
	}
	return nil
}
func (kvm ItemKv) String() string {
	nb := bytes.NewBuffer(nil)
	i := 0
	for k, v := range kvm {
		if i > 0 {
			nb.WriteString(",")
		}
		nb.WriteString(fmt.Sprintf("%d=[id=%d,cnt=%d,%p]", k, v.ItemId, v.ItemCnt, v))
		i++
	}
	return nb.String()
}

type GDItemExchange struct {
	Sid      int    `json:",string"`
	Order    int16  `json:"Order,string"`
	DESC     string `json:"DESC"`
	ItemsReq ItemKv `json:"ItemsReq"` //需要道具
	ItemsGet ItemKv `json:"ItemsGet"` //合成道具
}

func (s *GDItemExchange) GetSid() int {
	return s.Sid
}
func (s *GDItemExchange) Clone() samplefactory.Sample {
	//这里做深层拷贝
	cp := &GDItemExchange{ItemsReq: make(ItemKv), ItemsGet: make(ItemKv)}
	//	samplefactory.DeepGobCopy(cp, s)
	samplefactory.DeepJsonCopy(cp, s)
	return cp
}

type GDPreItem struct {
	ItemId  int16
	ItemCnt int32
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
	err1 := csvconfig.Load([]string{"kv", "itemexchange"})
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
	err = factory.Store("./tmp", 5, func() samplefactory.Sample { return new(KVAttr) })
	if err != nil {
		panic(err)
	}
	//浅层拷贝
	vsample := factory.GetSample(1).(*KVAttr)
	fmt.Printf("vsample=%+v,%+v,%+v,【%+v】,%+v\n", &vsample.Sid, &vsample.Value, &vsample.DESC, vsample.Obj.String(), vsample.Obj.Attr1)
	vsample.Value = 1
	vsample.Obj.Attr1 = 2
	fmt.Println("----")
	fmt.Printf("vsample=%+v,%+v,%+v,【%+v】,%+v\n", &vsample.Sid, &vsample.Value, &vsample.DESC, vsample.Obj.String(), vsample.Obj.Attr1)
	fmt.Println("----")

	vsample = factory.NewSample(1).(*KVAttr)
	vsample.Value = 2
	vsample.Obj.Attr1 = 3
	fmt.Printf("vsample=%+v,%+v,%+v,【%+v】,%+v\n", &vsample.Sid, &vsample.Value, &vsample.DESC, vsample.Obj.String(), vsample.Obj.Attr1)
	fmt.Println("----")
	vsample = factory.GetSample(1).(*KVAttr)
	fmt.Printf("vsample=%+v,%+v,%+v,【%+v】,%+v\n", &vsample.Sid, &vsample.Value, &vsample.DESC, vsample.Obj.String(), vsample.Obj.Attr1)

	factory, err = samplefactory.CreateSampleFacotry("itemexchange", func() samplefactory.Sample {
		return &GDItemExchange{ItemsReq: make(ItemKv), ItemsGet: make(ItemKv)}
	})
	if err != nil {
		panic(err)
	}
	err = factory.Store("./tmp", 5, func() samplefactory.Sample { return new(GDItemExchange) })
	if err != nil {
		panic(err)
	}

	//深层拷贝
	vsample1 := factory.GetSample(1).(*GDItemExchange)
	fmt.Printf("vsample1=%+v,%+v,%+v,【%+v】,%+v\n", &vsample1.Sid, &vsample1.Order, &vsample1.DESC, vsample1.ItemsReq.String(), vsample1.ItemsReq[572].ItemCnt)
	vsample1.ItemsReq[572].ItemCnt = 7
	fmt.Println("----")
	fmt.Printf("vsample1=%+v,%+v,%+v,【%+v】,%+v\n", &vsample1.Sid, &vsample1.Order, &vsample1.DESC, vsample1.ItemsReq.String(), vsample1.ItemsReq[572].ItemCnt)
	fmt.Println("----")

	vsample1 = factory.NewSample(1).(*GDItemExchange)
	vsample1.ItemsReq[572].ItemCnt = 1
	fmt.Printf("vsample1=%+v,%+v,%+v,【%+v】,%+v\n", &vsample1.Sid, &vsample1.Order, &vsample1.DESC, vsample1.ItemsReq.String(), vsample1.ItemsReq[572].ItemCnt)
	fmt.Println("----")
	vsample1 = factory.GetSample(1).(*GDItemExchange)
	fmt.Printf("vsample1=%+v,%+v,%+v,【%+v】,%+v\n", &vsample1.Sid, &vsample1.Order, &vsample1.DESC, vsample1.ItemsReq.String(), vsample1.ItemsReq[572].ItemCnt)
}
