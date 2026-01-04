# cashflow-assesment-task

---

## Run the stack

```bash
docker compose up -d
```

This starts the API, workers, RabbitMQ, and Postgres.

---

## Database Migrations

1. Install goose if not already installed:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

2. Navigate to migrations folder:

```bash
cd payment-service/migrations
```

3. Run migrations: using localhost for convenient migration

```bash
goose postgres "postgres://postgres:strongpassword@localhost:5432/app_db?sslmode=disable" up
```

- Use `db` as the host when running via Docker Compose.
- Use `localhost` for a local Postgres instance.

---

## Scale Payment Workers

Start multiple workers:

```bash
docker compose up -d --scale payment-worker=3
```

RabbitMQ load-balances messages across all workers.
Redilevery happens in cases of transient errors.

---

## Retry testing

Go to the env file in payment worker to set max retry for DLX and duration gaps to lower cpu spike, ideally 10 seconds.
Logs are used extensively so please inspect container logs to check states

---

## API

### Create Payment

```http
POST /payment
Content-Type: application/json

{
  "ref": "123",
  "amount": "20000.50",
  "currency": "ETB"
}
```

- `amount` must be a string with `.` as the decimal separator; commas are not allowed, and it's parsed to int64 considering minor cent values.
- `ref` can be any random string.
- Payment is processed asynchronously; initial status is usually `PENDING`.

Status values:

| Status | Meaning |
| ------ | ------- |
| -1     | FAILED  |
| 0      | PENDING |
| 1      | SUCCESS |

Response example:

```json
{
  "payment_id": "83071c68-31b0-423c-92be-dc7882c58f4b",
  "status": 0
}
```

### Get Payment

```http
GET /payment/83071c68-31b0-423c-92be-dc7882c58f4b
```

Response example:

```json
{
  "amount": "20000.50",
  "currency": "ETB",
  "ref": "123",
  "status": 1,
  "created_at": "2026-01-04T16:20:51.082206Z"
}
```
