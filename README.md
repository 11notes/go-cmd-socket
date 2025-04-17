![banner](https://github.com/11notes/defaults/blob/main/static/img/banner.png?raw=true)

# SYNOPSIS 📖
**What can I do with this?** cmd-socket is a very slim and ultra simple socket http server that will execute commands on the host systems it runs on. Its purpose is to give other containers the ability to execute commands on other systems by simply mounting the socket. You can for instance run a scheduled command in another container without having to open access to the Docker socket. The [compose.yaml](https://github.com/11notes/docker-postgres/blob/master/compose.yaml) of the [11notes/postgres](https://github.com/11notes/docker-postgres) image has a great example for this use.

# COMMAND LINE 📟
* **-s** - path to socket file, by default */run/cmd/cmd.sock*

# JSON DATA STRUCTURE 📦
The http server will parse any json object in this format.
```json
{
  "bin":"df",
  "arguments":["-h"],
  "environment":{}
}
```
* **bin** - The name of the binary to be executed by the socket, **required**
* **arguments** - An array ([]string) with the parameters passed to the binary, *optional*
* **environment** - An object ([]string) with the environment passed to the binary, *optional*

# EXAMPLES
```shell
curl --unix-socket /run/cmd/cmd.sock http:/cmd -H 'Content-Type: application/json' -d '{"bin":"df", "arguments":["-h"]}'
curl --unix-socket /run/cmd/cmd.sock http:/cmd -H 'Content-Type: application/json' -d '{"bin":"env", "environment":{"FOO":"bar"}}'
```