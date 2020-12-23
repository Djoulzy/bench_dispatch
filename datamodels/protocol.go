package datamodels

// DataParams : Message generique
type DataParams interface{}

// Request : Format d'une requete
type Request struct {
	ID     int        `json:"id"`
	Method string     `json:"method"`
	Params DataParams `json:"params"`
	Status Error      `json:"status"`
}

// Response : Format d'une reponse
type Response struct {
	ID     int        `json:"id"`
	Result DataParams `json:"result"`
}

// Error : gestion d'une erreur de requete
type Error struct {
	ID      int    `json:"errorCode"`
	Message string `json:"errorMessage"`
}

// VehicleType : type de voiture du chauffeur
type VehicleType string

const (
	green    VehicleType = "GREEN"
	berline  VehicleType = "BERLINE"
	van      VehicleType = "VAN"
	prestige VehicleType = "PRESTIGE"
	medical  VehicleType = "MEDICAL"
	other    VehicleType = "OTHER"
)

// RideFlowType : Provenace de la demande
type RideFlowType int

const (
	defaut RideFlowType = iota
	leTaxi
)

// Coordinates : geolocalisation
type Coordinates struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

// AddressRide : adresse de prise en charge
type AddressRide struct {
	Address string      `json:"address"`
	Coord   Coordinates `json:"coordinates"`
}

// OptionsRide : Options à ajouter à la course
type OptionsRide struct {
	Luggages   int         `json:"numberOfLuggages"`
	Passengers int         `json:"numberOfPassengers"`
	Vehicle    VehicleType `json:"vehicleType"`
}

// AcceptRide : Message du chauffeur pour accepter la course
type AcceptRide struct {
	RideID string `json:"rideId"`
}

// RideState : Status d'une course
type RideState int

const (
	started RideState = iota
	approach
	pickUpPassenger
	ended
	cancelled
	booked
	pending
)

/*
{
   "id":1,
   "method":"new_ride",
   "params":{
      "origin":0,
      "id":"7698765",
      "date":"",
      "toAddress":{
         "address":"arrivee adresse 14019 Ubcqoeu",
         "coordinates":{
            "longitude":5.37436,
            "latitude":43.29539
         }
      },
      "validUntil":"",
      "options":{
         "numberOfLuggages":1,
         "numberOfPassengers":2
      },
      "fromAddress":{
         "address":"depart adresse 14019 Ubcqoeu",
         "coordinates":{
            "longitude":5.4925626895974045,
            "latitude":43.471590283851015
         }
      },
      "isImmediate":true
   }
}
*/
