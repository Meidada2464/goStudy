package Builder

import "fmt"

// 建造者模式属于对象创建型模式
// 建造者模式可以一步一步地建造一个对象，允许用户只通过指定复杂对象的类型和内容就可以构建它们，无需知道内部的具体构造细节

// 模式结构包含一下结构
// 1、Builder: 抽象建造者	2、ConcreteBuilder: 具体建造者  3、Director： 指挥官  4、Product：产品角色
// 指挥官是关键，其属性是抽象建造者，指挥官有其方法可以规定组装顺序很重要

type Builder interface {
	SetBrand(s string)
	SetEngine(s string)
	SetPrice(s string)
}

type Director struct {
	// 结构体里面的属性可以是借口类型的
	builder Builder
}

func (d *Director) Construct(brand, engine, price string) {
	// 规定建造顺序
	d.builder.SetBrand(brand)
	d.builder.SetEngine(engine)
	d.builder.SetPrice(price)
}

// NewDirector 创建一个方法，可以创建不同类型的指挥官
func NewDirector(bu Builder) Director {
	return Director{
		builder: bu,
	}
}

// Product1 Product 创建具体z制造者
type Product1 struct {
	// 比如这个制造商有不同的属性
	Brand  string
	Engine string
	Price  string
}

func NewProduct1() *Product1 {
	return &Product1{}
}

func (p *Product1) getProduct() Product1 {
	return *p
}

func (p *Product1) SetBrand(s string) {
	fmt.Println("loading P1 Brand:", s)
	p.Brand = s
}

func (p *Product1) SetEngine(s string) {
	fmt.Println("loading P1 Engine:", s)
	p.Engine = s
}

func (p *Product1) SetPrice(s string) {
	fmt.Println("loading P1 Price:", s)
	p.Price = s
}

// Product2 Product1 Product 创建具体z制造者
type Product2 struct {
	// 比如这个制造商有不同的属性
	Brand  string
	Engine string
	Price  string
}

func NewProduct2() *Product2 {
	return &Product2{}
}

func (p *Product2) getProduct() Product2 {
	return *p
}

func (p *Product2) SetBrand(s string) {
	fmt.Println("loading P2 Brand:", s)
	p.Brand = s
}

func (p *Product2) SetEngine(s string) {
	fmt.Println("loading P2 Engine:", s)
	p.Engine = s
}

func (p *Product2) SetPrice(s string) {
	fmt.Println("loading P2 Price:", s)
	p.Price = s
}
