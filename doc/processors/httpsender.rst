.. _httpsender:

HTTPSender
==========

The HTTPSender sends a specific entry of the message content via HTTP POST to
an HTTP host.

Configuration
-------------
The HTTPSender requires the following configuration entries:

.. code-block:: ini

    [httpsender-config]
    type = HTTPSender
    read.from = body
    url = http://some.endpoint/receive_msg

type
    To configure a HTTPSender, use ``HTTPSender`` as ``type``.

read.from
    The path of the message's content to send to the host.

url
    The URL to send the content to.
