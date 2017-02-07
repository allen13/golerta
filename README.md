Golerta
-------

A simplified reimplementation of [alerta](https://github.com/guardian/alerta) in golang.

Rethinkdb is used instead of Mongodb for operational simplicity and real-time data features.

features
--------

   * All-in-one server - static content, api, and continuous queries
   * LDAP Authentication
   * OAuth Authenticaion
   * Alert routing via plugins - current plugins: file, pagerduty
   * Optional alert time outs - if an alert has a timeout field it will escalate to critical if a new alert isn't received within a given amount of time
   * Timed acknowledgements - alerts will reopen after a specified amount of time
   * Flap detection - alerts that continually change severity state will be marked as flapping
   * Select and update alerts from the list view
   * IN DEVELOPMENT: Web sockets on the client side - web browser no longer polls and new alerts are received immediately
   * IN DEVELOPMENT: Fully distributed - Currently the continuous queries can only be run on a single node without duplicating effort.


docker image
------------
Example: Run golerta minimal

    docker run -p 5608:5608 -e RETHINKDB_ADDRESS=your-rethinkdb-host:28015 -e allen13/golerta:latest

Example: Run golerta by linking to both the rethinkdb and smtp service containers.

    docker run -p 5608:5608 --link rethinkdb:rethinkdb --link smtp:smtp -e RETHINKDB_ADDRESS=rethinkdb:28015 -e EMAIL_NOTIFIER_ENABLED=true -e EMAIL_SMTP_SERVER=smtp golerta

development environment
-----------------------

Run rethinkdb in a container on localhost

    docker run -d --name rethinkdb -p 8080:8080 -p 28015:28015 rethinkdb

Run postfix in a container on localhost
		docker run -p 25:25 --name smtp -e maildomain=localhost -e smtp_user=test1:password,test2:password -d catatnight/postfix

Get golerta code and run it with the example config.

    go get github.com/allen13/golerta
    cd $GOPATH/github.com/allen13/golerta
    go run golerta.go server --config golerta.toml

Log in using credentials from the [forumsys test ldap server](http://www.forumsys.com/en/tutorials/integration-how-to/ldap/online-ldap-test-server/)

    username: gauss
    password: password

Run all unit tests:

    go test ./app/...


