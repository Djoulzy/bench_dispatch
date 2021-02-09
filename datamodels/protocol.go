package datamodels

// DataParams : Message generique
type DataParams interface{}

// Request : Format d'une requete
type Request struct {
	ID     int        `json:"id"`
	Method string     `json:"method"`
	Params DataParams `json:"params"`
}

// Response : Format d'une reponse
type Response struct {
	ID     int        `json:"id"`
	Method string     `json:"method"`
	Params DataParams `json:"params"`
	Status Error      `json:"status"`
}

// Error : gestion d'une erreur de requete
type Error struct {
	ID      int    `json:"errorCode"`
	Message string `json:"errorMessage"`
}

// Success : Simplification d'ecriture en cas de succes
var Success = Error{
	ID:      0,
	Message: "OK",
}

// Login : Structure d'info pour le login Driver
type Login struct {
	Token string      `json:"token"`
	State DriverState `json:"state"`
	ID    int         `json:"id"`
	Name  string      `json:"name"`
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

// Origine de la demande
const (
	Defaut RideFlowType = iota
	LeTaxi
)

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

// OptionsRide : Options à ajouter à la course
type OptionsRide struct {
	Luggages   int         `mapstructure:"numberOfLuggages" json:"numberOfLuggages"`
	Passengers int         `mapstructure:"numberOfPassengers" json:"numberOfPassengers"`
	Vehicle    VehicleType `mapstructure:"vehicleType" json:"vehicleType"`
}

// AcceptRide : Message du chauffeur pour accepter la course
type AcceptRide struct {
	RideID string `mapstructure:"rideId" json:"rideId"`
}

// RideState : Status d'une course
type RideState int

// Etats des course
const (
	Pending RideState = iota
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

// Ride : modele de donnée pour une course
type Ride struct {
	Origin      RideFlowType `mapstructure:"origin" json:"origin"`
	ID          string       `mapstructure:"id" json:"id"`
	Date        string       `mapstructure:"date" json:"date"`
	State       RideState    `mapstructure:"state" json:"state"`
	ToAddress   Address      `mapstructure:"toAddress" json:"toAddress"`
	ValidUntil  string       `mapstructure:"validUntil" json:"validUntil"`
	IsImmediate bool         `mapstructure:"isImmediate" json:"isImmediate"`
	FromAddress Address      `mapstructure:"fromAddress" json:"fromAddress"`
	Options     OptionsRide  `mapstructure:"options" json:"options"`
}

// RideUpdate : Modifie l'état de la course
type RideUpdate struct {
	ID    string    `mapstructure:"rideId" json:"rideId"`
	State RideState `mapstructure:"state" json:"state"`
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
	Err
)

// DriverStateChange : Changement du status d'un Driver
type DriverStateChange struct {
	State DriverState `mapstructure:"state" json:"state"`
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
         "address":"250 plage de l'Estaque 13016 Marseille",
         "coordinates":{
            "longitude":5.2934914462,
            "latitude":43.3586309109
         }
      },
      "validUntil":"",
      "options":{
         "numberOfLuggages":1,
         "numberOfPassengers":2
      },
      "fromAddress":{
         "address":"Quai de Rive Neuve 13007 Marseille",
         "coordinates":{
            "longitude":5.36785419454433,
            "latitude":43.2924901379708
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
