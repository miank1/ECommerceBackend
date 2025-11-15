package service

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"sync"
	"time"
)

type notification struct {
	ID    string                 `json:"id"`
	Type  string                 `json:"type"`
	To    string                 `json:"to"`
	Title string                 `json:"title"`
	Body  string                 `json:"body"`
	Meta  map[string]interface{} `json:"meta"`
	Tries int                    `json:"tries"`
}

type NotificationService struct {
	queue       chan notification
	queueSize   int
	smtpAddr    string
	from        string
	wg          sync.WaitGroup
	shutdownCh  chan struct{}
	clientMutex sync.Mutex
}

func NewNotificationService(smtpAddr, from string) *NotificationService {
	return &NotificationService{
		queueSize:  100,
		queue:      make(chan notification, 100),
		smtpAddr:   smtpAddr,
		from:       from,
		shutdownCh: make(chan struct{}),
	}
}

// EnqueueNotification pushes event to internal queue
func (s *NotificationService) EnqueueNotification(typ, to, title, body string, meta map[string]interface{}) error {
	if to == "" {
		return errors.New("missing to")
	}
	n := notification{
		ID:    fmt.Sprintf("%d", time.Now().UnixNano()),
		Type:  typ,
		To:    to,
		Title: title,
		Body:  body,
		Meta:  meta,
		Tries: 0,
	}
	select {
	case s.queue <- n:
		return nil
	default:
		// queue full — return transient error so caller can retry if needed
		return errors.New("notification queue full")
	}
}

func (s *NotificationService) StartWorkers(n int) {
	for i := 0; i < n; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}
}

func (s *NotificationService) worker(id int) {
	defer s.wg.Done()
	log.Printf("notification worker %d started", id)
	for {
		select {
		case n := <-s.queue:
			s.processWithRetry(n)
		case <-s.shutdownCh:
			log.Printf("worker %d shutting down", id)
			return
		}
	}
}

func (s *NotificationService) processWithRetry(n notification) {
	maxTries := 3
	backoff := time.Second * 2

	for {
		err := s.process(n)
		if err == nil {
			return
		}
		n.Tries++
		if n.Tries >= maxTries {
			log.Printf("notification %s failed after %d tries: %v", n.ID, n.Tries, err)
			return
		}
		log.Printf("notification %s failed, retrying in %v: %v", n.ID, backoff, err)
		time.Sleep(backoff)
		backoff *= 2
	}
}

func (s *NotificationService) process(n notification) error {
	// choose sending method by Type (for now only email)
	switch n.Type {
	case "order_paid":
		return s.sendEmail(n)
	default:
		// log fallback
		log.Printf("unhandled notification type=%s body=%s", n.Type, n.Body)
		return nil
	}
}

func (s *NotificationService) sendEmail(n notification) error {
	// smtpAddr is like "mailhog:1025" or "smtp.example.com:587"
	// For local dev MailHog listens on 1025 (no auth)
	msg := bytes.Buffer{}
	subject := n.Title
	fmt.Fprintf(&msg, "From: %s\r\n", s.from)
	fmt.Fprintf(&msg, "To: %s\r\n", n.To)
	fmt.Fprintf(&msg, "Subject: %s\r\n", subject)
	fmt.Fprintf(&msg, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&msg, "Content-Type: text/plain; charset=\"utf-8\"\r\n")
	fmt.Fprintf(&msg, "\r\n")
	fmt.Fprintf(&msg, "%s\r\n", n.Body)

	// No auth example (MailHog)
	err := smtp.SendMail(s.smtpAddr, nil, s.from, []string{n.To}, msg.Bytes())
	if err != nil {
		return err
	}
	log.Printf("email sent to=%s subject=%s", n.To, subject)
	return nil
}

// Shutdown to gracefully stop workers
func (s *NotificationService) Shutdown() {
	close(s.shutdownCh)
	s.wg.Wait()
}
