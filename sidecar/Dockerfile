FROM ubuntu:latest

RUN addgroup --system --gid 50000 data \
  && adduser --system --uid 50000 --gid 50000 data 

RUN apt-get update -y \
  && apt-get install -y curl unzip

RUN mkdir /data \
  && chown data:data /data \
  && chmod 755 /data

RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" \
  && unzip awscliv2.zip \
  && ./aws/install

COPY entrypoint.sh /entrypoint.sh

RUN chown data:data /entrypoint.sh \
  && chmod 700 /entrypoint.sh

WORKDIR "/data"

USER data

VOLUME [ "/data" ]

ENTRYPOINT ["/entrypoint.sh"]
