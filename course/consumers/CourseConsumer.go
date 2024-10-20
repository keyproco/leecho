package consumers

import (
	"course/config"
	"course/models"
	"encoding/json"
	"log"

	"gorm.io/gorm"
)

func StartCourseCreatedConsumer(rabbitMQConfig *config.RabbitMQConfig, db *gorm.DB) {
	msgs, err := rabbitMQConfig.Channel.Consume(
		"course_created",
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

			var course models.Course
			err := json.Unmarshal(msg.Body, &course)
			if err != nil {
				log.Printf("Failed to unmarshal course data: %s", err)
				continue
			}
			if err := models.CreateCourse(db, &course); err != nil {
				log.Printf("Failed to insert course into database: %s", err)
				continue
			}

			log.Printf("Course %s inserted into the database successfully!", course.Title)
		}
	}()
	log.Println("Waiting for course_created messages.")
}
