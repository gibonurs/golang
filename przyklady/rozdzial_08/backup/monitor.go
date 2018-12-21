package backup

import (
  "fmt"
  "path/filepath"
  "time"
)

// Monitor sprawdza ścieżki i archiwuzuje te, które uległy zmianie.
type Monitor struct {
  Paths map[string]string
  Archiver Archiver
  Destination string
}


// Now sprawdza wszystie katalogu w mapie Paths, porównując je 
// z ostatnim wyznaczonym skrótem. Dla każdej ścieżki, której 
// bieżący skrót różni się od zapamiętanego, zostanie wywołana
// metoda Archive.
func (m *Monitor) Now() (int, error) {
  var counter int
  for path, lastHash := range m.Paths {
    newHash, err := DirHash(path)
    if err != nil {
      return 0, err
    }
    if newHash != lastHash {
      err := m.act(path)
      if err != nil {
        return counter, err
      }
      m.Paths[path] = newHash // aktualizacja skrótu
      counter++
    }
  }
  return counter, nil
}

func (m *Monitor) act(path string) error {
  dirname := filepath.Base(path)
  filename := fmt.Sprintf(m.Archiver.DestFmt(), time.Now().UnixNano())
  return m.Archiver.Archive(path, filepath.Join(m.Destination, dirname, filename))
}
