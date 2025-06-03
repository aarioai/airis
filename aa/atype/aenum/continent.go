package aenum

type Continent uint8

const (
	NilContinent Continent = 0 // no, or invalid continent
	Asia         Continent = 1
	Europe       Continent = 2
	NorthAmerica Continent = 3
	SouthAmerica Continent = 4
	Oceania      Continent = 5
	Africa       Continent = 6
	Antarctica   Continent = 7
)
