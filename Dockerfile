FROM docker.arvancloud.ir/golang:1.25 AS builder
WORKDIR /app
# downloading dependencies ignored until vendoring is an option
# COPY go.* ./
# RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/shortener .

FROM scratch
COPY --from=builder /app/shortener /app/shortener
ENTRYPOINT [ "/app/shortener" ]
CMD [ "server" ]
