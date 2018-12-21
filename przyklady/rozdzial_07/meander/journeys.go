package meander

import (
  "strings"
)

// j reprezentuje szablon podróży.
type j struct {
  Name       string
  PlaceTypes []string
}

// Public zwraca publiczną postać tej podróży gets a public view of this Journey.
func (j j) Public() interface{} {
  return map[string]interface{}{
    "name": j.Name,
    "journey": strings.Join(j.PlaceTypes, "|"),
  }
}

// Journeys reprezentuje wstępnie zdefiniowane dane podróży.
var Journeys = []interface{}{
  j{Name: "romantyczny dzień", PlaceTypes: []string{"park", "bar", "movie_theatre", "restaurant", "florist", "taxi_stand"}},
  j{Name: "szał zakupów", PlaceTypes: []string{"department_store", "cafe", "clothing_store", "jewelry_store", "shoe_store"}},
  j{Name: "nocna eksapada", PlaceTypes: []string{"bar", "casino", "food", "bar", "night_club", "bar", "bar", "hospital"}},
  j{Name: "pełna kultura", PlaceTypes: []string{"museum", "cafe", "cemetery", "library", "art_gallery"}},
  j{Name: "chwila dla siebie", PlaceTypes: []string{"hair_care", "beauty_salon", "cafe", "spa"}},
}
