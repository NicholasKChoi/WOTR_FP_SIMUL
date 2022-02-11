package main

import "fmt"

func main() {
	fmt.Println("Hello world")

	var (
		charR CharDiceNeededTable
		corrR CorruptionInflictedTable
		reavR RevealsTable
	)

	fmt.Printf("%#v\n", charR)
	fmt.Printf("%#v\n", corrR)
	fmt.Printf("%#v\n", reavR)
}
