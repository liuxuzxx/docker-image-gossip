FROM golang:1.18

COPY ./cmd/gossiper/gossiper /home/gossiper/

WORKDIR /home/gossiper/

ENTRYPOINT ["/home/gossiper/gossiper"]