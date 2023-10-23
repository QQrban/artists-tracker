package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"groupie-tracker/helpers"
	"groupie-tracker/types"
	"html/template"
	"log"
	"net/http"
	"slices"
	"time"

	//"slices"
	"strconv"
	"strings"
)

var ApiBaseUrl = "https://groupietrackers.herokuapp.com/api"

var Artists []types.ArtistFull

func getData(baseUrl string, endpoint string) (interface{}, error) {
	apiUrl := fmt.Sprintf("%s/%s", baseUrl, endpoint)
	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch endpoint {
	case "artists":
		// decode an array value (Artists)
		var artists []types.Artist
		err = json.NewDecoder(resp.Body).Decode(&artists)
		return artists, err
	case "dates":
		// decode a map value (Dates)
		var dates map[string][]types.Dates
		err = json.NewDecoder(resp.Body).Decode(&dates)
		return dates, err
	case "locations":
		// decode a map value (Locations)
		var locations map[string][]types.Locations
		err = json.NewDecoder(resp.Body).Decode(&locations)
		return locations, err
	case "relation":
		// decode a map value (Relations)
		var relation map[string][]types.Relations
		err = json.NewDecoder(resp.Body).Decode(&relation)
		return relation, err
	default:
		return nil, fmt.Errorf("unknown data type: %s", endpoint)
	}
}

// Get integrated data for all artists
func getArtists(w http.ResponseWriter, r *http.Request) {
	artistData, err := getData(ApiBaseUrl, "artists")
	if err != nil {
		log.Printf("Error getting artists data: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	artists, ok := artistData.([]types.Artist)
	if !ok {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	dateData, err := getData(ApiBaseUrl, "dates")
	if err != nil {
		log.Printf("Error getting artists data: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	dates, ok := dateData.(map[string][]types.Dates)
	if !ok {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	/*
		locationData, err := getData(ApiBaseUrl, "locations")
		if err != nil {
			log.Printf("Error getting artists data: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		locations, ok := locationData.(map[string][]types.Locations)
		if !ok {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	*/
	relationData, err := getData(ApiBaseUrl, "relation")
	if err != nil {
		log.Printf("Error getting artists data: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	relations, ok := relationData.(map[string][]types.Relations)
	if !ok {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create a slice with integrated Artist data
	for _, a := range artists {
		var artist types.ArtistFull
		artist.Id = a.Id
		artist.Image = a.Image
		artist.Name = a.Name
		artist.Members = a.Members
		artist.CreationDate = a.CreationDate
		artist.FirstAlbum = a.FirstAlbum
		// Extract the dates for the artist
		for _, d := range dates["index"] {
			if d.Id == a.Id {
				artist.ConcertDates = d.Dates
			}
		}
		// Extract the locations longitude and latitude for the artist
		artist.Locations = make(map[string][]float64)

		// Extract the relations for the artist
		artist.Relations = make(map[string][]string)
		for _, r := range relations["index"] {
			if r.Id == a.Id {
				// Edit keys to be more readable
				for k, v := range r.DatesLocations {
					k = normalize(k)
					//k += ":"
					artist.Relations[k] = v
				}
			}
		}
		// Add artist data to the slice
		Artists = append(Artists, artist)
	}
}

func normalize(k string) string {
	k = strings.ReplaceAll(k, "-", ", ")
	k = strings.ReplaceAll(k, "_", " ")
	k = strings.Title(k)
	if strings.HasSuffix(k, "Usa") {
		k = k[:len(k)-3] + "USA"
	}
	if strings.HasSuffix(k, "Uk") {
		k = k[:len(k)-2] + "UK"
	}
	return k
}

var Filter0 []types.Field

// Get data for main page (root path)
func ArtistsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.ServeFile(w, r, "templates/404.html")
		return
	}

	// Get data only if it's not already loaded
	if Artists == nil {
		getArtists(w, r)
	}

	var spec = []types.Field{
		helpers.NewField("CreationDate", "creation year", "number", true, false, false),
		helpers.NewField("FirstAlbum", "first album", "number", true, false, false),
		helpers.NewField("LocationDates", "concert dates", "date", true, false, false),
		helpers.NewField("CountMembers", "members", "number", false, true, false),
		helpers.NewField("Countries", "concert countries", "text", false, false, true),
		helpers.NewField("Cities", "concert cities", "text", false, false, true),
	}

	Filter0 = makeFilter(spec)

	data := struct {
		Artists []types.ArtistFull
		Filter  []types.Field
	}{Artists, Filter0}

	// Get and fill in template
	template, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := template.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Get data for individual artist page
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	artistId := r.URL.Path[len("/artist/"):]
	if len(artistId) == 0 {
		http.ServeFile(w, r, "templates/404.html")
		return
	}

	// Get data for all artists only if it's not already loaded, i.e. artist page is visited first
	if Artists == nil {
		getArtists(w, r)
	}

	var artist types.ArtistFull

	for i, a := range Artists {
		id, _ := strconv.Atoi(artistId)
		if a.Id == id {
			artist = a
			break
		}

		// If artist is not found, return 404
		if i == len(Artists)-1 {
			http.ServeFile(w, r, "templates/404.html")
			return
		}
	}
	if len(artist.Locations) == 0 {
		for location := range artist.Relations {
			if _, ok := artist.Locations[location]; !ok {
				lonLat, err := helpers.Geolocate(location)
				if err != nil {
					log.Printf("Error getting location: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				artist.Locations[location] = lonLat
			}
		}
	}

	template, err := template.ParseFiles("templates/base.html", "templates/artist.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := template.Execute(w, artist); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func autocompleteHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("query"))
	if len(query) < 1 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]map[string]string{})
		return
	}
	var suggestions []map[string]string
	for _, artist := range Artists {
		// Check the artist's name
		if strings.Contains(strings.ToLower(artist.Name), query) {
			suggestions = append(suggestions, map[string]string{
				"name": fmt.Sprintf("%s (band)", artist.Name),
				"id":   strconv.Itoa(artist.Id),
			})
		}
		// Check each member's name
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), query) {
				suggestions = append(suggestions, map[string]string{
					"name": fmt.Sprintf("%s (%s) - member", member, artist.Name),
					"id":   strconv.Itoa(artist.Id),
				})
			}
		}
		// Check each location
		for location := range artist.Relations {
			if strings.Contains(strings.ToLower(location), query) {
				suggestions = append(suggestions, map[string]string{
					"name": fmt.Sprintf("%s %s", location, artist.Name),
					"id":   strconv.Itoa(artist.Id),
				})
			}
		}
		// Check the creation date
		if strings.Contains(strings.ToLower(strconv.Itoa(artist.CreationDate)), query) {
			suggestions = append(suggestions, map[string]string{
				"name": fmt.Sprintf("Creation date: %d (%s)", artist.CreationDate, artist.Name),
				"id":   strconv.Itoa(artist.Id),
			})
		}
		// Check the first album
		if strings.Contains(strings.ToLower(artist.FirstAlbum), query) {
			formatted := artist.FirstAlbum[6:len(artist.FirstAlbum)]
			suggestions = append(suggestions, map[string]string{
				"name": fmt.Sprintf("First Album: %s (%s)", formatted, artist.Name),
				"id":   strconv.Itoa(artist.Id),
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suggestions)
}

func filterHandler(w http.ResponseWriter, r *http.Request) {
	filter, err := getFilter(w, r)
	if err != nil {
		return
	}

	registerChoices(filter)

	artists, err := filterArtists(filter, w)
	if err != nil {
		return
	}

	data := struct {
		Artists []types.ArtistFull
		Filter  []types.Field
	}{artists, Filter0}

	template, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := template.Execute(w, data); err != nil {
		report(err, "Error executing template", w)
		return
	}
}

func report(val any, message string, w http.ResponseWriter) {
	log.Printf(message+": %v", val)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func registerChoices(filter map[string]map[string]any) {
	//fmt.Printf("\nfilter: %v\n\n", filter)
	//fmt.Printf("Filter0: %v\n", Filter0)
	for i, spec := range Filter0 {
		if spec.Minmax {
			//fmt.Printf("spec: %v\n\n", spec)
			if min, ok := filter[spec.Name]["min"]; ok {
				if fmt.Sprintf("%v", Filter0[i].Attributes["min"]) != fmt.Sprintf("%v", min) {
					Filter0[i].MMValues["min"] = fmt.Sprintf("%v", min)
				}
			} else {
				delete(Filter0[i].MMValues, "min")
			}
			if max, ok := filter[spec.Name]["max"]; ok {
				if fmt.Sprintf("%v", Filter0[i].Attributes["max"]) != fmt.Sprintf("%v", max) {
					Filter0[i].MMValues["max"] = fmt.Sprintf("%v", max)
				}
			} else {
				delete(Filter0[i].MMValues, "max")
			}
		}
		if spec.Checkbox || spec.Options {
			if values, ok := filter[spec.Name]["values"]; ok {
				//Filter0[i].Values = make(map[string]bool)
				for name := range spec.Values {
					Filter0[i].Values[name] = slices.Contains(values.([]string), name)
				}
			} else {
				for name := range spec.Values {
					Filter0[i].Values[name] = false
				}
			}
		}
	}
}

func makeFilter(specs []types.Field) []types.Field {
	Min := make(map[string]any)
	Max := make(map[string]any)
	Check := make(map[string]map[string]bool)

	// Initialize with first artist
	a := Artists[0]
	for _, spec := range specs {
		spec.Attributes = make(map[string]string)
		val := a.Get(spec.Name)
		switch fmt.Sprintf("%T", val) {
		case "int":
			//spec.Attributes["type"] = "number"
			if spec.Minmax {
				Min[spec.Name] = val.(int)
				Max[spec.Name] = val.(int)
			}
			if spec.Checkbox || spec.Options {
				Check[spec.Name] = make(map[string]bool)
				v := fmt.Sprintf("%d", val)
				Check[spec.Name][v] = false
			}
		case "string":
			//spec.Attributes["type"] = "text"
			v := val.(string)
			if spec.Minmax {
				Min[spec.Name] = v[6:] + "-" + v[3:5] + "-" + v[:2] // Only dates are strings in min-max filters
				Max[spec.Name] = v[6:] + "-" + v[3:5] + "-" + v[:2]
			}
			if spec.Checkbox || spec.Options {
				Check[spec.Name] = make(map[string]bool)
				Check[spec.Name][v] = false
			}
		case "[]string":
			//spec.Attributes["type"] = "text"
			vs := val.([]string)
			if spec.Minmax {
				Min[spec.Name] = vs[0]
				Max[spec.Name] = vs[len(vs)-1]
			}
			if spec.Checkbox || spec.Options {
				Check[spec.Name] = make(map[string]bool)
				for _, v := range vs {
					Check[spec.Name][v] = false
				}
			}
		}
	}
	for _, a := range Artists[1:] {
		for _, spec := range specs {
			val := a.Get(spec.Name)
			switch fmt.Sprintf("%T", val) {
			case "int":
				v := val.(int)
				if spec.Minmax {
					if Min[spec.Name].(int) > v {
						Min[spec.Name] = v
					}
					if Max[spec.Name].(int) < v {
						Max[spec.Name] = v
					}
				}
				if spec.Checkbox || spec.Options {
					i := fmt.Sprintf("%d", v)
					Check[spec.Name][i] = false
				}
			case "string":
				v := val.(string)
				layout := "02-01-2006"
				parsedTime, err := time.Parse(layout, v)
				if err != nil {
					fmt.Printf("Error parsing date: %v\n", err)
					continue
				}
				if spec.Minmax {
					minDate, err := time.Parse(layout, fmt.Sprintf("%v", Min[spec.Name]))
					if err != nil {
						fmt.Println("Error getting date")
					}
					maxDate, err := time.Parse(layout, fmt.Sprintf("%v", Max[spec.Name]))
					if err != nil {
						fmt.Println("Error getting date")
					}
					if parsedTime.Before(minDate) {
						Min[spec.Name] = v
					}
					if parsedTime.After(maxDate) {
						Max[spec.Name] = v
					}
				}
				if spec.Checkbox || spec.Options {
					Check[spec.Name] = make(map[string]bool)
					Check[spec.Name][v] = false
				}
			case "[]string":
				vs := val.([]string)
				if spec.Minmax {
					vmin := vs[0]
					vmax := vs[len(vs)-1]
					if fmt.Sprintf("%v", Min[spec.Name]) > vmin {
						Min[spec.Name] = vmin
					}
					if fmt.Sprintf("%v", Max[spec.Name]) < vmax {
						Max[spec.Name] = vmax
					}
				}
				if spec.Checkbox || spec.Options {
					if _, ok := Check[spec.Name]; !ok {
						Check[spec.Name] = make(map[string]bool)
					}
					for _, v := range vs {
						Check[spec.Name][v] = false
					}
				}
			}
		}
	}
	for i, spec := range specs {
		if spec.Minmax {
			var mi, ma string
			switch fmt.Sprintf("%T", Min[spec.Name]) {
			case "int":
				mi = strconv.Itoa(Min[spec.Name].(int))
				ma = strconv.Itoa(Max[spec.Name].(int))
			case "string":
				mi = fmt.Sprintf("%v", Min[spec.Name])
				ma = fmt.Sprintf("%v", Max[spec.Name])
				specs[i].MMValues["min"] = mi
				specs[i].MMValues["max"] = ma
			}
			specs[i].Attributes["min"] = mi
			specs[i].Attributes["max"] = ma
		}
		if spec.Checkbox || spec.Options {
			boxes := make([]string, len(Check[spec.Name]))
			for c := range Check[spec.Name] {
				boxes = append(boxes, c)
			}
			boxes = boxes[len(boxes)/2:]
			slices.Sort(boxes)
			specs[i].Boxes = boxes
			specs[i].Values = Check[spec.Name]
		}
	}
	return specs
}

func getFilter(w http.ResponseWriter, r *http.Request) (map[string]map[string]any, error) {
	filter := make(map[string]map[string]any) // Filtering criteria will be saved here from form
	r.ParseForm()
	for key, values := range r.Form { // range over input fields of filter-form
		if strings.HasSuffix(key, "_min") || strings.HasSuffix(key, "_max") {
			k := strings.Index(key, "_")
			name := key[:k]
			dim := key[k+1:]
			var f int
			var field types.Field
			for f, field = range Filter0 {
				if field.Name == name {
					break
				}
			}
			if _, ok := filter[name]; !ok {
				filter[name] = make(map[string]any) // In case map is not yet created for this key
			}

			v := values[0] // Values for key is always in array; but for min/max we have just single element
			switch len(v) {
			case 10: // Full date, 1990-10-31  //e.g. 01-02-1990
				filter[name][dim] = v //v[6:] + v[3:5] + v[:2] // Rearrange into <year><mo><dy>
				Filter0[f].MMValues[dim] = v
			case 0: // No date, set min or max to default value
				delete(Filter0[f].MMValues, dim)
				delete(filter[name], dim)
			default: // In default case we expect number (year or number of members)
				if strings.Contains(v, "-") {
					report(v, "Error parsing values", w)
					err := errors.New("error parsing values")
					return nil, err
				}
				Filter0[f].MMValues[dim] = v
				val, err := strconv.Atoi(v)
				if err != nil {
					report(v, "Error parsing values", w)
					return nil, err
				}
				filter[name][dim] = val
			}
		} else {
			//fmt.Println(key, values)
			if _, ok := filter[key]; !ok {
				filter[key] = make(map[string]any) // In case map is not yet created for this key
			}
			filter[key]["values"] = values
		}
	}
	return filter, nil
}

func filterArtists(filter map[string]map[string]any, w http.ResponseWriter) ([]types.ArtistFull, error) {
	var artists []types.ArtistFull
	for _, a := range Artists {
		pass := true // By default all artists are included
		for field, dims := range filter {
			//fmt.Println(field, dims)
			val := a.Get(field)             // Get value from artist for this field (see type ArtistFull.Get())
			if min, ok := dims["min"]; ok { // If we have min field defined
				switch min.(type) {
				case int: // min was either default 0 or single number (wo '-')
					if min.(int) > 0 { // it was specified
						switch val.(type) {
						case int: // val in artist data is number
							if val.(int) < min.(int) {
								pass = false
								break
							}
						case string: // val in artist data is string
							// Extract year
							k := strings.LastIndex(val.(string), "-")
							v, _ := strconv.Atoi(val.(string)[k+1:])
							if v < min.(int) {
								pass = false
								break
							}
						case []string:
							// value in artist data is slice (hopefully sorted)
							val0 := val.([]string)
							va := val0[len(val0)-1] // Last value
							v, err := strconv.Atoi(va[:4])
							if err != nil {
								report(val0, "Error parsing values", w)
								return nil, err
							}
							if v < min.(int) {
								pass = false
								break
							}
						}
					}
				case string: // min was string
					switch val.(type) {
					case string:
						mi := min.(string)
						va := val.(string)
						switch len(mi) {
						case 8: // it is transformed already into form <year><mo><dy>
							if mi > va[6:]+va[3:5]+va[:2] {
								pass = false
								break
							}
						case 6: // it is transformed already into form <year><mo>
							if mi > va[6:]+va[3:5] {
								pass = false
								break
							}
						}
					case []string:
						mi := min.(string)
						va := val.([]string)
						if mi > va[len(va)-1] {
							pass = false
							break
						}
					}
				}
			}
			if max, ok := dims["max"]; ok { // Same for max as for min
				switch max.(type) {
				case int:
					if max.(int) < 3000 {
						switch val.(type) {
						case int:
							if val.(int) > max.(int) {
								pass = false
								break
							}
						case string:
							k := strings.LastIndex(val.(string), "-")
							v, _ := strconv.Atoi(val.(string)[k+1:])
							if v > max.(int) {
								pass = false
								break
							}
						case []string:
							// value in artist data is slice (hopefully sorted)
							val0 := val.([]string)
							va := val0[0] // First value
							v, err := strconv.Atoi(va[:4])
							if err != nil {
								report(val0, "Error parsing values", w)
								return nil, err
							}
							if v > max.(int) {
								pass = false
								break
							}
						}
					}
				case string:
					switch val.(type) {
					case string:
						va := val.(string)
						ma := max.(string)
						switch len(ma) {
						case 8:
							if ma < va[6:]+va[3:5]+va[:2] {
								pass = false
								break
							}
						case 6:
							if ma < va[6:]+va[3:5] {
								pass = false
								break
							}
						}
					case []string:
						va := val.([]string)
						ma := max.(string)
						if ma < va[0] {
							pass = false
							break
						}
					}
				}
			}
			if values, ok := dims["values"]; ok && len(values.([]string)) > 0 {
				//va := val.([]string)
				switch val.(type) {
				case int:
					v := fmt.Sprintf("%d", val.(int))
					if !slices.Contains(values.([]string), v) {
						pass = false
						break
					}
				case string:
					v := val.(string)
					if !slices.Contains(values.([]string), v) {
						pass = false
						break
					}
				case []string:
					va := val.([]string)
					pass0 := false
					for _, v := range va {
						if slices.Contains(values.([]string), v) {
							pass0 = true
							break
						}
					}
					if !pass0 {
						pass = false
						break
					}
				}
			}
		}
		// artist has passed all checks
		if pass {
			artists = append(artists, a)
		}

	}
	return artists, nil
}
