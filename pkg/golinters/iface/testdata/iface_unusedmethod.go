//golangcitest:args -Eiface
//golangcitest:config_path testdata/iface_unusedmethod.yml
package testdata

type MessageSender interface {
	SendMessage(username, msg string) error
	MustSendMessage(username, msg string) // want "method 'MustSendMessage\\(\\)' is declared on interface 'MessageSender' but not used within the package"

	// Close the resource.
	Close() error // want "method 'Close\\(\\)' is declared on interface 'MessageSender' but not used within the package"
}

type Service struct {
	sender MessageSender
}

func NewService(sender MessageSender) *Service {
	return &Service{
		sender: sender,
	}
}

func (s *Service) GreetMorning(username string) error {
	return s.sender.SendMessage(username, "Good morning")
}
