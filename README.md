# TI1

The best thing to happen since yesterday at 3 pm

## Usage

Start with getting Docker then do the following:

### Create the setup files
Create a `docker-compose.yml`
```yaml
services:
  db:
    image: postgres:17.2
    container_name: postgres-db
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: Root Password
      POSTGRES_DB: ti1
    ports:
      - "5432:5432"
    volumes:
      - ./postgres_data:/var/lib/postgresql/data  # Store data in the current directory
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql:ro
    networks:
      - app-network
    healthcheck:
      test: ["CMD", "pg_isready", "-U", "postgres", "-d", "ti1", "-h", "db"]
      interval: 10s
      retries: 5
    restart: always  # Ensure the container always restarts

  ti1-container:
    image: pigwin1/ti1:v0.1.1
    container_name: ti1-container
    environment:
      DB_HOST: db
      DB_PORT: 5432
      DB_USER: ti1
      DB_PASSWORD: ti1 password
      DB_NAME: ti1
      DB_SSLMODE: disable
    depends_on:
      db:
        condition: service_healthy  # Wait until the db service is healthy
    networks:
      - app-network
    restart: always  # Ensure the container always restarts

networks:
  app-network:
    driver: bridge

volumes:
  postgres_data:
    driver: local
```

Create `init.sql`
```sql
-- Check if 'post' user exists; create if not
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'post') THEN
        CREATE ROLE post WITH LOGIN PASSWORD 'post password';
        GRANT ALL PRIVILEGES ON DATABASE ti1 TO post;
        ALTER ROLE post WITH SUPERUSER;
    END IF;
END
$$;

-- Check if 'ti1' user exists; create if not
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'ti1') THEN
        CREATE ROLE ti1 WITH LOGIN PASSWORD 'ti1 password';
        GRANT ALL PRIVILEGES ON DATABASE ti1 TO ti1;
        GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ti1;
        GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ti1;
        GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO ti1;
-- Grant the ti1 user the necessary permissions on the public schema
GRANT USAGE, CREATE ON SCHEMA public TO ti1;

-- Grant all permissions (SELECT, INSERT, UPDATE, DELETE, etc.) on all existing tables in the public schema
GRANT ALL ON ALL TABLES IN SCHEMA public TO ti1;

-- Grant all permissions on all existing sequences in the public schema
GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO ti1;

-- Grant all permissions on all functions in the public schema
GRANT ALL ON ALL FUNCTIONS IN SCHEMA public TO ti1;

-- Ensure that the ti1 user will have access to new tables, sequences, and functions created in the public schema
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO ti1;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO ti1;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO ti1;

-- Optionally, grant full permissions on the entire database to ti1 (if needed)
-- GRANT ALL PRIVILEGES ON DATABASE ti1 TO ti1;

    END IF;
END
$$;
```

Remember to change the password values

### Run the Docker Containers
```sh
docker compose up -d
```

### edit the postgress config (optinal)
open the config file
```sh
nano postgres_data/postgresql.conf
```
Change the following values
```conf
listen_addresses = '*'
max_connections = 100
shared_buffers = 16GB
work_mem = 256MB
maintenance_work_mem = 2GB
dynamic_shared_memory_type = posix
max_wal_size = 1GB
min_wal_size = 80MB
```
set these to what makes most sense for you

These values should also be set bet not necessarily changed
```conf
log_timezone = 'Etc/UTC'
datestyle = 'iso, mdy'
timezone = 'Etc/UTC'
lc_messages = 'en_US.utf8'
lc_monetary = 'en_US.utf8'
lc_numeric = 'en_US.utf8'
lc_time = 'en_US.utf8'
default_text_search_config = 'pg_catalog.english'
```

### Docker Hub Repository
You can find the Docker image on Docker Hub at the following link:

https://hub.docker.com/repository/docker/pigwin1/ti1/general

