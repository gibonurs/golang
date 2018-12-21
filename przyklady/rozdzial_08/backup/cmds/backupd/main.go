package main

import (
  "encoding/json"
  "errors"
  "flag"
  "fmt"
  "log"
  "os"
  "os/signal"
  "syscall"
  "time"

  "github.com/matryer/filedb"
  "../../../backup"
)

type path struct {
  Path string
  Hash string
}

func main() {
  var fatalErr error
  defer func() {
    if fatalErr != nil {
      log.Fatalln(fatalErr)
    }
  }()
  var (
    interval = flag.Duration("interval", 10 * time.Second, "odstęp pomiędzy kolejnymi sprawdzeniami")
    archive = flag.String("archive", "archive", "ścieżka dostępu do katalogu archiwum")
    dbpath = flag.String("db", "db", "ścieżka do bazy danych")
  )
  flag.Parse()

  m := &backup.Monitor{
    Destination: *archive,
    Archiver: backup.ZIP,
    Paths: make(map[string]string),
  }

  db, err := filedb.Dial(*dbpath)
  if err != nil {
    fatalErr = err
    return
  }
  defer db.Close()
  col, err := db.C("paths")
  if err != nil {
    fatalErr = err
    return
  }

  var path path
  col.ForEach(func(_ int, data []byte) bool {
    if err := json.Unmarshal(data, &path); err != nil {
      fatalErr = err
      return true
    }
    m.Paths[path.Path] = path.Hash
    return false // dalej
  })
  if fatalErr != nil {
    return
  }
  if len(m.Paths) < 1 {
    fatalErr = errors.New("Brak ścieżek - użyj programu backup, by dodać przynajmniej jedną!")
    return
  }

  check(m, col)
  signalChan := make(chan os.Signal, 1)
  signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
  for {
    select {
      case <-time.After(*interval):
        check(m, col)
      case <-signalChan:
        // zatrzymanie
        fmt.Println()
        log.Printf("Zatrzymywanie...")
        return
    }
  }
}


func check(m *backup.Monitor, col *filedb.C) {
  log.Println("Sprawdzanie...")
  counter, err := m.Now()
  if err != nil {
    log.Fatalln("Nie udało się zarchiwizować katalogów: ", err)
  }
  if counter > 0 {
    log.Printf("  Liczba zarchiwizowanych katalogów: %d.\n", counter)
    // aktualizacja skrótów
    var path path
    col.SelectEach(func(_ int, data []byte) (bool, []byte, bool) {
      if err := json.Unmarshal(data, &path); err != nil {
        log.Println("Błąd odtwarzania danych (pomijam):", err)
        return true, data, false
      }
      path.Hash, _ = m.Paths[path.Path]
      newdata, err := json.Marshal(&path)
      if err != nil {
        log.Println("Błąd szeregowania danych (pomijam):", err)
        return true, data, false
      }
      return true, newdata, false
    })
  } else {
    log.Println("  Brak zmian.")
  }
}
