## Send email microservices example

### Usage

```docker
docker-compose up -d
```

Example api call

```bash
curl --request POST \
  --url http://localhost:8000 \
  --header 'authorization: API_KEY' \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data To=email@example.com \
  --data From=email@example.com \
  --data 'Subject=Hello World' \
  --data 'Body=Flameo Hotmans!' \
  --data Host=smtp.example.com \
  --data Port=587 \
  --data Username=email@example.com\
  --data 'very secure password'
```

Make sure to not call api from client side application without refactoring the producer with necessary security.

### Replication

You can replicate the micro services if necessary.

```bash
docker-compose up --scale sendemail-consumer=5
```
