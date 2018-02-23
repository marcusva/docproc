.. _proc:

docproc.proc
============

Processing content is mainly done by docproc.proc, with some minor exception
for docproc.fileinput, which feeds file-based data into the processing queues
being handled by docproc.proc.

When invoking docproc.proc, you may specify any of the following options::

    docproc.proc [-hv] [-c file] [-l file]

.. option:: -c <file>

   Load the configuration from the passed file.

.. option:: -h

   Print a short description of all command-line options and exit.

.. option:: -l <file>

   Log information to the passed file.

.. option:: -v

   Print version details and information about the build configuration and exit.
