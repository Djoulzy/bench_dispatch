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

// Vehicule type list
const (
	Green    VehicleType = "GREEN"
	Berline  VehicleType = "BERLINE"
	Van      VehicleType = "VAN"
	Prestige VehicleType = "PRESTIGE"
	Medical  VehicleType = "MEDICAL"
	Other    VehicleType = "OTHER"
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

// Ride : modele de donnée pour une course
type Ride struct {
	Origin      RideFlowType `json:"origin"`
	ID          string       `json:"id"`
	Date        string       `json:"date"`
	State       RideState    `json:"state"`
	ToAddress   AddressRide  `json:"toAddress"`
	ValidUntil  string       `json:"validUntil"`
	IsImmediate bool         `json:"isImmediate"`
	FromAddress AddressRide  `json:"fromAddress"`
	Options     OptionsRide  `json:"options"`
}

/*

////////// UpdateDriverLocation //////////

{
   "status" : {
      "errorMessage" : "status",
      "errorCode" : 987
   },
   "method" : "UpdateDriverLocation",
   "id" : 1,
   "params" : {
      "longitude" : 4.9878669999999996,
      "latitude" : 43.987867000000001
   }
}

////////// NewRide //////////
{
   "id":1,
   "method":"NewRide",
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

////////// AcceptRide //////////
{
   "status" : {
      "errorMessage" : "status",
      "errorCode" : 987
   },
   "method" : "AcceptRide",
   "id" : 1,
   "params" : {
      "rideId" : "8976987986"
   }
}

////////// AcceptRideResponse //////////
{
   "status" : {
      "errorMessage" : "status",
      "errorCode" : 987
   },
   "method" : "AcceptRideResponse",
   "id" : 1,
   "params" : {
      "options" : {
      "numberOfPassengers" : 2,
      "numberOfLuggages" : 1
      },
      "isImmediate" : true,
      "origin" : 0,
      "id" : "8976987986",
      "validUntil" : "2020-12-16T17:20:00.000Z",
      "date" : "2020-12-16T17:20:00.000Z",
      "toAddress" : {
      "coordinates" : {
         "longitude" : 5.4925626895974045,
         "latitude" : 43.471590283851015
      },
      "address" : "arrivée adresse 14019 Ubcqoeu"
      },
      "fromAddress" : {
      "coordinates" : {
         "longitude" : 5.4925626895974045,
         "latitude" : 43.471590283851015
      },
      "address" : "départ adresse 14019 Ubcqoeu"
      }
   }
}

////////// StartRide //////////
{
   "status" : {
      "errorMessage" : "status",
      "errorCode" : 987
   },
   "method" : "StartRide",
   "id" : 1,
   "params" : {
      "memo" : "Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Nulla vitae elit libero, a pharetra augue.",
      "reference" : "Chambre 208A",
      "ride" : {
      "options" : {
         "numberOfPassengers" : 2,
         "numberOfLuggages" : 1
      },
      "state" : 0,
      "origin" : 0,
      "isImmediate" : true,
      "id" : "8976987986",
      "validUntil" : "2020-12-16T17:20:00.000Z",
      "date" : "2020-12-16T17:20:00.000Z",
      "toAddress" : {
         "coordinates" : {
            "longitude" : 5.452597832139817,
            "latitude" : 43.52645372148015
         },
         "address" : "arrivée adresse 14019 Ubcqoeu"
      },
      "fromAddress" : {
         "coordinates" : {
            "longitude" : 5.53859787072443,
            "latitude" : 43.47865284174063
         },
         "address" : "départ adresse 14019 Ubcqoeu"
      }
      },
      "passenger" : {
      "picture" : "https:\/\/media-exp1.licdn.com\/dms\/image\/C4D03AQHAPi4WceJ6rA\/profile-displayphoto-shrink_400_400\/0\/1516561939955?e=1614211200&v=beta&t=Mk1eA5tDgOODt3V9cLqITaWj9TAelHZTHDAFXVhx4vE",
      "phone" : "+330987654321",
      "firstname" : "Jérôme",
      "id" : "8976987986",
      "lastname" : "TONNELIER"
      }
   }
}

////////// ChangeRideState //////////
{
   "status" : {
      "errorMessage" : "status",
      "errorCode" : 987
   },
   "method" : "ChangeRideState",
   "id" : 1,
   "params" : {
      "rideId" : "8976987986",
      "state" : 1
   }
}

////////// RideStateChanged //////////
{
   "status" : {
      "errorMessage" : "status",
      "errorCode" : 987
   },
   "method" : "RideStateChanged",
   "id" : 1,
   "params" : {
      "memo" : "Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Nulla vitae elit libero, a pharetra augue.",
      "reference" : "Chambre 208A",
      "ride" : {
      "options" : {
         "numberOfPassengers" : 2,
         "numberOfLuggages" : 1
      },
      "state" : 0,
      "origin" : 0,
      "isImmediate" : true,
      "id" : "8976987986",
      "validUntil" : "2020-12-16T17:20:00.000Z",
      "date" : "2020-12-16T17:20:00.000Z",
      "toAddress" : {
         "coordinates" : {
            "longitude" : 5.452597832139817,
            "latitude" : 43.52645372148015
         },
         "address" : "arrivée adresse 14019 Ubcqoeu"
      },
      "fromAddress" : {
         "coordinates" : {
            "longitude" : 5.53859787072443,
            "latitude" : 43.47865284174063
         },
         "address" : "départ adresse 14019 Ubcqoeu"
      }
      },
      "passenger" : {
      "picture" : "https:\/\/media-exp1.licdn.com\/dms\/image\/C4D03AQHAPi4WceJ6rA\/profile-displayphoto-shrink_400_400\/0\/1516561939955?e=1614211200&v=beta&t=Mk1eA5tDgOODt3V9cLqITaWj9TAelHZTHDAFXVhx4vE",
      "phone" : "+330987654321",
      "firstname" : "Jérôme",
      "id" : "8976987986",
      "lastname" : "TONNELIER"
      }
   }
}
*/
