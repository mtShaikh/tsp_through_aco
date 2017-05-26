/*
	Ant Colony Optimization for the Traveling Salesman Problem
	Author: Taha Shaikh
*/
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"math/rand"
	"strings"
	//"text/scanner"
	//"github.com/twmb/algoimpl/go/graph"
)

type ant_t struct {
	tabulist    []int
	currentCity int
	nextCity    int
	tour        []int
	tourIndex   int
	tourlength  int
}

type city struct {
	x, y float64
}

var adjMatrix [][]float64
var tauMatrix [][]float64
var cities []city
var ants []ant_t
var rho = 0.6
var qval = 1
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
	numCities := len(cities)
	for i := 0; i < numAnts; i++ {
		
	}
}

//initialize ants
func initAnts() {
	
}

//choosing next city
func goToNewCity(ant *ant_t) {
	numCities := len(cities)
	var from, to int
	var p float64
	d := 0.0
	from = ant.currentCity
	for to = 0; to < numCities; to++ {
		if ant.tabulist[to] == 0 {
			d += math.Pow(tauMatrix[from][to], alpha) * math.Pow((1.0/adjMatrix[from][to]), beta)
		}
	}
	to = 0
	for {
		if ant.tabulist[to] == 0 {
			p = (math.Pow(tauMatrix[from][to], alpha) * math.Pow((1.0/adjMatrix[from][to]), beta)) / d
			if (rand.Float64() < p){
				break
			}
			to = ((to + 1) % numCities)
		}
		ant.nextCity = to
		ant.tabulist[ant.nextCity] = 1
		ant.tour[ant.tourIndex++] = ant.nextCity
		ant.tourlength += adjMatrix[ant.currentCity][ant.nextCity]
		if ant.tourIndex == numCities {
			ant.tourlength += adjMatrix[ant.tour[numCities-1]][ant.tour[0]]
		}
		ant.currentCity = ant.nextCity
	}

}

//reads from file and creates a list of all the cities coordinates
func readFile(name string) {
	var dim, i int
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
					//fmt.Println(cities[i-1])
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
}

//intensifying pheromone levels
func intensifyTrail() {
	var from, to, i, c int
	numCities := len(cities)
	for i = 0; i < numAnts; i++ {
		for c = 0; c < numCities; c++ {
			from = ants[i].tour[city]
			to = ants[i].tour[((city + 1) % numCities)]
			tauMatrix[from][to] += ((qval / ants[i].tour_length) * rho)
			tauMatrix[to][from] = tauMatrix[from][to]
		}
	}
}

//making graph
func initGraph(name string) {
	readFile(name)
	adjMatrix = make([][]float64, len(cities))
	for i := range adjMatrix {
		adjMatrix[i] = make([]float64, len(cities))
		for j := range adjMatrix[i] {
			adjMatrix[i] = make([]float64, len(cities))
			adjMatrix[i][j] = calEdge(cities[i], cities[j])
			fmt.Println(adjMatrix[i][j])
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
