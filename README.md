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

## Features
- Schedule messages for private chats or groups
- Set exact send times
- Retry send

## Installation

### Prerequisites
- Docker
- Docker Compose
- WhatsApp account

### Quick Start

1. Run this

```bash
git clone https://github.com/ghazlabs/wa-scheduler.git

make run
  ```

2. Open http://localhost:9865 to access the dashboard
3. Login with username `admin` and password `admin`
4. Scan the QR code to connect your WhatsApp account
5. When the dashboard shows, you already can schedule message

## Contributing
<a name="contributing"></a>

Thank you for considering contributing to WA Scheduler!

Whether you’re fixing bugs, improving documentation, or building new features, **your help makes this project better**.

You can also solve our [issues](https://github.com/ghazlabs/wa-scheduler/issues) by code or simply join the discussion to start contribute.

### Getting Started
1. The first step is to fork and clone this repository. You can perform these steps manually by clicking the ["Fork"](https://github.com/ghazlabs/wa-scheduler/fork) button on the GitHub repository page, or by using [Github CLI](https://cli.github.com/)
2. Clone your fork locally
3. Create a branch
```bash
git checkout -b your-feature-name
```
4. Make your changes
5. Commit and push
```bash
git commit -m "Add: your message"
git push origin your-feature-name
```
6. Open a pull request
