package main

import "testing"

func testPopulateCities(t *testing.T) {
	var testState = new(Invasion)
	testState.initInvasion()

	// Test populateCities
}

func TestPopulateAliens(t *testing.T) {
	for i := 0; i < 10000; i++ {
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

			testState.uniqueCities = []string{"NewYork", "Boston", "Miami",
				"Portland", "Stamford", "Houston", "TwinPeaks"}

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

func TestRunSimulation(t *testing.T) {
	var testState = new(Invasion)
	testState.initInvasion()

	testState.populateCities()
}
