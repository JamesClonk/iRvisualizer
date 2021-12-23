package web

import "strconv"

func isDriverMarked(drivers []string, driverID int) bool {
	for _, driver := range drivers {
		if strconv.Itoa(driverID) == driver {
			return true
		}
	}
	return false
}
