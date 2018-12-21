package vault

import (
  "testing"

  "golang.org/x/net/context"
)

func TestHasherService(t *testing.T) {
  srv := NewService()
  ctx := context.Background()
  h, err := srv.Hash(ctx, "hasło")
  if err != nil {
    t.Errorf("Skrót: %s", err)
  }
  ok, err := srv.Validate(ctx, "hasło", h)
  if err != nil {
    t.Errorf("Poprawne: %s", err)
  }
  if !ok {
    t.Error("Metoda Valid miała zwrócić wartość true!")
  }
  ok, err = srv.Validate(ctx, "złe hasło", h)
  if err != nil {
    t.Errorf("Poprawne: %s", err)
  }
  if ok {
    t.Error("Metoda Valid miała zwrócić wartość false!")
  }
}
