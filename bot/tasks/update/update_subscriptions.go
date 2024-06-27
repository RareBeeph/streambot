package update

import (
	"streambot/models"
	"streambot/query"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func UpdateSubscriptions(s *discordgo.Session, minHealth int, maxHealth int) {
	m := query.Message
	qs := query.Subscription

	// Load all active subscriptions
	subscriptions, err := qs.
		Preload(qs.Messages.Order(m.PostOrder.Asc())).
		GetByHealth(minHealth, maxHealth)
	if err != nil {
		log.Err(err).Msg("Failed to find subscriptions")
	}

	var wg sync.WaitGroup
	wg.Add(len(subscriptions))
	for _, sub := range subscriptions {
		go func(sub *models.Subscription) {
			defer wg.Done()
			UpdateSubscription(s, sub)
		}(sub)
	}

	wg.Wait()
}

func UpdateSubscription(s *discordgo.Session, sub *models.Subscription) {
	qs := query.Subscription

	actions, err := msgActionsForSub(sub)
	if err != nil {
		return
	}

	for idx, action := range actions {
		// if errored, keep performing actions, but don't unset the error
		// another option is to just break on the first error of the subscription
		actionErr := action.perform(s, sub, idx)
		if err == nil {
			err = actionErr
		}
	}

	if err != nil {
		// propagate failure count
		qs.Where(qs.ID.Eq(sub.ID)).Update(qs.TimesFailed, qs.TimesFailed.Add(1))
		return
	}

	// on success, reset failure count
	qs.Where(qs.ID.Eq(sub.ID)).Update(qs.TimesFailed, 0)
}
