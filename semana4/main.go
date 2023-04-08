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
	numContagios   int
}

type Traveler struct {
	userId  User
	origin  string
	destiny string
}

type ListaSolicitudes struct {
	solicitudes []Traveler
	users       map[string]User
}

type ArbolUsuarios struct {
	raiz *Nodo
}

func (ls *ListaSolicitudes) AddNewTravel(user User, origin string, destiny string) {
	newTraveler := Traveler{userId: user, origin: origin, destiny: destiny}
	ls.solicitudes = append(ls.solicitudes, newTraveler)
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
				if st.userId.id == uid {
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
	nivel   int
}

func (a *ArbolUsuarios) insertaUsuario(u *User) {
	nuevoNodo := &Nodo{usuario: u}
	if a.raiz == nil {
		a.raiz = nuevoNodo
		fmt.Println("Usuario insertado:", u)
		return
	}
	nodoActual := a.raiz
	for nodoActual != nil {
		if u.nombre == nodoActual.usuario.nombre {
			fmt.Println("El usuario ya existe")
			return
		}
		if u.nombre < nodoActual.usuario.nombre {
			if nodoActual.izq == nil {
				nodoActual.izq = nuevoNodo
				fmt.Println("Usuario insertado:", u)
				return
			}
			nodoActual = nodoActual.izq
		} else {
			if nodoActual.der == nil {
				nodoActual.der = nuevoNodo
				fmt.Println("Usuario insertado:", u)
				return
			}
			nodoActual = nodoActual.der
		}
	}
}

func (a *ArbolUsuarios) contagio(nombre string) {
	nodoActual := a.raiz
	for nodoActual != nil {
		if nombre == nodoActual.usuario.nombre {
			nodoActual.usuario.numContagios++
			if nodoActual.usuario.numContagios == 5 {
				fmt.Println("El usuario", nombre, "ha recibido 5 contagios y debe ser notificado.")
			}
			return
		}
		if nombre < nodoActual.usuario.nombre {
			nodoActual = nodoActual.izq
		} else {
			nodoActual = nodoActual.der
		}
	}
	fmt.Println("El usuario", nombre, "no se encontró en el árbol.")
}

func (a *ArbolUsuarios) usuariosVulnerables() []string {
	vulnerableUsers := make([]string, 0)
	if a.raiz == nil {
		return vulnerableUsers
	}
	queue := make([]*Nodo, 0)
	queue = append(queue, a.raiz)
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		if n.usuario.numerodequejas >= 2 {
			vulnerableUsers = append(vulnerableUsers, n.usuario.nombre)
		}
		if n.izq != nil {
			queue = append(queue, n.izq)
		}
		if n.der != nil {
			queue = append(queue, n.der)
		}
	}
	return vulnerableUsers
}

func (a *ArbolUsuarios) remueveUsuario(id string) {
	a.raiz = remueveNodo(a.raiz, id)
	fmt.Print("El usuario ", id, "fue eliminado", "\n")
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

	fmt.Println("Buscando usuario: ", "\n", nombre)
	current := tree.raiz
	for current != nil {
		if current.usuario.nombre == nombre {
			fmt.Println("Usuario encontrado: ", "\n", current.usuario)
			return current.usuario, nil
		} else if current.usuario.nombre > nombre {
			fmt.Println(current.izq)
			current = current.izq
		} else {
			fmt.Println(current.der)
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
	fmt.Println("Usuario muestra:", "\n", nodo.usuario)
	nodo.der.muestra()
}

func (tree *ArbolUsuarios) muestraNivel() {
	if tree.raiz == nil {
		fmt.Println("El árbol está vacío")
		return
	}
	cola := []*Nodo{{usuario: tree.raiz.usuario, nivel: 0}}

	for len(cola) > 0 {
		current := cola[0]
		cola = cola[1:]

		fmt.Printf("Nivel %d: %v\n", current.nivel, current.usuario)

		if current.izq != nil {
			cola = append(cola, &Nodo{usuario: current.izq.usuario, nivel: current.nivel + 1})
		}
		if current.der != nil {
			cola = append(cola, &Nodo{usuario: current.der.usuario, nivel: current.nivel + 1})
		}
	}
}

func (tree *ArbolUsuarios) quejas(nombre string) error {
	user, err := tree.encuentraUsuario(nombre)
	if err != nil {
		return err
	}
	fmt.Println("El número de quejas del usuario", user.nombre, "es: ", "\n", user.numerodequejas)
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
	var traverse func(nodo *Nodo)
	traverse = func(nodo *Nodo) {
		if nodo == nil {
			return
		}

		if nodo.usuario.numerodequejas > maxQuejas {
			maxQuejas = nodo.usuario.numerodequejas
			maxUsuarios = []string{nodo.usuario.nombre}
		} else if nodo.usuario.numerodequejas == maxQuejas && !contains(maxUsuarios, nodo.usuario.nombre) {
			maxUsuarios = append(maxUsuarios, nodo.usuario.nombre)
		}

		if nodo.usuario.numerodequejas < minQuejas {
			minQuejas = nodo.usuario.numerodequejas
			minUsuarios = []string{nodo.usuario.nombre}
		} else if nodo.usuario.numerodequejas == minQuejas && !contains(minUsuarios, nodo.usuario.nombre) {
			minUsuarios = append(minUsuarios, nodo.usuario.nombre)
		}

		traverse(nodo.izq)
		traverse(nodo.der)
	}
	traverse(tree.raiz)

	fmt.Printf("Usuarios con más quejas (%d): %v\n", maxQuejas, maxUsuarios)
	fmt.Printf("Usuarios con menos quejas (%d): %v\n", minQuejas, minUsuarios)

	return maxUsuarios, minUsuarios, nil
}

func contains(slice []string, val string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
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
	ls := &ArbolUsuarios{}
	user1 := User{
		id:             uuid.New().String(),
		nombre:         "Maided",
		edad:           20,
		sexo:           "Mujer",
		numerodequejas: 20,
	}
	user2 := User{
		id:             uuid.New().String(),
		nombre:         "Guadalupe",
		edad:           20,
		sexo:           "Mujer",
		numerodequejas: 3,
	}
	user3 := User{
		id:             uuid.New().String(),
		nombre:         "Lupita",
		edad:           20,
		sexo:           "Mujer",
		numerodequejas: 8,
	}
	ls.insertaUsuario(&user1)
	fmt.Printf("\n")
	ls.insertaUsuario(&user2)
	fmt.Printf("\n")
	ls.encuentraUsuario(user1.nombre)
	fmt.Printf("\n")
	ls.remueveUsuario(user3.id)
	fmt.Printf("\n")
	ls.muestra()
	fmt.Printf("\n")
	fmt.Println("Niveles de usuarios:")
	ls.muestraNivel()
	fmt.Printf("\n")
	ls.quejas(user1.nombre)
	fmt.Printf("\n")
	ls.usuariosCriticos()

	// Contagio
	ls.contagio("Maided")
	ls.contagio("Maided")
	ls.contagio("Maided")
	ls.contagio("Maided")
	ls.contagio("Maided")

	// Usuarios vulnerables
	fmt.Println("Usuarios vulnerables:")
	vulnerables := ls.usuariosVulnerables()
	fmt.Println(vulnerables)

	// Ordena
	ls2 := ListaSolicitudes{}
	newTravel1 := Traveler{
		userId:  user1,
		origin:  "Medellín",
		destiny: "Cartagena",
	}
	newTravel12 := Traveler{
		userId:  user2,
		origin:  "Medellín",
		destiny: "Cartagena",
	}
	ls2.AddNewTravel(newTravel1.userId, newTravel1.origin, newTravel1.destiny)
	ls2.AddNewTravel(newTravel12.userId, newTravel12.origin, newTravel12.destiny)
	fmt.Printf("\n")
	fmt.Println("Segunda lista de solicitudes:", "\n", ls2.solicitudes)
	fmt.Printf("\n")
	ls2.Ordena(1)
	fmt.Println("Ordena:", ls2.solicitudes)

}
