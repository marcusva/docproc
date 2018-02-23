.. _inputhandlers:

File Input Handlers
===================

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

.. toctree::
    :maxdepth: 1
    :caption: Available Input Transformers:

    rditransformer
    csvtransformer