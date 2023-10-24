package repository



type SalaEstado int16

const (
	SalaAvailable SalaEstado = iota
	SalaUnAvailable
	SalaReserved 
)

type NotificationType int8

const (
	//"0"
	NotificationMessageGroup NotificationType = 0
	//"1"
	NotificationMessageComplejo = 1
	NotificationSalaCreation = 2
	NotificationSalaReservationConflict = 3
	NotificationSalaHasBeenReserved = 4
	NotificationBilling = 5
)




type ConsumoType int8

const (
	ReservaInstalacion ConsumoType = 0
	ReservaSala = 1
)