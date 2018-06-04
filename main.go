package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "my_db"
)

type Donor struct {
	Id       sql.NullInt64  `json:"id"`
	Email    sql.NullString `json:"email"`
	Name     sql.NullString `json:"name"`
	Gender   sql.NullString `json:"gender"`
	Age      sql.NullInt64  `json:"age"`
	City     sql.NullString `json:"city"`
	Password sql.NullString `json:"password"`
}
type DonorFmt struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Age      int64  `json:"age"`
	City     string `json:"city"`
	Password string `json:"password"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/donors", GetDonors).Methods("GET")
	router.HandleFunc("/donors/{email}", GetDonor).Methods("GET")
	router.HandleFunc("/donors", CreateDonor).Methods("POST")
	router.HandleFunc("/donors/{email}", UpdateDonor).Methods("POST")
	router.HandleFunc("/donors/{email}", DeleteDonor).Methods("DELETE")

	router.HandleFunc("/hospitals", GetHospitals).Methods("GET")
	router.HandleFunc("/hospitals/{email}", GetHospital).Methods("GET")
	router.HandleFunc("/hospitals", CreateHospital).Methods("POST")
	router.HandleFunc("/hospitals/{email}", UpdateHospital).Methods("POST")
	router.HandleFunc("/hospitals/{email}", DeleteHospital).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func GetDonors(w http.ResponseWriter, r *http.Request) {

	var donors []DonorFmt
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM donors")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var donor Donor
		err = rows.Scan(&donor.Id, &donor.Email, &donor.Name, &donor.Gender, &donor.Age, &donor.City, &donor.Password)
		if err != nil {
			panic(err)
		}
		// fmt.Println(donor)
		donors = append(donors, DonorFmt{Id: donor.Id.Int64, Email: donor.Email.String, Name: donor.Name.String, Gender: donor.Gender.String, Age: donor.Age.Int64, City: donor.City.String, Password: donor.Password.String})
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	fmt.Println(donors)
	json.NewEncoder(w).Encode(donors)
}

func GetDonor(w http.ResponseWriter, r *http.Request) {

	var donor Donor
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	params := mux.Vars(r)
	theEmail := params["email"]
	sqlStatement := `SELECT * FROM donors WHERE email=$1;`
	row := db.QueryRow(sqlStatement, theEmail)
	err = row.Scan(&donor.Id, &donor.Email, &donor.Name, &donor.Gender, &donor.Age, &donor.City, &donor.Password)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	donorFmt := DonorFmt{Id: donor.Id.Int64, Email: donor.Email.String, Name: donor.Name.String, Gender: donor.Gender.String, Age: donor.Age.Int64, City: donor.City.String, Password: donor.Password.String}
	json.NewEncoder(w).Encode(donorFmt)
}

func CreateDonor(w http.ResponseWriter, r *http.Request) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("ERR 1", err.Error)
		panic(err)
	}
	defer db.Close()

	email := r.Form.Get("email")
	var theEmail string
	fmt.Println(theEmail)

	sqlStatement := `SELECT email FROM donors WHERE email=$1;`
	row := db.QueryRow(sqlStatement, email)
	err = row.Scan(&theEmail)
	if err == nil {
		fmt.Println("ERR 2", err.Error)
		fmt.Println("User Already Exist")
	} else {
		r.ParseForm()
		email = r.Form.Get("email")
		name := r.Form.Get("name")
		gender := r.Form.Get("gender")
		age := r.Form.Get("age")
		city := r.Form.Get("city")
		password := r.Form.Get("password")
		fmt.Println(email, name, gender, age, city, password)

		sqlStatement := `
		INSERT INTO donors(
			email, name, gender, age, city, password)
			VALUES  ($1, $2, $3, $4, $5, $6)`

		db.QueryRow(sqlStatement, email, name, gender, age, city, password)
	}
}

func DeleteDonor(w http.ResponseWriter, r *http.Request) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("ERR 1", err.Error)
		panic(err)
	}
	defer db.Close()

	params := mux.Vars(r)
	theEmail := params["email"]

	sqlStatement := `
	DELETE FROM donors
	WHERE email = $1`

	db.QueryRow(sqlStatement, theEmail)
}
func UpdateDonor(w http.ResponseWriter, r *http.Request) {
	var id sql.NullInt64
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("ERR 1", err.Error)
		panic(err)
	}
	defer db.Close()

	params := mux.Vars(r)
	theEmail := params["id"]

	sqlStatement := `SELECT id FROM donors WHERE email=$1;`
	row := db.QueryRow(sqlStatement, theEmail)
	err = row.Scan(&id)
	if err != nil {
		fmt.Println("ERR 2", err.Error)
		fmt.Println("User Does Not Exist!")
	} else {
		r.ParseForm()
		email := r.Form.Get("email")
		name := r.Form.Get("name")
		gender := r.Form.Get("gender")
		age := r.Form.Get("age")
		city := r.Form.Get("city")
		password := r.Form.Get("password")
		fmt.Println(email, name, gender, age, city, password)

		sqlStatement := `
		UPDATE donors
		SET email=$1, name=$2, gender=$3, age=$4, city=$5, password=$6
		WHERE id=$7;`

		db.QueryRow(sqlStatement, email, name, gender, age, city, password, id)
	}
}

func GetHospitals(w http.ResponseWriter, r *http.Request) {

	var donors []DonorFmt
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM donors")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var donor Donor
		err = rows.Scan(&donor.Id, &donor.Email, &donor.Name, &donor.Gender, &donor.Age, &donor.City, &donor.Password)
		if err != nil {
			panic(err)
		}
		// fmt.Println(donor)
		donors = append(donors, DonorFmt{Id: donor.Id.Int64, Email: donor.Email.String, Name: donor.Name.String, Gender: donor.Gender.String, Age: donor.Age.Int64, City: donor.City.String, Password: donor.Password.String})
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	fmt.Println(donors)
	json.NewEncoder(w).Encode(donors)
}

func GetHospital(w http.ResponseWriter, r *http.Request) {

	var donor Donor
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	params := mux.Vars(r)
	theEmail := params["email"]
	sqlStatement := `SELECT * FROM donors WHERE email=$1;`
	row := db.QueryRow(sqlStatement, theEmail)
	err = row.Scan(&donor.Id, &donor.Email, &donor.Name, &donor.Gender, &donor.Age, &donor.City, &donor.Password)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	donorFmt := DonorFmt{Id: donor.Id.Int64, Email: donor.Email.String, Name: donor.Name.String, Gender: donor.Gender.String, Age: donor.Age.Int64, City: donor.City.String, Password: donor.Password.String}
	json.NewEncoder(w).Encode(donorFmt)
}

func CreateHospital(w http.ResponseWriter, r *http.Request) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("ERR 1", err.Error)
		panic(err)
	}
	defer db.Close()

	email := r.Form.Get("email")
	var theEmail string
	fmt.Println(theEmail)

	sqlStatement := `SELECT email FROM donors WHERE email=$1;`
	row := db.QueryRow(sqlStatement, email)
	err = row.Scan(&theEmail)
	if err == nil {
		fmt.Println("ERR 2", err.Error)
		fmt.Println("User Already Exist")
	} else {
		r.ParseForm()
		email = r.Form.Get("email")
		name := r.Form.Get("name")
		gender := r.Form.Get("gender")
		age := r.Form.Get("age")
		city := r.Form.Get("city")
		password := r.Form.Get("password")
		fmt.Println(email, name, gender, age, city, password)

		sqlStatement := `
		INSERT INTO donors(
			email, name, gender, age, city, password)
			VALUES  ($1, $2, $3, $4, $5, $6)`

		db.QueryRow(sqlStatement, email, name, gender, age, city, password)
	}
}

func DeleteHospital(w http.ResponseWriter, r *http.Request) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("ERR 1", err.Error)
		panic(err)
	}
	defer db.Close()

	params := mux.Vars(r)
	theEmail := params["email"]

	sqlStatement := `
	DELETE FROM donors
	WHERE email = $1`

	db.QueryRow(sqlStatement, theEmail)
}
func UpdateHospital(w http.ResponseWriter, r *http.Request) {
	var id sql.NullInt64
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("ERR 1", err.Error)
		panic(err)
	}
	defer db.Close()

	params := mux.Vars(r)
	theEmail := params["id"]

	sqlStatement := `SELECT id FROM donors WHERE email=$1;`
	row := db.QueryRow(sqlStatement, theEmail)
	err = row.Scan(&id)
	if err != nil {
		fmt.Println("ERR 2", err.Error)
		fmt.Println("User Does Not Exist!")
	} else {
		r.ParseForm()
		email := r.Form.Get("email")
		name := r.Form.Get("name")
		gender := r.Form.Get("gender")
		age := r.Form.Get("age")
		city := r.Form.Get("city")
		password := r.Form.Get("password")
		fmt.Println(email, name, gender, age, city, password)

		sqlStatement := `
		UPDATE donors
		SET email=$1, name=$2, gender=$3, age=$4, city=$5, password=$6
		WHERE id=$7;`

		db.QueryRow(sqlStatement, email, name, gender, age, city, password, id)
	}
}
