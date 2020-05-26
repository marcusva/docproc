Installing
==========

This section provides an overview and guidance for installing docproc.

Binary Releases
---------------

You can download pre-built binary releases of docproc for different platforms
from https://github.com/marcusva/docproc/releases. If your platform is not
listed, you can also build docproc from source. Take a look at :ref:`building`
for more details, if necessary.

Installation
------------

Unpack the matching distribution package for your operating system and copy the
required binaries into the desired target location.

Example for Windows:

.. code-block:: batch

    > unzip docproc-<version>-windows-amd64.zip
    > cd docproc-<version>-windows-amd64
    > copy docproc*.exe C:\docproc\bin

Example for Linux:

.. code-block:: console

    $ unzip docproc-<version>-linux-amd64.zip
    $ cd docproc-<version>-linux-amd64
    $ cp docproc*. /usr/local/bin

Set up the configuration files as appropriate and you are good to go.

.. _building:

Building From Source
--------------------

You can download source snapshots of docproc from
https://github.com/marcusva/docproc/tags. Besides the source distribution, you
also will need the following tools:

* Golang 1.13 or newer (https://golang.org/)

docproc relies on a message queue implementation. It currently supports the
following message queues:

* beanstalk - http://beanstalkd.github.io/beanstalkd/
* NSQ - http://nsq.io/

On Unix and Linux, unpack docproc-|version|.tar.gz into :envvar:`$GOPATH/src`, then
run

.. code-block:: console

    $ cd $GOPATH/src/github.com/marcusva/docproc
    $ make install

This will install the docproc binaries into ``/usr/local/bin`` by default. You
can change the :envvar:`PREFIX` as well as :envvar:`DESTDIR` for your own
installation scheme.

On Windows, unpack docproc-|version|.zip into :envvar:`%GOPATH%\\src`, then run

.. code-block:: batch

    > cd %GOPATH%/src/github.com/marcusva/docproc
    > make.bat

Those commands will build docproc and put the binaries, documentation and
examples into the `dist` folder. Copy them into the desired locations.
