package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Color struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
}

type Point struct {
	X      int    `json:"x"`
	Y      int    `json:"y"`
	UserID string `json:"user_id"`
	Color  Color  `json:"color"`
}

var (
	connections = make(map[string]*websocket.Conn)
	points      = make(map[string]Point)
	mu          sync.Mutex
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	userID := uuid.New().String()

	color := Color{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
	}

	mu.Lock()
	connections[userID] = conn
	log.Println("New connection:", userID)

	points[userID] = Point{UserID: userID, Color: color}

	existingPoints := make([]Point, 0, len(points))
	for _, p := range points {
		// if p.X != 0 || p.Y != 0 {
		existingPoints = append(existingPoints, p)
		// }
	}
	data, err := json.Marshal(existingPoints)
	if err == nil {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("Error sending initial data:", err)
		}
	} else {
		log.Println("Error encoding initial objects:", err)
	}
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(connections, userID)
		delete(points, userID) // Remove cursor on disconnect
		mu.Unlock()
		conn.Close()
		broadcastPoints()
	}()

	// Handle incoming messages (points)
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Connection lost:", userID, err)
			return
		}

		var point Point
		if err := json.Unmarshal(p, &point); err != nil {
			log.Println("Error decoding JSON:", err)
			continue
		}

		// Ignore points at (0, 0)
		if point.X == 0 && point.Y == 0 {
			continue
		}

		// Assign the correct user ID and retain the original color
		mu.Lock()
		if existingPoint, exists := points[userID]; exists {
			point.UserID = userID
			point.Color = existingPoint.Color // Preserve the original color
		}
		points[userID] = point
		mu.Unlock()

		broadcastPoints()
	}
}

// Broadcast all points to every connected client
func broadcastPoints() {
	mu.Lock()
	defer mu.Unlock()

	updatedPoints := make([]Point, 0, len(points))
	for _, p := range points {
		// Filter out points at (0, 0)
		if p.X != 0 || p.Y != 0 {
			updatedPoints = append(updatedPoints, p)
		}
	}

	data, _ := json.Marshal(updatedPoints)

	for userID, conn := range connections {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("Error broadcasting to", userID, ":", err)
			conn.Close()
			delete(connections, userID)
			delete(points, userID) // Remove cursor on disconnect
		}
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Server started on :8080")
	http.ListenAndServe(":8080", nil)
}
