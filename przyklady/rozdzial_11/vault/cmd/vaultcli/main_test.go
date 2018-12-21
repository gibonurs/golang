package main

import "testing"

func TestPop(t *testing.T) {
  args := []string{"jeden", "dwa", "trzy"}
  var s string
  s, args = pop(args)
  if s != "jeden" {
    t.Errorf("nieoczekiwana wartość: \"%s\"", s)
  }
  s, args = pop(args)
  if s != "dwa" {
    t.Errorf("nieoczekiwana wartość: \"%s\"", s)
  }
  s, args = pop(args)
  if s != "trzy" {
    t.Errorf("nieoczekiwana wartość: \"%s\"", s)
  }
  s, args = pop(args)
  if s != "" {
    t.Errorf("nieoczekiwana wartość: \"%s\"", s)
  }
}
