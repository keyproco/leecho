package consumers

import (
	"class/config"
	"class/models"
	"encoding/json"
	"log"

	"gorm.io/gorm"
)

type ClassEvent struct {
	EventType string       `json:"event_type"`
	Class     models.Class `json:"class"`
	ID        uint         `json:"id"`
}

func StartClassEventConsumer(rabbitMQConfig *config.RabbitMQConfig, db *gorm.DB) {
	msgs, err := rabbitMQConfig.Channel.Consume(
		"class_events",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)

			var classEvent ClassEvent
			err := json.Unmarshal(msg.Body, &classEvent)
			if err != nil {
				log.Printf("Failed to unmarshal class event data: %s", err)
				continue
			}

			switch classEvent.EventType {
			case "class.created":
				log.Printf("Handling class created event for class: %s", classEvent.Class.Title)
				if err := models.CreateClass(db, &classEvent.Class); err != nil {
					log.Printf("Failed to insert class into the database: %s", err)
					continue
				}
				log.Printf("Class '%s' inserted into the database successfully!", classEvent.Class.Title)

			case "class.updated":
				log.Printf("Handling class updated event for class: %s", classEvent.Class.Title)
				if classEvent.Class.ID == 0 {
					log.Printf("No Class ID provided for update event")
					continue
				}

				if err := models.UpdateClass(db, classEvent.Class.ID, &classEvent.Class); err != nil {
					log.Printf("Failed to update class in the database: %s", err)
					continue
				}
				log.Printf("Class '%s' updated in the database successfully!", classEvent.Class.Title)

			case "class.deleted":
				log.Printf("Handling class deleted event for class ID: %d", classEvent.ID)
				if err := models.DeleteClass(db, classEvent.ID); err != nil {
					log.Printf("Failed to delete class from the database: %s", err)
					continue
				}
				log.Printf("Class with ID %d deleted from the database successfully!", classEvent.ID)

			default:
				log.Printf("Unknown event type: %s", classEvent.EventType)
			}
		}
	}()

	log.Println("Waiting for class event messages.")
}
