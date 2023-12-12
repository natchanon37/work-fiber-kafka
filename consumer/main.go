package main

import (
	"consumer/repositories"
	"consumer/services"
	"context"
	"events"
	"fmt"
	"strings"

	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func initDatabase() *gorm.DB {
	//postgres dsn
	dsn := fmt.Sprintf("%v://%v:%v@%v:%v/%v",
		viper.GetString("db.driver"),
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.host"),
		viper.GetInt("db.port"),
		viper.GetString("db.name"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}
	return db
}

func main() {
	consumer, err := sarama.NewConsumerGroup(viper.GetStringSlice("kafka.server"), viper.GetString("kafka.group"), nil)
	if err != nil {
		panic(err)
	}

	defer consumer.Close()
	db := initDatabase()
	accountRepo := repositories.NewAccountRepository(db)
	accountEventHandler := services.NewEventHandler(accountRepo)
	accountConsumerHandler := services.NewConsumerHandler(accountEventHandler)

	fmt.Println("Consumer is running")
	for {
		consumer.Consume(context.Background(), events.Topic, accountConsumerHandler)
	}
}
