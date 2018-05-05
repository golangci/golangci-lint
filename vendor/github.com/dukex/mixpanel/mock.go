package mixpanel

import (
	"errors"
	"fmt"
	"time"
)

// Mocked version of Mixpanel which can be used in unit tests.
type Mock struct {
	// All People identified, mapped by distinctId
	People map[string]*MockPeople
}

func NewMock() *Mock {
	return &Mock{
		People: map[string]*MockPeople{},
	}
}

func (m *Mock) String() string {
	str := ""
	for id, p := range m.People {
		str += id + ":\n" + p.String()
	}
	return str
}

// Identifies a user. The user will be added to the People map.
func (m *Mock) people(distinctId string) *MockPeople {
	p := m.People[distinctId]
	if p == nil {
		p = &MockPeople{
			Properties: map[string]interface{}{},
		}
		m.People[distinctId] = p
	}

	return p
}

func (m *Mock) Track(distinctId, eventName string, e *Event) error {
	p := m.people(distinctId)
	p.Events = append(p.Events, MockEvent{
		Event: *e,
		Name:  eventName,
	})
	return nil
}

type MockPeople struct {
	Properties map[string]interface{}
	Time       *time.Time
	IP         string
	Events     []MockEvent
}

func (mp *MockPeople) String() string {
	timeStr := ""
	if mp.Time != nil {
		timeStr = mp.Time.Format(time.RFC3339)
	}

	str := fmt.Sprintf("  ip: %s\n  time: %s\n", mp.IP, timeStr)
	str += "  properties:\n"
	for key, val := range mp.Properties {
		str += fmt.Sprintf("    %s: %v\n", key, val)
	}
	str += "  events:\n"
	for _, event := range mp.Events {
		str += "    " + event.Name + ":\n"
		str += fmt.Sprintf("      IP: %s\n", event.IP)
		if event.Timestamp != nil {
			str += fmt.Sprintf(
				"      Timestamp: %s\n", event.Timestamp.Format(time.RFC3339),
			)
		} else {
			str += "      Timestamp:\n"
		}
		for key, val := range event.Properties {
			str += fmt.Sprintf("      %s: %v\n", key, val)
		}
	}
	return str
}

func (m *Mock) Update(distinctId string, u *Update) error {
	p := m.people(distinctId)

	if u.IP != "" {
		p.IP = u.IP
	}
	if u.Timestamp != nil && u.Timestamp != IgnoreTime {
		p.Time = u.Timestamp
	}

	switch u.Operation {
	case "$set":
		for key, val := range u.Properties {
			p.Properties[key] = val
		}
	default:
		return errors.New("mixpanel.Mock only supports the $set operation")
	}

	return nil
}

func (m *Mock) Alias(distinctId, newId string) error {
	return nil
}

type MockEvent struct {
	Event
	Name string
}
