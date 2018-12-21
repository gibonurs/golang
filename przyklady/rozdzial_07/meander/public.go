package meander

// Facade reprezentuje obiekty które udostępniają publiczny
// widok samych siebie.
type Facade interface {
  Public() interface{}
}

// Public pobiera publiczny widok określonego obiektu, o ile
// ten implementuje interfejs Facade. W przeciwny razie zwraca
// pierwotny obiekt.
func Public(o interface{}) interface{} {
  if p, ok := o.(Facade); ok {
    return p.Public()
  }
  return o
}