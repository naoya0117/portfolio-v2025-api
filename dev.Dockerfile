FROM golang:1.24.5-bookworm

ARG UID=1000
ARG GID=1000


# 非特権ユーザの設定
RUN (getent passwd ${UID} && /usr/sbin/userdel -r $(getent passwd ${UID} | cut -d: -f1) || true) && \
  (getent group ${GID} || groupadd -g ${GID} nonroot) && \
  /usr/sbin/useradd -u ${UID} -g ${GID} -m -s /bin/bash nonroot

RUN go install github.com/air-verse/air@latest

# Go module cache用のディレクトリを作成し、nonrootユーザに権限を付与
RUN mkdir -p /go/pkg/mod/cache && \
  chown -R ${UID}:${GID} /go/pkg/mod

USER nonroot
WORKDIR /api

EXPOSE 8080

CMD ["air"]
