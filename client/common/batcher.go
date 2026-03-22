package common

import (
	"strings"
)

const batchSeparator = "#"

type Batcher struct {
	batch     string
	maxAmount int
}

func NewBatcher(maxAmount int) *Batcher {
	return &Batcher{
		batch:     "",
		maxAmount: maxAmount,
	}
}

func (b *Batcher) CanAddBet(bet Bet) bool {
	if b.batch == "" {
		return 1 <= b.maxAmount
	}
	current := strings.Count(b.batch, batchSeparator) + 1
	return current+1 <= b.maxAmount
}

func (b *Batcher) AddBet(bet Bet) {
	serializeBet := SerializeBet(bet)
	if b.batch == "" {
		b.batch = serializeBet
	} else {
		b.batch = strings.Join([]string{b.batch, serializeBet}, batchSeparator)
	}
}

func (b *Batcher) Clean(bet Bet) {
	b.batch = SerializeBet(bet)
}

func (b *Batcher) GetBatch() string {
	return b.batch
}

func (b *Batcher) IsEmpty() bool {
	return len(b.batch) == 0
}
