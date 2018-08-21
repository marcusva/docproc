.. _fileinput:

docproc.fileinput
=================

docproc supports processing content from files, such as CSV or SAP RDI, via
the *docproc.fileinput* application.

When invoking docproc.fileinput, you may specify any of the following options::

    docproc.fileinput [-hv] [-c file] [-l file]

.. option:: -c <file>

   Load the configuration from the passed file. If not provided,
   docproc.fileinput will try to read from the configuration file
   ``docproc-fileinput.conf`` in the current working directory.

.. option:: -h

   Print a short description of all command-line options and exit.

.. option:: -l <file>

   Log information to the passed file.

.. option:: -v

   Print version details and information about the build configuration and exit.

Configuration
-------------

The configuration file uses an INI-style layout and contains several sections,
some of them being optional and some of them being mandatory.

.. code-block:: ini

    [log]
    # log to a specific file instead of stdout
    # file=<path/to/the/file>
    # level can be one of Emergency,Alert,Critical,Error,Warning,Notice,Info,Debug
    level = Info

    # Queue to write the read messages to
    [out-queue]
    type = nsq
    host = 127.0.0.1:4150
    topic = input

    # Enabled file input handlers
    [input]
    handlers = rdi-in, csv-in

    # SAP RDI file handler
    [rdi-in]
    format = rdi
    transformer = RDITransfomer
    folder.in = data
    pattern = *.gz
    interval = 2

    # CSV file handler
    [csv-in]
    format = csv
    transformer = CSVTransformer
    delim = ;
    folder.in = data
    pattern = *.csv
    interval = 2

Logging
^^^^^^^

Logging is configured via the ``[log]`` section. It can contain two entries.
The ``[log]`` section is optional and, if omitted, logging will happen on
STDERR with the log level ``Error``.

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

.. note::

    It is recommended to use the log level ``Error`` in a production environment
    to spot message processing issues (e.g. a queue being not reachable
    anymore). In rare situations, docproc.fileinput may use a more severe log
    level to indicate critical internal application problems.

Output Queue
^^^^^^^^^^^^

The output queue to write messages, generated from the input files, to, is
configured via the ``[out-queue]`` section. Configuration entries for the queue
may vary slightly, depending on the used message queue provider. The following
entries are required nevertheless.

.. code-block:: ini

    [out-queue]
    type = nsq
    host = 127.0.0.1:4150
    topic = input

type
    The message queue type to use. This can be one of

    * ``beanstalk``
    * ``nats``
    * ``nsq``

host
    The host or URI to use for connecting to the queue. The exact connection
    string to use varies, depending on the queue type and your service layout.

topic
    The message queue topic to write to. Consumers, such as docproc.proc can
    use the same topic to receive and process the incoming messages from
    docproc.fileinput.

File Input
^^^^^^^^^^

File input handlers are activated in the ``[input]`` section and configured in
an own, user-defined section. The ``[input]`` section tells docproc.fileinput,
which other sections it shall read to configure the appropriate handlers.

The currently supported handlers are explained in file input handlers.

handlers
    A comma-separated list of sections to use for configuring and activating
    input handlers. The entries must match a section within the configuration
    file.

.. code-block:: ini

    [input]
    # Set up two handlers, which are configured in [rdi-in] and [csv-in]
    handlers = rdi-in, csv-in

    [rdi-in]
    ...

    [csv-in]
    ...

File Input Handlers
^^^^^^^^^^^^^^^^^^^

:ref:`fileinput` comes with support for converting different file formats and
file content into processable messages, which can be individually activated
and configured.

It currently supports the conversion of the following input formats:

* SAP RDI spool files via the :ref:`rditransformer`
* CSV data via the :ref:`csvtransformer`

Each individual file handler shares a common set of configuration entries:

.. code-block:: ini

    [your-config]
    format = <format-name>
    folder.in = <directory to check>
    pattern = <file pattern to check>
    interval = <check interval>
    transformer = <the relevant input transformer>
    # additional, transformer-specific configuration entries

format
    The input file format. This is mainly used for informational purposes within
    the message's metadata and does not have any effect on the message
    processing.

folder.in
    The directory to watch for RDI files to be processed.

pattern
    The file pattern to use for identifying RDI files. This can be a wildcard
    pattern, strict file name matching or regular expression that identifies
    those files, that shall be picked up by the ``RDITransformer``.

interval
    The time interval in seconds to use for checking for new files. This must
    be a positive integer.

transformer
    The input transformer to use. See below for a list of currently available
    input transformers.

Input Transformers
------------------

:ref:`rditransformer`
    Processes SAP RDI spool files and transforms the contained documents into
    messages.

:ref:`csvtransformer`
    Processes CSV files and transforms the contained rows into messages.

.. toctree::
    :hidden:
    :maxdepth: 1

    input/rditransformer
    input/csvtransformer

.. _RFC-5424: http://www.rfc-base.org/txt/rfc-5424.txt