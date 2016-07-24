# Lambda function manager
Gourmet runs your code "as a service". Lambda function manager written in
golang and based on docker.

## Implementation
gourmet is a cli application that start an http server to start new build.

```bash
./gourmet api
```
After this command your server is ready to go on port 8000.

### Env vars
`GOURMET_REGISTRY_URL` allow registry push and pull in order to
manage more gourmet.

## API

```json
POST /func

{
    "img": "gourmet/php",
    "source": "https://ramdom-your-source.net/gourmet.zip",
}
```
* `img` is the name of docker image to use how started point
* `source` is the artifact of your script, it should be contain an executable console entrypoint `bin/console`
This function return the function's id `FuncId`
```
{
    "FuncId": "34gaw23t2"
}
```


```json
POST /func/{FuncId}

{
    "env": [
        "AWS_KEY=EXAMPLE",
        "AWS_SECRET=",
        "AWS_QUEUE=https://sqs.eu-west-1.amazonaws.com/test"
    ]
}
```
* `evn` are environment varaibles, you can use them to configure your script
This function return status of our function
```
At the moment caos but I am working on it
```

## How it works
gourmet prepares your container, downloads source and executes this steps:
```bash
wget <your-source-zip>
unzip  <your-source-zip>.zip -d .
bin/console
```

`bin/console` is the console entrypoint of your scirpt, it should be executable.

### Troubleshooting
* In this repository you can try an example of docker image (PHP7) build it and go!
* During my test I'm using this [php-script](https://github.com/gianarb/gourmet-php-example), it is very easy require 3 env var
    * AWS_KEY, AWS_SECRET, AWS_QUEUE and push a message in queue

### TODO
This project has a log todo list you can follow it and help me if you like this idea.
Here only a little list of possibile improvement:

* Download source with authentication
* Increase application config
* other things on [gourmet/issues](https://github.com/gianarb/gourmet/issues)

