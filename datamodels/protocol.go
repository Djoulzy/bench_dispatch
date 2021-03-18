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

// Login : Structure d'info pour le login Driver
type Login struct {
	Token string      `json:"token"`
	State DriverState `json:"state"`
	ID    int         `json:"id"`
	Name  string      `json:"name"`
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

// OptionsRide : Options à ajouter à la course
type OptionsRide struct {
	Luggages   int         `mapstructure:"numberOfLuggages" json:"numberOfLuggages"`
	Passengers int         `mapstructure:"numberOfPassengers" json:"numberOfPassengers"`
	Vehicle    VehicleType `mapstructure:"vehicleType" json:"vehicleType"`
}

// AcceptRide : Message du chauffeur pour accepter la course
type AcceptRide struct {
	ID int64 `mapstructure:"rideId" json:"rideId"`
}

// AcceptRideResponse : Retour pour course acceptée
type AcceptRideResponse struct {
	Memo      string    `mapstructure:"memo" json:"memo"`
	Reference string    `mapstructure:"reference" json:"reference"`
	Ride      Ride      `mapstructure:"ride" json:"ride"`
	Passenger Passenger `mapstructure:"passenger" json:"passenger"`
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

// Ride : modele de donnée pour une course
type Ride struct {
	ID          int64        `mapstructure:"id" json:"id"`
	Origin      RideFlowType `mapstructure:"origin" json:"origin"`
	ExternalID  string       `mapstructure:"externalId" json:"externalId"`
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
	ID    int64     `mapstructure:"rideId" json:"rideId"`
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
	Payment
	Err
)

// DriverStateChange : Changement du status d'un Driver
type DriverStateChange struct {
	State DriverState `mapstructure:"state" json:"state"`
}

// RideStats : Détails de la course pour payement
type RideStats struct {
	Value            float32 `mapstructure:"value" json:"value"`
	Unit             string  `mapstructure:"unit" json:"unit"`
	AdditionnalValue float32 `mapstructure:"additionnalValue" json:"additionnalValue"`
	Type             int     `mapstructure:"type" json:"type"`
}

// Payment : Payement d'une course
type PaymentResponse struct {
	Ride          Ride        `mapstructure:"ride" json:"ride"`
	PickUpAddress Address     `mapstructure:"pickUpAddress" json:"pickUpAddress"`
	Stats         []RideStats `mapstructure:"stats" json:"stats"`
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
