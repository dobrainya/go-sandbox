package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

func main() {
	fmt.Println(add(42, 13))

	a, b, c, d := true, 12.2, "hello", 1

	fmt.Println(a, b, c, d)

	var i interface{} = 123.45

	val, ok := i.(float64)

	fmt.Println(val, ok)

	val1, ok1 := i.(string)

	fmt.Println(val1, ok1)
	//mathFn()

	//cycles()
	//whileCytcle()
	//arrays()
	//slices()
	//maps()
	switchFn()

	//fmt.Println(returnMathFn("/")(5, 2))
	//fmt.Println(returnMathFn("+")(5, 2))

	//pointers()
}

func pointers() {
	a := 12

	pointer := &a

	fmt.Println(pointer, *pointer)

	changer := func(str string, value string) {
		str = value
	}

	str := "Hello"

	changer(str, "World")

	fmt.Println(str)

	changerV2 := func(str *string, value string) {
		*str = value
	}

	changerV2(&str, "World")

	fmt.Println(str)

	var pointer2 *string = &str

	changerV2(pointer2, "People")

	fmt.Println(str)
}

func add(x int, y int) int {
	return x + y
}

func returnMathFn(action string) func(x1, x2 float64) (result float64) {
	return func(x1, x2 float64) (result float64) {
		switch action {
		case "/":
			result = x1 / x2
		case "*":
			result = x1 * x2
		case "+":
			result = x1 + x2
		case "-":
			result = x1 - x2
		}

		return
	}
}

func mathFn() {
	a := 5
	b := 2

	fmt.Println(a / b)

	var c float64 = 5
	var d float64 = 2

	fmt.Println(c / d)

	f := 2.501

	fmt.Println(math.Ceil(f))
	fmt.Println(math.Floor(f))
	fmt.Println(math.Round(f))
}

func cycles() {
	for i := 0; i < 10; i++ {
		fmt.Println(i)
	}
}

func whileCytcle() {
	var i int
	rand.Seed(time.Now().UnixNano())

	for i != 5 {
		fmt.Println(i)

		i = rand.Intn(100)
	}
}

func arrays() {
	var arr1 [3]int

	fmt.Println(arr1)

	arr := []int{1, 2, 3, 4, 5}

	for i := 0; i < len(arr); i++ {
		fmt.Println(arr[i])
	}

	for idx, elem := range arr {
		fmt.Printf("Index: %d, Element: %d\n", idx, elem)
	}

	matrix := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}

	for x1, arr := range matrix {
		for x2, value := range arr {
			fmt.Printf("x1: %d, x2: %d, Value: %d\n", x1, x2, value)
		}
	}

}

func slices() {
	slice0 := []int{5, 6, 2, 1}

	fmt.Println(slice0)

	slice0 = append(slice0, 0)
	fmt.Println(slice0)

	sort.Ints(slice0)
	fmt.Println(slice0)

	sort.Ints(slice0)
	fmt.Println(slice0)

}

func maps() {
	var regCodes map[int]string = map[int]string{
		777: "Msc",
		116: "Tat",
	}

	fmt.Println(regCodes)

	fmt.Println(regCodes[116])

	regCodes0 := map[int]string{
		777: "Msc",
		116: "Tat",
		100: "Undefined",
	}

	fmt.Println(regCodes0)

	delete(regCodes0, 100)
	fmt.Println(regCodes0)
	fmt.Println(regCodes[100])

	el, status := regCodes[100]

	fmt.Println(el, status)

	el1, status1 := regCodes[116]

	fmt.Println(el1, status1)
}

func switchFn() {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(5)

	switch i {
	case 1:
		fmt.Println("1")
	case 2:
		fmt.Println("2")
	case 3:
		fmt.Println("3")
	case 4:
		fmt.Println("4")
	default:
		fmt.Println("0")
	}

	a := 10

	switch {
	case a > 1:
		fmt.Println("var a greater than 1")
	case a > 5:
		fmt.Println("var a greater than 5")
	}

	switch {
	case a > 1:
		fmt.Println("var a greater than 1")
		fallthrough
	case a > 5:
		fmt.Println("var a greater than 5")
	}

	//var ii interface{} = 123.45
	var ii interface{} = "text"

	switch v := ii.(type) {
	case float64:
		fmt.Printf("value is %T\n", v)
	case int:
		fmt.Printf("value is %T\n", v)
	case string:
		fmt.Printf("value is %T\n", v)
	default:
		fmt.Println("the type of value is unknown")
	}
}
