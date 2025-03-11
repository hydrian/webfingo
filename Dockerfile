FROM golang:1.23

ARG DOCKERCOMPOSE_USERID
ARG DOCKERCOMPOSE_GROUPID

ENV DOCKERCOMPOSE_USERID=${DOCKERCOMPOSE_USERID}
ENV DOCKERCOMPOSE_GROUPID=${DOCKERCOMPOSE_GROUPID}

RUN addgroup --gid ${DOCKERCOMPOSE_GROUPID} --system thedockeruser
RUN adduser --system --shell /bin/bash -u ${DOCKERCOMPOSE_USERID} thedockeruser
RUN usermod -aG thedockeruser thedockeruser
RUN usermod -aG sudo thedockeruser

RUN mkdir /app
RUN chown -R thedockeruser:thedockeruser /app

RUN mkdir /nonexistent
RUN chown -R thedockeruser:thedockeruser /nonexistent

USER thedockeruser

WORKDIR /app/

CMD /bin/bash
