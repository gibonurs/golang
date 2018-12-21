package main

import (
  "bufio"
  "fmt"
  "log"
  "os"

  "../thesaurus"
)

func main() {
  apiKey := os.Getenv("BHT_APIKEY")
  thesaurus := &thesaurus.BigHugh{APIKey: apiKey}
  s := bufio.NewScanner(os.Stdin)
  for s.Scan() {
    word := s.Text()
    syns, err := thesaurus.Synonyms(word)
    if err != nil {
      log.Fatalln("Nie udało się pobrać synonimów słowa \""+word+"\"", err)
    }
    if len(syns) == 0 {
      log.Fatalln("Nie udało się znaleźć synonimów słowa \"" + word + "\"")
    }
    for _, syn := range syns {
      fmt.Println(syn)
    }
  }
}
