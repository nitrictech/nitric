package service

type Status int

const (
	Status_Init Status = iota
	Status_Start
	Status_Restart
	Status_Stop
	Status_Fatal
)
