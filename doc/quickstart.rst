.. highlight:: sh

Quick Start
===========

The following steps will run a small docproc content processing pipeline
on your local machine and transform simple CSV input into HTML files.

We are using `NSQ`_ as message queue implementation for your local system.
Follow its `installation instructions`_ to get it installed.

#. Start nsqlookupd in a separate shell::

    $ nsqlookupd

#. Start nsqd in a separate shell::

    $ nsqd -lookupd-tcp-address 127.0.0.1:4160

#. Create two new directories for in- and output named ``data`` and ``output``::

    $ mkdir data
    $ mkdir output

#. Create the NSQ topic. *Note: this is an optional step to avoid delays in the
   service discovery.*::

    $ curl -X POST http://127.0.0.1:4161/topic/create?topic=input

#. Start docproc.fileinput in a separate shell using the NSQ example
   configuration. docproc.fileinput will watch the directory named ``data`` in
   the current directory::

    $ docproc.fileinput -c examples/docproc-fileinput.conf

#. Start docproc.proc in a separate shell using the docproc-proc example
   configuration. docproc.proc will write processed contents into a directory
   named ``output`` in the current directory::

        $ docproc.proc -c examples/docproc-proc.conf

#. Copy the the test records CSV file into the data directory to start content
   processing::

        $ cp examples/data/testrecords.csv data/testrecords.csv

#. To verify that everything worked as expected, check the ``data`` directory
   for the now processed ``testrecords.csv.DONE`` file and the ``output``
   directory for a set of new HTML files.

For a more sophisticated setup, take a look at the docker configuration for the
integration tests, which can be found in the test folder of the *source* distribution.

.. _nsq: https://nsq.io/
.. _installation instructions: https://nsq.io/deployment/installing.html
