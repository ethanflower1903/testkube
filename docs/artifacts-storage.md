# Artifacts Storage

Testkube supports test artifacts collection.

Currently, only the Cypress executor job produces test artifacts. The executor will scrape the files and store them in [Minio](https://min.io/).  The executor will create a bucket named by execution ID and collect all files that are stored in the Cypress artifacts location `Cypress/`.

The available configuration parameters in Helm charts are:

| Parameter                          | Is optional | Default                              | Default                                              |
| ---------------------------------- | ----------- | ------------------------------------ | ---------------------------------------------------- |
| api-server.storage.endpoint        | yes         | testkube-minio-service-testkube:9000 | URL of the S3 bucket                                 |
| api-server.storage.accessKeyId     | yes         | minio                                | Access Key ID                                        |
| api-server.storage.accessKey       | yes         | minio123                             | Access Key                                           |
| api-server.storage.location        | yes         |                                      | Region                                               |
| api-server.storage.token           | yes         |                                      | S3 Token                                             |
| api-server.storage.SSL             | yes         | false                                | Indicates whether SSL communication is to be enabled |
| api-server.storage.scrapperEnabled | yes         | true                                 | Indicates whether executors should scrape artifacts  |

The API Server accepts the following environment variables:

```sh
STORAGE_ENDPOINT
STORAGE_ACCESSKEYID
STORAGE_SECRETACCESSKEY
STORAGE_LOCATION
STORAGE_TOKEN 
STORAGE_SSL
SCRAPPERENABLED
```

Which can be set while installing with Helm:

```bash
helm install --create-namespace my-testkube testkube/testkube --set STORAGE_ENDPOINT=custom_value
```

Alternatively, these values can be read from Kubernetes secrets and set:

```yaml
- env:
 - name: STORAGE_ENDPOINT
   secret:
  secretName: test-secret
```
