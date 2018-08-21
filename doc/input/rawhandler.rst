.. _rawhandler:

RawHandler
==========

Content transformation via HTTP requests can be configured by setting up the
``RawHandler`` within the input handler configuration.

.. code-block:: ini

    [your-config]
    endpoint = /receive
    type = RawHandler
    maxsize = 128
