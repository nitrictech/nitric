package service

type Status int

const (
	Status_Init Status = iota
	Status_Starting
	Status_Restarting
	Status_Running
	Status_Stopping
	Status_Stopped
	Status_Fatal
)
