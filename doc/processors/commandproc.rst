.. _commandproc:

CommandProc
===========

The CommandProc executes a command on a specific content part of a message and
stores the result in the message. This allows docproc to perform more complex
tasks with external utilities.

Configuration
-------------

The CommandProc requires the following configuration entries:

.. code-block:: ini

    [commandproc-config]
    type  = CommandProc
    read.from = inputfield
    store.in = outputfield
    store.base64 = false
    exec = xsltproc --novalid myxslt.xslt

type
    To configure a CommandProc, use ``CommandProc`` as ``type``.

read.from
    The path of the message's content to pass to the command. The content will
    be stored in a temporary file and is passed as **last** argument to the
    ``exec`` command line.

store.in
    The path to use on the message's content to store the output result in.

store.base64
    Indicates, if the output result must be base64-encoded. This is necessary,
    if e.g. the output is binary content, which must not be converted to a
    string.

    The entry is optional and, if unset, assumed to be ``false``.

exec
    The command to execute. The configured content of the message will be
    stored in a temporary file and passed as last argument to the command.
    The ``exec`` configuration entry::

        exec = xsltproc --novalid myxslt.xslt

    thus effectively will be executed like::

        xsltproc --novalid myxslt.xslt /tmp/tmp-12347-docproc-cmdproc
