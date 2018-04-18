.. _contentvalidator:

ContentValidator
================

The ContentValidator checks, if the contents of a message conform to a set
of predefined rules. This allows docproc to process only those messages, which
are considered functionally valid.

Configuration
-------------

The ContentValidator requires the following configuration entries:

.. code-block:: ini

    [contentvalidator-config]
    type  = ContentValidator
    rules = /path/to/a/rules/set

type
   To configure a ContentValidator, use ``ContentValidator`` as ``type``.

rules
   ``rules`` refers to a file on disk containing the rules to be executed on
   the message content. A rule consists of one or more conditions.

Defining Rules
--------------

The rules to be executed are kept in a simple JSON-based list.

.. code-block:: json

    [
        {
            "name": "First rule",
            "path": "NET",
            "op": "less than",
            "value": 0,
        },
        {
            "name": "Second rule",
            "path": "ZIP",
            "op": "exists",
        }
    ]

See :ref:`rulesengine` for more details about how to configure rules.

Exampe Usage
------------

Only messages containing a customer number shall be processed. So let's assume
the following input content with a missing customer number: ::

    CUSTNO;FIRSTNAME;LASTNAME;STREET;ZIP;CITY;NET;GROSS;DATE
    100112;John;Doe;Example Lane 384;10006;New York;10394.00;12386.86;2017-04-07
    ;Jane;Doeanne;Another Lane 384;10009;New York;-376.00;-405.88;2017-05-18
    194227;Max;Mustermann;Musterstraße 123;12345;Berlin;0.00;0.00;2017-12-04

Define a rule in a ``validate.json`` file:

.. code-block:: json

    [
        {
            "name": "Check for a customer number",
            "path": "CUSTNO",
            "op": "!=",
            "value": "",
        }
    ]


Set up the ContentValidator in your docproc.proc configuration:

.. code-block:: ini

    [execute]
    handlers = validate-data

    [validate-data]
    type  = ContentValidator
    rules = /app/docproc/rules/validate.json

Messages with an empty customer produce an error within the log now and are put
into the error queue, if configured.
