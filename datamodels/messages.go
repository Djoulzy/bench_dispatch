package datamodels

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

/////////////////////////////////////////
/////////// OUTGOING MESSAGES ///////////
/////////////////////////////////////////

// Login : Structure d'info pour le login Driver
// Response is error code
type Login struct {
	Token string         `json:"token"`
	State TaximeterState `json:"state"`
	ID    int            `json:"id"`
	Name  string         `json:"name"`
}

// CreateRide : Proposition de course
type CreateRide struct {
	Ride          RideData      `mapstructure:"ride" json:"ride"`
	SearchOptions SearchOptions `mapstructure:"searchOptions" json:"searchOptions"`
	Passenger     Passenger     `mapstructure:"passenger" json:"passenger"`
	Proposal      Proposal      `mapstructure:"proposal" json:"proposal"`
}

// RideUpdate : Modifie l'état de la course
// Response is the same with error code
type ChangeRideState struct {
	ID    int64     `mapstructure:"rideId" json:"rideId"`
	State RideState `mapstructure:"state" json:"state"`
}

// UpdateDriverLocation: Mise à jour de la position du chauffeur
type UpdateDriverLocation struct {
	Coord          Coordinates     `mapstructure:"coordinates" json:"coordinates"`
	VehicleOptions []VehicleOption `mapstructure:"vehicleOptions" json:"vehicleOptions"`
	VehicleType    VehicleType     `mapstructure:"vehicleType" json:"vehicleType"`
}

// AcceptRide : Message du chauffeur pour accepter la course
type AcceptRide struct {
	ID      int64   `mapstructure:"rideId" json:"rideId"`
	Vehicle Vehicle `mapstructure:"vehicle" json:"vehicle"`
}

// ChangeTaximeterState : Changement du status d'un Driver
// Response is the same with error code
type ChangeTaximeterState struct {
	State TaximeterState `mapstructure:"state" json:"state"`
}

/////////////////////////////////////////
/////////// INCOMING MESSAGES ///////////
/////////////////////////////////////////

// MonitorConfig : Structure d'info pour le login Driver
type MonitorConfig struct {
	Config Globals `json:"config"`
}

// Payment : Payement d'une course
type PendingPaymentResponse struct {
	Ride               RideData `mapstructure:"ride" json:"ride"`
	PickUpAddress      Address  `mapstructure:"pickUpAddress" json:"pickUpAddress"`
	Payment            Payment  `mapstructure:"payment" json:"payment"`
	CancellationReason string   `mapstructure:"cancellationReason" json:"cancellationReason"`
}

// AcceptRideResponse : Retour pour course acceptée
type AcceptRideResponse struct {
	Ride          RideData      `mapstructure:"ride" json:"ride"`
	Passenger     Passenger     `mapstructure:"passenger" json:"passenger"`
	Vehicle       Vehicle       `mapstructure:"vehicle" json:"vehicle"`
	SearchOptions SearchOptions `mapstructure:"searchOptions" json:"searchOptions"`
}

/*

{
   "method" : "Login",
   "id" : 1,
   "params" : {
      "token" : "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTIyMDkwMDksImlhdCI6MTYxMjIwNjAwOSwidXNlclV1aWQiOiIwMDQwYjUzMi1lMmM1LTEwMzktOGNiMC02NWFkMTQ1ODZjMTUifQ.rl237_kl6-lj2nsNth5oBZ4fvW2UuapbdmW2NhmxPAaOJDTEcObtjHxvxuo0VxO6EvmnMa-lQs9JpA2Zn7ZfGqripx3zUYyHWrOgjL9zKLfy0QOb7NqXqwryn2HiMgqXmd0ZpwrXNjFeSr2jBZT2BWslWIO_oN3fJpFiORtf8384y6SvjjquZO4Jkwv8m44fDJyKXRFIq-koQJh5nAHj0dP7LAwEpBMFMf_6pnzUqOMvzNfVyEtmnKuK6jwSxqy98IMCJjp2UiCitjGIU88_yHJA5ZLAAmOj1yfKUJeNNtDVUdkdGTrGCaBAIgHHSBvRdk4X4M4079AcfFuerIw9yQ"
   }
}

{
   "status" : {
      "errorMessage" : "OK",
      "errorCode" : 0
   },
   "method" : "LoginResponse",
   "id" : 1,
   "params" : {
      "id" : 12345,
      "name" : "Joe le Taxi"
   }
}

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
    "id":256,
    "method":"NewRide",
    "params":{

        "ride": {
            "id" : 987,
            "externalId": "skhdf455",
            "origin" : "BOOKER",
            "state" : 1,
            "memo" : "mémo optionnel",
            "reference" : "reference optionnelle",
            "isImmediate" : true,
            "startDate" : "2020-12-16T17:20:00.00Z",
            "validUntil" : "2020-12-16T17:20:00.00Z",
            "fromAddress" : {
                "name" : "nom cours de l'adresse",
                "address" : "adresse complete",
                "coordinates" : {
                    "latitude" : 42.9867,
                    "longitude" : 4.9867
                }
            },
            "toAddress" : {
                "name" : "nom cours de l'adresse optionelle",
                "address" : "adresse complete optionelle",
                "coordinates" : {
                    "latitude" : 42.9867,
                    "longitude" : 4.9867
                }
            },
            "numberOfPassengers" : 1,
            "numberOfLuggages" : 0,
            "vehicleOptions" : [1, 3, 5],
            "vehicleType" : 2,

            "passenger" : {
                "id" : 89767,
                "firstname" : "prénom",
                "lastname" : "lastname",
                "phone" : "phone",
                "picture" : "url de la photo optionnelle"
            },

            "vehicle" : {
                "id" : 8976,
                "brand" : "BMW",
                "model" : "Série 3",
                "vehicleType" : 2,
                "color" : "WHITE",
                "plate" : "TY-496-CZ",
                "numberOfSeats" : 6
            },

            "pickUpAddress" :  {
                "name" : "nom cours de l'adresse optionelle",
                "address" : "adresse complete optionelle",
                "coordinates" : {
                    "latitude" : 42.9867,
                    "longitude" : 4.9867
                }
            },

            "vatValue" : 20.0,
            "stats" : [
                {
                    "value" : 0.0,
                    "unit" : "€",
                    "type" : 0,
                    "additionnalValue": 10.0
                },
                {
                    "value" : 25.0,
                    "unit" : "km",
                    "type" : 1
                },
                {
                    "value" : 40.0,
                    "unit" : "min",
                    "type" : 2
                }
            ]
        },

        "proposal": {
            "saveForMe": true,
            "shareGroups": "Nom du groupe"
        }
    },

    "status" : {
        "errorMessage" : "status",
        "errorCode" : 987
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
      "rideId" : 8976987986
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
      "memo" : "Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Nulla vitae elit libero, a pharetra augue.",
      "reference" : "Chambre 208A",
      "ride" : {
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

////////// StartOngoingRide //////////
{
   "status" : {
      "errorMessage" : "status",
      "errorCode" : 987
   },
   "method" : "StartOngoingRide",
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
      "id" : 12,
      "lastname" : "TONNELIER"
      }
   }
}

////////// ChangeRideStateResponse //////////

{
   "status" : {
      "errorMessage" : "status",
      "errorCode" : 987
   },
   "method" : "ChangeRideStateResponse",
   "id" : 1,
   "params" : {
      "rideId" : 8976987986,
      "state" : 1
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
      "rideId" : 8976987986,
      "state" : 1
   }
}

////////// PendingPaymentResponse //////////

{
   "status" : {
      "errorMessage" : "status",
      "errorCode" : 987
   },
   "method" : "PendingPaymentResponse",
   "id" : 1,
   "params" : {
      "ride" : {
         "options" : {
            "numberOfPassengers" : 9,
            "numberOfLuggages" : 5
         },
         "state" : 7,
         "origin" : 0,
         "isImmediate" : true,
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
      },
      "pickUpAddress" : {
         "coordinates" : {
            "longitude" : 5.53858787072443,
            "latitude" : 43.46865284174063
         },
         "address" : "Acception de la course address"
      },
      "stats" : [
         {
            "value" : 0.0,
            "unit" : "€",
            "type" : 0,
            "additionnalValue" = 10.0
         },
         {
               "value" : 25.0,
               "unit" : "km",
               "type" : 1
         },
         {
               "value" : 40.0,
               "unit" : "min",
               "type" : 2
         }
      ]
   }
}

////////// Ended //////////

{
   "status" : {
      "errorMessage" : "status",
      "errorCode" : 987
   },
   "method" : "RideEndedResponse",
   "id" : 1,
   "params" : {
      "rideId" : "JHGCUYGC-HCS"
   }
}


Option Véhicule
Paroi de séparation COVID
Chauffeur parlant Anglais
MKIDS : 1 réhausseur
MKIDS : 2 réhausseurs
MKIDS : 1 Siège bébé
MKIDS : 1 Siège bébé+1réhausseur
Animal accepté.
-ACCESS (personne à mobilité réduite).

case cpam = 1, covidShield, englishSpoken, mkids1, mkids2, mkids3, mkids4, pets, access

{
  "method" : "CreateRide",
  "id" : 1,
  "params" : {
    "start" : { "coordinates" : { "latitude" : 42.2324, "longitude" : 42.2324 }, "address" : "l'adresse en texte" },
    "end" : { "coordinates" : { "latitude" : 42.2324, "longitude" : 42.2324 }, "address" : "l'adresse d'arrivée est optionnelle" },
    "vehicleType" : 1,
    "options" : [1, 4],
    "date" : "2020-12-16T17:20:00.000Z",
    "shareGroups" : [123, 98327],
    "driverId" : 87586
        var passengerName: String?
    var passengerPhone: String?
  }
}
*/
