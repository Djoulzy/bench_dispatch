package datamodels

// Globals : Partie globale du fichier de conf
type Globals struct {
	ServerID     int
	LogLevel     int
	StartLogging bool
	MaxUsers     int
	Workers      int
	QueueSize    int
}

// Bench : Parametre des tests
type Bench struct {
	NbClients int
}

// WSserver : Configuration des servers
type WSserver struct {
	Addr string
}

// RideConfig : paramètres d'une course
type RideConfig struct {
	TimeBeetwinSteps int
}

// ConfigData : Data structure du fichier de conf
type ConfigData struct {
	Globals
	Bench
	WSserver
	RideConfig
}
