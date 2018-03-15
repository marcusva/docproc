.. _processors:

Processors
==========

docproc's core command, :ref:`proc` offers a set of different, simple
processing tools, which can enhance, change, transform or send message contents.

The following pages provide in-depth information about the different processors
and their usage.

:ref:`contentvalidator`
    Validates the message contents against a predefined set of rules.

:ref:`valueenricher`
    Enables docproc to add new content to a message or to modify
    existing content of the message.

:ref:`templatetransformer`
    Provides templating support via Go's ``text/template`` package.

:ref:`htmlrenderer`
    Provides templating support via Go's ``html/template`` package. It is similar
    to the :ref:`templatetransformer`, except that ``html/template`` contains
    some builtin safety nets for HTML content.

:ref:`filewriter`
    Writes a specific entry of the message content to a file on disk.

:ref:`httpsender`
    Sends a specific entry of the message content via HTTP POST to an HTTP host.

.. toctree::
    :maxdepth: 1
    :hidden:

    contentvalidator
    valueenricher
    templatetransformer
    htmlrenderer
    filewriter
    httpsender
