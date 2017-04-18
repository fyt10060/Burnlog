FROM golang
RUN mkdir /app
ADD burnlog /app/burnlog
ADD conf /app/conf
WORKDIR /app
EXPOSE 10086
EXPOSE 6379
EXPOSE 16379
ENTRYPOINT /app/burnlog

CMD ["bee", "run"]