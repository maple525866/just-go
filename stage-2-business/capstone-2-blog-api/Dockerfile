FROM golang:1.24 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /out/capstone-2 ./stage-2-business/capstone-2-blog-api

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/capstone-2 /app/capstone-2
EXPOSE 8080
ENTRYPOINT ["/app/capstone-2"]
