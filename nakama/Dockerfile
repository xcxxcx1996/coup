FROM heroiclabs/nakama-pluginbuilder:3.10.0 AS builder

# ENV GO111MODULE on
ENV GO111MODULE=on 

ENV CGO_ENABLED 1
# ENV GOPRIVATE "github.com/heroiclabs/nakama-project-template"

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOPROXY="https://goproxy.cn,direct"

WORKDIR /backend
COPY . .
RUN go mod download
RUN go mod vendor
RUN go build --trimpath --mod=vendor --buildmode=plugin -o ./backend.so

FROM heroiclabs/nakama:3.10.0

COPY --from=builder /backend/backend.so /nakama/data/modules
COPY --from=builder /backend/local.yml /nakama/data/
