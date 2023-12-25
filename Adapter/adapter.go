package Adapter

// 适配器模式是为了解决现有接口\类结构不符合用户的期望，我们再此基础上封装一层。且用户不能直接调用原接口

// 该模式的重点是：
// 1、将一个接口转换成用户希望的另一个接口

// 单例模式包含的角色： 1、Target:目标抽象类 2、Adapter:适配器类 3、Adapted:适配者类 4、Client：客户类

// Adapted 创建一个现有的接口
type Adapted interface {
	// SpecificRequest 接口定义一个需要实现的方法实现方法
	SpecificRequest() string
}

// 创建一个实现类
type adaptedImpl struct {
}

// SpecificRequest 实现Adapted接口
func (a *adaptedImpl) SpecificRequest() string {
	return "this is adapted"
}

// NewAdaptedFactory 创建一个创建被适配者的生产工厂
func NewAdaptedFactory() Adapted {
	return &adaptedImpl{}
}

// Target 创建目标接口，即我们需要转化后的接口
type Target interface {
	Request() string
}

// NewAdapter 为了做装换，我们在原来的基础上加一层，入参是被适配者的类型(可以接受接口的实现类)，出参是目标接口,
//
//	同时需要定义一个被装换的类
func NewAdapter(ad Adapted) Target {
	return &adapter{
		Adapted: ad,
	}
}

// 这个类做中间的转换，需要实现目标类的方法
type adapter struct {
	// 类中组合了一个接口，以接口为属性
	Adapted
}

func (a *adapter) Request() string {
	return "this is adapter ,complement target class"
}
