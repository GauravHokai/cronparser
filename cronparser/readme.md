# Cron Expression Parser in Go

A command-line application, written in Go, that parses a standard cron string and expands each field to show the times at which it will run.

## Features

* Parses standard cron expressions with five time fields (minute, hour, day of month, month, day of week) and a command.
* Supports all standard cron field operators:
    * `*` (wildcard)
    * `,` (list separator)
    * `-` (range)
    * `/` (step/interval)
* Handles complex combinations of operators.
* Provides clear, table-formatted output for each expanded field.
* Includes a full suite of automated tests.

## Requirements

* [Go](https://go.dev/doc/install) version 1.18 or newer.

## Setup

1.  Clone the repository to your local machine:
    ```sh
    git clone <your-repository-url>
    ```
2.  Navigate into the project directory:
    ```sh
    cd cronapp
    ```
3.  Run the application:

    * **Method 1 (With Go Installed):**
        ```sh
        go run main.go "*/15 0 1,15 * 1-5 /usr/bin/find"
        ```
    * **Method 2 (Using the pre-compiled MacOS binary):**
        ```sh
        sudo ./cronParserBuild "*/15 0 1,15 * 1-5 /usr/bin/find"
        ```
4.  Run the tests (Requires Go):
    ```sh
    go test -v ./cronapp
    ```

## Current Limitations

* No support for special time strings (e.g., "@yearly")
* No support for special characters ('?' or 'L')
* No validation for varying days in months (e.g., February)