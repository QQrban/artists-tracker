package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

//os.Setenv("GEOAPIFY_KEY", "3c6389ab57f04024b13cd2d76ff57b84")

func Geolocate(place string) ([]float64, error) {
	//example := "38%20Upper%20Montagu%20Street%2C%20Westminster%20W1H%201LJ%2C%20United%20Kingdom&"
	//example := "Westminster%20UK"

	const key = "3c6389ab57f04024b13cd2d76ff57b84"
	const base = "https://api.geoapify.com/v1/geocode/search?apiKey=" + key

	url := base + "&lang=en&limit=10&text=" + url.QueryEscape(place)
	method := "GET"
	//&type=city" + "
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	/*body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}*/

	var result map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	coordinates, ok := result["features"].([]any)[0].(map[string]any)["geometry"].(map[string]any)["coordinates"]
	if !ok {
		fmt.Println("Error while parsing JSON", place, result["query"])
		return nil, err
	}
	coord := []float64{coordinates.([]any)[0].(float64), coordinates.([]any)[1].(float64)}
	return coord, nil
}
