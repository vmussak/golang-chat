package main

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

//log --> console.log()
//net/http --> http
//github.com/gorilla/websocket --> socket.io

var clients = make(map[*websocket.Conn]bool) // clients que conectaram
var broadcast = make(chan Message)           // canal de comunicação

// Configurando o 'atualizador'
var upgrader = websocket.Upgrader {
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Definindo o objeto de mensagem
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	// Criando um servidor de arquivos
	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	// Configurando as rotas do Socket
	http.HandleFunc("/ws", handleConnections)

	// Inicia a escuta para novas mensagens
	go handleMessages()

	// Inicia o servidor na porta 8000
	log.Println("server rodando na porta 8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Faz a primeira requisição GET se tornar um WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Definindo para fechar a conecção quando a função terminar
	defer ws.Close()

	// Registrando nosso novo client
	clients[ws] = true

	for {
		var msg Message
		// Lendo o JSON da mensagem e colocando no objeto
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// Coloca a nova mensagem para os demais clients
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Pega a próxima mensagem do canal de comunicação
		msg := <-broadcast
		// Manda a mensagem para todos os clients conectados
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
