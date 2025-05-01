# REST API

For authenticating the call, client is expected to submit basic authentication using a
 predefined username and password. You can lookup the value of this key from `CLIENT_USERNAME` and `CLIENT_PASSWORD` environment variable.

**Table of contents:**

- [REST API](#rest-api)
  - [Get All Messages](#get-all-messages)
  - [Send Message](#send-message)
  - [Retry Message](#retry-message)

## Get All Messages

GET: `/messages`

This endpoint is used to get all messages from the system including the status of the message.

**Headers:**

| Field           | Type   | Required | Description                                           |
| --------------- | ------ | -------- | ----------------------------------------------------- |
| `Authorization` | String | Yes      | The Basic Authentication for authenticating the call. |
| `Content-Type`  | String | Yes      | The only accepted value is `application/json`.        |

**Example Call:**

```json
GET /messages
Authorization: Basic admF6bGFicy5jb206cGFzc3dvcmQ=
Content-Type: application/json
```

**Success Response:**

```json
HTTP/1.1 200 OK
Content-Type: application/json

{
    "ok": true,
    "messages": [
        // scheduled message, should be sent at `scheduled_sending_at`
        {
            "id": "1da2f3e4-5b6c-7d8e-9a0b-c1d2e3f4g5h6",
            "message": "Job alert for Software Engineer at Invertase...",
            "scheduled_sending_at": 1735432224,
            "sent_at": null,
            "retried_count": 0,
            "status": "scheduled",
            "reason": null,
            "created_at": 1735432224,
            "updated_at": 1735432224
        },
        // successfully sent message, `sent_at` is set
        {
            "id": "2b3c4d5e-6f7g-8h9i-0j1k-l2m3n4o5p6q7",
            "message": "Job alert for Software Engineer at dev.to...",
            "scheduled_sending_at": 1735432224,
            "sent_at": 1735432224,
            "retried_count": 0,
            "status": "sent",
            "reason": null,
            "created_at": 1735432224,
            "updated_at": 1735432224
        },
        // has been retried and expected to be sent
        {
            "id": "2b3c4d5e-6f7g-8h9i-0j1k-l2m3n4o5p6q7",
            "message": "Job alert for Software Engineer at dev.to...",
            "scheduled_sending_at": 1735432224,
            "sent_at": null,
            "retried_count": 1,
            "status": "scheduled",
            "reason": null,
            "created_at": 1735432224,
            "updated_at": 1735432224
        },
        // failed message, 
        {
            "id": "2b3c4d5e-6f7g-8h9i-0j1k-l2m3n4o5p6q7",
            "message": "Job alert for Software Engineer at dev.to...",
            "scheduled_sending_at": 1735432224,
            "sent_at": null,
            "retried_count": 3,
            "status": "failed",
            "reason": "session expired",
            "created_at": 1735432224,
            "updated_at": 1735432224
        }
    ],
    "ts": 1735432224
}
```

[Back to Top](#rest-api)

## Send Message

POST: `/messages`

This endpoint is used to send a scheduled message to Whatsapp.

**Headers:**

| Field           | Type   | Required | Description                                           |
| --------------- | ------ | -------- | ----------------------------------------------------- |
| `Authorization` | String | Yes      | The Basic Authentication for authenticating the call. |
| `Content-Type`  | String | Yes      | The only accepted value is `application/json`.        |

**Body Payload:**

| Field               | Type            | Required | Description                                            |
| ------------------- | --------------- | -------- | ------------------------------------------------------ |
| `recipient_numbers` | Array of string | Yes      | The list of recipient numbers.                         |
| `message`           | String          | Yes      | The message to be sent.                                |
| `scheduled_sending_at`        | Number          | Yes      | The Unix timestamp of when the message should be sent. |

**Example Call:**

```json
POST /vacancies
Authorization: Basic admF6bGFicy5jb206cGFzc3dvcmQ=
Content-Type: application/json

{
    "recipient_numbers": [
        "120363352351961274@g.us",
        "120363352351961275@g.us"
    ],
    "message": "Job alert for Software Engineer at Invertase...",
    "scheduled_sending_at": 1735432224
}
```

**Success Response:**

```json
HTTP/1.1 200 OK
Content-Type: application/json

{
  "ok": true,
  "ts": 1735432224
}
```

[Back to Top](#rest-api)

## Retry Message

POST: `/messages/{id}/retry`

This endpoint is used to retry a failed and halted message. The message will be retried immediately or based on the `scheduled_sending_at` timestamp.

**Headers:**

| Field           | Type   | Required | Description                                           |
| --------------- | ------ | -------- | ----------------------------------------------------- |
| `Authorization` | String | Yes      | The Basic Authentication for authenticating the call. |
| `Content-Type`  | String | Yes      | The only accepted value is `application/json`.        |

**Body Payload:**

| Field        | Type   | Required | Description                                                                                          |
| ------------ | ------ | -------- | ---------------------------------------------------------------------------------------------------- |
| `scheduled_sending_at` | Number | No       | The Unix timestamp of when the message should be sent. If not provided, it will be sent immediately. |

**Example Call:**

```json
POST /messages/2b3c4d5e-6f7g-8h9i-0j1k-l2m3n4o5p6q7/retry
Authorization: Basic admF6bGFicy5jb206cGFzc3dvcmQ=
Content-Type: application/json

{
    "scheduled_sending_at": 1735432224
}
```

**Success Response:**

```json
HTTP/1.1 200 OK
Content-Type: application/json

{
  "ok": true,
  "ts": 1735432224
}
```

[Back to Top](#rest-api)
