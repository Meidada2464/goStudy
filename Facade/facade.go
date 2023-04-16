package Facade

import "fmt"

// 这个简单的门面模式类似于我们买车，我们知道一个汽车真正怎么造出来的。
// 我们不知道里面的细节，只知道我们需要购买的汽车的参数，对参数（外观）进行选择即可得到一个真正的示例
// 该例子中，经过层层包装，直接调用门面的方法CreateNewCompleteCar即可得到实例
// 类似于提供了一个set方法供外部调用，能对结构体里面的私有参数赋值

type CarModel struct {
	model string
}

func NewCarModel() *CarModel {
	return &CarModel{}
}

func (m *CarModel) SetCarModel(s string) {
	m.model = s
	fmt.Println("loading car Model ", s)
}

type CarEngine struct {
	engine string
}

func NewCarEngine() *CarEngine {
	return &CarEngine{}
}

func (e *CarEngine) SetCarEngine(s string) {
	e.engine = s
	fmt.Println("then loading car engine ", s)
}

type CarBody struct {
	body string
}

func NewCarBody() *CarBody {
	return &CarBody{}
}

func (b *CarBody) SetCarBody(s string) {
	b.body = s
	fmt.Println("loading car body ", s)
}

// CarFacade 组合汽车的零部件
type CarFacade struct {
	Model  CarModel
	Engine CarEngine
	Body   CarBody
}

// NewCarFacade 进行初始化用的
func NewCarFacade() *CarFacade {
	return &CarFacade{
		Model:  CarModel{},
		Engine: CarEngine{},
		Body:   CarBody{},
	}
}

func (f *CarFacade) CreateNewCompleteCar(m, e, b string) {
	f.Model.SetCarModel(m)
	f.Engine.SetCarEngine(e)
	f.Body.SetCarBody(b)
}
