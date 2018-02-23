.. highlight:: sh

Quick Start
===========

The following steps will run a small docproc content processing pipeline
on your local machine and transform simple CSV input into HTML files.

We are using `NATS`_ as message queue implementation for your local system.
Follow its `installation instructions`_ to get it installed.

#. Start gnatsd in a separate shell::

    $ gnatsd

#. Create two new directories for in- and output named ``data`` and ``output``::

    $ mkdir data
    $ mkdir output

#. Start docproc.fileinput in a separate shell using the NATS example
   configuration. docproc.fileinput will watch the directory named ``data`` in
   the current directory::

        $ docproc.fileinput -c examples/docproc-fileinput-natsio.conf

#. Start docproc.proc in a separate shell using the NATS example
   configuration. docproc.proc will write processed contents into a directory
   named ``output`` in the current directory::

        $ docproc.proc -c examples/docproc-proc-natsio.conf

#. Copy the the test records CSV file into the data directory to start content
   processing::

        $ cp examples/data/testrecords.csv data/testrecords.csv

#. To verify that everything worked as expected, check the ``data`` directory
   for the now processed ``testrecords.csv.DONE`` file and the ``output``
   directory for a set of new HTML files.

For a more sophisticated setup, take a look at the :ref:`docker-setup` section.

.. _NATS: https://nats.io/
.. _installation instructions: https://nats.io/documentation/tutorials/gnatsd-install/
