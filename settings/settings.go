package settings

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Settings struct {
	Rabbit RabbitMQSettings
}

type RabbitMQSettings struct {
	User string
	Password string
	Host string
	Port int
	Vhost string
	Exhange string
	RoutingKey string
}

func (s *Settings) GetSettings() *Settings {
	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file")
	}
	s.Rabbit.User = os.Getenv("RABBITMQ_USER")
	s.Rabbit.Password = os.Getenv("RABBITMQ_PASSWORD")
	s.Rabbit.Host = os.Getenv("RABBITMQ_HOST")
	s.Rabbit.Port = s.StrToIntParseOrGetDefault("RABBITMQ_PORT", 5672)
	s.Rabbit.Vhost = os.Getenv("RABBITMQ_VHOST")
	s.Rabbit.Exhange = os.Getenv("RABBITMQ_EXCHANGE")
	s.Rabbit.RoutingKey = os.Getenv("RABBITMQ_ROUTING_KEY")
	return s

}

func (s *Settings) GetRabbitmqUrl() string {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d%s",
	s.Rabbit.User,
	s.Rabbit.Password,
	s.Rabbit.Host,
	s.Rabbit.Port,
	s.Rabbit.Vhost,
	)
	return url
}

func (s *Settings) StrToIntParseOrGetDefault(envName string, defaultValue int) int {
	// Method parse env variable to int. If its not posible or env is not set it returns default int value
	strValue := os.Getenv(envName)
	if strValue == "" {
		log.Println("Variable is empty. Using default.")
		return defaultValue
	}
	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		log.Println("Can't use variable. Using default.")
		return defaultValue
	} 
	return intValue
}