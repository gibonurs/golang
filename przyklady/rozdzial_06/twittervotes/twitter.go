package main

import (
  "encoding/json"
  "io"
  "log"
  "net"
  "net/http"
  "net/url"
//  "os"
//  "os/signal"
  "strconv"
  "strings"
//  "syscall"
  "time"
  "sync"

  "github.com/joeshaw/envdecode"
  "github.com/matryer/go-oauth/oauth"
//  "github.com/nsqio/go-nsq"
//  "gopkg.in/mgo.v2"
)

var conn net.Conn
func dial(netw, addr string) (net.Conn, error) {
  if conn != nil {
    conn.Close()
    conn = nil
  }
  netc, err := net.DialTimeout(netw, addr, 5*time.Second)
  if err != nil {
    return nil, err
  }
  conn = netc
  return netc, nil
}

var reader io.ReadCloser
func closeConn() {
  if conn != nil {
    conn.Close()
  }
  if reader != nil {
    reader.Close()
  }
}


var (
  authClient *oauth.Client
  creds *oauth.Credentials
)

func setupTwitterAuth() {
  var ts struct {
    ConsumerKey string `env:"SP_TWITTER_KEY,required"`
    ConsumerSecret string `env:"SP_TWITTER_SECRET,required"`
    AccessToken string `env:"SP_TWITTER_ACCESSTOKEN,required"`
    AccessSecret string `env:"SP_TWITTER_ACCESSSECRET,required"`
  }
  if err := envdecode.Decode(&ts); err != nil {
    log.Fatalln(err)
  }
  creds = &oauth.Credentials{
    Token: ts.AccessToken,
    Secret: ts.AccessSecret,
  }
  authClient = &oauth.Client{
    Credentials: oauth.Credentials{
      Token: ts.ConsumerKey,
      Secret: ts.ConsumerSecret,
    },
  }
}


var (
  authSetupOnce sync.Once
  httpClient *http.Client
)
func makeRequest(req *http.Request, params url.Values) (*http.Response, error) {
  authSetupOnce.Do(func() {
    setupTwitterAuth()
    httpClient = &http.Client{
      Transport: &http.Transport{
        Dial: dial,
      },
    }
  })
  formEnc := params.Encode()
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  req.Header.Set("Content-Length", strconv.Itoa(len(formEnc)))
  req.Header.Set("Authorization", authClient.AuthorizationHeader(creds, "POST", req.URL, params))
  return httpClient.Do(req)
}


type tweet struct {
  Text string
}

func readFromTwitter(votes chan<- string) {
  options, err := loadOptions()
  if err != nil {
    log.Println("Nie udało się wczytać opcji:", err)
    return
  }
  u, err := url.Parse("https://stream.twitter.com/1.1/statuses/filter.json")
  if err != nil {
    log.Println("Nie udało się utworzyć żądania filtra:", err)
    return
  }
  query := make(url.Values)
  query.Set("track", strings.Join(options, ","))
  req, err := http.NewRequest("POST",u.String(),strings.NewReader(query.Encode()))
  if err != nil {
    log.Println("Nie udało się utworzyć żądania filtra:", err)
    return
  }
  resp, err := makeRequest(req, query)
  if err != nil {
    log.Println("Nie udało się wykonać żądania:", err)
    return
  }
  reader := resp.Body
  decoder := json.NewDecoder(reader)
  for {
    var t tweet
    if err := decoder.Decode(&t); err != nil {
      break
    }
    for _, option := range options {
      if strings.Contains(
        strings.ToLower(t.Text),
        strings.ToLower(option),
      ) {
        log.Println("głos:", option)
        votes <- option
      }
    }
  }
}    


func startTwitterStream(stopchan <-chan struct{}, votes chan<- string) <-chan struct{} {
  stoppedchan := make(chan struct{}, 1)
  go func() {
    defer func() {
      stoppedchan <- struct{}{}
    }()
    for {
      select {
      case <-stopchan:
        log.Println("Zamykanie połączenia z Twitterem...")
        return
      default:
        log.Println("Przeszukiwanie Twittera...")
        readFromTwitter(votes)
        log.Println(" (czekam...)")
        time.Sleep(10 * time.Second) // czekamy chwilę przed ponownym nawiązaniem połączenia
      }
    }
  }()
  return stoppedchan
}