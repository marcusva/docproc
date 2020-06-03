.. _httpsender:

HTTPSender
==========

The HTTPSender sends a specific entry of the message content via HTTP POST to
an HTTP host. The HTTPSender can utilize both, HTTP and HTTPS connections.

Configuration
-------------
The HTTPSender requires the following configuration entries:

.. code-block:: ini

    [httpsender-config]
    type = HTTPSender
    read.from = body
    url = http://some.endpoint/receive_msg
    timeout = 10
    headers = "http-headers"

type
    To configure a HTTPSender, use ``HTTPSender`` as ``type``.

read.from
    The path of the message's content to send to the host.

url
    The URL to send the content to.

timeout
    The time limit for an individual request in seconds. 

    The entry is optional and, if unset, set to 5 seconds.

headers
    The path of the message's content to read additional HTTP headers
    from. docproc always will set the ``Content-Length`` 
    The entry is optional and, if unset or empty,
    ``"Content-Type": "text/plain"`` will be added to the HTTP headers.

