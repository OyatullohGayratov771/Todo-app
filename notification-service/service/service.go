package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"

	notifactionpb "notification-service/protos/notification"

	"github.com/IBM/sarama"
)

type NotificationServiceServer struct {
	notifactionpb.UnimplementedNotificationServiceServer
	Consumer sarama.Consumer
	Topic    string
}

func NewNotificationServer(consumer sarama.Consumer, topic string) *NotificationServiceServer {
	return &NotificationServiceServer{
		Consumer: consumer,
		Topic:    topic,
	}
}

func (s *NotificationServiceServer) SendNotification(ctx context.Context, req *notifactionpb.NotificationReq) (*notifactionpb.NotifactionRes, error) {
	// Your logic here
	return &notifactionpb.NotifactionRes{Message: "Notification sent"}, nil
}

func (s *NotificationServiceServer) ListenForKafkaMessages() {
	partitionConsumer, err := s.Consumer.ConsumePartition(s.Topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Partition consumer yaratishda xatolik: %v", err)
	}
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
		log.Printf("Olingan Kafka xabari: %s\n", string(msg.Value))

		var emailMsg notifactionpb.NotificationReq
		if err := json.Unmarshal(msg.Value, &emailMsg); err != nil {
			log.Printf("Xabarni parsing qilishda xatolik: %v", err)
			continue
		}
		

		log.Println("!!!!!!!!!!!!!!!",emailMsg.UserEmail, emailMsg.Subject, emailMsg.Message,"!!!!!!!!!!!!!!!")
		err = SendEmail(emailMsg.UserEmail, emailMsg.Subject, emailMsg.Message)
		log.Println(err)
		if err != nil {
			log.Printf("❌ Email yuborishda xatolik: %v", err)
		} else {
			log.Printf("📩 Email yuborildi: %s", emailMsg.UserEmail)
		}
	}
}

func SendEmail(to, subject, body string) error {
	from := "tcode696@gmail.com"
	pass := "umel yswa aahx hmjj"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, pass, smtpHost)

	// RFC 5322 formatida xabar
	msg := []byte(
		"From: Todo App <" + from + ">\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
			"Content-Transfer-Encoding: 8bit\r\n" +
			"\r\n" +
			body + "\r\n" +
			"\r\n") // <- SMTP uchun yakuniy CRLF

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
