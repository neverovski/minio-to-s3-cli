# minio-to-s3-cli

## Usage

```bash
go run main.go -bucket <string> [-minio-region <string>] [-minio-endpoint <string>] [-minio-access-key <string>] [-minio-secret-key <string>] [-s3-region <string>] [-s3-endpoint <string>] [-s3-access-key <string>] [-s3-secret-key <string>]

  -bucket [string]
        The name of the S3 bucket to list objects from.
  -minio-region [string]
        The region of the MinIO server.
  -minio-endpoint [string]
        The endpoint of the MinIO server.
  -minio-access-key [string]
        The access key for the MinIO server.
  -minio-secret-key [string]
        The secret key for the MinIO server.
  -s3-region [string]
        The region of the S3 server.
  -s3-endpoint [string]
        The endpoint of the S3 server.
  -s3-access-key [string]
        The access key for the S3 server.
  -s3-secret-key [string]
        The secret key for the S3 server.
```
