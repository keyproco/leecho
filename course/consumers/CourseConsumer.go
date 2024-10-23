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

type CoursePathEvent struct {
	EventType  string            `json:"event_type"`
	CoursePath models.CoursePath `json:"coursePath"`
	ID         uint              `json:"id"`
}

func StartCourseEventConsumer(rabbitMQConfig *config.RabbitMQConfig, db *gorm.DB) {
	// Consumer for course_events
	consumeCourseEvents(rabbitMQConfig, db)
	// Consumer for coursePath_events
	consumeCoursePathEvents(rabbitMQConfig, db)

	log.Println("Waiting for course and course path event messages.")
}

func consumeCourseEvents(rabbitMQConfig *config.RabbitMQConfig, db *gorm.DB) {
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
		log.Fatalf("Failed to register a consumer for course_events: %s", err)
	}

	go func() {
		for msg := range msgs {
			log.Printf("Received a message from course_events: %s", msg.Body)

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
					log.Printf("Failed to insert course into the database: %s", err)
					continue
				}
				log.Printf("Course '%s' inserted into the database successfully!", courseEvent.Course.Title)

			case "course.updated":
				log.Printf("Handling course updated event for course: %s", courseEvent.Course.Title)
				if courseEvent.Course.ID == 0 {
					log.Printf("No Course ID provided for update event")
					continue
				}

				if err := models.UpdateCourse(db, courseEvent.Course.ID, &courseEvent.Course); err != nil {
					log.Printf("Failed to update course in the database: %s", err)
					continue
				}
				log.Printf("Course '%s' updated in the database successfully!", courseEvent.Course.Title)

			case "course.deleted":
				log.Printf("Handling course deleted event for course ID: %d", courseEvent.ID)
				if err := models.DeleteCourse(db, courseEvent.ID); err != nil {
					log.Printf("Failed to delete course from the database: %s", err)
					continue
				}
				log.Printf("Course with ID %d deleted from the database successfully!", courseEvent.ID)

			default:
				log.Printf("Unknown event type: %s", courseEvent.EventType)
			}
		}
	}()
}

func consumeCoursePathEvents(rabbitMQConfig *config.RabbitMQConfig, db *gorm.DB) {
	msgs, err := rabbitMQConfig.Channel.Consume(
		"coursePath_events",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer for coursePath_events: %s", err)
	}

	go func() {
		for msg := range msgs {
			log.Printf("Received a message from coursePath_events: %s", msg.Body)

			var coursePathEvent CoursePathEvent
			err := json.Unmarshal(msg.Body, &coursePathEvent)
			if err != nil {
				log.Printf("Failed to unmarshal course path event data: %s", err)
				continue
			}

			switch coursePathEvent.EventType {
			case "coursePath.created":
				log.Printf("Handling course path created event for course path: %s", coursePathEvent.CoursePath.Title)
				if err := models.CreateCoursePath(db, &coursePathEvent.CoursePath); err != nil {
					log.Printf("Failed to insert course path into the database: %s", err)
					continue
				}
				log.Printf("Course Path '%s' inserted into the database successfully!", coursePathEvent.CoursePath.Title)

			case "coursePath.updated":
				log.Printf("Handling course path updated event for course path: %s", coursePathEvent.CoursePath.Title)
				if coursePathEvent.CoursePath.ID == 0 {
					log.Printf("No Course Path ID provided for update event")
					continue
				}

				if err := models.UpdateCoursePath(db, coursePathEvent.CoursePath.ID, &coursePathEvent.CoursePath); err != nil {
					log.Printf("Failed to update course path in the database: %s", err)
					continue
				}
				log.Printf("Course Path '%s' updated in the database successfully!", coursePathEvent.CoursePath.Title)

			case "coursePath.deleted":
				log.Printf("Handling course path deleted event for course path ID: %d", coursePathEvent.ID)
				if err := models.DeleteCoursePath(db, coursePathEvent.ID); err != nil {
					log.Printf("Failed to delete course path from the database: %s", err)
					continue
				}
				log.Printf("Course Path with ID %d deleted from the database successfully!", coursePathEvent.ID)

			default:
				log.Printf("Unknown event type: %s", coursePathEvent.EventType)
			}
		}
	}()
}
