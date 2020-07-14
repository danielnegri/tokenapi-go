# Copyright 2020 The Ledger Authors
#
# Licensed under the AGPL, Version 3.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.gnu.org/licenses/agpl-3.0.en.html
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


FROM golang:1.14-buster as builder
MAINTAINER Daniel Negri <danielgomesnegri@gmail.com>

ENV GO111MODULE=on

RUN set -x \
    && apt-get update \
    && apt-get install -y build-essential ca-certificates git-core zip \
    && rm -rf /var/lib/apt/lists/*

RUN set -x \
   && go get github.com/AlekSi/gocov-xml \
   && go get github.com/axw/gocov/gocov \
   && go get github.com/t-yuki/gocover-cobertura \
   && go get github.com/tebeka/go2xunit

COPY . /go/src/github.com/danielnegri/adheretech
WORKDIR /go/src/github.com/danielnegri/adheretech

RUN set -x \
    && make testall \
    && make release-binary \
    && mkdir -p /usr/share/adheretech \
    && cp -r ./release/bin /usr/share/adheretech/. \
    && cp -r ./results /usr/share/adheretech/. \
    && echo "Build complete."

# Release
FROM debian:buster
MAINTAINER Daniel Negri <danielgomesnegri@gmail.com>

ENV ENVIRONMENT=prod
ENV GIN_MODE=release

RUN set -x \
    && apt-get update \
    && apt-get install -y ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/share/adheretech /usr/share/adheretech
RUN ln -s /usr/share/adheretech/bin/ledger /usr/bin/ledger

WORKDIR /usr/share/adheretech

EXPOSE 8080

ENTRYPOINT ["ledger"]
CMD ["serve"]
