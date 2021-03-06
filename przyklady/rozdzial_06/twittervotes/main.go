package main

import (
  "log"
  "os"
  "os/signal"
  "sync"
  "syscall"  
  "time"

//  "github.com/joeshaw/envdecode"
//  "github.com/matryer/go-oauth/oauth"
  "github.com/nsqio/go-nsq"
  "gopkg.in/mgo.v2"
)


func main(){
  var stoplock sync.Mutex // zabezpieczenie zatrzymania
  stop := false
  stopChan := make(chan struct{}, 1)
  signalChan := make(chan os.Signal, 1)
  go func() {
    <-signalChan
    stoplock.Lock()
    stop = true
    stoplock.Unlock()
    log.Println("Zatrzymywanie...")
    stopChan <- struct{}{}
    closeConn()
  }()
  signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

  if err := dialdb(); err != nil {
    log.Fatalln("Nie udało się nazwiązać połączenia z MongoDB:", err)
  }
  defer closedb()


  // zaczynamy wszystko
  votes := make(chan string) // kanał do przekazywania głosów
  publisherStoppedChan := publishVotes(votes)
  twitterStoppedChan := startTwitterStream(stopChan, votes)
  go func() {
    for {
      time.Sleep(1 * time.Minute)
      closeConn()
      stoplock.Lock()
      if stop {
        stoplock.Unlock()
        return
      }
      stoplock.Unlock()
    }
  }()
  <-twitterStoppedChan
  close(votes)
  <-publisherStoppedChan

}

var db *mgo.Session
func dialdb() error {
  var err error
  log.Println("dialing mongodb: localhost")
  db, err = mgo.Dial("localhost")
  return err
}
func closedb() {
  db.Close()
  log.Println("Zamknięto połączenie z bazą danych")
}



type poll struct {
  Options []string
}
func loadOptions() ([]string, error) {
  var options []string
  iter := db.DB("ballots").C("polls").Find(nil).Iter()
  var p poll
  for iter.Next(&p) {
    options = append(options, p.Options...)
  }
  iter.Close()
  return options, iter.Err()
}    


func publishVotes(votes <-chan string) <-chan struct{} {
  stopchan := make(chan struct{}, 1)
  pub, _ := nsq.NewProducer("localhost:4150", nsq.NewConfig())
  go func() {
    for vote := range votes {
      pub.Publish("votes", []byte(vote)) // publikacja głosu
    }
    log.Println("Producent: Zatrzymywanie...")
    pub.Stop()
    log.Println("Producent: Zatrzymany")
    stopchan <- struct{}{}
  }()
  return stopchan
}