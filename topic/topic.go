package topic

import (
	"errors"
	"fmt"
	"github.com/zentures/message"
)

const (
	// MWC is the multi-level wildcard
	MWC = "#"

	// SWC is the single level wildcard
	SWC = "+"

	// SEP is the topic level separator
	Message SEP = "/"

	// SYS is the starting character of the system level topics
	SYS = "$"

	// Both wildcards
	_WC = "#+"
)

var (
	// ErrAuthFailure is returned when the user/pass supplied are invalid
	ErrAuthFailure = errors.New("auth: Authentication failure")

	// ErrAuthProviderNotFound is returned when the requested provider does not exist.
	// It probably hasn't been registered yet.
	ErrAuthProviderNotFound = errors.New("auth: Authentication provider not found")

	providers = make(map[string]TopicsProvider)
)

// TopicsProvider
type TopicsProvider interface {
	Subscribe(topic []byte, qos byte, subscriber interface{}) (byte, error)
	Unsubscribe(topic []byte, subscriber interface{}) error
	Subscribers(topic []byte, qos byte, subs *[]interface{}, qoss *[]byte) error
	Retain(msg *message.PublishMessage) error
	Retained(topic []byte, msgs *[]*message.PublishMessage) error
	Close() error
}

func Register(name string, provider TopicsProvider) {
	if provider == nil {
		panic("topics: Register provide is nil")
	}

	if _, dup := providers[name]; dup {
		panic("topics: Register called twice for provider " + name)
	}

	providers[name] = provider
}

func Unregister(name string) {
	delete(providers, name)
}

type Manager struct {
	p TopicsProvider
}

func NewManager(providerName string) (*Manager, error) {
	p, ok := providers[providerName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provider %q", providerName)
	}

	return &Manager{p: p}, nil
}

func (this *Manager) Subscribe(topic []byte, qos byte, subscriber interface{}) (byte, error) {
	return this.p.Subscribe(topic, qos, subscriber)
}

func (this *Manager) Unsubscribe(topic []byte, subscriber interface{}) error {
	return this.p.Unsubscribe(topic, subscriber)
}

func (this *Manager) Subscribers(topic []byte, qos byte, subs *[]interface{}, qoss *[]byte) error {
	return this.p.Subscribers(topic, qos, subs, qoss)
}

func (this *Manager) Retain(msg *message.PublishMessage) error {
	return this.p.Retain(msg)
}

func (this *Manager) Retained(topic []byte, msgs *[]*message.PublishMessage) error {
	return this.p.Retained(topic, msgs)
}

func (this *Manager) Close() error {
	return this.p.Close()
}
