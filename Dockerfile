FROM golang:1.22.0-alpine

WORKDIR /simaku-elearning

COPY go.mod /simaku-elearning 
COPY . /simaku-elearning

RUN go mod tidy 
RUN go build -o /simaku-elearning/bin/main /simaku-elearning/main.go 


LABEL authors="Ridho Galih Pambudi"

LABEL org.opencontainers.image.authors="Ridho Galih Pambudi <rneko2006@gmail.com>"
LABEL org.opencontainers.image.title="simaku-api"
LABEL org.opencontainers.image.description="An API service for Simaku E-Learning."
LABEL org.opencontainers.image.vendor="Simaku"


EXPOSE 8082

ENTRYPOINT ["/simaku-elearning/bin/main"]
CMD ["http-gw-srv", "--port", "8082"]

