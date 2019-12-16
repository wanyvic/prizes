#!/bin/bash
wget -qO - https://www.mongodb.org/static/pgp/server-4.2.asc | sudo apt-key add -
echo "deb [ arch=amd64 ] https://repo.mongodb.org/apt/ubuntu xenial/mongodb-org/4.2 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-4.2.listsudo apt-get update
sudo apt-get install -y mongodb-org
echo -e "security:\n\tauthorization: enabled"  | sudo tee /etc/mongod.conf
use docker
db.createUser(
  {
    user: "xxxx",
    pwd: "xxxxx",
    roles: [ { role: "dbOwner", db: "docker" }]
  }
)
sudo service mongod start