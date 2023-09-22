# sharp-cert-manager
This project aims to provide a simple tool to monitor certificate validity. It is composed of a golang backend API built using GO http server and a frontend build using [Svelte](https://svelte.dev/).

![Demo frontend image](/docs/demo.jpeg)

Additionally, the app can be configured to run a job at a given schedule. The job will check the configured websites and send a message to a Webhook with a summary of the websites and their certificate validity.

Teams message:
![Demo teams message](/docs/TeamsDemo.jpg)

Slack message:
![Demo slack message](/docs/SlackDemo.jpg)

# Getting started
The easiest way to get started is to run the Docker image published to [Docker Hub](https://hub.docker.com/repository/docker/jlucaspains/sharp-cert-manager/general). Replace the `SITE_1` parameter value with a website to monitor. To add other websites, just add parameters `SITE_n` where `n` is an integer.

```bash
docker run -it -p 8000:8000 \
    --env ENV=DEV \
    --env SITE_1=https://expired.badssl.com/ \
    jlucaspains/sharp-cert-manager
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

## Running in Azure
### Azure Container Instance
Create an ACI resource via Azure CLI. The following parameters may be adjusted
1. `--resource-group`: resource group to be used
2. `--name`: name of the ACI resource
3. `--dns-name-label`: DNS to expose the ACI under
4. `--environment-variables`
   1. `SITE_1..SITE_N`: monitored websites.

```bash
az container create \
    --resource-group rg-sharpcertmanager-001 \
    --name aci-sharpcertmanager-001 \
    --image jlucaspains/sharp-cert-manager \
    --dns-name-label sharp-cert-manager \
    --ports 8000 \
    --environment-variables ENV=DEV SITE_1=https://expired.badssl.com/
```

### Azure Container App
> While more expensive, an ACA is a better option for production environments as it provides a more robust and scalable environment.

First, create an ACA environment using Azure CLI:

```bash
az containerapp env create \
    --name ace-sharpcertmanager-001 \
    --resource-group rg-experiments-soutchcentralus-001
```

Now, create the actual ACA. The following parameters may be adjusted:
1. `-g`: resource group to be used
2. `-n`: name of the app
4. `--env-vars`
   1. `SITE_1..SITE_N`: monitored websites.

```bash
az containerapp create \
    -n aca-sharpcertmanager-001 \
    -g rg-experiments-soutchcentralus-001 \
    --image jlucaspains/sharp-cert-manager \
    --environment ace-sharpcertmanager-001 \
    --ingress external --target-port 8000 \
    --env-vars ENV=DEV SITE_1=https://expired.badssl.com/ \
    --query properties.configuration.ingress.fqdn
```

## Jobs and Webhook Notifications
The app can be configured to run a job at a given schedule. The job will check the configured websites and send a message to a Webhook with a summary of the websites and their certificate validity. Currently, Teams and Slack are supported.

Adjust the `CHECK_CERT_JOB_SCHEDULE` cron to run at the desired schedule.

The `WEBHOOK_URL` is the URL of the Teams/Slack Webhook to send the message to. Generate a webhook URL for Teams following [this guide](https://docs.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/add-incoming-webhook#add-an-incoming-webhook-to-a-teams-channel) and for Slack following [this guide](https://api.slack.com/messaging/webhooks).

```bash
docker run -it -p 8000:8000 `
    --env ENV=DEV `
    --env SITE_1=https://expired.badssl.com/ `
    --env CHECK_CERT_JOB_SCHEDULE=* * * * * `
    --env WEBHOOK_URL=ReplaceWithWebhookUrl `
    jlucaspains/sharp-cert-manager
```

## All environment options
| Environment variable              | Description                                                                     | Default value                                 |
|-----------------------------------|---------------------------------------------------------------------------------|-----------------------------------------------|
| ENV                               | Environment name. Used to configure the app to run in different environments.   |                                               |
| SITE_1..SITE_N                    | Websites to monitor.                                                            |                                               |
| CHECK_CERT_JOB_SCHEDULE           | Cron schedule to run the job that checks the certificates.                      |                                               |
| WEBHOOK_URL                       | Webhook URL to send the message to.                                             |                                               |
| MESSAGE_URL                       | URL to be used as Card action                                                   |                                               |
| MESSAGE_TITLE                     | message  title                                                                  | Sharp Cert Manager Summary                    |
| MESSAGE_BODY                      | Message body body                                                               | The following certificates were checked on %s |
| WEB_HOST_PORT                     | host and port the web server will listen on                                     | :8000                                         |
| TLS_CERT_FILE                     | Certificate used for TLS hosting                                                |                                               |
| TLS_CERT_KEY_FILE                 | Certificate key used for TLS hosting                                            |                                               |
| CERT_WARNING_VALIDITY_DAYS        | Defines how many days from today a cert need to have to prevent a warning       | 30                                            |
| CHECK_CERT_JOB_NOTIFICATION_LEVEL | Defines minimum notification level for jobs. Values are Info, Warning, or Error | Warning                                       |

## Security considerations
This app is intended to run in private environments or at a minimum be behind a secure gateway with proper TLS and authentication to ensure it is not improperly used.

The app will allow unsecured requests to the configured websites. It will perform a get and discard any data returned. All information used is derived from the connection and certificate negotiated between the http client and the web server being monitored.

## Features
Below features are currentl being evaluated and/or built. If you have a suggestion, please create an issue.

- [x] Display list of monitored certificates
- [x] Display certificate details
- [x] Monitor certificate in background
- [x] Teams WebHook integration
- [x] Slack WebHook integration
- [ ] Monitoring