.. _filehandler:

FileHandler
===========

Uploading files via HTTP requests can be configured by setting up the
``FileHandler`` within the input handler configuration.

.. code-block:: ini

    [your-config]
    endpoint = /upload
    maxsize = 5000
    type = FileHandler
    folder.out = /app/data
    file.prefix = out-
    file.suffix = .csv

folder.out
    The directory to use for storing the file.

file.prefix
    File names are created randomly by the ``FileHandler``. To allow other
    software to automatically recognize files provided by this specific
    instance, a fixed prefix can be added to the random part.

file.suffix
    File names are created randomly by the ``FileHandler``. To allow other
    software to automatically recognize files provided by this specific
    instance, a fixed suffix can be added to the random part.
