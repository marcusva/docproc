Installing
==========

This section provides an overview and guidance for installing docproc.

Binary Releases
---------------

You can download pre-built binary releases of docproc for different platforms
from https://github.com/marcusva/docproc/releases. If your platform is not
listed, you can also build docproc from source.

Building From Source
--------------------

You can download source snapshots of docproc from
https://github.com/marcusva/docproc/tags. Besides the source distribution, you
also will need the following tools:

* Golang 1.8 or newer (https://golang.org/)
* dep (https://golang.github.io/dep/)

docproc relies on a message queue implementation. It currently supports the
following message queues:

* beanstalk - http://kr.github.io/beanstalkd/
* NSQ - http://nsq.io/
* NATS - https://nats.io/

Unpack the source snapshot into your `GOPATH`, run the `dep` command and
build docproc.

On Unix and Linux run

.. code-block:: console

    $ tar xzvf docproc-.tar.gz $GOPATH
    $ cd $GOPATH/github.com/marcusva/docproc
    $ dep ensure
    $ build-release.sh

On Windows run

.. code-block:: batch

    > unzip docproc-.zip %GOPATH%
    > cd %GOPATH%/github.com/marcusva/docproc
    > dep ensure
    > build-release.bat

Those commands will build a set o docproc release distributions in the `dist`
folder.

Installation
------------

Unpack the matching distribution package for your operating system and copy the
required binaries into the desired target location.

Example for Windows:

.. code-block:: batch

    > unzip docproc-0.0.1-windows-amd64.zip
    > cd docproc-0.0.1-windows-amd64
    > copy docproc*.exe C:\docproc\bin

Example for Linux:

.. code-block:: console

    $ unzip docproc-0.0.1-linux-amd64.zip
    $ cd docproc-0.0.1-linux-amd64
    $ cp docproc*. /usr/local/bin

Set up the configuration files as appropriate and you are good to go.
