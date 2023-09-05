# sharp-cert-checker
This project aims to provide a simple tool to monitor certificate validity. It is composed of a golang backend API built using [Fiber](https://gofiber.io/) and a frontend build using [Svelte](https://svelte.dev/).

![Demo image](/docs/demo.jpeg)

## Getting started
Run docker image

## Running locally
### Prerequisites
* Go 1.16+
* NodeJS

### CLone the repo
```bash
git clone https://github.com/jlucaspains/sharp-cert-checker.git
```

### Run the frontend
```bash
cd sharp-cert-checker/frontend
npm install
npm run dev
```

### Run the backend
First, Install the dependencies:

```bash
cd sharp-cert-checker/backend
go mod download
```

Create a dev `.env` file:
```bash
echo "ENV=local\nSITE_1=https://blog.lpains.net" > .env
```

Finally, run the app:
```bash
go run main.go
```

## Deploying
Build and deploy 

### Security considerations
This app is intended to run in private environments or at a minimum be behind a secure gateway with proper TLS and authentication to ensure it is not improperly used.

The app will allow unsecured requests to the configured websites. It will perform a get and discard any data returned. All information used is derived from the connection and certificate negotiated between the http client and the web server being monitored.
