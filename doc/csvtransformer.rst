.. _csvtransformer:

CSVTransformer
==============

CSV files can be read using the ``CSVTransformer`` within the input handler
configuration.

.. code-block:: ini

    [your-config]
    format = csv
    ...
    transformer = CSVTransformer
    delim = ;

delim
    The column separator to use.

When processing a CSV file, the CSVTransformer creates one or more messages,
depending on the amount of rows within the CSV file.
The first row is considered the header row and its column values are used as
field names for the message content.

The CSV contents

.. code-block:: none

    CUSTNO;FIRSTNAME;LASTNAME
    100112;John;Doe
    194228;Manuela;Mustermann

would result in two messages to be created:

Message 1
    .. code-block:: json

        {
            "metadata": {
                "format": "csv",
                "batch": 1517607828,
                "created": "2018-02-02T22:43:48.0220047+01:00"
            },
            "content": {
                "CUSTNO": "100112",
                "FIRSTNAME": "John",
                "LASTNAME": "Doe"
            }
        }

Message 2
    .. code-block:: json

        {
            "metadata": {
                "format": "csv",
                "batch": 1517607828,
                "created": "2018-02-02T22:43:48.0220047+01:00"
            },
            "content": {
                "CUSTNO": "194228",
                "FIRSTNAME": "Manuela",
                "LASTNAME": "Mustermann"
            }
        }