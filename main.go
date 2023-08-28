package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ServiceWeaver/weaver"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type JSONReturn struct {
	Id        int       `json:"id"`
	Nama      string    `json:"nama,omitempty"`
	Kondisi   bool      `json:"kondisi"`
	Tanggal   time.Time `json:"tanggal"`
	Timestamp time.Time `json:"timestamp"`
	Email     string    `json:"email,omitempty"`
}

type Test struct {
	Nama, Email        string
	Kondisi            bool
	Tanggal, Timestamp time.Time
	Id                 int
}

type Response struct {
	Pesan  string `json:"pesan"`
	Status int    `json:"status"`
}

func main() {
	if err := weaver.Run(context.Background(), run); err != nil {
		log.Fatal(err)
	}
}

type app struct {
	weaver.Implements[weaver.Main]
	listen weaver.Listener `weaver:"API_CRUD"`
}

func run(ctx context.Context, app *app) error {
	fmt.Printf("Listener alamat %s:", app.listen)

	r := mux.NewRouter()
	http.HandleFunc("/read", read)
	r.HandleFunc("/insert", insert).Methods("POST")
	r.HandleFunc("/delete/{id:[0-9]+}", delete).Methods("DELETE")
	r.HandleFunc("/update", update).Methods("PUT")
	http.Handle("/", r)

	return http.Serve(app.listen, nil)
}

func update(w http.ResponseWriter, r *http.Request) {
	// Konek DB
	dbQuery := "root:admin@tcp(localhost:3306)/golang?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dbQuery), &gorm.Config{})
	if err != nil {
		log.Fatal("Query Failed", err)
	}

	// Variabel Form
	id := r.FormValue("idEdit")
	nama := r.FormValue("namaEdit")
	email := r.FormValue("emailEdit")
	tanggal := r.FormValue("tglEdit")
	kondisi := r.FormValue("kondisiEdit")

	// Convert Id
	idConv, errIdConv := strconv.Atoi(id)
	if errIdConv != nil {
		http.Error(w, "Error Konversi Id"+errIdConv.Error(), http.StatusInternalServerError)
	}

	// Convert Bool
	kondisiBool, errKondisiBool := strconv.ParseBool(kondisi)
	if errKondisiBool != nil {
		http.Error(w, "Error Konversi Bool"+errKondisiBool.Error(), http.StatusInternalServerError)
	}

	fmt.Println(kondisiBool)

	// Convert tanggal
	layout := "2006-01-02"
	tglConv, errTglConv := time.Parse(layout, tanggal)
	if errTglConv != nil {
		http.Error(w, "Error Konversi Tanggal"+errTglConv.Error(), http.StatusInternalServerError)
	}

	// Interface
	testUpdate := Test{
		Nama:      nama,
		Email:     email,
		Tanggal:   tglConv,
		Kondisi:   kondisiBool,
		Timestamp: time.Now(),
	}

	// Exec DB
	db.Table("test").Where("id=?", idConv).Updates(&testUpdate)

	// Pesan JSON
	pesan := Response{
		Pesan:  "Berhasil Update",
		Status: 200,
	}

	// Return JSON
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	errEncoder := encoder.Encode(pesan)
	if errEncoder != nil {
		http.Error(w, "Encoder Error"+errEncoder.Error(), http.StatusInternalServerError)
		return
	}
}

func delete(w http.ResponseWriter, r *http.Request) {
	// Konek DB
	dbQuery := "root:admin@tcp(localhost:3306)/golang?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dbQuery), &gorm.Config{})
	if err != nil {
		log.Fatal("Query Failed", err)
	}

	// Get Variabel ID
	vars := mux.Vars(r)
	id := vars["id"]

	// Convert Id to int
	idConv, _ := strconv.Atoi(id)

	// Exec DB
	deleteUser := []Test{{Id: idConv}}
	db.Table("test").Delete(&deleteUser)

	// Pesan JSON
	pesanJson := Response{
		Pesan:  "Berhasil Hapus",
		Status: 200,
	}

	// Return JSON
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	errEncoder := encoder.Encode(pesanJson)
	if errEncoder != nil {
		http.Error(w, "Encoder Error"+errEncoder.Error(), http.StatusInternalServerError)
		return
	}
}

func insert(w http.ResponseWriter, r *http.Request) {

	// Konek DB
	dbQuery := "root:admin@tcp(localhost:3306)/golang?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dbQuery), &gorm.Config{})
	if err != nil {
		log.Fatal("Query Failed", err)
	}

	// Get Form
	nama := r.FormValue("nama")
	tanggal := r.FormValue("tanggal")
	email := r.FormValue("email")
	kondisi := r.FormValue("kondisi")

	// Convert Kondisi
	kondisiBool, errConvertKondisi := strconv.ParseBool(kondisi)
	if errConvertKondisi != nil {
		http.Error(w, "Error Convert"+errConvertKondisi.Error(), http.StatusInternalServerError)
	}

	// Convert tanggal
	layout := "2006-01-02"
	tanggalTime, errConvertTanggal := time.Parse(layout, tanggal)
	if errConvertTanggal != nil {
		http.Error(w, "Error Convert"+errConvertTanggal.Error(), http.StatusInternalServerError)
	}

	test := Test{
		Nama:      nama,
		Kondisi:   kondisiBool,
		Tanggal:   tanggalTime,
		Email:     email,
		Timestamp: time.Now(),
	}

	// Exec DB
	db.Table("test").Create(&test)

	// Pesan JSON
	pesan := Response{
		Pesan:  "Berhasil",
		Status: 200,
	}

	// Return JSON
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	errEncoder := encoder.Encode(pesan)
	if errEncoder != nil {
		http.Error(w, "Encoder Error"+errEncoder.Error(), http.StatusInternalServerError)
		return
	}
}

func read(w http.ResponseWriter, r *http.Request) {
	// Konek DB
	dbQuery := "root:admin@tcp(localhost:3306)/golang?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dbQuery), &gorm.Config{})
	if err != nil {
		log.Fatal("Query Failed", err)
	}

	var test []Test
	db.Table("test").Find(&test)
	if db.Error != nil {
		fmt.Println("Error", db.Error)
		return
	}

	// JSON
	var kembaliJSON []JSONReturn

	// Print Value
	for _, row := range test {

		kembaliJSON = append(kembaliJSON,
			JSONReturn{
				Id:        row.Id,
				Nama:      row.Nama,
				Email:     row.Email,
				Tanggal:   row.Tanggal,
				Timestamp: row.Timestamp,
				Kondisi:   row.Kondisi,
			})
	}

	// Return JSON
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	errEncoder := encoder.Encode(kembaliJSON)
	if errEncoder != nil {
		http.Error(w, "Encoder Error"+errEncoder.Error(), http.StatusInternalServerError)
		return
	}
}

func selectSpesifik(w http.ResponseWriter, r *http.Request) {

}
