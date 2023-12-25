package ActualCombat

// 创建衣服工厂接口

type IClothesFactory interface {
	MakeClothes(name string) IClothes
}

type IClothes interface {
	setName(name string)
	setSize(size string)
	GetName() string
	GetSize() string
}

// GYClothesStore 工厂类
type GYClothesStore struct {
}

// MakeClothes 实现工厂方法
func (g *GYClothesStore) MakeClothes(name string) IClothes {
	// 根据不同的输入，工厂返回不同的产品
	switch name {
	case "peak":
		return NewPEAK()
	case "anta":
		return NewANTA()
	}
	return nil
}

// 创建具体的产品类
type clothes struct {
	name string
	size string
}

func (c *clothes) setName(name string) {
	c.name = name
}

func (c *clothes) setSize(size string) {
	c.size = size
}

func (c *clothes) GetName() string {
	return c.name
}

func (c *clothes) GetSize() string {
	return c.size
}

// PEAK 若要创建具体的产品，直接包一层,这样PEAK直接实现了产品的接口
type PEAK struct {
	clothes
}

type ANTA struct {
	clothes
}

// NewPEAK 创建不同种类的产品
func NewPEAK() IClothes {
	var p PEAK
	p.setSize("xxl")
	p.setName("peak")
	return &p
}

// NewANTA 创建不同种类的产品
func NewANTA() IClothes {
	var p ANTA
	p.setSize("xl")
	p.setName("anta")
	return &p
}
