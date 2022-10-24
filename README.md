# Fabrica Coding Challenge
by Tim Garthwaite
2022-10-24

## Getting Started

#### Tools
To use this tutorial, you must have an installation of Go installed. The example usage below also assumes that you
have `curl` installed on your system and available in your path.

There is a `.tool-versions` file included, which is pre-set with the `golang` version that the author used when
developing this code. The author uses the [asdf](https://asdf-vm.com) version manager to manage different tool
versions for each repository on his system. `asdf` shims the shell commands for the proper tool versions for each
repository as he changes the present working directory.

You can install Go for your OS directly from [the website](https://golang.google.cn/dl/). Or, with the `asdf` version
manager installed, change to the root directory in your shell (or any subdirectory) and run the following commands:
```shell
asdf plugin add golang
asdf install
```

### Running the REST API server locally
From the root directory of the repository, run:
```shell
asdf plugin add golang
asdf install
go run .
```

#### Customizing the bound hostname and port
The API server, by default, will start up bound to `localhost` on the port `8080`. You may change either or both of
these by setting the `SERVER_HOST` or `SERVER_PORT` environment variables, respectively. The author uses the `direnv`
utility to keep a different set of environment variables for each local code repository; as such, the `.envrc` file
is included in the `.gitignore` file for the repository.

### API Usage

#### Creating an order
Here is an example of creating a new order from the command line using cURL:
```shell
curl http://localhost:8080/orders \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"userId": "A", "items": [1, 3, 5]}'
```

Here is the expected output:
```
HTTP/1.1 201 Created
Content-Type: application/json; charset=utf-8
Date: Mon, 24 Oct 2022 13:07:48 GMT
Content-Length: 125

{
    "id": "4eb6c67a-6403-431e-8e50-8b58eafa640b",
    "userId": "A",
    "items": [
        1,
        3,
        5
    ],
    "status": "queued"
```

#### Driver picking up an order from the queue
```shell
curl http://localhost:8080/orders \
    --include \
    --header "Content-Type: application/json" \
    --request "PUT" \
    --data '{"orderId": "4eb6c67a-6403-431e-8e50-8b58eafa640b", "status": "en-route"}'
```

Expected output if the order with the ID is found:
```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Mon, 24 Oct 2022 14:00:44 GMT
Content-Length: 151

{
    "id": "d724215a-fc96-4b56-992f-2603a8eba81c",
    "userId": "A",
    "items": [
        1,
        3,
        5
    ],
    "status": "en-route"
}
```

Expected output if the order is not found:
```
HTTP/1.1 404 Not Found
Content-Type: application/json; charset=utf-8
Date: Mon, 24 Oct 2022 13:57:55 GMT
Content-Length: 93

{"code":"NOT_FOUND","message":"Order with id 69ecf8c0-bf80-469a-ab56-9923a682d304 not found"}
```

Expected output if the status is not allowed:
```
HTTP/1.1 422 Unprocessable Entity
Content-Type: application/json; charset=utf-8
Date: Mon, 24 Oct 2022 14:09:12 GMT
Content-Length: 97

{"code":"UNPROCESSABLE_ENTITY","message":"Change of status from queued to closed is not allowed"}
```