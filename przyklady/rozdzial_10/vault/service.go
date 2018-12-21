package vault

import (
   "encoding/json"
   "errors"
   "net/http"
 
   "github.com/go-kit/kit/endpoint"
   "golang.org/x/crypto/bcrypt"
   "golang.org/x/net/context"
)

// Service udostępnia możliwości związane z generacją i sprawdzaniem skrótu hasła.
type Service interface {
  Hash(ctx context.Context, password string) (string, error)
  Validate(ctx context.Context, password, hash string) (bool, error)
}

// NewService tworzy nową instancję Service.
func NewService() Service {
  return vaultService{}
}

type vaultService struct{}

func (vaultService) Hash(ctx context.Context, password string) (string, error) {
  hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  if err != nil {
    return "", err
  }
  return string(hash), nil
}

func (vaultService) Validate(ctx context.Context, password, hash string) (bool, error) {
  err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
  if err != nil {
    return false, nil
  }
  return true, nil
}

type hashRequest struct {
  Password string `json:"password"`
}
type hashResponse struct {
  Hash string `json:"hash"`
  Err  string `json:"err,omitempty"`
}

type validateRequest struct {
  Password string `json:"password"`
  Hash     string `json:"hash"`
}
type validateResponse struct {
  Valid bool   `json:"valid"`
  Err   string `json:"err,omitempty"`
}

func decodeHashRequest(ctx context.Context, r *http.Request) (interface{}, error) {
  var req hashRequest
  err := json.NewDecoder(r.Body).Decode(&req)
  if err != nil {
    return nil, err
  }
  return req, nil
}

func decodeValidateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
  var req validateRequest
  err := json.NewDecoder(r.Body).Decode(&req)
  if err != nil {
    return nil, err
  }
  return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
  return json.NewEncoder(w).Encode(response)
}

func MakeHashEndpoint(srv Service) endpoint.Endpoint {
  return func(ctx context.Context, request interface{}) (interface{}, error) {
    req := request.(hashRequest)
    v, err := srv.Hash(ctx, req.Password)
    if err != nil {
      return hashResponse{v, err.Error()}, nil
    }
    return hashResponse{v, ""}, nil
  }
}

func MakeValidateEndpoint(srv Service) endpoint.Endpoint {
  return func(ctx context.Context, request interface{}) (interface{}, error) {
    req := request.(validateRequest)
    v, err := srv.Validate(ctx, req.Password, req.Hash)
    if err != nil {
      return validateResponse{false, err.Error()}, nil
    }
    return validateResponse{v, ""}, nil
  }
}

// Endpoints representuje wszystkie punkty końcowe naszej usługi.
type Endpoints struct {
  HashEndpoint     endpoint.Endpoint
  ValidateEndpoint endpoint.Endpoint
}

// Hash używa HashEndpoint do wyznaczenia skrótu hasła.
func (e Endpoints) Hash(ctx context.Context, password string) (string, error) {
  req := hashRequest{Password: password}
  resp, err := e.HashEndpoint(ctx, req)
  if err != nil {
    return "", err
  }
  hashResp := resp.(hashResponse)
  if hashResp.Err != "" {
    return "", errors.New(hashResp.Err)
  }
  return hashResp.Hash, nil
}

// Validate używa ValidateEndpoint do sprawdzania poprawności hasła na podstawie skrótu
func (e Endpoints) Validate(ctx context.Context, password,
  hash string) (bool, error) {
  req := validateRequest{Password: password, Hash: hash}
  resp, err := e.ValidateEndpoint(ctx, req)
  if err != nil {
    return false, err
  }
  validateResp := resp.(validateResponse)
  if validateResp.Err != "" {
    return false, errors.New(validateResp.Err)
  }
  return validateResp.Valid, nil
}
