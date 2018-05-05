package analytics

import "context"

type EventName string

type trackingContextKeyType string

const trackingContextKey trackingContextKeyType = "tracking context"

func ContextWithTrackingProps(ctx context.Context, props map[string]interface{}) context.Context {
	return context.WithValue(ctx, trackingContextKey, props)
}

func getTrackingProps(ctx context.Context) map[string]interface{} {
	tp := ctx.Value(trackingContextKey)
	if tp == nil {
		return map[string]interface{}{}
	}

	return tp.(map[string]interface{})
}

func ContextWithEventPropsCollector(ctx context.Context, name EventName) context.Context {
	return context.WithValue(ctx, name, map[string]interface{}{})
}

func SaveEventProp(ctx context.Context, name EventName, key string, value interface{}) {
	ec := ctx.Value(name).(map[string]interface{})
	ec[key] = value
}

func SaveEventProps(ctx context.Context, name EventName, props map[string]interface{}) {
	ec := ctx.Value(name).(map[string]interface{})

	for k, v := range props {
		ec[k] = v
	}
}
