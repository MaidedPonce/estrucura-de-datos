package main

import (
	"fmt"
	"sort"

	"github.com/google/uuid"
)

type User struct {
	id             string
	nombre         string
	edad           int
	sexo           string
	numerodequejas int
}

type Traveler struct {
	userId  string
	origin  string
	destiny string
}

type ListaSolicitudes struct {
	solicitudes []Traveler
	users       map[string]User
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

func (ls *ListaSolicitudes) Ordena(orden int) ListaSolicitudes {
	result := ListaSolicitudes{}
	if orden == 1 {
		sort.SliceStable(ls.solicitudes, func(i, j int) bool {
			return ls.solicitudes[i].origin < ls.solicitudes[j].origin
		})
	} else if orden == 2 {
		sort.SliceStable(ls.solicitudes, func(i, j int) bool {
			return ls.solicitudes[i].destiny < ls.solicitudes[j].destiny
		})
	} else {
		userIds := make([]string, 0, len(ls.users))
		for k := range ls.users {
			userIds = append(userIds, k)
		}
		sort.Strings(userIds)
		for _, uid := range userIds {
			t := Traveler{}
			for _, st := range ls.solicitudes {
				if st.userId == uid {
					t = st
					break
				}
			}
			result.solicitudes = append(result.solicitudes, t)
		}
	}
	return result
}

type Nodo struct {
	usuario *User
	izq     *Nodo
	der     *Nodo
}

type ArbolUsuarios struct {
	raiz *Nodo
}

func (a *ArbolUsuarios) insertaUsuario(u *User) {
	nuevoNodo := &Nodo{usuario: u}
	if a.raiz == nil {
		a.raiz = nuevoNodo
		return
	}
	nodoActual := a.raiz
	for nodoActual != nil {
		if u.id == nodoActual.usuario.id {
			fmt.Println("El usuario ya existe")
			return
		}
		if u.id < nodoActual.usuario.id {
			if nodoActual.izq == nil {
				nodoActual.izq = nuevoNodo
				return
			}
			nodoActual = nodoActual.izq
		} else {
			if nodoActual.der == nil {
				nodoActual.der = nuevoNodo
				return
			}
			nodoActual = nodoActual.der
		}
	}
}

func (a *ArbolUsuarios) remueveUsuario(id string) {
	a.raiz = remueveNodo(a.raiz, id)
}

func remueveNodo(nodo *Nodo, id string) *Nodo {
	if nodo == nil {
		return nil
	}

	if id < nodo.usuario.id {
		nodo.izq = remueveNodo(nodo.izq, id)
		return nodo
	}

	if id > nodo.usuario.id {
		nodo.der = remueveNodo(nodo.der, id)
		return nodo
	}

	if nodo.izq == nil && nodo.der == nil {
		return nil
	}

	if nodo.izq == nil {
		return nodo.der
	}

	if nodo.der == nil {
		return nodo.izq
	}

	sucesor := nodo.der
	for sucesor.izq != nil {
		sucesor = sucesor.izq
	}

	nodo.usuario = sucesor.usuario
	nodo.der = remueveNodo(nodo.der, sucesor.usuario.id)
	return nodo
}

func (tree *ArbolUsuarios) encuentraUsuario(nombre string) (*User, error) {
	if tree.raiz == nil {
		return nil, fmt.Errorf("El árbol está vacío")
	}

	current := tree.raiz
	for current != nil {
		if current.usuario.nombre == nombre {
			return current.usuario, nil
		} else if current.usuario.nombre > nombre {
			current = current.izq
		} else {
			current = current.der
		}
	}

	return nil, fmt.Errorf("Usuario no encontrado")
}

func (tree *ArbolUsuarios) muestra() {
	if tree.raiz == nil {
		fmt.Println("El árbol está vacío")
		return
	}

	tree.raiz.muestra()
}

func (nodo *Nodo) muestra() {
	if nodo == nil {
		return
	}

	nodo.izq.muestra()
	fmt.Println(nodo.usuario)
	nodo.der.muestra()
}

func (tree *ArbolUsuarios) muestraNivel() {
	if tree.raiz == nil {
		fmt.Println("El árbol está vacío")
		return
	}
	cola := []*Nodo{tree.raiz}

	for len(cola) > 0 {
		current := cola[0]
		cola = cola[1:]

		fmt.Println(current.usuario)

		if current.izq != nil {
			cola = append(cola, current.izq)
		}
		if current.der != nil {
			cola = append(cola, current.der)
		}
	}
}

func (tree *ArbolUsuarios) quejas(nombre string) error {
	user, err := tree.encuentraUsuario(nombre)
	if err != nil {
		return err
	}
	user.numerodequejas++
	if user.numerodequejas == 5 {
		tree.raiz = remueveNodo(tree.raiz, nombre)
	}
	return nil
}

func (tree *ArbolUsuarios) usuariosCriticos() ([]string, []string, error) {
	if tree.raiz == nil {
		return nil, nil, fmt.Errorf("El árbol está vacío")
	}
	maxQuejas := tree.raiz.usuario.numerodequejas
	minQuejas := tree.raiz.usuario.numerodequejas
	maxUsuarios := []string{tree.raiz.usuario.nombre}
	minUsuarios := []string{tree.raiz.usuario.nombre}
	updateLists := func(user *User, usuarios *[]string, quejas int) {
		if user.numerodequejas == quejas {
			*usuarios = append(*usuarios, user.nombre)
		} else if user.numerodequejas > quejas {
			*usuarios = []string{user.nombre}
		}
	}
	var traverse func(nodo *Nodo)
	traverse = func(nodo *Nodo) {
		if nodo == nil {
			return
		}

		updateLists(nodo.usuario, &maxUsuarios, maxQuejas)
		updateLists(nodo.usuario, &minUsuarios, minQuejas)

		if nodo.usuario.numerodequejas > maxQuejas {
			maxQuejas = nodo.usuario.numerodequejas
		} else if nodo.usuario.numerodequejas < minQuejas {
			minQuejas = nodo.usuario.numerodequejas
		}

		traverse(nodo.izq)
		traverse(nodo.der)
	}

	traverse(tree.raiz)

	return maxUsuarios, minUsuarios, nil
}

func imprimirArbol(nodo *Nodo) {
	if nodo == nil {
		return
	}
	imprimirArbol(nodo.izq)
	fmt.Println(nodo.usuario)
	imprimirArbol(nodo.der)
}

func main() {
	ls := ListaSolicitudes{}
	ls.AddNewTraveler("Bogotá", "Medellín")
	ls.AddNewTraveler("Bogotá", "Cali")
	fmt.Println("Lista de solicitudes:", ls.solicitudes)

	ls2 := ListaSolicitudes{}
	ls2.AddNewTraveler("Bogotá", "Cali")
	ls2.AddNewTraveler("Medellín", "Cartagena")
	fmt.Println("Segunda lista de solicitudes:", ls2.solicitudes)

	ls3 := ls.fusionaSolicitudes(ls2)
	fmt.Println("Lista fusionada de solicitudes:", ls3.solicitudes)

	ls4 := ls.compartir(ls2)
	fmt.Println("Solicitudes compartidas:", ls4.solicitudes)

	ls5 := ls.Ordena(1)
	fmt.Println("Solicitudes ordenadas por origen:", ls5.solicitudes)

	tree := ArbolUsuarios{}
	user1 := User{id: "1", nombre: "Juan", edad: 25, sexo: "M", numerodequejas: 0}
	user2 := User{id: "2", nombre: "Ana", edad: 30, sexo: "F", numerodequejas: 1}
	user3 := User{id: "3", nombre: "Pedro", edad: 20, sexo: "M", numerodequejas: 2}

	tree.insertaUsuario(&user1)
	tree.insertaUsuario(&user2)
	tree.insertaUsuario(&user3)

	fmt.Println("Usuarios del árbol:")
	imprimirArbol(tree.raiz)

	tree.remueveUsuario("2")

	fmt.Println("Usuarios del árbol después de remover el usuario con id=2:")
	imprimirArbol(tree.raiz)
}
