package main

import (
	"fmt"
	"strconv"
)

func main() {

	//FOR Wersja 1
	i := 0

	for {
		if i > 20 {
			break
		}
		i++
		fmt.Println(i)
	}
	fmt.Println("koniec petli FOR Wersja 1, i=" + strconv.Itoa(i))

	//FOR Wersja 2
	var j int
	for j = 0; j < 20; j++ {
		fmt.Println(j)
	}
	fmt.Println("koniec petli FOR Wersja 2, j=" + strconv.Itoa(j))

	//FOR Wersja 3
	var k1 int
	for k := 0; k < 20; k++ {
		fmt.Println(k)
		k1 = k
	}
	fmt.Println("koniec petli FOR Wersja 3, k=" + strconv.Itoa(k1))
	fmt.Printf("%d \n", k1)

	//FOR Wersja 4
	suma := 0
	for suma < 1000 {
		suma += 5
	}
	fmt.Printf("%d \n", suma)

	//FOR Wersja 5
	imiona := []string{"Jan", "Adam", "Witek", "Grazyna", "Bolek"}

	for index, value := range imiona {
		fmt.Printf("%d %s \n", index, value)
	}

}
