package gss

type SubscribeEvent struct {
	Channel string
	ConnID  string
	User    string
}

type UnsubscribeEvent struct {
	ConnID string
	User   string
}

func NewSubscribeEvent(channel, connid, user string) *SubscribeEvent {
	return &SubscribeEvent{Channel: channel, ConnID: connid, User: user}
}

func NewUnsubscribeEvent(connid, user string) *UnsubscribeEvent {
	return &UnsubscribeEvent{ConnID: connid, User: user}
}
