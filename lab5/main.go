package main

import(
	"fmt"
	"math"
)

type Shape interface{
	Area() float64
}

type Person struct{
	name string
	age int
}

func (p Person) Info(){
	fmt.Println(p.name, p.age)
}

func (p *Person) Birthday(){
	p.age++
}

type Circle struct{
	radius float64
}

func (c Circle) Area() float64{
	return math.Pi * c.radius * c.radius
}

type Rectangle struct{
	height, width float64
}

func (s Rectangle) Area() float64{
	return s.height * s.width
}

func printAreas(shapes[] Shape){
	for _, shape := range shapes{
		fmt.Printf("%.2f\n", shape.Area())
	}
}

type Book struct{
	title string
	author string
	year int
}

func (b Book) String() string {
	return fmt.Sprintf("«%s», автор: %s (%d)", b.title, b.author, b.year)
}

func main(){

	var p1 Person = Person{"Kirill", 21}
	p1.Info()
	p1.Birthday()
	p1.Info()

	var c Circle = Circle{5}
	fmt.Printf("%.2f\n", c.Area())

	var r Rectangle = Rectangle{10, 5}

	var shapes []Shape = []Shape{c, r}

	printAreas(shapes)

	var book Book = Book{
		title: "Путь к успеху",
		author: "Николай Соболев",
		year: 2018,
	}

	fmt.Println(book)
}