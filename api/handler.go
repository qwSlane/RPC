package api

import "log"

type Handler struct{}

func (s *Handler) Welcome(str string) {
	log.Println("Hello", str)
}
