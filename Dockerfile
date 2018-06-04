FROM lwolf/golang-glide

ENV APP_PATH=/go/src/github.com/b0ralgin/geo-ip-finder

RUN mkdir -p $APP_PATH
WORKDIR $APP_PATH

COPY glide.yaml glide.yaml
COPY glide.lock glide.lock
RUN glide install -v

ADD . $APP_PATH
RUN CGO_ENABLED=0 go build -o service
CMD ./service