/*
*	Author: Taha Shaikh
*	=== Ant Colony Optimization for the Traveling Salesman Problem ===
*
*	=== Usage ===
*	Comment out whichever file/data not needed and uncomment file/data at
*	line numbers 60 and 90 to execute ACO on
*
*	=== Implementation ===
*	- Reads the DIMENSION field in the file provided for the total number of cities
*	- Reads the coordinates of each city in the NODE_COORD_SECTION and saves them in
*	  a list of cities
*	- Creates an adjacency matrix for the cities, computing edge weights by the
*	  calculating the euclidean distance between every two cities
*	- Initialize the tau matrix (pheromone levels) for each edge
*	- initialize each ant and provide them a random city to start
*	- Each ant traveses the path and chooses each city using a probability which
*	  is computed using a formula
*	- Update the tau matrix (pheromone levels) after one tour has end
*	- Run 500 tours 10 times to obtain average and best tour length values
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

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
)

//ant implementation
type ant_t struct {
	tabulist    []int
	currentCity int
	nextCity    int
	tour        []int
	tourIndex   int
	tourlength  float64
}

//city implementation
type city struct {
	x, y float64
}

//declaration and initialization
var ants []ant_t
var cities []city
var tauMatrix [][]float64
var adjMatrix [][]float64
var besttour []int
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
	//comment other to use the other one
	//initGraph("TSP_D")
	initGraph("TSP_WS") //other TSP file

	avgavg := make([]float64, 10)
	avgbest := make([]float64, 10)
	besttourlens := make([]float64, 10)
	avgtourlens := make([]float64, 10)
	iterations := 0
	i := 0
	//initializes pheromone levels with base pheromone i.e. 1
	initTrail()
	for i < 10 {
		iterations = 0
		besttour = nil
		for iterations < 500 {
			//initializes each ant
			initAnts()
			//tour for each ant
			moveAnts()
			//intensify pheromone levels
			intensifyTrail()
			//compute best and avg tour length for every tour in
			//500 and add them for average
			besttourlens[i] += calculateBest()
			avgtourlens[i] += calculateAvg()
			iterations++
		}
		//evaporate pheromone to obtain better results
		evaporatePheromone()
		fmt.Println("Iteration: ", i)
		fmt.Println("Optimal Path: ", besttour)
		avgavg[i] = avgtourlens[i] / 500.0
		avgbest[i] = besttourlens[i] / 500.0
		i++
	}
	fmt.Println("Average Average:", avgavg)
	fmt.Println("Average Best:", avgbest)

	p, err := plot.New()
	if err != nil {
		printError(err)
	}

	//comment other to use the other one
	//p.Title.Text = "Dijibouti TSP"
	p.Title.Text = "Western Sahara TSP"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	avgpts := make(plotter.XYs, 10)
	for i := range avgpts {
		avgpts[i].Y = avgavg[i]
		avgpts[i].X = float64(i)
	}

	bestpts := make(plotter.XYs, 10)
	for i := range bestpts {
		bestpts[i].Y = avgbest[i]
		bestpts[i].X = float64(i)
	}
	err = plotutil.AddLinePoints(p,
		"Average So Far", avgpts,
		"Best So Far", bestpts)
	if err != nil {
		printError(err)
	}

	// Save the plot to a PNG file.
	/*if err := p.Save(4*vg.Inch, 4*vg.Inch, "djibouti.png"); err != nil {
		printError(err)
	}*/
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "westernsahara.png"); err != nil {
		printError(err)
	}
}

//calculate average tour length of all length for one tour
func calculateAvg() float64 {
	avglength := ants[0].tourlength
	total := 0.0
	for i := range ants {
		total += ants[i].tourlength
	}
	avglength = total / 10.0
	return avglength
}

//calculate best tour length of all length for one tour
func calculateBest() float64 {
	var bestlength float64
	bestlength = ants[0].tourlength
	if besttour == nil {
		besttour = make([]int, numCities)
		copy(besttour, ants[0].tour)
	}
	for i := range ants {
		if ants[i].tourlength < bestlength {

			bestlength = ants[i].tourlength
			copy(besttour, ants[i].tour)
		}
	}
	return bestlength
}

//initialize pheromone levels
func initTrail() {
	tauMatrix = make([][]float64, numCities)
	for i := range tauMatrix {
		tauMatrix[i] = make([]float64, numCities)
		for j := range tauMatrix[i] {
			if i != j {
				tauMatrix[i][j] = 1.0 //initialize to base pheromone = 1
			} else {
				tauMatrix[i][j] = 0.0
			}
		}
	}
}

//initialize ants
func initAnts() {
	ants = nil
	ants = make([]ant_t, numAnts)
	for i := range ants {
		ants[i].tabulist = make([]int, numCities)
		ants[i].currentCity = rand.Intn(numCities) //randomly assigns ant a city to s
		ants[i].nextCity = 0
		ants[i].tour = make([]int, numCities)
		ants[i].tourIndex = 0
		ants[i].tourlength = 0.0
	}
}

//move all ants to visit the whole graph
func moveAnts() {
	for i := range ants {
		currentIndex = 0
		for currentIndex < numCities {
			goToNewCity(&ants[i])
			currentIndex++
		}
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

//evaporating pheromone after each iteration of the algorithm
func evaporatePheromone() {
	var from, to int
	for from = 0; from < numCities; from++ {
		for to = 0; to < numCities; to++ {
			tauMatrix[from][to] = tauMatrix[from][to] * (1.0 - rho)
			if tauMatrix[from][to] < 0.0 {
				tauMatrix[from][to] = 1.0
			}
		}
	}
}

//intensifying pheromone levels
func intensifyTrail() {
	var from, to, i, c int
	for i = 0; i < numAnts; i++ {
		for c = 0; c < numCities; c++ {
			from = ants[i].tour[c]
			to = ants[i].tour[((c + 1) % numCities)]
			deltatau := (qval / ants[i].tourlength)
			tauMatrix[from][to] = tauMatrix[from][to] + deltatau
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
