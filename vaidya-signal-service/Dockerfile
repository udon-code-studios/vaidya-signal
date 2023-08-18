FROM golang:1.21

# set port to serve app
ENV PORT 3000
EXPOSE ${PORT}

# get alpaca variables as build arguments and set env variables
ARG APCA_API_KEY_ID
ARG APCA_API_SECRET_KEY
ARG APCA_API_BASE_URL
ENV APCA_API_KEY_ID=$APCA_API_KEY_ID
ENV APCA_API_SECRET_KEY=$APCA_API_SECRET_KEY
ENV APCA_API_BASE_URL=$APCA_API_BASE_URL

WORKDIR /app

# copy app
COPY ./src .

# install go packages
RUN go install

# build app
RUN go build -o app

# start app
CMD ./app