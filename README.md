<div >
    <img src="assets/owl.jpg" align="left" height="40px" width="40px"/>
    <img src="assets/medusa.png" align="right" height="40px" width="40px"/>
    <h1 align="center" > Go-SimpleHTTPS </h1>
</div>

## About

Simple HTTP/S alternative to python's http.server in go

## Features
- Mode selection between HTTPS and HTTP.
- No need to use openssl to generate the .crt and .key file, when HTTPS mode is selected the server will generate them for you without writing any file.
- Prints a detailed activity log including IP, time and Method

## Why?

I use python's 'http.server module' a lot so I decided to make my own version adding the possibility to use an encrypted connection just as a fun side project.

## Usage

- Clone the repo to compile it.
```bash
git clone https://github.com/Alpharivs/go-simplehttps.git
```
- Compile and name it whatever suits you.
```bash
❯ go build -o [NAME] main.go
```
- Set the options that you want and execute!
```bash
❯ goserver -h
Usage of ./goserver:
  -d string
    	The directory to serve files from. (Default: current dir) (default ".")
  -p string
    	Listening port. (Default 80 or 443 if using HTTPS)
  -s	Use HTTPS.
```
## Example

HTTPS mode:
```bash
# Start server with the '-s' flag
❯ ./goserver -s
[!] Started HTTPS server on port :443

# Send request
❯ curl https://127.0.0.1 -k
<pre>
<a href="go.mod">go.mod</a>
<a href="go.sum">go.sum</a>
<a href="goserver">goserver</a>
<a href="main.go">main.go</a>
</pre>

# Log
127.0.0.1 - - [04/Feb/2024:18:24:08 +0100] "GET / HTTP/2.0" 200 163
```
HTTP mode:
```bash
# Start server
❯ goserver
[!] Started HTTP server on port 80

# Send request
❯ curl http://127.0.0.1
<pre>
<a href="go.mod">go.mod</a>
<a href="go.sum">go.sum</a>
<a href="goserver">goserver</a>
<a href="main.go">main.go</a>
</pre>

# Log
127.0.0.1 - - [04/Feb/2024:18:24:37 +0100] "GET / HTTP/1.1" 200 163
```
<h2 align="center" > LVX-SIT</h2>
<h3 align="center" > MMDCCLXXV -- Ab urbe condita </h3>