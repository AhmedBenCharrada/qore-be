package entities

// Person represents the person entity.
type Person struct {
	ID   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// TableName ..
func (Person) TableName() string {
	return "person"
}

// Phone represents the phone entity.
type Phone struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	PersonID int    `json:"person_id"`
	Number   string `json:"number"`
}

// TableName ..
func (Phone) TableName() string {
	return "phone"
}

// Address represents the address entity.
type Address struct {
	ID      int    `json:"id" gorm:"primaryKey;autoIncrement"`
	City    string `json:"city"`
	State   string `json:"state"`
	Street1 string `json:"street1"`
	Street2 string `json:"street2"`
	Zip     string `json:"zip_code"`
}

// TableName ..
func (Address) TableName() string {
	return "address"
}

// PersonAddress join table for person and address.
type PersonAddress struct {
	ID        int `json:"id" gorm:"primaryKey;autoIncrement"`
	PersonID  int `json:"person_id"`
	AddressID int `json:"address_id"`
}

// TableName ..
func (PersonAddress) TableName() string {
	return "address_join"
}
