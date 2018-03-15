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
