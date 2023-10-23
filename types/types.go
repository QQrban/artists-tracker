package types

import (
	"slices"
	"strconv"
	"strings"
)

// This will hold integrated data about the artist
type ArtistFull struct {
	Id           int
	Image        string
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
	ConcertDates []string
	Locations    map[string][]float64
	Relations    map[string][]string
}

func (a *ArtistFull) Get(field string) any {
	switch field {
	case "Id":
		return a.Id
	case "Image":
		return a.Image
	case "Name":
		return a.Name
	case "Members":
		return a.Members
	case "CountMembers":
		return len(a.Members)
	case "CreationDate":
		return a.CreationDate
	case "FirstAlbum":
		year := a.FirstAlbum[6:]
		y, _ := strconv.Atoi(year)
		return y
	case "ConcertDates":
		return a.ConcertDates
	case "Locations":
		return a.Locations
	case "Relations":
		return a.Relations
	case "Countries":
		var countries = make([]string, len(a.Relations))
		for location := range a.Relations {
			i := strings.LastIndex(location, ", ")
			country := location[i+2:]
			countries = append(countries, country)
		}
		slices.Sort(countries)
		return countries
	case "Cities":
		var cities = make([]string, len(a.Relations))
		for location := range a.Relations {
			cities = append(cities, location)
		}
		slices.Sort(cities)
		return cities
	case "LocationDates":
		var dates = make([]string, 0)
		for _, locationDates := range a.Relations {
			for _, date := range locationDates {
				dates = append(dates, date[6:]+"-"+date[3:5]+"-"+date[:2])
			}
		}
		slices.Sort(dates)
		return dates
	}
	return 0
}

// Following types will hold the data received from the API
type Artist struct {
	Id              int      `json:"id"`
	Image           string   `json:"image"`
	Name            string   `json:"name"`
	Members         []string `json:"members"`
	CreationDate    int      `json:"creationDate"`
	FirstAlbum      string   `json:"firstAlbum"`
	ConcertDatesUrl string   `json:"concertDates"`
	LocationsUrl    string   `json:"locations"`
	RelationsUrl    string   `json:"relations"`
}

type Dates struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Locations struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Relations struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type Field struct {
	Name       string
	Label      string
	Type       string
	Minmax     bool
	Checkbox   bool
	Options    bool
	Attributes map[string]string
	Boxes      []string
	Values     map[string]bool
	MMValues   map[string]string
}
