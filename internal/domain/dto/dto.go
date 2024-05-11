package dto

// PersonDTO represents the person DTO.
type PersonDTO struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Number  string `json:"phone_number"`
	City    string `json:"city"`
	State   string `json:"state"`
	Street1 string `json:"street1"`
	Street2 string `json:"street2"`
	Zip     string `json:"zip_code"`
}
