
FROM golang:1.20.3-alpine3.17 AS goBuilder
WORKDIR /app

RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

ENV USER=appuser
ENV UID=10001 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"
    
COPY backend/go.mod go.mod
COPY backend/go.sum go.sum
RUN go mod download
COPY backend/. .
RUN go build -ldflags "-s -w" -o ./certChecker ./main.go

FROM node:16-alpine AS svelteBuiler
WORKDIR /app
COPY frontend/ ./
RUN npm install --ignore-scripts
RUN echo "PUBLIC_API_BASE_PATH=http://localhost:3000/api" > .env
RUN npm run build

FROM scratch AS runner
COPY --from=goBuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=goBuilder /etc/passwd /etc/passwd
COPY --from=goBuilder /etc/group /etc/group
WORKDIR /app
COPY --from=goBuilder /app/certChecker .
COPY --from=svelteBuiler /app/build/ ./public/
USER appuser:appuser
EXPOSE 3000
ENTRYPOINT ["./certChecker"]
