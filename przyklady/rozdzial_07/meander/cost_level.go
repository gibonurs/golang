package meander

import (
  "errors"
  "strings"
)

// Cost reprezentuje poziom cen
type Cost int8

// Cost to wyliczenie reprezentuje różne poziomu cen
const (
  _ Cost = iota
  Cost1
  Cost2
  Cost3
  Cost4
  Cost5
)

var costStrings = map[string]Cost{
  "$":     Cost1,
  "$$":    Cost2,
  "$$$":   Cost3,
  "$$$$":  Cost4,
  "$$$$$": Cost5,
}

func (l Cost) String() string {
  for s, v := range costStrings {
    if l == v {
      return s
    }
  }
  return "invalid"
}

// ParseCost przetwarza łańcuch znaków na wartość typu Cost.
func ParseCost(s string) Cost {
  return costStrings[s]
}

// CostRange reprezentuje zakres wartości typu Cost.
type CostRange struct {
  From Cost
  To   Cost
}

func (r CostRange) String() string {
  return r.From.String() + "..." + r.To.String()
}

// ParseCostRange przetwarza łańcuch zakresu cen na wartość typu CostRange.
func ParseCostRange(s string) (CostRange, error) {
  var r CostRange
  segs := strings.Split(s, "...")
  if len(segs) != 2 {
    return r, errors.New("nieprawidłowy zakres cen")
  }
  r.From = ParseCost(segs[0])
  r.To = ParseCost(segs[1])
  return r, nil
}
