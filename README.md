## DESCRIPTION

```text
This application takes commands from HR system, such as CREATE_USER, UPDATE_USER and stores those 
users in keycloak service. Keycloak plays role of IAM service for alif-bank employees.

If one wants to have mentioned users in his database, then there are few options:
- implement OAuth2 protocol and you are done
- create your own queue in RabbitMQ and stay up to day in async way (see User sync)
```

## User sync

```ts
// NOTE:
// USER_REPRESENTATION = keycloak.org/docs-api/18.0/rest-api/index.html#_userrepresentation
TOPIC_NAME: `keycloak.user_created`
PAYLOAD: USER_REPRESENTATION

EMITS: `keycloak.user_updated`
PAYLOAD: USER_REPRESENTATION
```

## Getting started

```shell
# clone repo
mkdir -p $GOPATH/src/github.com/alifcapital/
cd $GOPATH/src/github.com/alifcapital/
git clone git@github.com:alifcapital/keycloack_module.git
cd keycloack_module

# install go dependencies
go mod tidy

# init app dependencies
make docker-up

# start listening to command
make start

```

## Contributing

Pull requests are welcome. For any changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.
