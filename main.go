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

// QUESTIONS
// 1. Is numAliens < numCities? Any restriction on these?
// 2. Is a road always a 2-way connection? Meaning, if NY says it is connected
// to Boston, is does Boston say it is connected to NY?

// CONCERNS
// 1. Does city need to be destroyed as soon as move is made?
// 2. What if 2 aliens begin in same city, but additional alien
// moves to city prior to iterator handling what happens for that city?
// We will have 3 aliens in a city in this case

// Invasion represents an active simulation and its state
type Invasion struct {
	cities       map[string][]string
	aliens       map[string][]int
	uniqueCities []string
}

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

	state.printMapState()
}

func (invasion *Invasion) initInvasion() {
	invasion.cities = make(map[string][]string)
	invasion.aliens = make(map[string][]int)
	invasion.uniqueCities = make([]string, 0)
}

func (invasion *Invasion) printMapState() {
	f, err := os.Create("result")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	for city := range invasion.cities {
		_, err := w.WriteString(fmt.Sprintf("%v %v\n", city, invasion.cities[city]))
		if err != nil {
			log.Fatal(err)
		}
	}

	w.Flush()
}

func (invasion *Invasion) formatNeighbors(city string) {
	// asd
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
		c := strings.Split(scanner.Text(), " ")

		// If this city does not already exist
		if len(invasion.cities[c[0]]) == 0 {
			// Create new city with name of current line
			invasion.cities[c[0]] = make([]string, 0)

			// Add to uniqueCities for use in alien assignments
			invasion.uniqueCities = append(invasion.uniqueCities, c[0])
		}

		// Regex to filter out direction and equals sign
		// reg, err := regexp.Compile(`^(.*?)=`)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// Add neighbors of current city
		// TODO Probably good to make sure this isn't the same as current city
		for _, neighbor := range c[1:] {
			// filtered := reg.ReplaceAllString(neighbor, "")

			invasion.cities[c[0]] = append(invasion.cities[c[0]], neighbor)
		}
	}
}

func (invasion *Invasion) populateAliens(numAliens int) {
	if len(invasion.uniqueCities) == 0 || len(invasion.cities) == 0 {
		log.Fatalf("populateAliens: must populate cities first")
	}

	rand.Seed(time.Now().Unix())

	for i := 0; i < numAliens; i++ {
		city := invasion.uniqueCities[rand.Intn(len(invasion.uniqueCities))]

		// Ensure no cities have more than 2 aliens
		for len(invasion.aliens[city]) == 2 {
			city = invasion.uniqueCities[rand.Intn(len(invasion.uniqueCities))]
		}

		invasion.aliens[city] = append(invasion.aliens[city], i)
	}
}

func (invasion *Invasion) destroyCity(city string) {
	fmt.Printf("%v has been destroyed by alien %v and alien %v!\n", city, invasion.aliens[city][0], invasion.aliens[city][1])

	reg, err := regexp.Compile(`^(.*?)=`)
	if err != nil {
		log.Fatal(err)
	}

	// Go to neighbors of deleted and delete itself from their lists
	for _, neighbor := range invasion.cities[city] {
		neighbor = reg.ReplaceAllString(neighbor, "")
		for i, n := range invasion.cities[neighbor] {
			if strings.Contains(n, city) {
				invasion.cities[neighbor][i] = invasion.cities[neighbor][len(invasion.cities[neighbor])-1]
				invasion.cities[neighbor] = invasion.cities[neighbor][:len(invasion.cities[neighbor])-1]
			}
		}
	}
	delete(invasion.aliens, city)
	delete(invasion.cities, city)
}

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
