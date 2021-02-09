FROM ubuntu:18.04
RUN apt-get update \
 && apt-get install -y make \
 && apt-get install -y wget \
 && apt-get install -y build-essential \
 && apt-get install -y clang-format

# Install golang
RUN wget -q https://golang.org/dl/go1.15.7.linux-amd64.tar.gz -O /tmp/go1.15.7.linux-amd64.tar.gz \
 && tar -C /usr/local -xzf /tmp/go1.15.7.linux-amd64.tar.gz

# Add go lang into path.
ENV PATH="/usr/local/go/bin:${PATH}"
# Set gopath
ENV GOPATH="/go"
WORKDIR /go/src/govmtouch