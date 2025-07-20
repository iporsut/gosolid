package main

import "github.com/gin-gonic/gin"

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
		// Violating LSP by using a type switch
		// This is not a good practice as it breaks the Liskov Substitution Principle
		// Ideally, we should not check the type of notifier here
		// Instead, we should rely on the interface contract
		switch notifier.(type) {
		case *EmailNotifier:
			// Specific logic for EmailNotifier if needed
			if err := notifier.NotifyPostUpdated(post, ActionUpdate); err != nil {
				c.Error(err) // Handle error appropriately
				return
			}
		}
	}
	c.JSON(200, gin.H{"status": "post updated"})
}
