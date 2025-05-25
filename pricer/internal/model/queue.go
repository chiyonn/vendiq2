package model


type QueueInfo struct {
	Name               string `json:"name"`
	Messages           int    `json:"messages"`
	MessagesReady      int    `json:"messages_ready"`
	MessagesUnack      int    `json:"messages_unacknowledged"`
}

