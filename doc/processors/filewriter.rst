.. _filewriter:

FileWriter
==========

The FileWriter writes a specific entry of the message content to a file on disk.

Configuration
-------------
The FileWriter requires the following configuration entries:

.. code-block:: ini

    [filewrite-config]
    type = FileWriter
    read.from = htmlresult
    read.base64 = true
    rules = /app/rules/output/file-html.json
    filename = filename
    path = /app/output

type
    To configure a FileWriter, use ``FileWriter`` as ``type``.

read.from
    The path of the message's content to save to the file.

read.base64
    Indicates, if the message content path is base64-encoded and needs to
    be decoded beforehand (necessary, if e.g. the ``read.from`` entry
    originally contained binary data).

    The entry is optional and, if unset, assumed to be ``false``.

filename
    The path of the message's content containing the filename to use.

path
    The directory to use for storing the file.

rules
    The set of rules to use to decide, if the file shall be written or not.
