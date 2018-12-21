package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

var server = "localhost"
var user = "sa"
var port = 1433
var password = "P@ssw0rd"
var database = "testdb"

type ludzie struct {
	id       int
	name     string
	location string
}

var czlowiek []ludzie

func main() {

	//Odcczyt z bazy
	constr := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;", server, user, password, port, database)
	fmt.Println(constr)
	RowCount, ErrRead := readdb(constr)
	fmt.Printf("Record count = %d error %s\n", RowCount, ErrRead)

	//dodanie nowego do bazy
	idnowego, errordodawania := writedb(constr, "Franek", "Mars")
	if errordodawania != nil {
		log.Println("Drivers error ", errordodawania.Error())
	}
	fmt.Printf("ID Nowego = %d\n", idnowego)

	// //aktualizacja bazy
	// errorupdate := updatedb(constr, "Blob", "Katowice", 17)
	// if errorupdate != nil {
	// 	log.Println("Update error", errorupdate.Error())
	// }

	//aktualizacja bazy v2
	// var oldLoc string
	// var oldName string
	oldLoc, oldName, errorupdatev2 := updatedbv2(constr, "Blobsrrrrr", "Katowicefffff", 17)
	if errorupdatev2 != nil {
		log.Println("Update error", errorupdatev2.Error())
	}
	fmt.Printf("Old Location=%s old Name=%s\n", oldLoc, oldName)

	//delete from db
	errordelete := deletedb(constr, 10)
	if errordelete != nil {
		log.Println("Update error", errordelete.Error())
	}
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

	for rows.Next() {
		var id int
		var name, location string

		errqw := rows.Scan(&id, &name, &location)
		if errqw != nil {
			log.Println("Query Error: ", errqw.Error())
			return -1, errqw
		}
		// fmt.Printf("%d %s %s \n", id, name, location)

		czlowiek = append(czlowiek, ludzie{id: id, name: name, location: location})

		RowCount++
	}
	fmt.Println(czlowiek)
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
