package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"testing"
)

// NOTES:
// 1. Not sure how to test runSimulation since it is non-deterministic

// TestPopulateCities verifies that number of cities and neighbors
// matches the input file
func TestPopulateCities(t *testing.T) {
	var testState = new(Invasion)
	testState.initInvasion()

	testState.populateCities()

	data, err := os.Open("cities")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanLines)

	cityCount := 0
	neighborCount := 0

	for scanner.Scan() {
		cityCount++

		currentCity := strings.Split(scanner.Text(), " ")

		for range currentCity[1:] {
			neighborCount++
		}
	}

	createdNeighborCount := 0

	for _, neighbors := range testState.cities {
		for range neighbors {
			createdNeighborCount++
		}
	}

	if len(testState.cities) != cityCount {
		t.Errorf("Number of created cities %v does not match number in cities file %v", len(testState.cities), cityCount)
	} else if createdNeighborCount != neighborCount {
		t.Errorf("Number of created neighbors %v does not match number in cities file %v", createdNeighborCount, neighborCount)
	}
}

// Test creating aliens and placing them in cities
// Assumption: Max number of aliens can only be # of cities * 2
func TestPopulateAliens(t *testing.T) {
	uniqueCities := []string{"NewYork", "Boston", "Miami",
		"Portland", "Stamford", "Houston", "TwinPeaks"}
	for i := 0; i < len(uniqueCities)*2; i++ {
		curAliens := i
		// Parallelize
		go func() {
			var testState = new(Invasion)
			testState.initInvasion()

			testState.cities["NewYork"] = []string{"north=Boston", "south=Miami", "east=Stamford"}
			testState.cities["Boston"] = []string{"south=NewYork", "north=Portland"}
			testState.cities["Miami"] = []string{"north=NewYork", "west=Houston"}
			testState.cities["Portland"] = []string{"south=Boston"}
			testState.cities["Stamford"] = []string{"west=NewYork"}
			testState.cities["Houston"] = []string{"east=Miami", "north=TwinPeaks"}
			testState.cities["TwinPeaks"] = []string{"south=Houston"}

			testState.uniqueCities = uniqueCities

			testState.populateAliens(curAliens)

			totalAliens := make([]int, 0)
			for _, v := range testState.aliens {
				for alien := range v {
					totalAliens = append(totalAliens, alien)
				}
			}

			if len(totalAliens) != curAliens {
				t.Errorf("Created aliens %v does not match input %v", len(totalAliens), curAliens)
			}
		}()
	}
}

// TestDestroyCity puts 2 aliens in every city, then goes through cities
// and destroys each
// Verifies cities map and alien locations map are empty as a result
func TestDestroyCity(t *testing.T) {
	var testState = new(Invasion)
	testState.initInvasion()

	testState.populateCities()
	alienID := 0
	for k := range testState.cities {
		testState.aliens[k] = append(testState.aliens[k], alienID)
		alienID++
		testState.aliens[k] = append(testState.aliens[k], alienID)
		alienID++
	}

	for city := range testState.aliens {
		testState.destroyCity(city)
	}

	if len(testState.aliens) != 0 && len(testState.cities) != 0 {
		t.Errorf("Not all cities destroyed!\nAlien occupied cities: %v\nCities: %v", testState.aliens, testState.cities)
	}
}

func TestWriteMapState(t *testing.T) {
	var testState = new(Invasion)
	testState.initInvasion()

	testState.populateCities()

	data, err := os.Open("cities")
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		//
	}
}
