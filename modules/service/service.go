package service

var services = make(Generators)

// Generator
type Generator func() (Service, error)

// Generators
type Generators map[string]Generator

// List
type List map[string]Service

type Service interface {
	Name() string
}

// Register 將參數加入services中
func Register(k string, gen Generator) {
	if _, ok := services[k]; ok {
		panic("service已被加入")
	}
	services[k] = gen
}

// GetServices 初始化List
func GetServices() List {
	var (
		l   = make(List)
		err error
	)
	for k, gen := range services {
		if l[k], err = gen(); err != nil {
			panic("初始化Service失敗")
		}
	}
	return l
}

// Get 取得匹配的Service
func (g List) Get(k string) Service {
	if s, ok := g[k]; ok {
		return s
	}
	panic("找不到匹配的Service")
}

// Add 加入新的Service至List
func (g List) Add(k string, service Service) {
	if _, ok := g[k]; ok {
		panic("service存在")
	}
	g[k] = service
}
