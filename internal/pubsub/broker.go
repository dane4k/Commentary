package pubsub

import (
	"Commentary/internal/graph/model"
	"github.com/sirupsen/logrus"
	"sync"
)

type Broker struct {
	subs map[int][]chan *model.Comment
	mu   sync.RWMutex
}

func NewBroker() *Broker {
	return &Broker{
		subs: make(map[int][]chan *model.Comment),
	}
}

func (b *Broker) Subscribe(postID int) <-chan *model.Comment {
	logrus.Debugf("subscribing to post %v", postID)

	b.mu.Lock()

	defer b.mu.Unlock()

	ch := make(chan *model.Comment, 10)
	b.subs[postID] = append(b.subs[postID], ch)
	logrus.Debugf("subscribed to post %v", postID)

	return ch
}

func (b *Broker) Unsubscribe(postID int, ch <-chan *model.Comment) {
	logrus.Debugf("unsubscribing from post %v", postID)

	b.mu.Lock()

	defer b.mu.Unlock()
	defer logrus.Debugf("unsubscribed from post %v", postID)

	subs := b.subs[postID]
	for i, sub := range subs {
		if sub == ch {
			close(sub)
			b.subs[postID] = append(subs[:i], subs[i+1:]...)
			break
		}
	}
}

func (b *Broker) Publish(postID int, comment *model.Comment) {
	logrus.Debugf("publishing comment %s to post %v", comment.Content, postID)

	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, ch := range b.subs[postID] {
		select {
		case ch <- comment:
			logrus.Debugf("published comment %s to post %v", comment.Content, postID)
		default:
			logrus.Errorf("failed to publish comment %s to post %v", comment.Content, postID)
		}
	}
}
