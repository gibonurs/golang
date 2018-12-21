package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
)

var server = "localhost"
var user = "sa"
var port = 1433
var password = "P@ssw0rd"
var database = "testdb"

type Person struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Location string `json:"location,omitempty"`
}

var people []Person

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/people/{name}+{location}", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people/mod/{id}+{name}+{location}", UpdatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePersonEndpoint).Methods("DELETE")
	log.Fatal(http.ListenAndServeTLS(":8000", "server.cert", "server.key", router))

}

func readdb(cnstr string) (int, error) {

	db, err := sql.Open("sqlserver", cnstr)

	if err != nil {
		log.Println("Drivers error ", err.Error())
		return 0, err
	}

	ctx := context.Background()
	err1 := db.PingContext(ctx)

	if err1 != nil {
		log.Println("Problem with connection pool", err1.Error())
		return 0, err1
	}

	sqlq := "select * from TestSchema.Employees"
	rows, errq := db.QueryContext(ctx, sqlq)

	if errq != nil {
		log.Println("Context error", errq.Error())
		return 0, errq
	}

	RowCount := 0
	people = nil
	for rows.Next() {
		var id int
		var name, location string

		errqw := rows.Scan(&id, &name, &location)
		if errqw != nil {
			log.Println("Query Error: ", errqw.Error())
			return -1, errqw
		}
		// fmt.Printf("%d %s %s \n", id, name, location)

		people = append(people, Person{ID: strconv.Itoa(id), Name: name, Location: location})

		RowCount++
	}
	fmt.Println(people)
	defer rows.Close()
	return RowCount, nil
}

func writedb(cnstr string, name string, location string) (int64, error) {
	db, err := sql.Open("sqlserver", cnstr)

	if err != nil {
		log.Println("Drivers error ", err.Error())
		return 0, err
	}

	ctx := context.Background()
	err1 := db.PingContext(ctx)

	if err1 != nil {
		log.Println("Problem with connection pool", err1.Error())
		return 0, err1
	}

	//weryfikacja poprawnosci składni należy ja zamknąc po zakończeniu
	var zapis = "insert into TestSchema.Employees(name,location) values(@name,@location); select @@identity;"

	skladnia, errsk := db.Prepare(zapis)
	if errsk != nil {
		fmt.Println("Error skladnia")
		return 0, errsk
	}
	defer skladnia.Close()

	row := skladnia.QueryRowContext(ctx,
		sql.Named("name", name),
		sql.Named("location", location))

	var noweID int64
	errid := row.Scan(&noweID)
	if errid != nil {
		fmt.Println("cos sie wyje......o")
		return 0, errid
	}

	return noweID, nil
}

func updatedb(cnstr string, name string, location string, id int64) error {
	db, err := sql.Open("sqlserver", cnstr)

	if err != nil {
		log.Println("Drivers error ", err.Error())
		return err
	}

	ctx := context.Background()
	err1 := db.PingContext(ctx)

	if err1 != nil {
		log.Println("Problem with connection pool", err1.Error())
		return err1
	}

	//weryfikacja poprawnosci składni należy ja zamknąc po zakończeniu
	var update = fmt.Sprintf("update TestSchema.Employees set name=@name,location=@location where id=@id")

	skladnia, errsk := db.Prepare(update)
	if errsk != nil {
		fmt.Println("Error skladnia")
		return errsk
	}
	defer skladnia.Close()

	_, errupd := skladnia.ExecContext(ctx,
		update,
		sql.Named("name", name),
		sql.Named("location", location),
		sql.Named("id", id))

	if errupd != nil {
		fmt.Println("Error skladnia")
		return errupd
	}

	return nil
}

func updatedbv2(cnstr string, name string, location string, id int64) (string, string, error) {
	db, err := sql.Open("sqlserver", cnstr)

	if err != nil {
		log.Println("Drivers error ", err.Error())
		return "", "", err
	}

	ctx := context.Background()
	err1 := db.PingContext(ctx)

	if err1 != nil {
		log.Println("Problem with connection pool", err1.Error())
		return "", "", err1
	}

	//weryfikacja poprawnosci składni należy ja zamknąc po zakończeniu
	var update = fmt.Sprintf("select * from TestSchema.Employees where id=@id;update TestSchema.Employees set name=@name,location=@location where id=@id")

	skladnia, errsk := db.Prepare(update)
	if errsk != nil {
		fmt.Println("Error skladnia")
		return "", "", errsk
	}
	defer skladnia.Close()

	oldrow := skladnia.QueryRowContext(ctx,
		update,
		sql.Named("name", name),
		sql.Named("location", location),
		sql.Named("id", id))

	var oldName string
	var oldLocation string
	errupd := oldrow.Scan(&id, &oldName, &oldLocation)

	if errupd != nil {
		fmt.Println("Error skladnia")
		return "", "", errupd
	}

	return oldLocation, oldName, nil
}

func deletedb(cnstr string, id int64) error {
	db, err := sql.Open("sqlserver", cnstr)

	if err != nil {
		log.Println("Drivers error ", err.Error())
		return err
	}

	ctx := context.Background()
	err1 := db.PingContext(ctx)

	if err1 != nil {
		log.Println("Problem with connection pool", err1.Error())
		return err1
	}

	//weryfikacja poprawnosci składni należy ja zamknąc po zakończeniu
	var delete = fmt.Sprintf("delete from TestSchema.Employees where id=@id")

	skladnia, errsk := db.Prepare(delete)
	if errsk != nil {
		fmt.Println("Error skladnia")
		return errsk
	}
	defer skladnia.Close()

	_, errdel := skladnia.ExecContext(ctx,
		delete,
		sql.Named("id", id))

	if errdel != nil {
		fmt.Println("Error skladnia")
		return errdel
	}

	return nil
}

//import JSON
func GetPersonEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}
func GetPeopleEndpoint(w http.ResponseWriter, r *http.Request) {
	constr := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;", server, user, password, port, database)
	fmt.Println(constr)

	RowCount, ErrRead := readdb(constr)
	fmt.Printf("Record count = %d error %s\n", RowCount, ErrRead)

	json.NewEncoder(w).Encode(people)
}
func CreatePersonEndpoint(w http.ResponseWriter, r *http.Request) {
	constr := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;", server, user, password, port, database)
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.Name = params["name"]
	person.Location = params["location"]

	idnowego, errordodawania := writedb(constr, params["name"], params["location"])
	if errordodawania != nil {
		log.Println("Drivers error ", errordodawania.Error())
	}

	person.ID = strconv.FormatInt(idnowego, 10)
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
	fmt.Println(people)
}

func UpdatePersonEndpoint(w http.ResponseWriter, r *http.Request) {
	constr := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;", server, user, password, port, database)
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	person.Name = params["name"]
	person.Location = params["location"]

	modid, _ := strconv.ParseInt(params["id"], 10, 32)
	oldLoc, oldName, errorupdatev2 := updatedbv2(constr, params["name"], params["location"], modid)
	if errorupdatev2 != nil {
		log.Println("Update error", errorupdatev2.Error())
	}
	fmt.Printf("Old Location=%s old Name=%s\n", oldLoc, oldName)
	RowCount, _ := readdb(constr)
	fmt.Printf("Modified record = %d form %d records, Name=%s, Location=%s\n", modid, RowCount, params["name"], params["location"])
	json.NewEncoder(w).Encode(people)
	fmt.Println()
}

func DeletePersonEndpoint(w http.ResponseWriter, r *http.Request) {
	constr := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;", server, user, password, port, database)
	params := mux.Vars(r)
	var people2 []Person
	for _, item := range people {
		if item.ID == params["id"] {
			people2 = append(people2, Person{ID: item.ID, Name: item.Name, Location: item.Location})
			delindex, _ := strconv.ParseInt(item.ID, 10, 64)
			errordelete := deletedb(constr, delindex)
			if errordelete != nil {
				log.Println("Update error", errordelete.Error())
			}
			break
		}
	}
	fmt.Println(people2)

	json.NewEncoder(w).Encode(people2)
}
