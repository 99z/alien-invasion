package main

import "testing"

func testPopulateCities(t *testing.T) {
	var testState = new(Invasion)
	testState.initInvasion()

	// Test populateCities
}

func TestPopulateAliens(t *testing.T) {
	for i := 0; i < 10000; i++ {
		// Parallelize
		go func() {
			var testState = new(Invasion)
			testState.initInvasion()

			testState.populateCities()
			testState.populateAliens(i)

			totalAliens := make([]int, 0)
			for _, v := range testState.aliens {
				for alien := range v {
					totalAliens = append(totalAliens, alien)
				}
			}

			if len(totalAliens) != i {
				t.Errorf("Created aliens %v does not match input %v", len(totalAliens), i)
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
