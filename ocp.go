package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Action string

const (
	ActionCreate Action = "create"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

type PostUpdateNotifier interface {
	NotifyPostUpdated(post Post, action Action) error
}

type EmailService interface {
	SendEmail(sender string, recipient string, subject string, body string) error
}

type EmailNotifier struct {
	emailService EmailService
}

func (n *EmailNotifier) NotifyPostUpdated(post Post, action Action) error {
	subject := "Post Update Notification"
	body := "The post has been updated with the following details:\n" +
		"Title: " + post.Title + "\n" +
		"Body: " + post.Body + "\n" +
		"Action: " + string(action)
	return n.emailService.SendEmail("noreply@example.com", post.Author.Email, subject, body)
}

type PostHandler struct {
	notifiers []PostUpdateNotifier
}

func NewPostHandler(notifiers ...PostUpdateNotifier) *PostHandler {
	return &PostHandler{
		notifiers: notifiers,
	}
}

func (h *PostHandler) UpdateHandler(c gin.Context) {
	// Update logic for the post
	// ...
	post := Post{}

	// Notify all registered notifiers about the post update
	for _, notifier := range h.notifiers {
		if err := notifier.NotifyPostUpdated(post, ActionUpdate); err != nil {
			c.Error(err) // Handle error appropriately
			return
		}
	}
	c.JSON(200, gin.H{"status": "post updated"})
}

func main() {
	r := gin.Default()

	newGmailService := NewGmailService()
	emailNotifier := &EmailNotifier{emailService: newGmailService}
	lineNotifier := &LineNotifier{lineService: NewLineService()}
	postHandler := NewPostHandler(emailNotifier, lineNotifier)

	logReqMiddleware := func(c *gin.Context) {
		// Log the incoming request
		log.Printf("Incoming request: %s %s", c.Request.Method, c.Request.URL)
		c.Next()
	}

	r.PUT("/posts/:id", postHandler.UpdateHandler)
	r.Run(":8080") // Start the server
}
