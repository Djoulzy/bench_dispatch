package datamodels

// DataParams : Message generique
type DataParams interface{}

// JWTClaim : Contenue du token JWT
type JWTClaim struct {
	Exp      int    `json:"exp"`
	Iat      int    `json:"iat"`
	UserUuid string `json:"userUuid"`
}

// VehicleType : type de voiture du chauffeur
type VehicleType int

// Vehicule type list
const (
	Berline VehicleType = iota + 1
	Green
	Medical
	Other
	Prestige
	Van
)

// VehicleOption : Options list
type VehicleOption int

// Vehicule options list
const (
	CovidShield VehicleOption = iota + 1
	EnglishSpoken
	Mkids1
	Mkids2
	Mkids3
	Mkids4
	Pets
	Access
)

// Vehicle : Descriptif du véhicule du Driver
type Vehicle struct {
	ID           int    `mapstructure:"id" json:"id"`
	Brand        string `mapstructure:"brand" json:"brand"`
	Model        string `mapstructure:"model" json:"model"`
	Color        string `mapstructure:"color" json:"color"`
	Plate        string `mapstructure:"plate" json:"plate"`
	NumberOfSeat int    `mapstructure:"numberOfSeats" json:"numberOfSeats"`
}

// RideFlowType : Provenace de la demande
type RideFlowType int

// Origine de la demande
const (
	Defaut RideFlowType = iota
	LeTaxi
)

// Passenger : Détails du passager
type Passenger struct {
	ID        int    `mapstructure:"id" json:"id"`
	Picture   string `mapstructure:"picture" json:"picture"`
	Phone     string `mapstructure:"phone" json:"phone"`
	Firstname string `mapstructure:"firstname" json:"firstname"`
	Lastname  string `mapstructure:"lastname" json:"lastname"`
}

// Coordinates : geolocalisation
type Coordinates struct {
	Longitude float64 `mapstructure:"longitude" json:"longitude"`
	Latitude  float64 `mapstructure:"latitude" json:"latitude"`
}

// Address : adresse de prise en charge
type Address struct {
	Name  string      `mapstructure:"address" json:"address"`
	Coord Coordinates `mapstructure:"coordinates" json:"coordinates"`
}

// RideState : Status d'une course
type RideState int

// Etats des course
const (
	Pending RideState = iota + 1
	Booked
	Started
	Approach
	Delayed
	Waiting
	PickUpPassenger
	PendingPayment
	Ended
	Cancelled
)

// RideStats : Détails de la course pour payement
type RideStats struct {
	Value            float32 `mapstructure:"value" json:"value"`
	Unit             string  `mapstructure:"unit" json:"unit"`
	AdditionnalValue float32 `mapstructure:"additionnalValue" json:"additionnalValue"`
	Type             int     `mapstructure:"type" json:"type"`
}

type Proposal struct {
	SaveForMe   bool     `mapstructure:"saveForMe" json:"saveForMe"`
	ShareGroups []string `mapstructure:"shareGroups" json:"shareGroups"`
}

// Ride : modele de donnée pour une course
type RideData struct {
	ID             int64           `mapstructure:"id" json:"id"`
	Origin         RideFlowType    `mapstructure:"origin" json:"origin"`
	ExternalID     string          `mapstructure:"externalId" json:"externalId"`
	Memo           string          `mapstructure:"memo" json:"memo"`
	Reference      string          `mapstructure:"reference" json:"reference"`
	StartDate      string          `mapstructure:"date" json:"startDate"`
	State          RideState       `mapstructure:"state" json:"state"`
	ToAddress      Address         `mapstructure:"toAddress" json:"toAddress"`
	ValidUntil     string          `mapstructure:"validUntil" json:"validUntil"`
	IsImmediate    bool            `mapstructure:"isImmediate" json:"isImmediate"`
	FromAddress    Address         `mapstructure:"fromAddress" json:"fromAddress"`
	NbLuggages     int             `mapstructure:"numberOfLuggages" json:"numberOfLuggages"`
	NbPassengers   int             `mapstructure:"numberOfPassengers" json:"numberOfPassengers"`
	VehicleType    VehicleType     `mapstructure:"vehicleType" json:"vehicleType"`
	Vehicle        Vehicle         `mapstructure:"vehicle" json:"vehicle"`
	PickUpAddress  Address         `mapstructure:"pickUpAddress" json:"pickUpAddress"`
	Passenger      Passenger       `mapstructure:"passenger" json:"passenger"`
	Stats          []RideStats     `mapstructure:"stats" json:"stats"`
	VatValue       float32         `mapstructure:"vatValue" json:"vatValue"`
	VehicleOptions []VehicleOption `mapstructure:"vehicleOptions" json:"vehicleOptions"`
}

// DriverState : Etat du chauffeur
type DriverState int

// Liste des etats du chauffeur
const (
	Free DriverState = iota
	Occupied
	Offline

	Moving
	WaitOK
	WaitACK
	Payment
	Err
)
