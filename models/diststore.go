package models

type NotifyRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Gen   int    `json:"gen"`
}

type Value struct {
	Val string
	Gen int
}
