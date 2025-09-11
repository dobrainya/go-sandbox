package main

import "fmt"

type Animal struct {
	name string
	kind string
}

type PersonInterface interface {
	getName() string
	setName(value string)
	move()
}

type Person struct {
	name     string
	password string
	age      int
	pets     []*Animal
}

func (p Person) getName() string {
	return p.name
}

func (p Person) dump() {
	fmt.Println(p)
}

func (p *Person) setName(value string) {
	p.name = value
}

func (p Person) move() {
	fmt.Println("I am moving")
}

func main() {
	var pin PersonInterface
	var pets []*Animal
	pets = append(pets, &Animal{name: "Boby", kind: "dog"})
	var person1 Person = Person{name: "Kate", password: "1234", age: 25, pets: pets}
	pin = &person1

	pin.move()
	fmt.Println(pin.getName())

	person1.dump()

	var pin2 PersonInterface
	person2 := Person{"John", "12345", 30, pets}
	pin2 = &person2
	person2.dump()

	pin2.move()

	fmt.Println(person1.getName())

	person1.name = "Masha"
	person1.age = 27
	person1.dump()

	changeName := func(u Person, value string) {
		u.name = value
	}

	changeName(person1, "Kate")

	fmt.Println(person1.getName())

	changeNameV2 := func(u *Person, value string) {
		u.name = value
	}

	changeNameV2(&person1, "Kate")

	fmt.Println(person1.getName())

	person1.setName("Alice")

	fmt.Println(person1.getName())
}
