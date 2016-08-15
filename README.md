Golerta
-------

A simplified reimplementation of [alerta](https://github.com/guardian/alerta) in golang. It has an enterprise 
focus so the biggest change will be using LDAP for the authentication system. Additionally, I will be looking
to reimplement the backend in RethinkDB for simplified clustering and querying.

development environment
-----------------------

Make sure [godep](https://github.com/tools/godep) is installed and working.

    go get github.com/tools/godep

Run rethinkdb in a container on localhost

    docker run -d --name rethinkdb -p 8080:8080 -p 28015:28015 rethinkdb
    
Get golerta code and run it with the example config.

    go get github.com/allen13/golerta
    cd $GOPATH/github.com/allen13/golerta
    godep go run golerta.go --config example.toml

Log in using credentials from the [forumsys test ldap server](http://www.forumsys.com/en/tutorials/integration-how-to/ldap/online-ldap-test-server/) 

    username: gauss
    password: password
    
Run all unit tests:

    go test ./app/...

    