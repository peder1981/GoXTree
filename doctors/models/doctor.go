package models

type Doctor struct {
    ID                int64  `json:"id"`
    CRM               string `json:"crm"`
    UF                string `json:"uf"`
    DoctorName        string `json:"doctor_name"`
    TipoDeTalonario   string `json:"tipo_de_talonario"`
    SequencialInicial string `json:"sequencial_inicial"`
    Molde             []byte `json:"molde"`
}
