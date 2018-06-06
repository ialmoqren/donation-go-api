package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"os"
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

type Hospital struct {
	Id       sql.NullInt64  `json:"id"`
	Email    sql.NullString `json:"email"`
	Name     sql.NullString `json:"name"`
	City     sql.NullString `json:"city"`
	Password sql.NullString `json:"password"`
}
type HospitalFmt struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	City     string `json:"city"`
	Password string `json:"password"`
}

type Donation struct {
	Id         sql.NullInt64  `json:"id"`
	Type       sql.NullString `json:"type"`
	Notes      sql.NullString `json:"notes"`
	DonorEmail sql.NullString `json:"donor_email"`
}
type DonationFmt struct {
	Id         int64  `json:"id"`
	Type       string `json:"type"`
	Notes      string `json:"notes"`
	DonorEmail string `json:"donor_email"`
}

func main() {

	router := mux.NewRouter()

	headersOk := handlers.AllowedHeaders([]string{"POST"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

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

	router.HandleFunc("/donations", GetDonations).Methods("GET")
	router.HandleFunc("/donations/{email}", GetDonorDonations).Methods("GET")
	router.HandleFunc("/donations", CreateDonation).Methods("POST")
	router.HandleFunc("/donations/{id}", UpdateDonation).Methods("POST")
	router.HandleFunc("/donations/{id}", DeleteDonation).Methods("DELETE")
	fmt.Println("Listening at 8000")
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
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

	rows, err := db.Query("SELECT * FROM donors ORDER BY id;")
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
	sqlStatement := `SELECT * FROM donors WHERE email=$1 ORDER BY id;`
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
	fmt.Println("CreateDonor started")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("ERR 1")
		fmt.Println(err.Error())
		panic(err)
	}
	defer db.Close()

	email := r.Form.Get("email")
	var theEmail string
	fmt.Println("Form.Encode: ", r.Form.Encode())
	fmt.Println("Form.Get('email'): ", r.Form.Get("email"))
	fmt.Println("Body: ", r.Body)
	fmt.Println("URL: ", r.URL)

	sqlStatement := `SELECT email FROM donors WHERE email=$1;`
	row := db.QueryRow(sqlStatement, email)
	err = row.Scan(&theEmail)
	if err == nil {
		fmt.Println("ERR ")

		fmt.Println(err.Error())
		fmt.Println("User Already Exist")
	} else {
		fmt.Println(err.Error())
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
	theEmail := params["email"]

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
	var hospitals []HospitalFmt
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM hospitals ORDER BY id;")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var hospital Hospital
		err = rows.Scan(&hospital.Id, &hospital.Email, &hospital.Name, &hospital.City, &hospital.Password)
		if err != nil {
			panic(err)
		}
		// fmt.Println(donor)
		hospitals = append(hospitals, HospitalFmt{Id: hospital.Id.Int64, Email: hospital.Email.String, Name: hospital.Name.String, City: hospital.City.String, Password: hospital.Password.String})
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	fmt.Println(hospitals)
	json.NewEncoder(w).Encode(hospitals)
}
func GetHospital(w http.ResponseWriter, r *http.Request) {
	var hospital Hospital
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
	sqlStatement := `SELECT * FROM hospitals WHERE email=$1 ORDER BY id;`
	row := db.QueryRow(sqlStatement, theEmail)
	err = row.Scan(&hospital.Id, &hospital.Email, &hospital.Name, &hospital.City, &hospital.Password)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	hospitalFmt := HospitalFmt{Id: hospital.Id.Int64, Email: hospital.Email.String, Name: hospital.Name.String, City: hospital.City.String, Password: hospital.Password.String}
	json.NewEncoder(w).Encode(hospitalFmt)
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

	sqlStatement := `SELECT email FROM hospitals WHERE email=$1;`
	row := db.QueryRow(sqlStatement, email)
	err = row.Scan(&theEmail)
	if err == nil {
		fmt.Println("ERR 2", err.Error)
		fmt.Println("User Already Exist")
	} else {
		r.ParseForm()
		email = r.Form.Get("email")
		name := r.Form.Get("name")
		city := r.Form.Get("city")
		password := r.Form.Get("password")
		fmt.Println(email, name, city, password)

		sqlStatement := `
		INSERT INTO hospitals(
			email, name, city, password)
			VALUES  ($1, $2, $3, $4)`

		db.QueryRow(sqlStatement, email, name, city, password)
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
	DELETE FROM hospitals
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
	theEmail := params["email"]

	sqlStatement := `SELECT id FROM hospitals WHERE email=$1;`
	row := db.QueryRow(sqlStatement, theEmail)
	err = row.Scan(&id)
	if err != nil {
		fmt.Println("ERR 2", err.Error)
		fmt.Println("User Does Not Exist!")
	} else {
		r.ParseForm()
		email := r.Form.Get("email")
		name := r.Form.Get("name")
		city := r.Form.Get("city")
		password := r.Form.Get("password")
		fmt.Println(email, name, city, password)

		sqlStatement := `
		UPDATE hospitals
		SET email=$1, name=$2, city=$3, password=$4
		WHERE id=$5;`

		db.QueryRow(sqlStatement, email, name, city, password, id)
	}
}

func GetDonations(w http.ResponseWriter, r *http.Request) {
	var donations []DonationFmt
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM donations ORDER BY id;")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var donation Donation
		err = rows.Scan(&donation.Id, &donation.Type, &donation.Notes, &donation.DonorEmail)
		if err != nil {
			panic(err)
		}
		donations = append(donations, DonationFmt{Id: donation.Id.Int64, Type: donation.Type.String, Notes: donation.Notes.String, DonorEmail: donation.DonorEmail.String})
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	fmt.Println(donations)
	json.NewEncoder(w).Encode(donations)
}
func GetDonorDonations(w http.ResponseWriter, r *http.Request) {
	var donations []DonationFmt
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

	rows, err := db.Query(`SELECT * FROM donations WHERE donor_email=$1 ORDER BY id;`, theEmail)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var donation Donation
		err = rows.Scan(&donation.Id, &donation.Type, &donation.Notes, &donation.DonorEmail)
		if err != nil {
			panic(err)
		}
		donations = append(donations, DonationFmt{Id: donation.Id.Int64, Type: donation.Type.String, Notes: donation.Notes.String, DonorEmail: donation.DonorEmail.String})
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	fmt.Println(donations)
	json.NewEncoder(w).Encode(donations)
}
func CreateDonation(w http.ResponseWriter, r *http.Request) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("ERR 1", err.Error)
		panic(err)
	}
	defer db.Close()

	r.ParseForm()
	theType := r.Form.Get("type")
	notes := r.Form.Get("notes")
	donorEmail := r.Form.Get("donor_email")
	fmt.Println(theType, notes, donorEmail)

	sqlStatement := `
		INSERT INTO donations(
			type, notes, donor_email)
			VALUES  ($1, $2, $3)`

	db.QueryRow(sqlStatement, theType, notes, donorEmail)
}
func DeleteDonation(w http.ResponseWriter, r *http.Request) {
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
	theId := params["id"]

	sqlStatement := `
	DELETE FROM donations
	WHERE id = $1`

	db.QueryRow(sqlStatement, theId)
}
func UpdateDonation(w http.ResponseWriter, r *http.Request) {
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
	theId := params["id"]

	sqlStatement := `SELECT id FROM donations WHERE id=$1;`
	row := db.QueryRow(sqlStatement, theId)
	err = row.Scan(&id)
	if err != nil {
		fmt.Println("ERR 2", err.Error)
		fmt.Println("User Does Not Exist!")
	} else {
		r.ParseForm()
		theType := r.Form.Get("type")
		notes := r.Form.Get("notes")
		donorEmail := r.Form.Get("donor_email")
		fmt.Println(theType, notes, donorEmail)

		sqlStatement := `
		UPDATE donations
		SET type=$1, notes=$2
		WHERE id=$3;`

		db.QueryRow(sqlStatement, theType, notes, id)
	}
}
