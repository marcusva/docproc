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

In-, Output and Error Queue
^^^^^^^^^^^^^^^^^^^^^^^^^^^

The input queue to read messages from is configured via the ``[in-queue]``
section.  The output queue to write processed messages for other consumers is
configured via the ``[out-queue]`` section. If you want to preserve messages,
that failed to process, you can also coonfigure an error queue via the
``[error-queue]`` section.

Configuration entries for the queue(s) may vary slightly, depending on the used
message queue provider. The following entries are required nevertheless.

.. code-block:: ini

    [in-queue]
    type = nsq
    host = 127.0.0.1:4161
    topic = input

    [out-queue]
    type = nsq
    host = 127.0.0.1:4150
    topic = output

    [error-queue]
    type = nsq
    host = 127.0.0.1:4150
    topic = error

type
    The message queue type to use. This can be one of

    * ``beanstalk``
    * ``nsq``

host
    The host or URI to use for connecting to the queue. The exact connection
    string to use varies, depending on the queue type and your service layout.

topic
    ``[in-queue]``: The message queue topic to read messages from for
    processing.

    ``[out-queue]``: The message queue topic to write to. Consumers, such as
    following docproc.proc instances can use the same topic to receive and
    process the incoming messages.

    ``[error-queue]``: The message queue topic to write failed messages to.

Processors
^^^^^^^^^^

Processors are activated in the ``[execute]`` section and configured in
an own, user-defined section. The ``[execute]`` section tells docproc.proc,
which other sections it shall read to configure the appropriate handlers.

handlers
    A comma-separated list of sections to use for configuring and activating
    processors. The entries must match a section within the configuration
    file. The processors are executed in the order of appearance.

Processing the message stops immediately, if one of the configured processors
cannot sucessfully process the message. If an ``[error-queue]`` is configured,
docproc.proc will write the message in its current state to that queue.

.. code-block:: ini

    # Processors
    [execute]
    handlers = add-data, xml-transform

    # Processor of type "ValueEnricher"
    [add-data]
    type = ValueEnricher
    rules = /app/rules/preproc/testrules.json

    # Processor of type TemplateTransformer"
    [xml-transform]
    type = TemplateTransformer
    output = _xml_
    templates = /app/templates/preproc/*.tpl
    templateroot = main

The currently supported processors are explained in the chapter
:ref:`processors`.


.. _RFC-5424: http://www.rfc-base.org/txt/rfc-5424.txt