package example 

type Animal interface {
	Sound() 	string
	Age()		int
}


type Dog struct {
	age 	int
	name	string
}

func (d *Dog) Sound() string {
	return d.name + ":Bark"
}

func (d *Dog) Age() int {
	return d.age 
}

type Husky struct {
	Dog 
	fur 	int 
}

func (d *Husky) Sound() string {
	return d.name + ":Boof"
}

func MakeDog(name string, age int) *Dog {
	return &Dog{age:age, name:name}
}

func MakeHusky(name string, age int, fur int) *Husky {
	h := &Husky{fur:fur}
	h.age = age 
	h.name = name
	return h
}