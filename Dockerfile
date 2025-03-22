
FROM golang:1.24.1-alpine3.21 AS gobuilder
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

FROM node:18-alpine AS sveltebuilder
WORKDIR /app
COPY frontend/ ./
RUN npm install --ignore-scripts
RUN echo "PUBLIC_API_BASE_PATH=/api" > .env
RUN npm run build

FROM scratch AS runner
COPY --from=gobuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=gobuilder /etc/passwd /etc/passwd
COPY --from=gobuilder /etc/group /etc/group
WORKDIR /app
COPY --from=gobuilder /app/certChecker .
COPY --from=sveltebuilder /app/build/ ./public/
USER appuser:appuser
EXPOSE 8000
ENTRYPOINT ["./certChecker"]
