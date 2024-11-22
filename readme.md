# Receipt Processor Challenge

This repository contains my solution to the [Fetch Rewards Receipt Processor Challenge](https://github.com/fetch-rewards/receipt-processor-challenge). The challenge involves developing a RESTful API to process and store receipts, enabling retrieval of calculated points based on specific rules.

## Features

- **Receipt Submission**: Accepts receipt data and validates the structure.
- **Points Calculation**: Implements rules to calculate points based on receipt contents.
- **Retrieve Points**: Provides an API endpoint to fetch points associated with a specific receipt.

## Technologies Used

- **Programming Language**: Go
- **Framework**: [Echo](https://echo.labstack.com/) (or any other framework you used)
- **Database**: in-memory storage
- **Testing**: Standard Go testing library (`testing`)

## API Endpoints

1. **Submit a Receipt**  
   `POST /receipts/process`  
   Accepts a receipt in JSON format and returns a unique receipt ID.

2. **Retrieve Points**  
   `GET /receipts/:id/points`  
   Retrieves the points calculated for the receipt with the given ID.

## Setup and Usage

### Prerequisites

- Go 1.23+ installed on your machine
- (Optional) Docker installed if running in a containerized environment

### Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/receipt-processor
   cd receipt-processor
1. Build and run the application
    ```bash
      go run main.go
    ```
1. run tests
    ```bash
      go test ./...
    ```

### Docker Setup (if applicable)
Build the linux binary
```bash
make build
```
Build the Docker image:

```bash
docker build -t receipt-processor .
```
Run the container:

```bash
docker run -p 8080:8080 receipt-processor
```

### Improvements & Future Work
* Enhance validation logic for more complex receipt formats.
* Integrate a persistent database for production-ready storage.
* Add logging and error-handling for better observability.
* Add additional testing