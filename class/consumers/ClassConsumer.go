package consumers

import (
	"class/config"
	"class/models"
	"encoding/json"
	"log"

	"gorm.io/gorm"
)

func StartClassCreatedConsumer(rabbitMQConfig *config.RabbitMQConfig, db *gorm.DB) {
	msgs, err := rabbitMQConfig.Channel.Consume(
		"class_created",
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

			var class models.Class
			err := json.Unmarshal(msg.Body, &class)
			if err != nil {
				log.Printf("Failed to unmarshal class data: %s", err)
				continue
			}
			if err := models.CreateClass(db, &class); err != nil {
				log.Printf("Failed to insert class into database: %s", err)
				continue
			}

			log.Printf("Class %s inserted into the database successfully!", class.Title)
		}
	}()
	log.Println("Waiting for class_created messages. To exit press CTRL+C")
}
