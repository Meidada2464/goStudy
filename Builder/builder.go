package Builder

// 生成器模式
// 角色：1、生成器（builder） 2、具体生成器（ConcreteBuilder） 3、产品（Product） 4、主管（Director）

// Builder 创建生成器
type Builder interface {
	Build()
}

// ConcreteBuilder 具体生成器
type ConcreteBuilder struct {
	result Product
}

func (c *ConcreteBuilder) Build() {
	// 包含完整的创建对象的过程
	c.result = Product{
		name: "产品1",
		size: 25,
	}
}

func (c *ConcreteBuilder) GetResult() Product {
	// 返回在生成步骤中生成的产品对象
	return c.result
}

type Product struct {
	name string
	size int
}

// NewConcreteBuilder 初始化一个生成器
func NewConcreteBuilder() ConcreteBuilder {
	// 初始化时对生成器的内容进行初始化
	return ConcreteBuilder{result: Product{}}
}

// ========================================================

// Director 创建主管，用于控制生成最终产品对象算法的类
type Director struct {
	builder Builder
}

func (d *Director) Construct() {
	// 用于定义生成对象时的步骤
	d.builder.Build()
}

func NewDirector(builder Builder) Director {
	return Director{builder}
}
