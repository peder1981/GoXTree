package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strconv"

    "doctors/models"
    "doctors/db"

    "github.com/gorilla/mux"
)

// CreateDoctor creates a new doctor
func CreateDoctor(w http.ResponseWriter, r *http.Request) {
    var doctor models.Doctor
    json.NewDecoder(r.Body).Decode(&doctor)

    err := db.DB.QueryRow(
        "INSERT INTO doctors (crm, uf, doctor_name, tipo_de_talonario, sequencial_inicial, molde) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
        doctor.CRM, doctor.UF, doctor.DoctorName, doctor.TipoDeTalonario, doctor.SequencialInicial, doctor.Molde,
    ).Scan(&doctor.ID)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(doctor)
}

// GetDoctor retrieves a doctor by ID
func GetDoctor(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, _ := strconv.Atoi(params["id"])

    var doctor models.Doctor
    err := db.DB.QueryRow("SELECT id, crm, uf, doctor_name, tipo_de_talonario, sequencial_inicial, molde FROM doctors WHERE id = $1", id).Scan(
        &doctor.ID, &doctor.CRM, &doctor.UF, &doctor.DoctorName, &doctor.TipoDeTalonario, &doctor.SequencialInicial, &doctor.Molde,
    )

    if err == sql.ErrNoRows {
        http.Error(w, "Doctor not found", http.StatusNotFound)
        return
    } else if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(doctor)
}

// UpdateDoctor updates an existing doctor
func UpdateDoctor(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, _ := strconv.Atoi(params["id"])

    var doctor models.Doctor
    json.NewDecoder(r.Body).Decode(&doctor)

    _, err := db.DB.Exec(
        "UPDATE doctors SET crm = $1, uf = $2, doctor_name = $3, tipo_de_talonario = $4, sequencial_inicial = $5, molde = $6 WHERE id = $7",
        doctor.CRM, doctor.UF, doctor.DoctorName, doctor.TipoDeTalonario, doctor.SequencialInicial, doctor.Molde, id,
    )

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(doctor)
}

// DeleteDoctor deletes a doctor by ID
func DeleteDoctor(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    id, _ := strconv.Atoi(params["id"])

    _, err := db.DB.Exec("DELETE FROM doctors WHERE id = $1", id)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
