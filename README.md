# Augmented Country Information Service

## Description

This project implements a REST web service in Go that provides:

- General information about a country
- Currency exchange rates for neighbouring countries
- Status information about dependent third-party services

The service integrates data from:

- REST Countries API  
  http://129.241.150.113:8080/v3.1/

- Currency API  
  http://129.241.150.113:9090/currency/

The service follows the specification provided in the assignment description.

---

## Endpoints

Base path:

/countryinfo/v1/

### 1. Status

GET /countryinfo/v1/status/

Returns status information about:
- REST Countries API
- Currency API
- Service version
- Service uptime

---

### 2. Country Info

GET /countryinfo/v1/info/{two_letter_country_code}

Example:

GET /countryinfo/v1/info/no

Returns general country information.

---

### 3. Exchange Rates

GET /countryinfo/v1/exchange/{two_letter_country_code}

Example:

GET /countryinfo/v1/exchange/no

Returns exchange rates from the base country's currency to neighbouring countries.

---

## Running Locally

1. Clone the repository

git clone <repository-url>  
cd cloud-1

2. Initialize dependencies

go mod tidy

3. Run the service

go run .

The service will run on:

http://localhost:8080

---

## Deployment

The service will be deployed on Render.

The deployed URL will be added here after deployment.

---

## Notes

- Only Go standard library packages are used.
- The service queries third-party APIs in real time.
- No third-party libraries are included.

