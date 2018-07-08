FROM golang:latest

# Create the directory where the application will reside
RUN mkdir /app
COPY ./build/emit /app/emit
WORKDIR /app
EXPOSE 8080

CMD /app/emit