.. _proc:

docproc.proc
============

Processing content is mainly done by docproc.proc, with some minor exception
for docproc.fileinput, which feeds file-based data into the processing queues
being handled by docproc.proc.

When invoking docproc.proc, you may specify any of the following options::

    docproc.proc [-hv] [-c file] [-l file]

.. option:: -c <file>

   Load the configuration from the passed file. If not provided,
   docproc.proc will try to read from the configuration file
   ``docproc-proc.conf`` in the current working directory.

.. option:: -h

   Print a short description of all command-line options and exit.

.. option:: -l <file>

   Log information to the passed file.

.. option:: -v

   Print version details and information about the build configuration and exit.

Configuration
-------------

The configuration file uses an INI-style layout.

Logging
^^^^^^^

Logging is configured via the ``[log]`` section. It can contain two entries.

file
    The file to use for logging. This can be a file or writable socket.
    If omitted, STDERR will be used.

level
    The log level to use. The log level uses the severity values of `RFC-5424`_
    in either numerical (``0``, ``3``, ...) or textual form (``Error``,
    ``Info``, ...). If omitted, ``Error`` will be used.

.. code-block:: ini

    [log]
    file = /var/log/docproc-fileinput.log
    level = Info

The ``[log]`` section is optional and, if omitted, logging will happen on
STDERR with the log level ``Error``.

.. note::

    It is recommended to use the log level ``Error`` in a production environment
    to spot message processing issues (e.g. a queue being not reachable
    anymore). In rare situations, docproc.fileinput may use a more severe log
    level to indicate critical internal application problems.

.. _RFC-5424: http://www.rfc-base.org/txt/rfc-5424.txt