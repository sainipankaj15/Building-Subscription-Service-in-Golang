# About Project
- This is Subscription Service. This service is built with Golang, leveraging <b> Goroutines and Channels for concurrency</b>, and integrating modern technologies like <b> Redis, Docker, and PostgreSQL </b> for efficient performance and scalability.

---

# Prerequisites to run the Project

**Installing Golang**

- Before getting started with our Subscription Service, ensure you have Go installed on your local machine. If you haven't installed it yet, you can download and install it from the official Go website [here](https://go.dev/doc/install).

**Installing BeepKeeper**

- To use BeepKeeper with our Subscription Service, you'll need to install it on your local machine. You can download and install BeepKeeper from the official website [here](https://www.beekeeperstudio.io/get).

**Installing Docker**

- Before proceeding with our Subscription Service, make sure you have Docker installed on your local machine. If you haven't installed Docker yet, you can download and install it from the official website [here](https://www.docker.com/products/docker-desktop/).

---
# How to run the Project

Before running the project, ensure you have the following prerequisites:

- **Running Docker**: Make sure Docker is installed and running on your system. Our Docker image includes Redis, PostgreSQL, and Mailhog.

## Getting Started

1. **Navigate to the Main Directory**: Go to the main directory of the project where `docker-compose.yml` and `Dockerfile` are located.

2. **Start Docker Containers**:
   
    ```bash
    docker-compose up -d
    ```

    This command will start the Docker containers in detached mode.

3. **Stop Docker Containers** (Optional):

    If you need to stop the Docker containers, you can use the following command:

    ```bash
    docker-compose down
    ```

    This will stop and remove the Docker containers.


## Populating the Database

To populate the database with initial data, follow these steps:

1. **Copy SQL Commands**: Open the `db.sql` file and copy the SQL commands.

2. **Run Commands in BeepKeeper Query Section**: Launch BeepKeeper and navigate to the query section.

3. **Paste and Execute Commands**: Paste the copied SQL commands into the query section and execute them. This will create the necessary tables (`plans`, `user_plans`, and `users`) in the database.


### Starting the Application

If you're on a Linux system, you can start the application using the provided Makefile. However, if you're on Windows or macOS, you'll need to convert the Makefile commands accordingly using an online tool.

To start the application, follow these steps:

1. **Convert Makefile Commands (Windows/macOS)**:
   
   If you're on Windows or macOS, use an online tool to convert the Makefile commands to your system's equivalent commands.

2. **Run the Command**:
   
   Once you've converted the Makefile commands, open your terminal or command prompt and execute the following command:

   ```
   make start
    ```

### Managing the Main Application

You can manage the application using the following Makefile commands:

- **To start the application**: Execute the command `make start`.

- **To stop the application**: Execute the command `make stop`.

- **To restart the application**: Execute the command `make restart`.

## Accessing the Application

To access the application, please navigate to the following URL in your web browser:

[http://localhost:8080/](http://localhost:8080/)

This address points to the local server where the application is currently running. Simply paste the URL into your browser's address bar and press Enter to access the application's interface.

### Testing of this Application 
- Check out the testing branch of this repository for comprehensive test cases and results for each function.