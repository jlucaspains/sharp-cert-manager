# sharp-cert-manager
This project aims to provide a simple tool to monitor certificate validity. It is composed of a golang backend API built using GO http server and a frontend build using [Svelte](https://svelte.dev/).

![Demo image](/docs/demo.jpeg)

At the moment, the app doesn't actively monitor the configured websites. Instead, they are only available in the frontend for review.

## Getting started
The easiest way to get started is to run the Docker image published to [Docker Hub](https://hub.docker.com/repository/docker/jlucaspains/sharp-cert-manager/general). Replace the `SITE_1` parameter value with a website to monitor. To add other websites, just add parameters `SITE_n` where `n` is a integer.

> Remember to install Docker before running the docker run command.

```bash
docker run -it -p 8000:8000 --env ENV=DEV --env SITE_1=https://expired.badssl.com/ jlucaspains/sharp-cert-manager
```

## Running locally
### Prerequisites
* Go 1.16+
* NodeJS

### CLone the repo
```bash
git clone https://github.com/jlucaspains/sharp-cert-manager.git
```

### Run the frontend
```bash
cd sharp-cert-manager/frontend
npm install
npm run dev
```

### Run the backend
First, Install the dependencies:

```bash
cd sharp-cert-manager/backend
go mod download
```

Create a dev `.env` file:
```bash
echo "ENV=local\nSITE_1=https://expired.badssl.com/" > .env
```

Finally, run the app:
```bash
go run main.go
```

## Jobs and Teams Webhook
The app can be configured to run a job at a given schedule. The job will check the configured websites and send a message to a Teams Webhook with a summary of the websites and their certificate validity.

Adjust the `CHECK_CERT_JOB_SCHEDULE` cron to run at the desired schedule.

The `TEAMS_WEBHOOK_URL` is the URL of the Teams Webhook to send the message to. Generate a webhook URL following [this guide](https://docs.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/add-incoming-webhook#add-an-incoming-webhook-to-a-teams-channel).

```bash
docker run -it -p 8000:8000 `
    --env ENV=DEV `
    --env SITE_1=https://expired.badssl.com/ `
    --env CHECK_CERT_JOB_SCHEDULE=* * * * * `
    --env TEAMS_WEBHOOK_URL=ReplaceWithTeamsWebhookUrl `
    jlucaspains/sharp-cert-manager
```

## Security considerations
This app is intended to run in private environments or at a minimum be behind a secure gateway with proper TLS and authentication to ensure it is not improperly used.

The app will allow unsecured requests to the configured websites. It will perform a get and discard any data returned. All information used is derived from the connection and certificate negotiated between the http client and the web server being monitored.
