package main

func main() {
data := struct {
Message string `json:"message"`
}{Message: "Witaj, świecie!"}
err := json.NewEncoder(os.Stdout).Encode(data)
if err != nil {
log.Fatalln(err)
}
}
