FROM golang AS build

WORKDIR /src/

COPY . .

RUN CGO_ENABLED=0 go build -o mc

FROM alpine

WORKDIR /app/

COPY --from=build /src/mc .

ENTRYPOINT [ "./mc" ]
