.. _webinput:

docproc.webinput
================

docproc supports processing raw content provided as JSON via HTTP calls
through the *docproc.webinput* application. It can also receive file uploads for
:ref:`fileinput`.

When invoking docproc.webinput, you may specify any of the following options::

    docproc.webinput [-hv] [-c file] [-l file]

.. option:: -c <file>

   Load the configuration from the passed file. If not provided,
   docproc.webinput will try to read from the configuration file
   ``docproc-webinput.conf`` in the current working directory.

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

    # Enabled web input handlers
    [input]
    address = localhost:80
    handlers = web-in, file-in

    # Simple JSON message receiver
    [web-in]
    endpoint = /receive
    type = RawHandler
    maxsize = 128

    # File upload handler to pass files to docproc.fileinput
    [file-in]
    endpoint = /upload
    type = FileHandler
    folder.out = /app/data
    file.prefix = out-
    file.suffix = .csv
    maxsize = 5000

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
    file = /var/log/docproc-webinput.log
    level = Info

.. note::

    It is recommended to use the log level ``Error`` in a production environment
    to spot message processing issues (e.g. a queue being not reachable
    anymore). In rare situations, docproc.webinput may use a more severe log
    level to indicate critical internal application problems.

Output Queue
^^^^^^^^^^^^

The output queue to write messages, generated from the raw HTTP messages, to, is
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
    * ``nsq``

host
    The host or URI to use for connecting to the queue. The exact connection
    string to use varies, depending on the queue type and your service layout.

topic
    The message queue topic to write to. Consumers, such as docproc.proc can
    use the same topic to receive and process the incoming messages from
    docproc.webinput.

HTTP Input
^^^^^^^^^^

HTTP input handlers are activated in the ``[input]`` section and configured in
an own, user-defined section. The ``[input]`` section tells docproc.webinput,
which other sections it shall read to configure the appropriate handlers.

The currently supported handlers are explained the input handler section below.

address
    The host and port the webinput application shall listen on. It can be
    an IP, hostname, IP:port, hostname:port or just a port, e.g. ::

        localhost
        localhost:8001
        192.168.134.11
        :8080

handlers
    A comma-separated list of sections to use for configuring and activating
    input handlers. The entries must match a section within the configuration
    file.

.. code-block:: ini

    [input]
    address = 127.0.0.1:80
    # Set up two handlers, which are configured in [web-in] and [file-in]
    handlers = web-in, file-in

    [web-in]
    ...

    [file-in]
    ...

Input Handlers
^^^^^^^^^^^^^^

:ref:`webinput` comes with support for processing raw content provided as JSON
via HTTP as well as a basic file-upload mechanism, which allows for easy
HTTP integration of the :ref:`fileinput` file processing application.

Each individual HTTP handler shares a common set of configuration entries:

.. code-block:: ini

    [your-config]
    endpoint = <relative HTTP URL endpoint>
    type = <relevant input handler>
    maxsize = <maximum allowed message size in kB>

endpoint
    The URL endpoint to listen on, relative to the address part, including
    leading separators: ``http:/<host:port>[endpoint]``. The configured address
    (see above) and endpoint represent the full URL, the handler will listen on.

type
    The HTTP handler to use. See below for a list of currently available
    handlers.

maxsize
    The maximum allowed HTTP message size in kilobytes.

HTTP Handlers
-------------

:ref:`rawhandler`
    Processes HTTP requests messages, which provide the content to transform as
    JSON.

:ref:`filehandler`
    Allows to upload files for further processing by :ref:`fileinput`.

.. toctree::
    :hidden:
    :maxdepth: 1

    input/rawhandler
    input/filehandler

.. _RFC-5424: http://www.rfc-base.org/txt/rfc-5424.txt