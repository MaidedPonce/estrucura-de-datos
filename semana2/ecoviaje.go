package main

import (
	"fmt"

	"github.com/google/uuid"
)

type Traveler struct {
	userId  string
	origin  string
	destiny string
}

type ListaSolicitudes struct {
	solicitudes []Traveler
}

func (ls *ListaSolicitudes) AddNewTraveler(origin string, destiny string) {
	id := uuid.New().String()
	newTraveler := Traveler{userId: id, origin: origin, destiny: destiny}
	ls.solicitudes = append(ls.solicitudes, newTraveler)
}

func (ls *ListaSolicitudes) fusionaSolicitudes(lista2 ListaSolicitudes) ListaSolicitudes {
	result := ListaSolicitudes{}
	len1 := len(ls.solicitudes)
	len2 := len(lista2.solicitudes)
	if len1 >= len2 {
		for i := 0; i < len1; i++ {
			if i < len2 {
				result.solicitudes = append(result.solicitudes, ls.solicitudes[i], lista2.solicitudes[i])
			} else {
				result.solicitudes = append(result.solicitudes, ls.solicitudes[i])
			}
		}
	} else {
		for i := 0; i < len2; i++ {
			if i < len1 {
				result.solicitudes = append(result.solicitudes, ls.solicitudes[i], lista2.solicitudes[i])
			} else {
				result.solicitudes = append(result.solicitudes, lista2.solicitudes[i])
			}
		}
	}
	return result
}

func (ls *ListaSolicitudes) compartir(lista2 ListaSolicitudes) ListaSolicitudes {
	result := ListaSolicitudes{}
	for _, t1 := range ls.solicitudes {
		for _, t2 := range lista2.solicitudes {
			if t1.origin == t2.origin && t1.destiny == t2.destiny {
				result.solicitudes = append(result.solicitudes, t1, t2)
			}
		}
	}
	return result
}

func main() {
	ls1 := ListaSolicitudes{}
	ls1.AddNewTraveler("Puebla", "CDMX")
	ls1.AddNewTraveler("Guadalajara", "Monterrey")

	ls2 := ListaSolicitudes{}
	ls2.AddNewTraveler("CDMX", "Puebla")
	ls2.AddNewTraveler("Guadalajara", "Tijuana")

	ls3 := ls1.fusionaSolicitudes(ls2)
	fmt.Println("Fusión de solicitudes:")
	for _, t := range ls3.solicitudes {
		fmt.Printf("User ID: %s, Origen: %s, Destino: %s\n", t.userId, t.origin, t.destiny)
	}

	ls4 := ls1.compartir(ls2)
	fmt.Println("\nCompartir solicitudes:")
	for _, t := range ls4.solicitudes {
		fmt.Printf("User ID: %s, Origen: %s, Destino: %s\n", t.userId, t.origin, t.destiny)
	}
}
