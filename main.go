package main

import (
	"bytes"
	"flag"
	"io"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	minioConfig MinioConfig
	s3Config    S3Config
)

type MinioConfig struct {
	Bucket    string
	Region    string
	Endpoint  string
	AccessKey string
	SecretKey string
}

type S3Config struct {
	Bucket    string
	Region    string
	Endpoint  string
	AccessKey string
	SecretKey string
}

func init() {
	flag.StringVar(&minioConfig.Bucket, "minio-bucket", "", "Bucket of the MinIO server.")
	flag.StringVar(&minioConfig.Region, "minio-region", "", "Region of the MinIO server.")
	flag.StringVar(&minioConfig.Endpoint, "minio-endpoint", "", "Endpoint of the MinIO server.")
	flag.StringVar(&minioConfig.AccessKey, "minio-access-key", "", "Access key for the MinIO server.")
	flag.StringVar(&minioConfig.SecretKey, "minio-secret-key", "", "Secret key for the MinIO server.")

	flag.StringVar(&s3Config.Bucket, "s3-bucket", "", "Bucket of the S3 server.")
	flag.StringVar(&s3Config.Region, "s3-region", "", "Region of the S3 server.")
	flag.StringVar(&s3Config.Endpoint, "s3-endpoint", "", "Endpoint of the S3 server.")
	flag.StringVar(&s3Config.AccessKey, "s3-access-key", "", "Access key for the S3 server.")
	flag.StringVar(&s3Config.SecretKey, "s3-secret-key", "", "Secret key for the S3 server.")
}

func main() {
	flag.Parse()

	validateMinioConfig()
	validateS3Config()

	minioSession, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(minioConfig.AccessKey, minioConfig.SecretKey, ""),
		Region:           aws.String(minioConfig.Region),
		Endpoint:         aws.String(minioConfig.Endpoint),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	})
	minioParams := &s3.ListObjectsV2Input{Bucket: &minioConfig.Bucket}

	if err != nil {
		log.Fatalln("Failed to create AWS session:", err)
	}

	awsSession, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(s3Config.AccessKey, s3Config.SecretKey, ""),
		Region:           aws.String(s3Config.Region),
		Endpoint:         aws.String(s3Config.Endpoint),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	})

	if err != nil {
		log.Fatalln("Failed to create AWS session:", err)
	}

	// Create service client
	minioClient := s3.New(minioSession)
	awsClient := s3.New(awsSession)

	err = minioClient.ListObjectsV2Pages(minioParams, func(output *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range output.Contents {

			object, err := minioClient.GetObject(&s3.GetObjectInput{
				Bucket: aws.String(minioConfig.Bucket),
				Key:    obj.Key,
			})

			if err != nil {
				log.Println("Failed to get object from MinIO:", err)
				continue
			}

			bodyBytes, err := io.ReadAll(object.Body)

			if err != nil {
				log.Println("Failed to read object body:", err)
				continue
			}

			object.Body.Close()

			_, err = awsClient.PutObject(&s3.PutObjectInput{
				Bucket:      aws.String(s3Config.Bucket),
				Key:         obj.Key,
				Metadata:    object.Metadata,
				ContentType: object.ContentType,
				Body:        bytes.NewReader(bodyBytes),
			})
			if err != nil {
				log.Println("Failed to upload object to S3:", err)
				continue
			}

			log.Println("Object", *obj.Key, "transferred successfully from MinIO to S3.")
		}

		return true
	})

	if err != nil {
		log.Println("Failed to get object from MinIO:", err)
		return
	}

	log.Println("All objects transferred successfully from MinIO to S3.")
}

func validateMinioConfig() {
	if minioConfig.Bucket == "" {
		flag.PrintDefaults()
		log.Fatalf("Invalid parameters, minio-bucket required")
	}

	if minioConfig.Region == "" {
		flag.PrintDefaults()
		log.Fatalf("Invalid parameters, minio-region required")
	}

	if minioConfig.Endpoint == "" {
		flag.PrintDefaults()
		log.Fatalf("Invalid parameters, minio-endpoint required")
	}

	if minioConfig.AccessKey == "" {
		flag.PrintDefaults()
		log.Fatalf("Invalid parameters, minio-access-key required")
	}

	if minioConfig.SecretKey == "" {
		flag.PrintDefaults()
		log.Fatalf("Invalid parameters, minio-secret-key required")
	}
}

func validateS3Config() {
	if s3Config.Bucket == "" {
		flag.PrintDefaults()
		log.Fatalf("Invalid parameters, s3-bucket required")
	}

	if s3Config.Region == "" {
		flag.PrintDefaults()
		log.Fatalf("Invalid parameters, s3-region required")
	}

	if s3Config.Endpoint == "" {
		flag.PrintDefaults()
		log.Fatalf("Invalid parameters, s3-endpoint required")
	}

	if s3Config.AccessKey == "" {
		flag.PrintDefaults()
		log.Fatalf("Invalid parameters, s3-access-key required")
	}

	if s3Config.SecretKey == "" {
		flag.PrintDefaults()
		log.Fatalf("Invalid parameters, s3-secret-key required")
	}
}
