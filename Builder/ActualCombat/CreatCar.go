package ActualCombat

// InterfaceBuilder 创建一个生成器接口
type InterfaceBuilder interface {
	// 规定生成对象的步骤
	SetSeatsType()
	SetEngineType()
	SetNumber()
	GetCar() Car
}

type Car struct {
	seats  string
	Engine string
	Number int
}

// MpvBuilder 定义具体的生成器
type MpvBuilder struct {
	SeatsType  string
	EngineType string
	Number     int
}

func (m *MpvBuilder) SetSeatsType() {
	m.SeatsType = "七座"
}

func (m *MpvBuilder) SetEngineType() {
	m.EngineType = "V8"
}

func (m *MpvBuilder) SetNumber() {
	m.Number = 199999
}

func (m *MpvBuilder) GetCar() Car {
	return Car{
		seats:  m.SeatsType,
		Engine: m.EngineType,
		Number: m.Number,
	}
}

type SuvBuilder struct {
	SeatsType  string
	EngineType string
	Number     int
}

func (s *SuvBuilder) SetSeatsType() {
	s.SeatsType = "五座"
}

func (s *SuvBuilder) SetEngineType() {
	s.EngineType = "V12"
}

func (s *SuvBuilder) SetNumber() {
	s.Number = 299999
}

func (s *SuvBuilder) GetCar() Car {
	return Car{
		seats:  s.SeatsType,
		Engine: s.EngineType,
		Number: s.Number,
	}
}

// NewBuilder 获得生成器的方法
func NewBuilder(name string) InterfaceBuilder {
	switch name {
	case "mvp":
		return &MpvBuilder{}
	case "suv":
		return &SuvBuilder{}
	}
	return nil
}

// =================================================

type Director struct {
}
