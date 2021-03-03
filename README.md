# makephish

<p align="center">
  <img alt="makephish" src="https://github.com/andpalmier/makephish/blob/main/img/makephish.png?raw=true" />
  <p align="center">
    <a href="https://github.com/andpalmier/makephish/blob/master/LICENSE"><img alt="Software License" src="https://img.shields.io/badge/license-GPL3-brightgreen.svg?style=flat-square"></a>
    <a href="https://goreportcard.com/report/github.com/andpalmier/makephish"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/andpalmier/goransom?style=flat-square"></a>
    <a href="https://twitter.com/intent/follow?screen_name=andpalmier"><img src="https://img.shields.io/twitter/follow/andpalmier?style=social&logo=twitter" alt="follow on Twitter"></a>
  </p>
</p>


This is a proof of concept to automatically create phishing kits based on a specified URL, **please note that** `makephish` **will work exclusively on websites having a simple pages with** `<form>` **logins.**

The purpose of this project is purely educative: I wrote `makephish` to get familiar with Go, this also means that the code is poorly written and organized. The idea of the project is also to show how easy it is to clone a website and use it to create phishing pages.

## Usage

After downloading the repository, navigate into the directory and build the project:

```
$ make makephish
```

This will create a folder `build` with an executable called `makephish`. You can run the executable with the following flags:

- `-url`: URL of login page
- `-ua`: User Agent string, by defefault *"Mozilla/5.0 (X11; Linux x86_64; rv:83.0) Gecko/20100101 Firefox/83.0"*
- `-php`: Path to the PHP file to be used, by default *"phish.php"*
- `-kits`: Path used to store the kits, by default *"./kits"*


### Example

```
$ ./build/makephish -url "https://github.com/login/"

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

<p align="center">
  <img alt="fakeGH" src="https://github.com/andpalmier/makephish/blob/main/img/fakeGH.png?raw=true" />
</p>

If you enter some random credentials, you will note that you will be redirected to the real `github.com` login page:

<p align="center">
  <img alt="realGH" src="https://github.com/andpalmier/makephish/blob/main/img/realGH.png?raw=true" />
</p>

If you didn't modify the `phish.php` file, you can find the credentials you just enetered in `localhost:8000/log`:


<p align="center">
  <img alt="logs" src="https://github.com/andpalmier/makephish/blob/main/img/logs.png?raw=true" />
</p>

### PHP capabilities

A simple PHP file is provided in this repo, but you can easily adjust it to your needs. By default, the file will save username, password, User Agent string and IP of the victim in a log file, you can disable this option by removing the content of the variable `log`. You can also specify an email address to send these details via email everytime a new victim enters the credentials.

```
$exfilemail = ""; // -> enter an email address to send the details via email
$logpath = "log"; // -> make this variable empty to disable logging feature
```
