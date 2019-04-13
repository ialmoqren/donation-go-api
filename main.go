package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"os"

	"github.com/gorilla/handlers"
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
	Id      sql.NullInt64  `json:"id"`
	Type    sql.NullString `json:"type"`
	Notes   sql.NullString `json:"notes"`
	DonorId sql.NullString `json:"donor_id"`
	Gender  sql.NullString `json:"gender"`
	Email   sql.NullString `json:"email"`
	Age     sql.NullInt64  `json:"age"`
	City    sql.NullString `json:"city"`
}
type DonationFmt struct {
	Id      int64  `json:"id"`
	Type    string `json:"type"`
	Notes   string `json:"notes"`
	DonorId string `json:"donor_id"`
	Gender  string `json:"gender"`
	Email   string `json:"email"`
	Age     int64  `json:"age"`
	City    string `json:"city"`
}

func main() {

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"})
	credentialsOk := handlers.AllowCredentials()

	fmt.Println("Here")
	//db.Connect();
	router := mux.NewRouter()

	router.HandleFunc("/donors", GetDonors).Methods("GET")
	router.HandleFunc("/donors/{email}", GetDonor).Methods("GET")
	router.HandleFunc("/donorspass/{email}", DonorLogin).Methods("GET")
	router.HandleFunc("/donors", CreateDonor).Methods("POST")
	router.HandleFunc("/donors/{id}", UpdateDonor).Methods("POST")
	router.HandleFunc("/donors/{email}", DeleteDonor).Methods("DELETE")

	router.HandleFunc("/hospitals", GetHospitals).Methods("GET")
	router.HandleFunc("/hospitals/{email}", GetHospital).Methods("GET")
	router.HandleFunc("/hospitalspass/{email}", hospitalLogin).Methods("GET")
	router.HandleFunc("/hospitals", CreateHospital).Methods("POST")
	router.HandleFunc("/hospitals/{id}", UpdateHospital).Methods("POST")
	router.HandleFunc("/hospitals/{email}", DeleteHospital).Methods("DELETE")

	router.HandleFunc("/donations", GetDonations).Methods("GET")
	router.HandleFunc("/donations/{id}", GetDonorDonations).Methods("GET")
	router.HandleFunc("/donations", CreateDonation).Methods("POST")
	router.HandleFunc("/donations/{id}", UpdateDonation).Methods("POST")
	router.HandleFunc("/donations/{id}", DeleteDonation).Methods("DELETE")
	fmt.Println("Listening at 8000")
	log.Fatal(http.ListenAndServe(":8000"+os.Getenv("PORT"), handlers.CORS(originsOk, headersOk, methodsOk, credentialsOk)(router)))
}

func GetDonors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
func DonorLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	params := mux.Vars(r)
	r.ParseForm()
	theEmail := params["email"]
	enteredPassword := r.Form.Get("password")
	correctPassword := r.Form.Get("password")

	fmt.Println("EMAIL: ", theEmail)
	fmt.Println("PASSWORD: ", enteredPassword)

	sqlStatement := `SELECT password FROM donors WHERE email=$1;`

	row := db.QueryRow(sqlStatement, theEmail)
	err = row.Scan(&correctPassword)

	fmt.Println("CORRECT PASSWORD: ", correctPassword)

	result := false
	if enteredPassword == correctPassword {
		result = true
	}
	if err != nil {
		result = false
	}

	json.NewEncoder(w).Encode(result)

}
func CreateDonor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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

	sqlStatement := `SELECT email FROM donors WHERE email=$1;`
	row := db.QueryRow(sqlStatement, email)
	err = row.Scan(&theEmail)
	if err == nil {
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

		_, err := db.Exec(sqlStatement, email, name, gender, age, city, password)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			json.NewEncoder(w).Encode(true)
		}
	}
}
func DeleteDonor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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

	sqlStatement1 := `
	DELETE FROM donations
	WHERE donor_id = (	SELECT id
					FROM donors
					WHERE email = $1);`
	sqlStatement2 := `
	DELETE FROM donors
	WHERE email = $1;`
	fmt.Println("To be deleted:", theEmail)
	_, err1 := db.Exec(sqlStatement1, theEmail)
	_, err2 := db.Exec(sqlStatement2, theEmail)
	fmt.Println(err1.Error())
	fmt.Println(err2.Error())

}
func UpdateDonor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var id sql.NullInt64
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("ERR 1", err.Error())
		panic(err)
	}
	defer db.Close()

	params := mux.Vars(r)
	theId := params["id"]

	sqlStatement := `SELECT id FROM donors WHERE id=$1;`
	row := db.QueryRow(sqlStatement, theId)
	err = row.Scan(&id)
	if err != nil {
		fmt.Println("ERR 2", err.Error())
		fmt.Println("User Does Not Exist!")
	} else {
		r.ParseForm()
		email := r.Form.Get("email")
		name := r.Form.Get("name")
		gender := r.Form.Get("gender")
		age := r.Form.Get("age")
		city := r.Form.Get("city")
		password := r.Form.Get("password")

		sqlStatement := `
		UPDATE donors
		SET email=$1, name=$2, gender=$3, age=$4, city=$5, password=$6
		WHERE id=$7;`

		db.QueryRow(sqlStatement, email, name, gender, age, city, password, theId)
	}
}

func GetHospitals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
func hospitalLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	params := mux.Vars(r)
	r.ParseForm()
	theEmail := params["email"]
	enteredPassword := r.Form.Get("password")
	correctPassword := r.Form.Get("password")

	fmt.Println("EMAIL: ", theEmail)
	fmt.Println("PASSWORD: ", enteredPassword)

	sqlStatement := `SELECT password FROM hospitals WHERE email=$1;`

	row := db.QueryRow(sqlStatement, theEmail)
	err = row.Scan(&correctPassword)

	fmt.Println("CORRECT PASSWORD: ", correctPassword)

	result := false
	if enteredPassword == correctPassword {
		result = true
	}
	if err != nil {
		result = false
	}
	json.NewEncoder(w).Encode(result)
}
func CreateHospital(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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

	sqlStatement := `SELECT email FROM hospitals WHERE email=$1;`
	row := db.QueryRow(sqlStatement, email)
	err = row.Scan(&theEmail)
	if err == nil {
		fmt.Println("User Already Exist")
	} else {
		fmt.Println(err.Error())
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

		_, err := db.Exec(sqlStatement, email, name, city, password)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			json.NewEncoder(w).Encode(true)
		}
	}
}
func DeleteHospital(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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

	fmt.Println("To be deleted:", theEmail)
	_, err = db.Exec(sqlStatement, theEmail)
	fmt.Println(err.Error())
}
func UpdateHospital(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var id sql.NullInt64
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("ERR 1", err.Error())
		panic(err)
	}
	defer db.Close()

	params := mux.Vars(r)
	theId := params["id"]

	sqlStatement := `SELECT id FROM hospitals WHERE id=$1;`
	row := db.QueryRow(sqlStatement, theId)
	err = row.Scan(&id)
	if err != nil {
		fmt.Println("ERR 2", err.Error())
		fmt.Println("User Does Not Exist!")
	} else {
		r.ParseForm()
		email := r.Form.Get("email")
		name := r.Form.Get("name")
		city := r.Form.Get("city")
		password := r.Form.Get("password")

		sqlStatement := `
		UPDATE hospitals
		SET email=$1, name=$2, city=$3, password=$4
		WHERE id=$5;`

		db.QueryRow(sqlStatement, email, name, city, password, theId)
	}
}

func GetDonations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var donations []DonationFmt
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT donation.id, 
		donation.type, 
		donation.notes,
		donor.gender,
		donor.email,
		donor.age,
		donor.city
   FROM public.donations donation JOIN public.donors donor
	 ON donation.donor_id = donor.id
	 ORDER BY id;`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var donation Donation
		err = rows.Scan(&donation.Id, &donation.Type, &donation.Notes, &donation.Gender, &donation.Email, &donation.Age, &donation.City)
		if err != nil {
			panic(err)
		}
		donations = append(donations, DonationFmt{Id: donation.Id.Int64,
			Type:    donation.Type.String,
			Notes:   donation.Notes.String,
			DonorId: donation.DonorId.String,
			Gender:  donation.Gender.String,
			Email:   donation.Email.String,
			Age:     donation.Age.Int64,
			City:    donation.City.String})
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	fmt.Println(donations)
	json.NewEncoder(w).Encode(donations)
}
func GetDonorDonations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
	theId := params["id"]

	rows, err := db.Query(`SELECT donation.id, 
		donation.type, 
		donation.notes,
		donor.gender,
		donor.email,
		donor.age,
		donor.city
   FROM public.donations donation JOIN public.donors donor
	 ON donation.donor_id = donor.id
	 WHERE donor_id=$1 
	 ORDER BY id;`, theId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var donation Donation
		err = rows.Scan(&donation.Id, &donation.Type, &donation.Notes, &donation.Gender, &donation.Email, &donation.Age, &donation.City)
		if err != nil {
			panic(err)
		}
		donations = append(donations, DonationFmt{
			Id:      donation.Id.Int64,
			Type:    donation.Type.String,
			Notes:   donation.Notes.String,
			DonorId: donation.DonorId.String,
			Gender:  donation.Gender.String,
			Email:   donation.Email.String,
			Age:     donation.Age.Int64,
			City:    donation.City.String})
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	fmt.Println(donations)
	json.NewEncoder(w).Encode(donations)
}
func CreateDonation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
	donorEmail := r.Form.Get("donor_id")
	fmt.Println(theType, notes, donorEmail)

	sqlStatement := `
		INSERT INTO donations(
			type, notes, donor_id)
			VALUES  ($1, $2, $3)`

	db.QueryRow(sqlStatement, theType, notes, donorEmail)
}
func DeleteDonation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")

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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
