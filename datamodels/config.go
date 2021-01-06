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
	NbDrivers      int
	BaseTimer      int // Basde de temps
	SendPos        int // Nb de base de temps entre deux envois de position
	IdleDuration   int // Durée de la pause en BT
	PercentForIdle int // Pourcentage de chance de passer en Idle
	KmByBT         int // Nb de Km parcourus par BT
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
