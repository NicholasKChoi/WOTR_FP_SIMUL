package main

import (
	"fmt"

	"github.com/op/go-logging"
)

var (
	NumSimulationToRun = 1000
	GameLogLevel       = logging.WARNING
)

func main() {
	fmt.Println("Hello world")
	simulation := &Simulator{}
	simulation.Init()
	for i := 0; i < NumSimulationToRun; i++ {
		if err := simulation.RunSimulation(); err != nil {
			panic(err)
		}
	}
	simulation.PrintResult()
	fmt.Println("Analyzer complete")
}
