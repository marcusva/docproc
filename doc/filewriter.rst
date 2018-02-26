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
    identifier = htmlresult
    rules = /app/rules/output/file-html.json
    filename = filename
    path = /app/output

type
    To configure a FileWriter, use ``FileWriter`` as ``type``.

identifier
    The path of the message's content save to the file.

filename
    The path of the message's content containing the filename to use.

path
    The directory to use for storing the file.

rules
    The set of rules to use to decide, if the file shall be written or not.
