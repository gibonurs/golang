package meander_test

import (
  "testing"

  "github.com/cheekybits/is"
  "github.com/matryer/goblueprints/chapter7/meander"
)

type obj struct {
  value1 string
  value2 string
  value3 string
}

func (o *obj) Public() interface{} {
  return map[string]interface{}{"jedynka": o.value1, "trojka": o.value3}
}

func TestPublic(t *testing.T) {
  is := is.New(t)

  o := &obj{
    value1: "wartosc1",
    value2: "wartosc2",
    value3: "wartosc3",
  }

  v, ok := meander.Public(o).(map[string]interface{})
  is.Equal(true, ok) // powinno być OK
  is.Equal(v["jedynka"], "wartosc1")
  is.Nil(v["dwojka"]) // wartosc2
  is.Equal(v["trojka"], "wartosc3")

}
