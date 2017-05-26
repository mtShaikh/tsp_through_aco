/*
	Ant Colony Optimization for the Traveling Salesman Problem
	Author: Taha Shaikh
*/
package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type ant_t struct {
	tabulist    []int
	currentCity int
	nextCity    int
	tour        []int
	tourIndex   int
	tourlength  float64
}

type city struct {
	x, y float64
}

var ants []ant_t
var cities []city
var tauMatrix [][]float64
var adjMatrix [][]float64
var besttour []int
var besttourlength = 0.0
var avgtour []int
var currentIndex int
var numCities = 0
var rho = 0.6
var qval = 1.0
var alpha = 0.8
var beta = 0.8
var numAnts = 10

//prints error and exits on abnormal conditions
func printError(err error) {
	fmt.Print(err)
	os.Exit(2)
}

func main() {
	initGraph("read")
	iterations := 0
	initTrail()
	for iterations < 500 {
		initAnts()
		moveAnts()
		intensifyTrail()
		updateBestTour()
		iterations++
	}
	fmt.Println(besttour)
}

func updateBestTour() {
	if besttour == nil {
		besttour = make([]int, numCities)
		copy(besttour, ants[0].tour)
		besttourlength = ants[0].tourlength
	}
	for i := range ants {
		if ants[i].tourlength < besttourlength {
			besttourlength = ants[i].tourlength
			copy(besttour, ants[i].tour)
		}
	}
}

//initialize tauMatrix
func initTrail() {
	tauMatrix = make([][]float64, numCities)
	for i := range tauMatrix {
		tauMatrix[i] = make([]float64, numCities)
		for j := range tauMatrix[i] {
			if i != j {
				tauMatrix[i][j] = 1.0
			} else {
				tauMatrix[i][j] = 0.0
			}
		}
	}
}

//initialize ants
func initAnts() {
	currentIndex = -1
	ants = make([]ant_t, numAnts)
	for i := range ants {
		ants[i].tabulist = make([]int, numCities)
		ants[i].currentCity = rand.Intn(numCities)
		ants[i].nextCity = 0
		ants[i].tour = make([]int, numCities)
		ants[i].tourIndex = 0
		ants[i].tourlength = 0.0
	}
	currentIndex++
}

func moveAnts() {
	for currentIndex < numCities-1 {
		for i := range ants {
			goToNewCity(&ants[i])
		}
		currentIndex++
	}
}

//choosing next city
func goToNewCity(ant *ant_t) {
	var from, to int
	var p float64
	denom := 0.0
	from = ant.currentCity
	for to = 0; to < numCities; to++ {
		if from != to {
			if ant.tabulist[to] == 0 && tauMatrix[from][to] != 0 && adjMatrix[from][to] != 0 {
				denom += math.Pow(tauMatrix[from][to], alpha) * math.Pow((1.0/adjMatrix[from][to]), beta)
			}
		} else {
			continue
		}

	}
	to = 0
	for {
		if from != to {
			if ant.tabulist[to] == 0 {
				p = (math.Pow(tauMatrix[from][to], alpha) * math.Pow((1.0/adjMatrix[from][to]), beta)) / denom

				if rand.Float64() < p {
					break
				}
			}
		} else {
			to = ((to + 1) % numCities)
			continue
		}
		to = ((to + 1) % numCities)
	}
	ant.nextCity = to
	ant.tabulist[ant.nextCity] = 1
	ant.tour[ant.tourIndex] = ant.nextCity
	ant.tourIndex++
	ant.tourlength += adjMatrix[ant.currentCity][ant.nextCity]
	if ant.tourIndex == numCities {
		ant.tourlength += adjMatrix[ant.tour[numCities-1]][ant.tour[0]]
	}
	ant.currentCity = ant.nextCity
}

//reads from file and creates a list of all the cities coordinates
func readFile(name string) []city {
	var dim, i int
	var cities []city
	i, dim = 1, 0
	var startFlag bool
	startFlag = false
	if file, err := os.Open(name); err == nil {
		// make sure it gets closed
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			str := scanner.Text()
			if strings.Contains(str, "DIMENSION") {
				dim = getDim(str)
				break
			}
		}
		cities = make([]city, dim)
		for scanner.Scan() {
			str := scanner.Text()
			if strings.Contains(str, "EOF") {
				break
			} else if startFlag {
				x, y := tokenize(str)
				if i <= dim {
					cities[i-1] = city{x, y}
					i++
				} else {
					startFlag = false
				}
			} else if strings.Contains(str, "NODE_COORD_SECTION") {
				startFlag = true
			}
		}
		// check for errors
		if err = scanner.Err(); err != nil {
			printError(err)
		}
	} else {
		printError(err)
	}
	return cities
}

//intensifying pheromone levels
func intensifyTrail() {
	var from, to, i, c int
	for i = 0; i < numAnts; i++ {
		for c = 0; c < numCities; c++ {
			from = ants[i].tour[c]
			to = ants[i].tour[((c + 1) % numCities)]
			tauMatrix[from][to] += ((qval / ants[i].tourlength) * rho)
			tauMatrix[to][from] = tauMatrix[from][to]
		}
	}
}

//making graph
func initGraph(name string) {
	cities = readFile(name)
	numCities = len(cities)
	adjMatrix = make([][]float64, numCities)
	for i := range adjMatrix {
		adjMatrix[i] = make([]float64, numCities)
		for j := range adjMatrix[i] {
			adjMatrix[i][j] = calEdge(cities[i], cities[j])
		}
	}
}

//calculates edge weight (euclidiean distance)
func calEdge(c1, c2 city) float64 {
	return math.Pow((math.Pow((c2.y-c1.y), 2) + math.Pow((c2.y-c1.y), 2)), 0.5)
}

//tokenizes and converts to float
func tokenize(str string) (x, y float64) {
	s := strings.Split(str, " ")
	strX, strY := s[1], s[2]
	x, err := strconv.ParseFloat(strX, 64) //converts string to float64
	if err != nil {
		printError(err)
	}
	y, err = strconv.ParseFloat(strY, 64) //converts string to float64
	if err != nil {
		printError(err)
	}
	return x, y
}

//gets number of cities from the file
func getDim(str string) (dim int) {
	s := strings.Split(str, ":")
	num := strings.TrimLeft(s[1], " ")
	if dim, err := strconv.Atoi(num); err == nil {
		return dim
	} else {
		fmt.Print(err)
		os.Exit(2)
	}
	return 0
}
