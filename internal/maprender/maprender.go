package maprender

import (
	"fmt"
	"strings"

	"github.com/vsimakhin/web-logbook/internal/models"
)

type Line struct {
	Point1 []float64 `json:"point1"`
	Point2 []float64 `json:"point2"`
}

type Marker struct {
	Point     []float64 `json:"point"`
	Name      string    `json:"name"`
	CivilName string    `json:"civil_name"`
	City      string    `json:"city"`
	Country   string    `json:"country"`
	Elevation string    `json:"elevation"`
}

type MapRender struct {
	FlightRecords []models.FlightRecord
	AirportsDB    map[string]models.Airport

	FilterNoRoutes bool

	Lines   []Line
	Markers []Marker
}

// Render renders map with airport and routes markers
func (mr *MapRender) Render() {
	airportMarkers := make(map[string]struct{})
	routeLines := make(map[string]struct{})

	// parsing
	for _, fr := range mr.FlightRecords {
		// add to the list of the airport markers departure and arrival
		// it will be automatically a list of unique airports
		airportMarkers[fr.Departure.Place] = struct{}{}
		airportMarkers[fr.Arrival.Place] = struct{}{}

		// the same for the route lines
		if !mr.FilterNoRoutes {
			if fr.Departure.Place != fr.Arrival.Place {
				routeLines[fmt.Sprintf("%s-%s", fr.Departure.Place, fr.Arrival.Place)] = struct{}{}
			}
		}
	}

	// generate routes lines
	for route := range routeLines {
		places := strings.Split(route, "-")

		if airport1, ok := mr.AirportsDB[places[0]]; ok {
			if airport2, ok := mr.AirportsDB[places[1]]; ok {

				mr.Lines = append(mr.Lines, Line{Point1: []float64{airport1.Lon, airport1.Lat}, Point2: []float64{airport2.Lon, airport2.Lat}})
			}
		}
	}

	// generate airports markers
	for place := range airportMarkers {
		if airport, ok := mr.AirportsDB[place]; ok {
			marker := Marker{
				Point:     []float64{airport.Lon, airport.Lat},
				Name:      place,
				CivilName: airport.Name,
				City:      airport.City,
				Country:   airport.Country,
				Elevation: fmt.Sprintf("%d", airport.Elevation),
			}
			mr.Markers = append(mr.Markers, marker)
		}
	}
}
