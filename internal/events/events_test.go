package events

import (
	"context"
	"testing"
)

func TestDispatch(t *testing.T) {
	b := &Bus{handlers: make(map[string][]Handler)}

	var got []string
	b.Subscribe(TypeLevelUp, func(ctx context.Context, e Event) {
		got = append(got, "first")
	})
	b.Subscribe(TypeLevelUp, func(ctx context.Context, e Event) {
		got = append(got, "second")
	})
	b.Subscribe(TypeGithubSynced, func(ctx context.Context, e Event) {
		got = append(got, "wrong-type")
	})

	b.dispatch(context.Background(), Event{UserID: 1, Type: TypeLevelUp})

	if len(got) != 2 || got[0] != "first" || got[1] != "second" {
		t.Errorf("dispatched handlers = %v, want [first second]", got)
	}
}

func TestDispatch_NoSubscribers(t *testing.T) {
	b := &Bus{handlers: make(map[string][]Handler)}
	b.dispatch(context.Background(), Event{UserID: 1, Type: TypeEquipmentChanged})
}
