.. _rditransformer:

RDITransformer
==============

SAP RDI files can be read using the ``RDITransformer`` within the input handler
configuration:

.. code-block:: ini

    [your-config]
    format = rdi
    ...
    transformer = RDITransfomer

RDI files picked up by the ``RDITransfomer`` are assumed to be Gzip-compressed,
regardless of their suffix.

When processing an RDI file, the RDITransformer creates one or more messages.
A new message will be created, whenever a header ('H') window is found. All
data windows following the header will be put into the same message until the
next header window is found or if there are no more data windows to read. ::

    +---------------------------------------------------------+--------------+
    | RDI contents                                            | docproc      |
    |                                                         | Message      |
    +=========================================================+==============+
    | H0123456789789789789...                                 |   Message 1  |
    | S                                                       |              |
    | CCODEPAGE ...                                           |              |
    | C...                                                    |              |
    | CCODEPAGE 1100                                          |              |
    | DMAIN      SECTION_A                 ABCD ...           |   Content    |
    | DMAIN      SECTION_A                 FIELDX ...         |     of       |
    | DMAIN      SECTION_A                 FIELDQ ...         |  Message 1   |
    | ...                                                     |              |
    |                                                         +==============+
    | H0123456789789789789 ...                                | Message 2    |
    | S                                                       |              |
    | CCODEPAGE ...                                           |              |
    | C...                                                    |              |
    | CCODEPAGE 1100                                          |              |
    | DMAIN      SECTION_A                 ABCD ...           |   Content    |
    | DMAIN      SECTION_A                 FIELDX ...         |     of       |
    | DMAIN      SECTION_A                 FIELDQ ...         |  Message 2   |
    | C...                                                    |              |
    | DMAIN      SECTION_B                 FIELD99 ...        |              |
    | ...                                                     |              |
    +---------------------------------------------------------+--------------+

Control ('C') and sort ('S') windows will be skipped and have no effect on the
message order or content layout.

.. note::

    The ``RDITransformer`` follows an all-or-nothing approach when processing
    an RDI file. The created messages are only placed on the queue, if the
    whole RDI file can be read and transformed sucessfully.

The resulting message(s) consist of a content section, which contain one or more
``sections`` entries named after the data window that contains the fields.

The above example would produce the following messages.

**Message 1**

.. code-block:: json

    {
        "metadata": {
            "format": "rdi",
            "batch": 1517607828,
            "created": "2018-02-02T22:43:48.0220047+01:00"
        },
        "content": {
            "sections": [
                {
                    "name": "SECTION_A",
                    "content": {
                        "ABCD": "...",
                        "FIELDX": "...",
                        "FIELDQ": "...",
                    }
                }
            ]
        }
    }

**Message 2**

.. code-block:: json

    {
        "metadata": {
            "format": "rdi",
            "batch": 1517607828,
            "created": "2018-02-02T22:43:48.0220047+01:00"
        },
        "content": {
            "sections": [
                {
                    "name": "SECTION_A",
                    "content": {
                        "ABCD": "...",
                        "FIELDX": "...",
                        "FIELDQ": "...",
                    }
                },
                {
                    "name": "SECTION_B",
                    "content": {
                        "FIELD_99": "...",
                    }
                }
            ]
        }
    }

