<h1 align="center">Welcome to GoPassHunt 👋</h1>
<p>
</p>

>  Search drives for documents containing passwords

## required
- Install Golang on your computer

- if you want to hunt inside your Gdrive please, create a folder called `credentials` and put your oauth client ID inside it and rename the file `credentials.json`

- We are using go modules therefore you need to install all dependencies before running the program
``` sh
        go mod tidy
```
- check that the go modules are enabled
## Usage

```sh
USAGE
        GO111MODULE=on go run main.go <folderPath> [options]
OPTIONS
        -h, --help
                Display the program usage
        -v, --verbose
                Display additional logs
        -g, --gdrive
                Search in google drive
```

## Authors

👤 **hadi-ilies**

* Github: [@hadi-ilies](https://github.com/hadi-ilies)
* LinkedIn: [@https:\/\/www.linkedin.com\/in\/hadibereksi\]

👤 **Nicolas Barthere**

* Github: [@koff75](https://github.com/koff75)
* LinkedIn: [@https://www.linkedin.com/in/nicolas-barthere/]

## Show your support

Give a ⭐️ if this project helped you!
