# Gourmet
## script is served
Gourmet is a work in progress project written in go to run your code "as a service".  
The idea is easy, ping an endpoint and something runs your code.

At the moment it runs only PHP scripts, it prepare a docker container with PHP7 and 
it runs tree simple steps

```bash
wget <your-source-zip>
unzip  <your-source-zip>.zip -d .
bin/console
```
`bin/console` is the console entrypoint of your application, it should be executable from php.

### Implementation
gourmet is a cli application that start an http server to start new build.

```bash
./gourmet api
```
After this command your server is ready to go on port 8000.

### HTTP Api

Start build with a `POST` request on `/project`.
```
{
    "source": "http://site.net/static/my-script.zip
}
```

### Troubleshooting
* gourmet at the moment uses `gourmet/php` docker image. In this repository you can try my proposal.

### TODO
This project has a log todo list you can follow it and help me if you like this idea.  
Here only a little list of possibile improvement:

* Download source with authentication
* Support for more languages
* Increase application config
* Build container with environment variables (you will use them on your code)
* other things on [gourmet/issues](https://github.com/gianarb/gourmet/issues)

