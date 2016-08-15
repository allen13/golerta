Golerta
-------

A simplified reimplementation of [alerta](https://github.com/guardian/alerta) in golang. It has an enterprise 
focus which focuses on using LDAP for the authentication system.

Rethinkdb is used instead of Mongodb for operational simplicity.

Additionally golerta serves the api as well as the static user content. Alerta required the static content to be hosted separately.

docker image
------------

    docker run -p 5608:5608 -e RETHINKDB_ADDRESS=your-rethinkdb-host:28015 allen13/golerta:latest
    
development environment
-----------------------

Make sure [godep](https://github.com/tools/godep) is installed and working.

    go get github.com/tools/godep

Run rethinkdb in a container on localhost

    docker run -d --name rethinkdb -p 8080:8080 -p 28015:28015 rethinkdb
    
Get golerta code and run it with the example config.

    go get github.com/allen13/golerta
    cd $GOPATH/github.com/allen13/golerta
    godep go run golerta.go server --config example.toml

Log in using credentials from the [forumsys test ldap server](http://www.forumsys.com/en/tutorials/integration-how-to/ldap/online-ldap-test-server/) 

    username: gauss
    password: password
    
Run all unit tests:

    go test ./app/...

    