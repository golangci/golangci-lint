package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
)

type SessionBuilder interface {
	Config(*aws.Config) SessionBuilder
	Profile(string) SessionBuilder
	Options(*session.Options) SessionBuilder
	Endpoint(string) SessionBuilder
	S3ForcePathStyle(bool) SessionBuilder
	Build() *session.Session
}

type sessionBuilder struct {
	awsConfig      *aws.Config
	profile        string
	options        *session.Options
	endpoint       *string
	forcePathStyle *bool
}

func (sb *sessionBuilder) Config(c *aws.Config) SessionBuilder {
	sb.awsConfig = c
	return sb
}

func (sb *sessionBuilder) Profile(p string) SessionBuilder {
	sb.profile = p
	return sb
}

func (sb *sessionBuilder) Endpoint(e string) SessionBuilder {
	sb.endpoint = aws.String(e)
	return sb
}

func (sb *sessionBuilder) S3ForcePathStyle(b bool) SessionBuilder {
	sb.forcePathStyle = aws.Bool(b)
	return sb
}

func (sb *sessionBuilder) Options(o *session.Options) SessionBuilder {
	sb.options = o
	return sb
}

func (sb *sessionBuilder) Build() *session.Session {
	if sb.awsConfig == nil {
		sb.awsConfig = aws.NewConfig()
	}

	if sb.endpoint != nil {
		sb.awsConfig.Endpoint = sb.endpoint
		sb.awsConfig.S3ForcePathStyle = sb.forcePathStyle
	}

	sb.awsConfig.Credentials = credentials.NewChainCredentials([]credentials.Provider{
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{
			Profile: sb.profile,
		},
	})

	_, err := sb.awsConfig.Credentials.Get()
	if err == nil {
		return session.Must(session.NewSession(sb.awsConfig))
	}
	if sb.options == nil {
		sb.options = &session.Options{
			AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
			SharedConfigState:       session.SharedConfigEnable,
			Profile:                 sb.profile,
		}
	}

	return session.Must(session.NewSessionWithOptions(*sb.options))
}

func newSessionBuilder() SessionBuilder {
	return &sessionBuilder{}
}
