package server

import "strconv"

// ToString TODO Implement on an interface
func (r Server) ToString() string {
	return strconv.Itoa(r.Port)
}
