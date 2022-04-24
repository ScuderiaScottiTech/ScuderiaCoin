FROM golang
COPY . /go/src/github.com/ScuderiaScottiTech

WORKDIR /go/src/github.com/ScuderiaScottiTech/ScuderiaCoinMineAPI
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main .

# Second stage
FROM alpine

WORKDIR /app
COPY --from=0 /go/src/github.com/ScuderiaScottiTech/ScuderiaCoinMineAPI/main /app/main
CMD ["/app/main"]  