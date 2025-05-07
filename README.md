<h1 align="center" style="border-bottom: none">
  <div>
    WA Scheduler
  </div>
  Whatsapp Message Scheduler<br>
</h1>

<p align="center">
A simple message scheduling tool for WhatsApp private chats or groups. Built to make sure your messages are seen at the right time.
</p>

## Why We Built This

In our group, important messages often vanished into the noise — sent at odd hours, buried under a flood of chats, and seen too late (or not at all). We built this tool to change that. With manual scheduling, you control exactly when your message hits.

## Architecture

![High Level Architecture](./docs/architecture.drawio.svg)

Available services:

- `Server and Dashboard Service` => Handling dashboard and api requests from clients. For the API details see [this doc](./docs/rest_api.md).
- `WhatsApp Publisher` => Whatsapp publisher service. This service is responsible for sending messages to WhatsApp. Right now it only support using [go-whatsapp-web-multidevice](https://github.com/aldinokemal/go-whatsapp-web-multidevice).
- `Storage` => This service is responsible for storing all the message state. The database schema is available [here](./docs/db/schema.sql). Currently, it only supports MySQL.

## Features

- Schedule messages for private chats or groups
- Set exact send times
- Retry send

## Getting Started

### Locally (Docker)

1. Run the following commands:

    ```bash
    git clone https://github.com/ghazlabs/wa-scheduler.git

    make run
    ```

2. Open <http://localhost:9865> to access the dashboard
3. Log in with username `admin` and password `admin`
4. Scan the QR code to connect your WhatsApp account
5. When the dashboard shows, you already can schedule message
6. To get group recipients id, you need to check `List Groups` from the WhatsApp Publisher service on <http://localhost:3000>.

### Production

TBD

## Environment Variables

| Variable Name               | Required | Default | Description                                                                                                                                      |
| --------------------------- | -------- | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------ |
| `LISTEN_PORT`               | Yes      | `9865`  | Port number the server listens on                                                                                                                |
| `MYSQL_DSN`                 | Yes      | –       | MySQL Data Source Name                                                                                                                           |
| `DASHBOARD_CLIENT_USERNAME` | Yes      | –       | Username for dashboard authentication                                                                                                            |
| `DASHBOARD_CLIENT_PASSWORD` | Yes      | –       | Password for dashboard authentication                                                                                                            |
| `WA_DEFAULT_NUMBERS`        | No       | –       | Comma-separated list of default numbers could be private numbers or group id WhatsApp. E.g. `6287822334455@s.whatsapp.net,120363020892687898@g.us` |
| `WA_PUBLISHER_API_BASE_URL` | Yes      | –       | Base URL for WA Publisher API                                                                                                                    |
| `WA_PUBLISHER_USERNAME`     | Yes      | –       | Username for WA Publisher API                                                                                                                    |
| `WA_PUBLISHER_PASSWORD`     | Yes      | –       | Password for WA Publisher API                                                                                                                    |
| `WEB_CLIENT_PUBLIC_DIR`     | Yes      | `web`   | Directory for serving the web client                                                                                                             |

## Contributing

First and foremost, thank you for your interest in contributing to WA Scheduler 🙏

There are many ways to contribute, and most of them dont require writing code.

- [Spread the word](#spread-the-word)
- [Engage with the community](#engage-with-the-community)
- [Contribute code](#contribute-code)

### Spread the word

This might be the biggest help of all. Share WA Scheduler with your network or anyone who needs a simple way to schedule WhatsApp messages.

### Engage with the community

Every message, reaction, or bit of feedback counts. It keeps us motivated and reminds us that real people find this project useful.

### Contribute code

Code is just one piece of the puzzle—and contributing doesn’t always mean writing code. But if you do want to dive in, start small! Fix typos, report or squash bugs from the [issues page](https://github.com/ghazlabs/wa-scheduler/issues), polish up the docs, or add helpful features.

> [!TIP]
>
> Code matters, but it’s just one part of what makes a great product. Sometimes the easiest code fix isn’t the best choice overall. Don’t forget—there are plenty of other ways to contribute too!

#### Quick steps to contribute

1. Fork the repo via the ["Fork"](https://github.com/ghazlabs/wa-scheduler/fork) button
2. Clone your fork locally
3. Create a branch

    ```bash
    git checkout -b your-feature-name
    ```

4. Make your changes
5. Open a pull request
