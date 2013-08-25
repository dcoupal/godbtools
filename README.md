godbtools
=========

Tools to validate NOSQL/document databases.
It currently supports MongoDB, and partially CouchBase.
Documents stored in JSON files can also be used.

Main functionality are:
- verifying collections of database documents against schemas written in JSON Schema


Tool dependencies
-----------------

In order to run support all functionality of the tool, you will need to
have the following tools installed.

  Bazaar - to install 'mgo', the library to connect to MongoDB, no need at runtime
  CouchBase - only needed for CouchBase support and running the corresponding tests
  Git
  Go
  MongoDB - only needed for MongoDB support and running the corresponding tests

If you decide to not support some Database, you can configure the tool, so it does
not try to execute, or run tests, for the non-supported databases.


Source code dependencies
------------------------

  This project depends on the following open source projects.
  Run the following commands once your GOPATH is set.

    go get github.com/sigu-399/gojsonschema
    go get github.com/couchbaselabs/go-couchbase
    go get labix.org/v2/mgo
    

Running the tests
-----------------

  1) import the data in the different databases

  2) chmod +x ./runtests
    ./runtests
    

