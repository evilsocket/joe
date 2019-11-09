# Joe

Joe is a tool to automatically build a REST API, its documentation and charts around SQL queries and their resulting
 data. 

In a way it is an anti-[ORM](https://en.wikipedia.org/wiki/Object-relational_mapping): its purpose is to help backend engineers versioning
, annotating, exposing and charting the queries that can't be implemented or aren't worth implementing in the backend
 main business logic and that they would normally keep on .txt or .sql files.

## How to Install

```sh
go get -u https://github.com/evilsocket/joe/cmd/joe
```

## Example

First create an `/etc/joe/joe.conf` configuration file with the access credentials for the database:

```conf
# CHANGE THIS: use a complex secret for the JWT token generation
API_SECRET=02zygnJs5e0bBLJjaHCinWTjfRdheTYO

DB_HOST=joe-mysql
DB_DRIVER=mysql
DB_USER=joe
DB_PASSWORD=joe
DB_NAME=joe
DB_PORT=3306
```

Then create the `admin` user (this command will generate the file `/etc/joe/users/admin.yml`):

```sh
sudo mkdir -p /etc/joe/users
sudo joe -new-user admin -token-ttl 6 # JWT tokens for this user expire after 6 hours
```

For query and chart examples [you can check this repository](https://github.com/evilsocket/pwngrid-queries-joe). Once
 you have your `/etc/joe/queries` folder with your queries, you can:

Generate markdown documentation:

```sh
joe -doc /path/to/document.md
```

Start the joe server:

```sh
joe -conf /etc/joe/joe.conf -data /etc/joe/queries -users /etc/joe/users
```
        
## License

`joe` is made with ♥  by [@evilsocket](https://twitter.com/evilsocket) and it is released under the GPL3 license.