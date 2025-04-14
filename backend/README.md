# Running the Project with Docker

This section provides instructions to build and run the project using Docker.

## Prerequisites

- Ensure Docker and Docker Compose are installed on your system.
- Docker version: Compatible with syntax `docker/dockerfile:1`.
- Golang version: 1.24.2 (used in the build stage).

## Environment Variables

- The `database` service requires the following environment variables:
  - `POSTGRES_USER`: Database username (default: `user`).
  - `POSTGRES_PASSWORD`: Database password (default: `password`).
  - `POSTGRES_DB`: Database name (default: `appdb`).

## Build and Run Instructions

1. Clone the repository and navigate to the project directory.
2. Build and start the services using Docker Compose:
   ```bash
   docker-compose up --build
   ```
3. Access the application at `http://localhost:8080`.

## Exposed Ports

- `app` service: 8080 (mapped to host).
- `database` service: Not exposed to host.

## Notes

- Uncomment the `env_file` line in the `docker-compose.yml` file if an `.env` file is used for environment variables.
- The `db_data` volume persists database data across container restarts.

For further details, refer to the project documentation or contact the maintainers.