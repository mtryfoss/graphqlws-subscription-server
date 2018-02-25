package gss

type SubscribeEvent struct {
	Channel        string
	ConnID         string
	SubscriptionID string
	User           string
}

type UnsubscribeEvent struct {
	ConnID string
	User   string
}

func NewSubscribeEvent(channel, connid, subid, user string) *SubscribeEvent {
	return &SubscribeEvent{Channel: channel, ConnID: connid, SubscriptionID: subid, User: user}
}

func NewUnsubscribeEvent(connid, user string) *UnsubscribeEvent {
	return &UnsubscribeEvent{ConnID: connid, User: user}
}
