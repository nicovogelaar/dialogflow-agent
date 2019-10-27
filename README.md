# dialogflow-agent

## Build

```bash
go build -o dialogflow-agent .
```

## Usage

Import entities:
```bash
./dialogflow-agent \
  --project-id example-123 \
  --credentials-file ./credentials.json \
  entities import \
  -f examples/entities.yaml
```

Delete all entities:
```bash
./dialogflow-agent \
  --project-id example-123 \
  --credentials-file ./credentials.json \
  entities delete -a
```

Import intents:
```bash
./dialogflow-agent \
  --project-id example-123 \
  --credentials-file ./credentials.json \
  intent import \
  -f examples/intents.yaml
```

Delete all intents:
```bash
./dialogflow-agent \
  --project-id example-123 \
  --credentials-file ./credentials.json \
  intents delete -a
```
