package main

import (
  "strings"
)

// PathSeparator to znak używany do rozdzielania elementów
// ścieżek HTTP.
const PathSeparator = "/"

// Path reprezentuje ścieżkę żądania.
type Path struct {
  Path string
  ID string
}

// NewPath tworzy nowy obiekt Path na podstawie przekazanego 
// łańcucha znaków.
func NewPath(p string) *Path {
  var id string
  p = strings.Trim(p, PathSeparator)
  s := strings.Split(p, PathSeparator)
  if len(s) > 1 {
    id = s[len(s)-1]
    p = strings.Join(s[:len(s)-1], PathSeparator)
  }
  return &Path{Path: p, ID: id}
}

// HasID zwraca informację czy przekazana ścieżka
// zawiera ID czy nie.
func (p *Path) HasID() bool {
  return len(p.ID) > 0
}