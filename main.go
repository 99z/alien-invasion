package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Invasion represents an active simulation and its state
type Invasion struct {
	cities       map[string][]string
	aliens       map[string][]int
	uniqueCities []string
}

// Creates new Invasion, populates cities from a file, populates aliens
// from program arg, runs simulation, writes map state to a file
func main() {
	state := new(Invasion)
	state.initInvasion()

	state.populateCities()

	// Get number of aliens from program arg
	numAliens, err := strconv.Atoi(os.Args[1:2][0])
	if err != nil {
		log.Fatal(err)
	}

	state.populateAliens(numAliens)

	state.runSimulation()

	state.writeMapState()
}

// initInvasion initializes Invasion attributes to empty data structures
func (invasion *Invasion) initInvasion() {
	invasion.cities = make(map[string][]string)
	invasion.aliens = make(map[string][]int)
	invasion.uniqueCities = make([]string, 0)
}

// writeMapState writes the end state of the invasion simulation to a file
// in the same format as the input file
// Assumptions: can't specify result filename
func (invasion *Invasion) writeMapState() {
	f, err := os.Create("result")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	for city := range invasion.cities {
		var buffer bytes.Buffer
		buffer.WriteString(city + " ")

		for _, c := range invasion.cities[city] {
			buffer.WriteString(c + " ")
		}

		buffer.WriteString("\n")

		_, err := w.WriteString(buffer.String())
		if err != nil {
			log.Fatal(err)
		}
	}

	w.Flush()
}

// populateCities reads data from a file called "cities"
// and populates the cities map
// Assumptions: file is NOT specified by program argument - spec only mentions
// specifying # of aliens
func (invasion *Invasion) populateCities() {
	data, err := os.Open("cities")
	if err != nil {
		log.Fatal(err)
	}

	defer data.Close()

	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanLines)

	// Iterate over lines of cities file
	for scanner.Scan() {
		// Split city line by spaces
		currentCity := strings.Split(scanner.Text(), " ")

		// If this city does not already exist
		if len(invasion.cities[currentCity[0]]) == 0 {
			// Create new city with name of current line
			invasion.cities[currentCity[0]] = make([]string, 0)

			// Add to uniqueCities for use in alien assignments
			invasion.uniqueCities = append(invasion.uniqueCities, currentCity[0])
		}

		// Add neighbors of current city
		for _, neighbor := range currentCity[1:] {
			// Skip neighbor if it is the same as current city
			if strings.Contains(neighbor, currentCity[0]) {
				continue
			}

			invasion.cities[currentCity[0]] = append(invasion.cities[currentCity[0]], neighbor)
		}
	}

	// Throw error and exit if no cities were created
	if len(invasion.uniqueCities) == 0 || len(invasion.cities) == 0 {
		log.Fatalf("populateAliens: must populate cities first")
	}
}

// populateAliens creates numAliens and randomly places them in a city
// Assumption: no more than 2 aliens may begin in the same city
func (invasion *Invasion) populateAliens(numAliens int) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < numAliens; i++ {
		// Pick random city
		city := invasion.uniqueCities[rand.Intn(len(invasion.uniqueCities))]

		// Ensure no cities have more than 2 aliens
		for len(invasion.aliens[city]) == 2 {
			city = invasion.uniqueCities[rand.Intn(len(invasion.uniqueCities))]
		}

		// Append alien to city
		invasion.aliens[city] = append(invasion.aliens[city], i)
	}
}

// destroyCity deletes a city, specified by a string, from the cities map
// and alien locations map
// When destroying a city, visits its own neighbors and removes itself from
// their list of neighbors
func (invasion *Invasion) destroyCity(city string) {
	// Regex to filter out direction prefix
	reg, err := regexp.Compile(`^(.*?)=`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v has been destroyed by alien %v and alien %v!\n", city, invasion.aliens[city][0], invasion.aliens[city][1])

	// Go to neighbors of deleted and delete itself from their lists
	for _, neighbor := range invasion.cities[city] {
		neighbor = reg.ReplaceAllString(neighbor, "")
		for i, n := range invasion.cities[neighbor] {
			// Neighbors contain the direction prefix, so check that city
			// is a substring of the neighbor instead of direct comparison
			if strings.Contains(n, city) {
				invasion.cities[neighbor][i] = invasion.cities[neighbor][len(invasion.cities[neighbor])-1]
				invasion.cities[neighbor] = invasion.cities[neighbor][:len(invasion.cities[neighbor])-1]
			}
		}
	}

	delete(invasion.aliens, city)
	delete(invasion.cities, city)
}

// runSimulation runs the invasion simulation until one of two conditions is
// met: all aliens have died, or 10k turns have passed
func (invasion *Invasion) runSimulation() {
	turns := 0
	for len(invasion.aliens) != 0 && turns < 10000 {
		// First pass: destroy cities with 2 aliens in them
		// This prevents case of having 3 aliens in a city
		for city := range invasion.aliens {
			if len(invasion.aliens[city]) >= 2 {
				invasion.destroyCity(city)
			}
		}

		// Second pass: handle moving aliens
		for city := range invasion.aliens {
			// If this city has no neighbors to move to this alien
			// does nothing
			if len(invasion.cities[city]) == 0 {
				continue
			} else if len(invasion.aliens[city]) == 1 {
				rand.Seed(time.Now().Unix())

				// Regex to filter out direction prefix so we can properly move
				// This is because format of cities map is a non-prefixed
				// cites mapped to an array of direction-prefixed cities
				reg, err := regexp.Compile(`^(.*?)=`)
				if err != nil {
					log.Fatal(err)
				}

				// Get random neighbor of current city for alien to move to
				newCity := invasion.cities[city][rand.Intn(len(invasion.cities[city]))]
				filtered := reg.ReplaceAllString(newCity, "")

				// Move this alien
				invasion.aliens[filtered] = append(invasion.aliens[filtered], invasion.aliens[city][0])
				invasion.aliens[city] = invasion.aliens[city][1:]

				// Remove city from aliens if it was the only one there
				if len(invasion.aliens[city]) == 0 {
					delete(invasion.aliens, city)
				}
			}
		}

		turns++
	}
}
