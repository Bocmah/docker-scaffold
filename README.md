[![Build Status](https://github.com/Bocmah/phpdocker-gen/workflows/main/badge.svg)](https://github.com/Bocmah/phpdocker-gen/actions)
[![codecov](https://codecov.io/gh/Bocmah/phpdocker-gen/branch/master/graph/badge.svg)](https://codecov.io/gh/Bocmah/phpdocker-gen)
[![Go Report Card](https://goreportcard.com/badge/github.com/Bocmah/phpdocker-gen)](https://goreportcard.com/report/github.com/Bocmah/phpdocker-gen)

# Docker config generator for PHP projects

This tool abstracts away a tedious task of writing docker configuration files for PHP projects. It generates ready-to-use
docker config from high-level YAML file.

## The input YAML file
This tool consumes a YAML file, which contains a description of services, which are needed for your PHP application and
some information about your project. 

The following variables can/should be specified at the top indentation level:

| Name        | Type   | Required | Default value                                 | Description                                                  |
|-------------|--------|----------|-----------------------------------------------|--------------------------------------------------------------|
| appName     | string | yes      | -                                             | The name of your application. Can be anything.               |
| projectRoot | string | yes      | -                                             | Path to your project root.                                   |
| outputPath  | string | no       | ```.docker``` folder inside ```projectRoot``` | Path to folder where resulting configuration will be stored. |

Example:

```yaml
appName: awesome-app
projectRoot: /home/user/projects/awesome-project
outputPath: /home/user/awesome-project/custom-folder
```

You should also specify a list of services `services`, which will be discussed in the section below.

## Services
Services list is the core of the input YAML file. Each service you describe will be mapped to a single docker container.

Currently supported services:

```php``` - maps to a container with php-fpm.

Keys:

| Name       | Type    | Required | Default value                    | Description    |
|------------|---------|----------|----------------------------------|----------------|
| version    | numeric | no       | 7.4                              | PHP version    |
| extensions | list    | no       | [mbstring, zip, exif, pcntl, gd] | PHP extensions |

**Note**: ```extensions``` key is experimental. Not all extensions may install correctly.

Example:

```yaml
php:
  version: 7.4
  extensions:
    - redis
    - xdebug
```

```nginx``` - maps to a container with nginx.

Keys:

| Name                       | Type    | Required | Default value | Description                                                                                              |
|----------------------------|---------|----------|---------------|----------------------------------------------------------------------------------------------------------|
| httpPort                   | integer | no       | 80            | nginx will use this port for listening for HTTP requests                                                 |
| httpsPort                  | integer | no       | 443           | nginx will use this port for listening for HTTPS requests                                                |
| serverName                 | string  | yes      | -             | The hostname. The app will be navigable through the web using the value of server name followed by .test |
| fastCGI.passPort           | integer | no       | 9000          | This port will be used for connecting nginx and php-fpm                                                  |
| fastCGI.readTimeoutSeconds | integer | no       | 60            | How long nginx will wait for response from php-fpm before timing out with 504 error                      |

```yaml
nginx:
  httpPort: 81
  serverName: awesome-app
  fastCGI:
    passPort: 9001
    readTimeoutSeconds: 30
```

```nodejs``` - maps to a container with Node.js.

Keys:

| Name    | Type                | Required | Default value | Description     |
|---------|---------------------|----------|---------------|-----------------|
| version | numeric&#124;string | no       | latest        | Node.js version |

Example:

```yaml
nodejs:
  version: 10
```

```database``` - maps to a container with database.

Keys:

| Name         | Type                        | Required                      | Default value                                  | Description                                                                                                          |
|--------------|-----------------------------|-------------------------------|------------------------------------------------|----------------------------------------------------------------------------------------------------------------------|
| system       | enum(mysql&#124;postgresql) | yes                           | -                                              | Database system in use                                                                                               |
| version      | numeric                     | no                            | 8.0 for ```mysql``` 12.3 for ```postgresql```  | Database version                                                                                                     |
| name         | string                      | no                            | -                                              | If specified, database with ```name``` will be created on image startup                                              |
| port         | integer                     | no                            | 3306 for ```mysql``` 5432 for ```postgresql``` | Database port                                                                                                        |
| username     | string                      | no                            | -                                              | If specified, user with ```username``` will be created with superuser power                                          |
| password     | string                      | required for ```postgresql``` | -                                              | Sets the superuser password if system in use is ```postgresql``` or a password for username if system is ```mysql``` |
| rootPassword | string                      | required for ```mysql```      | -                                              | Sets the superuser password if system in use is ```mysql```                                                          |

Example:

```yaml
database:
  system: mysql
  version: 5.7
  name: test-db
  port: 3306
  username: joe
  password: test
  rootPassword: testRoot
```

Sometimes you may wish to use a service but all keys are optional, and you want to leave default values. In this
case you should specify a service as an empty object:

```yaml
nodejs: {}
```

## Full example file

```yaml
appName: awesome-app
projectRoot: /home/user/joe/projects/awesome-app
outputPath: /home/user/joe/projects/awesome-app/docker-config
services:
  php:
    version: 7.4
    extensions:
      - redis
      - xdebug
  nginx:
    httpPort: 81
    serverName: 
    fastCGI:
      passPort: 9001
      readTimeoutSeconds: 30
  nodejs:
    version: 10
  database:
    system: mysql
    version: 5.7
    name: test-db
    port: 3306
    username: bocmah
    password: test
    rootPassword: testRoot
```

Using the above file the tool will create a set of Docker configuration files inside
```/home/user/joe/projects/awesome-app/docker-config```. Final configuration will consist of four docker containers:

1. Php-fpm container with PHP 7.4
2. Nginx container listening on ports 81 and 9001 for HTTP and HTTPS ports respectively.
3. Container with Node.js v10.
4. MySQL container.

## Prerequisites

To run the tool, you should have Go 1.15+ installed.
You will also need Docker and Docker Compose for running output configuration.

## Installation

The following command will download and install the code as a binary. 

```$ go get github.com/Bocmah/phpdocker-gen/...```

**Note**: The location of Go binaries depends on Go environment variables ```GOPATH``` and ```GOBIN```. If you want
to use the binary without specifying its fullpath, you should export the location of Go binaries to your ```PATH```.
If you've just installed Go and didn't make any changes to Go env variables this command should work:

```$ export PATH=$PATH:$(go env GOPATH)/bin```

## Usage

```$ phpdocker-gen -file <path_to_input_file>```

You can use either an absolute path to input file or a path relative to current working directory.