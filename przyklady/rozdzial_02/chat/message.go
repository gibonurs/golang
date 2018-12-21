package main

import (
  "time"
)

// message reprezentuje pojedynczą wiadomość
type message struct {
  Name    string
  Message string
  When    time.Time
}
