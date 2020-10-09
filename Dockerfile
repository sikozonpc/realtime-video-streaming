FROM golang
  ENV GO111MODULE on
  # if left blank app will run with dev settings
  # to build production image run:
  # $ docker build ./api --build-args app_env=production
  ARG app_env
  ENV APP_ENV $app_env

  COPY go.mod go.sum /src/
  WORKDIR /src
  COPY . /src

  RUN go mod download
  RUN go build

  # if dev setting will use pilu/fresh for code reloading via docker-compose volume sharing with local machine
  # if production setting will build binary
  CMD if [ ${APP_ENV} = production ]; \
    then \
    api; \
    else \
    go get github.com/pilu/fresh && \
    fresh; \
    fi

  EXPOSE 8080