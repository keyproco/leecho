package consumers

import (
	"course/config"
	"course/models"
	"encoding/json"
	"log"

	"gorm.io/gorm"
)

type CourseEvent struct {
	EventType string        `json:"event_type"`
	Course    models.Course `json:"course"`
	ID        uint          `json:"id"`
}

func StartCourseEventConsumer(rabbitMQConfig *config.RabbitMQConfig, db *gorm.DB) {
	msgs, err := rabbitMQConfig.Channel.Consume(
		"course_events",
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

			var courseEvent CourseEvent
			err := json.Unmarshal(msg.Body, &courseEvent)
			if err != nil {
				log.Printf("Failed to unmarshal course event data: %s", err)
				continue
			}

			switch courseEvent.EventType {
			case "course.created":
				log.Printf("Handling course created event for course: %s", courseEvent.Course.Title)
				if err := models.CreateCourse(db, &courseEvent.Course); err != nil {
					log.Printf("Failed to insert course into database: %s", err)
					continue
				}
				log.Printf("Course %s inserted into the database successfully!", courseEvent.Course.Title)

			case "course.deleted":
				log.Printf("Handling course deleted event for course: %d", courseEvent.ID)
				if err := models.DeleteCourse(db, courseEvent.ID); err != nil {
					log.Printf("Failed to delete course from database: %s", err)
					continue
				}
				log.Printf("Course %s deleted from the database successfully!", courseEvent.Course.Title)

			default:
				log.Printf("Unknown event type: %s", courseEvent.EventType)
			}
		}
	}()

	log.Println("Waiting for course event messages.")
}
