# go-links
go-links is a custom url-shortener service that allows users to easily manage their own links. There are lots of options, but I wanted to build my own to make sure it supported my specific use cases.

Once installed, you will want to your DNS Search Suffix to include the domain you installed to. IE: If you're FQDN is `go.mysite.com`, you would want to set your DNS Search Suffix to `mysite.com`. Now when you type `http://go` into your browser, you should be directed to your site.

## SSL
Unless you want to sign all of your own certificates with your own root CA which then has to be trusted on all your devices, your SSL cert will most likely have to use an FQDN (`go.mysite.com`). To help with this, the service redirects all traffic from `http://go` to the value of the `FQDN` environment variable. This will allow you to use a certificate covered endpoint for things like SAML.

## Setup
### Kubernetes
```shell
helm repo add golinks https://imdevinc.github.io/go-links
helm repo update
helm install golinks -n golinks --create-namespace --set config.fqdn=go.mysite.com
```

## Config
There are some required and some optional config values depending on how you want to run the app.

| Config Option             | Environment Variable | Required | Description                                                                                                                                 | Example             | Default                   |
| ------------------------- | -------------------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------- | ------------------- | ------------------------- |
| `port`                    | `PORT`               | false    | Port used for app communication                                                                                                             | `8080`              | `8080`                    |
| `fqdn`                    | `FQDN`               | true     | FQDN used for redirects                                                                                                                     | `app.example.com`   | n/a                       |
| `storeType`               | `STORE_TYPE`         | false    | Type of storage to use. See [storeType](#storetype) for available options                                                                   | `mongo`             | `memory`                  |
| `mongo.username`          | `MONGO_USERNAME`     | false    | The username for the mongodb connection                                                                                                     | `mongoUser`         | n/a                       |
| `mongo.password`          | `MONGO_PASSWORD`     | false    | The password for the mongodb connection                                                                                                     | `mySecretPassword`  | n/a                       |
| `mongo.host`              | `MONGO_HOST`         | false    | The host used for the mongodb connection                                                                                                    | `mongodb.mongodb`   | n/a                       |
| `mongo.dbname`            | `MONGO_DB_NAME`      | false    | The database name used for the mongodb connection                                                                                           | `links`             | n/a                       |
| `postgres.username`       | `POSTGRES_USERNAME`  | false    | The username for the postgres connection                                                                                                    | `postgresUser`      | n/a                       |
| `postgres.password`       | `POSTGRES_PASSWORD`  | false    | The password for the postgres connection                                                                                                    | `mySecretPassword`  | n/a                       |
| `postgres.host`           | `POSTGRES_HOST`      | false    | The host used for the postgres connection                                                                                                   | `postgres.postgres` | n/a                       |
| `postgres.dbname`         | `POSTGRES_DB_NAME`   | false    | The database name used for the postgres connection                                                                                          | `links`             | n/a                       |
| `ssoEntityId`             | `SSO_ENTITY_ID`      | false    | The entity ID used for  SAML authentication                                                                                                 | `golinks`           | n/a                       |
| `ssoRequire`              | `SSO_REQUIRE`        | false    | If set to true and SAML auth is misconfigured, will not allow the service to startup                                                        | false               | `true`                    |
| `ssoMetadataFileContents` | n/a                  | false    | Sets the metadata XML file content (only used in Helm chart, see [SAML configuration](#saml-authentication) section below for more details) | false               | `<?xml version="1.0">...` |

## StoreType
go-links supports multiple storage types depending on your use case

### `memory`
> [!WARNING]
> Memory store will erase on every service restart, it is not recommended for production workloads and should only be used for testing

The memory store is the default store type and stores all data in memory. **When you restart the service, all existing data will be lost.**

### `file`
The file memory store keeps all data in a local file stored at `./links.json`. This store is not recommended for high traffic use cases as the query methods will be slower at high traffic levels.

### `mongo`
The mongo store type stores data in a [mongodb database](https://www.mongodb.com/). The following config values are required when using the mongo store type:
- `mongo.username`
- `mongo.password`
- `mongo.host`
- `mongo.dbname`
It is expected that the database will already exist, and the service will also create the appropriate indexes if they haven't been created yet.

### `postgres`
The postgres store type stores data in a postgres database. The following config values are required when using the postgres store type:
- `postgres.username`
- `postgres.password`
- `postgres.host`
- `postgres.dbname`
It is expected that the database `postgres.dbname` will already exist and the service will create the needed table if it does not exist.

## SAML Authentication
You can configure SAML authentication for the service to only allow specific actors. The following configuration options are required for SAML authentication to work properly:
- `ssoEntityId`
- `ssoMetadataFileContents`
You can also set `ssoRequire: true`  if you want the service to fail when SAML auth is not configured correctly.

Using SAML authentication also tracks who created each link for auditing. 