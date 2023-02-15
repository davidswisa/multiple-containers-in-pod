package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/davidswisa/multiple-containers-in-pod/pkg/orm"
	rsv "github.com/davidswisa/multiple-containers-in-pod/pkg/reservation"
)

// func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
// 	return kafka.NewReader(kafka.ReaderConfig{
// 		Brokers:  []string{kafkaURL},
// 		GroupID:  groupID,
// 		Topic:    topic,
// 		MinBytes: 10e3, // 10KB
// 		MaxBytes: 10e4, // 100KB
// 	})
// }

type Message struct {
	Key   []byte `json:"key"`
	Value []byte `json:"value"`
}

func main() {
	// get kafka reader using environment variables.
	client := orm.NewORMClient()

	// server := "kafka" //os.Getenv("kafka")
	// kafkaURL := server + ":9092"
	// topic := "topic1"
	// groupID := "group1"
	// reader := getKafkaReader(kafkaURL, topic, groupID)

	// defer reader.Close()
	c, err := net.Listen("unix", "/tmp/echo.sock")
	if err != nil {
		log.Fatalf("consumer: %v", err)
	}
	log.Println("consuming from socket...")

	for {
		conn, err := c.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("consumer: got connection from ", conn.RemoteAddr())
		buf := make([]byte, 256)
		b, err := conn.Read(buf[:])
		if err != nil {
			log.Fatal(err)
		}
		buf = buf[:b]
		var msg Message
		err = json.Unmarshal(buf, &msg)
		if err != nil {
			log.Fatal("consumer: ", err)
		}
		log.Printf("Reading message, operation : %v", msg.Key)
		switch string(msg.Key) {
		case rsv.OPNEW:
			{
				r, err := rsv.Decode(msg.Value)
				if err != nil {
					log.Printf("Error: unexpected decoding issue.\nreason: %v", err)
					continue
				}

				headers := http.Header{}
				headers.Add("Content-Type", "application/json")

				b, err := json.Marshal(r)
				if err != nil {
					log.Printf("Error: unexpected decoding issue.\nreason: %v", err)
					return
				}
				log.Printf("Reservation details : %s", string(b))

				body, res, err := client.Post("reservations", string(b), headers)
				if err != nil {
					log.Printf("Error: unexpected database issue.\nreason: %v", err)
					panic(err)
				}
				if res.StatusCode == 200 {
					log.Printf("Reservation (id: %d) inserted successfully", r.ID)
				}
				log.Printf("Reservation (id: %d) Status Code (code: %d) body: %s", r.ID, res.StatusCode, body)
			}
		case rsv.OPREM:
			{
				log.Println("delete flow")
				id := string(msg.Value)
				log.Printf("Reservation details : %s", id)

				body, res, err := client.Delete("reservations/"+id, http.Header{})
				if err != nil {
					log.Printf("Error: unexpected database issue.\nreason: %v", err)
					panic(err)
				}
				if res.StatusCode == 200 {
					log.Printf("Reservation (id: %d) inserted successfully", id)
				}
				log.Printf("Reservation (id: %d) Status Code (code: %d) body: %s", id, res.StatusCode, body)

			}
		case rsv.OPCHG:
			{
				fmt.Println("update flow")
				r, err := rsv.Decode(msg.Value)
				if err != nil {
					log.Printf("Error: unexpected decoding issue.\nreason: %v", err)
					continue
				}

				id := strconv.Itoa(r.ID)
				headers := http.Header{}
				headers.Add("Content-Type", "application/json")

				b, err := json.Marshal(r)
				if err != nil {
					log.Printf("Error: unexpected decoding issue.\nreason: %v", err)
					return
				}
				log.Printf("Reservation details:\nid: %s, body: %s", id, string(b))

				body, res, err := client.Put("reservations/"+id, string(b), headers)
				if err != nil {
					log.Printf("Error: unexpected database issue.\nreason: %v", err)
					panic(err)
				}
				if res.StatusCode == 200 {
					log.Printf("Reservation (id: %s) Updated successfully", id)
				}
				log.Printf("Reservation (id: %s) Status Code (code: %d) body: %s", id, res.StatusCode, body)
			}
		}
	}
}
