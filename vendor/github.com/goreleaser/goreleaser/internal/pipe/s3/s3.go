// Package s3 provides a Pipe that push artifacts to s3/minio
package s3

import (
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/goreleaser/goreleaser/internal/artifact"
	"github.com/goreleaser/goreleaser/internal/deprecate"
	"github.com/goreleaser/goreleaser/internal/pipe"
	"github.com/goreleaser/goreleaser/internal/semerrgroup"
	"github.com/goreleaser/goreleaser/internal/tmpl"
	"github.com/goreleaser/goreleaser/pkg/config"
	"github.com/goreleaser/goreleaser/pkg/context"
)

// Pipe for Artifactory
type Pipe struct{}

// String returns the description of the pipe
func (Pipe) String() string {
	return "S3"
}

// Default sets the pipe defaults
func (Pipe) Default(ctx *context.Context) error {
	for i := range ctx.Config.S3 {
		s3 := &ctx.Config.S3[i]
		if s3.Bucket == "" {
			continue
		}
		deprecate.Notice("s3")
		if s3.Folder == "" {
			s3.Folder = "{{ .ProjectName }}/{{ .Tag }}"
		}
		if s3.Region == "" {
			s3.Region = "us-east-1"
		}
		if s3.ACL == "" {
			s3.ACL = "private"
		}
	}
	return nil
}

// Publish to S3
func (Pipe) Publish(ctx *context.Context) error {
	if len(ctx.Config.S3) == 0 {
		return pipe.Skip("s3 section is not configured")
	}
	var g = semerrgroup.New(ctx.Parallelism)
	for _, conf := range ctx.Config.S3 {
		conf := conf
		g.Go(func() error {
			return upload(ctx, conf)
		})
	}
	return g.Wait()
}

func newS3Svc(conf config.S3) *s3.S3 {
	builder := newSessionBuilder()
	builder.Profile(conf.Profile)
	if conf.Endpoint != "" {
		builder.Endpoint(conf.Endpoint)
		builder.S3ForcePathStyle(true)
	}
	sess := builder.Build()

	return s3.New(sess, &aws.Config{
		Region: aws.String(conf.Region),
	})
}

func upload(ctx *context.Context, conf config.S3) error {
	var svc = newS3Svc(conf)

	template := tmpl.New(ctx)
	bucket, err := template.Apply(conf.Bucket)
	if err != nil {
		return err
	}

	folder, err := template.Apply(conf.Folder)
	if err != nil {
		return err
	}

	var filter = artifact.Or(
		artifact.ByType(artifact.UploadableArchive),
		artifact.ByType(artifact.UploadableBinary),
		artifact.ByType(artifact.Checksum),
		artifact.ByType(artifact.Signature),
		artifact.ByType(artifact.LinuxPackage),
	)
	if len(conf.IDs) > 0 {
		filter = artifact.And(filter, artifact.ByIDs(conf.IDs...))
	}

	var g = semerrgroup.New(ctx.Parallelism)
	for _, artifact := range ctx.Artifacts.Filter(filter).List() {
		artifact := artifact
		g.Go(func() error {
			f, err := os.Open(artifact.Path)
			if err != nil {
				return err
			}
			log.WithFields(log.Fields{
				"bucket":   bucket,
				"folder":   folder,
				"artifact": artifact.Name,
			}).Info("uploading")
			_, err = svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(filepath.Join(folder, artifact.Name)),
				Body:   f,
				ACL:    aws.String(conf.ACL),
			})
			return err
		})
	}
	return g.Wait()
}
