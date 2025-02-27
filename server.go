package main

import (
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
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

type DVector struct {
	DX float64 `json:"dx"` // Changed to float64
	DY float64 `json:"dy"` // Changed to float64
}

type Bullet struct {
	ToVector  DVector `json:"toVector"`
	IsHit     Point   `json:"isHit"`
	FromPoint Point   `json:"fromPoint"`
}

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

var (
	connections = make(map[string]*websocket.Conn)
	points      = make(map[string]Point)
	mu          sync.Mutex
	radius      = 100
	port        = os.Getenv("PORT")
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
		R: uint8(200 + rand.Intn(56)), // 200-255 (ensures brightness)
		G: uint8(rand.Intn(200)),      // 0-199 (adds contrast)
		B: uint8(rand.Intn(200)),      // 0-199 (adds contrast)
	}

	// Shuffle R, G, B to ensure different dominant colors
	perm := rand.Perm(3)
	shuffled := [3]uint8{color.R, color.G, color.B}

	color = Color{
		R: shuffled[perm[0]],
		G: shuffled[perm[1]],
		B: shuffled[perm[2]],
	}

	mu.Lock()
	connections[userID] = conn
	idData := map[string]string{"id": userID}

	log.Println("New connection:", userID)
	broadcastData("i", idData)
	points[userID] = Point{UserID: userID, Color: color}

	existingPoints := make([]Point, 0, len(points))
	for _, p := range points {
		existingPoints = append(existingPoints, p)
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
		delete(points, userID)
		mu.Unlock()
		conn.Close()
		broadcastPoints()
	}()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Connection lost:", userID, err)
			return
		}

		var msg Message

		if err := json.Unmarshal(p, &msg); err != nil {
			log.Println("Error decoding JSON:", err)
			continue
		}
		switch msg.Type {
		case "p":
			var point Point
			if err := json.Unmarshal(msg.Payload, &point); err != nil {
				log.Println("Error decoding point:", err)
				continue
			}
			if point.X == 0 && point.Y == 0 {
				continue
			}

			mu.Lock()
			if existingPoint, exists := points[userID]; exists {
				point.UserID = userID
				point.Color = existingPoint.Color
			}
			points[userID] = point
			mu.Unlock()
			broadcastPoints()

		case "b":
			var bullet Bullet
			var d DVector
			if err := json.Unmarshal(msg.Payload, &d); err != nil {
				log.Println("Error decoding bullet direction:", err)
				continue
			}
			mu.Lock()
			bullet.ToVector = d
			log.Printf("Bullet received: dx=%f, dy=%f", d.DX, d.DY)

			if point, exists := points[userID]; exists {
				bullet.FromPoint = point
				updatedPoints := make([]Point, 0, len(points))
				for _, p := range points {
					updatedPoints = append(updatedPoints, p)
				}
				checkHit := checkHit(point, d, updatedPoints)
				if checkHit == nil {
					log.Println("No points hit")
				} else {
					log.Println("Bullet hit:", checkHit.UserID)
					bullet.IsHit = *checkHit
					broadcastData("b", bullet)
				}
			} else {
				log.Println("No point found for userID:", userID)
			}

			mu.Unlock()

		}
	}
}

func checkHit(shooter Point, direction DVector, cursors []Point) *Point {
	x1 := float64(shooter.X)
	y1 := float64(shooter.Y)

	dx := direction.DX
	dy := direction.DY

	for _, cursor := range cursors {
		if cursor.UserID == shooter.UserID {
			continue
		}

		cx := float64(cursor.X)
		cy := float64(cursor.Y)

		// Vector from shooter to cursor
		vx := cx - x1
		vy := cy - y1

		// Project vector onto bullet dir
		projection := vx*dx + vy*dy

		// check if cursor behind shooter
		if projection < 0 {
			continue
		}

		//Calc perpendicular distance  -> cursor to trajectory
		distance := math.Abs(vx*dy - vy*dx)

		if distance <= float64(radius) {
			if cursor.X != 0 && cursor.Y != 0 {
				return &cursor
			} else {
				log.Println("point at 0,0")
			}
		}
	}
	return nil
}

func broadcastData(msgType string, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error encoding payload:", err)
		return
	}

	message := Message{
		Type:    msgType,
		Payload: data,
	}

	messageData, err := json.Marshal(message)
	if err != nil {
		log.Println("Error encoding message:", err)
		return
	}

	for userID, conn := range connections {
		if err := conn.WriteMessage(websocket.TextMessage, messageData); err != nil {
			log.Println("Error broadcasting to", userID, ":", err)
			conn.Close()
			delete(connections, userID)
			delete(points, userID)
		}
	}
}
func broadcastPoints() {
	mu.Lock()
	updatedPoints := make([]Point, 0, len(points))
	for _, p := range points {
		if p.X != 0 || p.Y != 0 {
			updatedPoints = append(updatedPoints, p)
		}
	}
	mu.Unlock()

	broadcastData("p", updatedPoints)
}

func main() {
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", handler)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
