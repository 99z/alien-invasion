package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var cities = make(map[string][]string)
var aliens = make(map[int]string)
var uniqueCities = make([]string, 0)

func main() {
	populateCities()

	// Get number of aliens from program arg
	numAliens, err := strconv.Atoi(os.Args[1:2][0])
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	runSimulation(numAliens)

	fmt.Println(cities)
	fmt.Println(aliens)
}

// populateCities reads data from a file called "cities"
// and populates the cities map
func populateCities() {
	data, err := os.Open("cities")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer data.Close()

	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanLines)

	// Iterate over lines of cities file
	for scanner.Scan() {
		// Split city line by spaces
		c := strings.Split(scanner.Text(), " ")

		// If this city does not already exist
		if len(cities[c[0]]) == 0 {
			// Create new city with name of current line
			cities[c[0]] = make([]string, 0)

			// Add to uniqueCities for use in alien assignments
			uniqueCities = append(uniqueCities, c[0])
		}

		// Regex to filter out direction and equals sign
		reg, err := regexp.Compile(`^(.*?)=`)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		// Add neighbors of current city
		// TODO Probably good to make sure this isn't the same as current city
		for _, neighbor := range c[1:] {
			filtered := reg.ReplaceAllString(neighbor, "")

			cities[c[0]] = append(cities[c[0]], filtered)
		}
	}
}

func runSimulation(numAliens int) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < numAliens; i++ {
		aliens[i] = uniqueCities[rand.Intn(len(uniqueCities))]
	}
}
