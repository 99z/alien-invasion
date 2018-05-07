# alien-invasion
Alien Invasion project simulation

## Running
The program assumes a file in the same directory exists named `cities`. See the file in this repo for an example.

`go run main.go numAliens`. It is assumed `numAliens` cannot be > 2x the number of cities. Crashes if this is the case.

The result of the map after the simulation has finished is written to a file named `result` in the same directory.

## Tests
The test coverage is 59.6%. The main game loop is not tested because it is non-deterministic - I wasn't sure the best way to handle this. All other functions are covered.
