 
## Sposób użycia

Uruchamianie usługi vaultd:

```
cd vault/cmd/vaultd
go run main.go
```

Wyznaczenie skrótu hasła:

```bash
curl -XPOST -d'{"password":"MySecretPassword123"}' localhost:8080/hash
```

```json
{"hash":"$2a$10$L/Riz9xbgTBDn7F6uLInq.9Tr67PvBCmxzrLgemitnRM53ht7LGpC"}
```

Weryfikacja hasła przy użyciu skrótu:

```bash
curl -XPOST -d'{"password":"MySecretPassword123","hash":"$2a$10$L/Riz9xbgTBDn7F6uLInq.9Tr67PvBCmxzrLgemitnRM53ht7LGpC"}' localhost:8080/validate
```

```json
{"valid":true}
```

lub jeśli hasło jest błędne:

```bash
curl -XPOST -d'{"password":"NOPE","hash":"$2a$10$L/Riz9xbgTBDn7F6uLInq.9Tr67PvBCmxzrLgemitnRM53ht7LGpC"}' localhost:8080/validate
```

```json
{"valid":false}
```

### Kompilacja protobuf

Instalacja proto3 ze źródeł:

```
brew install autoconf automake libtool
git clone https://github.com/google/protobuf
./autogen.sh ; ./configure ; make ; make install
```

Aktualizacja powiązań protoc Go:

```
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
```

Patrz także: https://github.com/grpc/grpc-go/tree/master/examples

Kompilacja protobuf (z poziomu katalogu `pb`):

```
protoc vault.proto --go_out=plugins=grpc:.
```
