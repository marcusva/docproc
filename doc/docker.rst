.. highlight:: sh
.. _docker-setup:

Docker Setup
=============

The following information will guide you through a simple configuration
scenario for creating your own `docker`_ setup. They can be summed up as

#. create the base image via docker
#. build the docproc images via docker-compose
#. run everything via docker-compose

.. note::
    You will need the source distribution of docproc, which you can find at
    https://github.com/marcusva/docproc/tags for stable snapshots.

Base Image
----------

The docproc base image contains all docproc applications as well as a
``nsqd`` binary to get docproc up and running for testing with the `NSQ`_
message queue system.

Create the base image with the following instructions::

    $ docker build -t docproc/base .

The base image is now registered in your local docker registry as
``docproc/base``.

.. note::
    The docproc applications of the base image will be built with nsq
    support only. To change this behaviour, you can tweak the ``BUILD_FLAGS``
    within ``Dockerfile`` as necessary or override the ``BUILD_FLAGS`` at the
    command line::

        $ docker build --build-arg BUILD_FLAGS="-tags beanstalk" -t docproc/base .

    The nsqd binary will be built and installed nevertheless, if ``Dockerfile``
    is not edited, though.

Build docproc Images
--------------------

Create all docproc images with the following instruction::

    $ docker-compose build

This creates the following set of docproc images:

* ``docproc/fileinput``
* ``docproc/preproc``
* ``docproc/renderer``
* ``docproc/postproc``
* ``docproc/output``

Each image can be run individually. Each image runs a local ``nsqd`` server to
be used by the individual docproc executable.

Run Everything
--------------

All services, including an nsqd, nsqlookupd and nsqadmin instance can be run via::

    $ docker-compose up

.. todo::
    document ports and directories properly.

.. _docker: https://docker.com
.. _NSQ: https://nsq.io
