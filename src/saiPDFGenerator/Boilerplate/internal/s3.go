package internal

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.uber.org/zap"
)

// init aws config from config
func (is *InternalService) awsConfig() {
	region := is.Context.GetConfig("specific.aws.region", "defaultRegion").(string)

	keyID := is.Context.GetConfig("specific.aws.key_id", "defaultKeyID").(string)

	secret := is.Context.GetConfig("specific.aws.secret", "123455").(string)

	token := is.Context.GetConfig("specific.aws.token", "").(string)

	config := &aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(keyID, secret, token),
	}

	is.Logger.Debug("AWS config initiated", zap.String("region", region), zap.String("keyID", keyID), zap.String("Token", token))

	is.AwsConfig = config
}

// upload file to s3, return link to file
func (is *InternalService) s3Upload(file []byte) (string, error) {
	s3Session, err := session.NewSession(is.AwsConfig)
	if err != nil {
		return "", fmt.Errorf("s3 - newSession - %w", err)
	}

	uploader := s3manager.NewUploader(s3Session)

	timestamp := strconv.Itoa(int(time.Now().Unix()))
	filename := fmt.Sprintf("%s.pdf", timestamp)

	input := &s3manager.UploadInput{
		Bucket:      aws.String(is.Context.GetConfig("specific.aws.bucket", "defaultBucket").(string)), // bucket's name
		Key:         aws.String(filename),                                                              // files destination location
		Body:        bytes.NewReader(file),                                                             // content of the file
		ContentType: aws.String("application/pdf"),                                                     // content type
	}
	output, err := uploader.UploadWithContext(context.Background(), input)
	if err != nil {
		return "", fmt.Errorf("s3 - upload - %w", err)
	}
	return output.Location, nil

}
