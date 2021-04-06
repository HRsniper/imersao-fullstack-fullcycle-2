package route

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"
)

// Route represents a request of new delivery request
type Route struct {
	ID        string     `json:"routeId"`
	ClientID  string     `json:"clientId"`
	Positions []Position `json:"position"`
}

// Position is a type which contains the lat and long
type Position struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

// PartialRoutePosition is the actual response which the system will return
type PartialRoutePosition struct {
	ID       string    `json:"routeId"`
	ClientID string    `json:"clientId"`
	Position []float64 `json:"position"`
	Finished bool      `json:"finished"`
}

// NewRoute creates a *Route struct
func NewRoute() *Route {
	return &Route{}
}

// LoadPositions loads from a .txt file all positions (lat and long) to the Position attribute of the struct
func (route *Route) LoadPositions() error {
	if route.ID == "" {
		return errors.New("route id not informed")
	}

	file, err := os.Open("destinations/" + route.ID + ".txt")
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), ",")

		lat, err := strconv.ParseFloat(data[0], 64)
		if err != nil {
			return nil
		}

		long, err := strconv.ParseFloat(data[1], 64)
		if err != nil {
			return nil
		}

		route.Positions = append(route.Positions, Position{
			Lat:  lat,
			Long: long,
		})
	}

	return nil
}

// ExportJsonPositions generates a slice of string in Json using PartialRoutePosition struct
func (route *Route) ExportJsonPositions() ([]string, error) {
	var routePartial PartialRoutePosition
	var result []string

	total := len(route.Positions)

	for key, value := range route.Positions {
		routePartial.ID = route.ID
		routePartial.ClientID = route.ClientID
		routePartial.Position = []float64{value.Lat, value.Long}
		routePartial.Finished = false

		if total-1 == key {
			routePartial.Finished = true
		}

		jsonRoute, err := json.Marshal(routePartial)
		if err != nil {
			return nil, err
		}

		result = append(result, string(jsonRoute))
	}

	return result, nil
}
