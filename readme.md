# Shorty Challenge
## Introduction
The challenge is to create a micro service to shorten urls,
in the style that TinyURL and bit.ly made popular.

## Getting Started
1.	Installation process
2.	Usage

### 1. Installation process
you can run the program with the following way:

a. by running the program by this command (golang way):
1. firstly, you can install go programming language.
2. go to the directory of the project, by this command
    ```shell
    cd test_amartha_muhammad_huzair
    ```
3. and install the dependencies by this command
    ```shell
    go mod download
    go mod tidy
    go mod verify
    ```
4. run the program by this command with the argument on which port are you using
    ```shell
    go run main.go -flag=3000
    ```
   or you can leave the option of the argument just to run the service by this exact command
   ```shell
   go run main.go
   ```
   the service will run on the default port on `8080`

b. by running the program by this command (bash way):
1. go to the directory of the project, by this command
     ```shell
     cd test_amartha_muhammad_huzair
     ```
2. run the program by this command with the argument on which port are you using
     ```shell
     ./bin/tiny_url -port=3000
     ```
   or you can leave the option of the argument just to run the service by this exact command
   ```shell
   ./bin/tiny_url
   ```
   the service will run on the default port on `8080`

## 2. Usage
all the url or feature provided by this microservice

### a. shorten the link
to use this endpoint use client http tools, such as postman, or insomnia and type this to the url format
`<address>:<PORT>/shorten` or example `http://localhost:8080/shorten` with the `POST` method
and the body request using the content type `application/json` is 
```json
{
    "url": "https://github.com",
    "shortcode": "prDeb7"
}
```
or you can just fill the url, so you dint have to think about the shortcode
like this example
```json
{
    "url": "https://github.com"
}
```

or using the cUrl on the terminal way
```shell
curl --location --request POST 'localhost:8080/shorten' \
--header 'Content-Type: application/json' \
--data-raw '{
    "url": "https://github.com"
}'
```

you should get the response

```json
{
    "shortcode": SHORTCODE_GENERATED_OR_SET_BY_YOU
}
```

like this example
```json
{
    "shortcode": "FJWEpm"
}
```

### b. get Redirect to the url you have been set
to use this endpoint use client http tools, such as postman, or insomnia and type this to the url format
`<address>:<PORT>/<SHORTCODE_GENERATED_OR_SET_BY_YOU>` 
or example `http://localhost:8080/<SHORTCODE_GENERATED_OR_SET_BY_YOU>` with the `GET` method

or using the cUrl on the terminal way 
```shell
curl --location --request GET 'localhost:8080/FJWEpm'
```

or you just can type the url on your favorite browser, `localhost:8080/FJWEpm` 
you should redirect to the url you have been set.

### c. statistic of your url
to use this endpoint use client http tools, such as postman, or insomnia and type this to the url format
`<address>:<PORT>/<SHORTCODE_GENERATED_OR_SET_BY_YOU>/stats`
or example `http://localhost:8080/FJWEpm/stats` with the `GET` method

or using the cUrl on the terminal way
```shell
curl --location --request GET 'localhost:8080/FJWEpm/stats'
```

you should get the response

```json
{
   "lastSeenDate": "2022-07-24T20:44:31+07:00",
   "redirectCount": 1,
   "startDate": "2022-07-24T20:36:31+07:00"
}
```