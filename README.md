# s3dbdump

A tool to dump a MariaDB (MySQL) database to a file and upload it to S3 or MinIO, with gzip compression.

## Usage

### Environment variables

| Environment Variable     | Required | Default Value             | Description                         |
| ------------------------ | -------- | ------------------------- | ----------------------------------- |
| `AWS_ACCESS_KEY_ID`      | Yes      | -                         | AWS access key ID                   |
| `AWS_SECRET_ACCESS_KEY`  | Yes      | -                         | AWS secret access key               |
| `AWS_REGION`             | Yes      | -                         | AWS region                          |
| `S3_BUCKET`              | Yes      | -                         | S3 bucket name                      |
| `S3_ENDPOINT`            | No       | -                         | Custom S3 endpoint (e.g. for MinIO) |
| `DB_HOST`                | Yes      | -                         | Database host                       |
| `DB_PORT`                | No       | 3306                      | Database port                       |
| `DB_USER`                | Yes      | -                         | Database user                       |
| `DB_PASSWORD`            | Yes      | -                         | Database password                   |
| `DB_NAME`                | Yes      | -                         | Database name to dump               |
| `DB_ALL_DATABASES`       | No       | 0                         | Set to 1 to dump all databases      |
| `DB_GZIP`                | No       | 1                         | Enable gzip compression             |
| `DB_DUMP_PATH`           | No       | ./dumps                   | Directory to store dumps            |
| `DB_DUMP_FILENAME`       | No       | %s-20060102T150405.sql.gz | Dump filename format                |
| `DB_DUMP_FILE_KEEP_DAYS` | No       | 7                         | Number of days to keep backups      |
