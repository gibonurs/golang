package main

import (
  "flag"
  "fmt"
  "os"
  "os/signal"
  "time"
  "log"
  "sync"
  "syscall"

  "github.com/nsqio/go-nsq"

  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
)

const updateDuration = 1 * time.Second

var fatalErr error

func fatal(e error) {
  fmt.Println(e)
  flag.PrintDefaults()
  fatalErr = e
}

func main() {
  
  defer func() {
    if fatalErr != nil {
      os.Exit(1)
    }
  }()

  log.Println("Nawiązywanie połączenia z bazą danych...")
  db, err := mgo.Dial("localhost")
  if err != nil {
    fatal(err)
    return
  }
  defer func() {
    log.Println("Zamykanie połączenia z bazą danych...")
    db.Close()
  }()
  pollData := db.DB("ballots").C("polls")

  var counts map[string]int
  var countsLock sync.Mutex

  log.Println("Nawiązywanie połączenia z nsq...")
  q, err := nsq.NewConsumer("votes", "counter", nsq.NewConfig())
  if err != nil {
    fatal(err)
    return
  }

  q.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
    countsLock.Lock()
    defer countsLock.Unlock()
    if counts == nil {
      counts = make(map[string]int)
    }
    vote := string(m.Body)
    counts[vote]++
    return nil
  }))

  if err := q.ConnectToNSQLookupd("localhost:4161");
   err != nil {
    fatal(err)
    return
  }

  log.Println("Czekamy na głosy w kolejce nsq...")
  ticker := time.NewTicker(updateDuration)
  termChan := make(chan os.Signal, 1)
  signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM,syscall.SIGHUP)
  for {
    select {
    case <-ticker.C:
      doCount(&countsLock, &counts,pollData) 
    case <- termChan:ticker.Stop()
      q.Stop()
    case <-q.StopChan:
      // koniec
      return
    }
  }

}  

func doCount(countsLock *sync.Mutex, counts *map[string]int, pollData *mgo.Collection) {
  countsLock.Lock()
  defer countsLock.Unlock()
  if len(*counts) == 0 {
    log.Println("Brak głosów, pominięcie aktualizacji.")
    return
  }
  log.Println("Aktualizacja bazy danych...")
  log.Println(*counts)
  ok := true
  for option, count := range *counts {
    sel := bson.M{"options": bson.M{"$in": []string{option}}}
    up := bson.M{"$inc": bson.M{"results." +  option:count}}
    if _, err := pollData.UpdateAll(sel, up); err != nil {
      log.Println("Błąd aktualizacji bazy danych:", err)
      ok = false
    }
  }
  if ok {
    log.Println("Kończenie aktualizacji bazy danych...")
    *counts = nil // przywrócenie początkowej wartości liczników
  }
}


