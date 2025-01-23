package services

import (
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var (
	store = make(map[string]string)
	mu    sync.Mutex
)

func GuardarSentencia(sentencia string) string {
	mu.Lock()
	defer mu.Unlock()

	id := generarID()
	store[id] = sentencia
	return id
}

func ObtenerSentenciaPorID(id string) string {
	mu.Lock()
	defer mu.Unlock()

	return store[id]
}

func generarID() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(999999)) // un número random
}
