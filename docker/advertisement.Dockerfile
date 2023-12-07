FROM golang:latest

WORKDIR /advertisement

COPY . .

RUN make build

RUN chmod ugo+x .bin/advertisement

CMD .bin/advertisement