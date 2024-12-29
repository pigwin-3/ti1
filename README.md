# TI1

The best thing to happen since yesterday at 3 pm

## Usage

To use this project, you can pull the Docker image from Docker Hub and run it using the following commands:

### Pull the Docker Image

```sh
docker pull pigwin1/ti1:latest
```

### Run the Docker Container
```sh
docker run -d --name ti1-container -e DB_HOST=<your_db_host> -e DB_PORT=<your_db_port> -e DB_USER=<your_db_user> -e DB_PASSWORD=<your_db_password> -e DB_NAME=<your_db_name> -e DB_SSLMODE=<your_db_sslmode> pigwin1/ti1:latest
```
Replace `<your_db_host>`, `<your_db_port>`, `<your_db_user>`, `<your_db_password>`, `<your_db_name>`, and `<your_db_sslmode>` with your actual database configuration values.

### Docker Hub Repository
You can find the Docker image on Docker Hub at the following link:

https://hub.docker.com/repository/docker/pigwin1/ti1/general

