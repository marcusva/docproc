ValueEnricher
=============

The ValueEnricher enables docproc to add new content to a message or to modify
existing content of the message. It uses a simple rules engine to conditionally
modify or add content.

Configuration
-------------
The ValueEnricher requires the following configuration entries:

.. code-block:: ini

    [valueenricher-config]
    type  = ValueEnricher
    rules = /path/to/a/rules/FileWriter

type
   To configure a ValueEnricher, use ``ValueEnricher`` as ``type``.

rules
   ``rules`` refers to a file on disk containing the rules to be executed on
   the message content. A rule consists of one or more conditions and a target
   value to set.

Defining Rules
--------------
The rules to be executed are kept in a simple JSON-based list.

.. code-block:: json-object

    [
        {
            "name": "First rule",
            "path": "NET",
            "op": "less than",
            "value": 0,
            "targetpath": "DOCTYPE",
            "targetvalue": "CREDIT_NOTE"
        },
        {
            "name": "Second rule",
            "path": "ZIP",
            "op": "exists",
            "targetpath": "HAS_ZIP",
            "targetvalue": true
        }
    ]

A rule to be used for the ValueEnricher consists of

name
   An optional name describing the rule. This is for maintenance purposes and
   does not have any effect on the rule, if provided or absent.

path
   The message's content element to check. Paths can be nested using a dotted
   notation.
   ``"path": "ZIP"`` refers to a content entry named "ZIP" on the topmost level
   of the message's content. Such an entry would e.g. match the following
   content::

        ...
        "content": {
            ...
            "CITY": "New York",
            "ZIP": "10006",
            ...
        }

   ``"path": "address.ZIP"`` refers to a content entry named "ZIP" within the
   "address" element of the message's content::

        ...
        "content": {
            ...
            "address": {
                "CITY": "New York",
                "ZIP": "10006",
                ...
            },
        ...
        }

   "path" can also contain brackets to access arrays, e.g.
   ``"path": "chargenumbers[2]"``::

        ...
        "content": {
            "chargenumbers": [
                14865,
                77896,
                12345
            ]
        }

value
    The value to compare the path's value against. **value** can be omitted, if
    the comparision operator is ``exists`` or ``not exists``. If it is provided
    for those operators, it will be ignored.

op
   The comparision operator to use. If not stated otherwise, the comparision
   will consider path being the left-hand and value the right-hand argument::

     value-of-path <op> rule-value

   See :ref:`rulesengine` for more details about the supported operators.
   docproc's rule engine currently understands the following operators:

targetpath
    Defines the path to use for writing the provided targetvalue. If the given
    path does not exist, it will be created. Similarily to the "path", the
    targetpath can be nested using a dotted notation.
    NOTE: Accessing arrays is currently not possible.

targetvalue
    The value to write into targetpath. if value or a part of it is surrounded
    by ``${}``, that specific part is treated as an existing path to be taken from
    the message's content.::

        ...
        "content": {
            ...
            "CITY": "New York",
            "ZIP": "10006",
            ...
        }

        {
            "path": "CITY",
            "op": "equals",
            "value": "New York",
            "targetpath": "PREFIXED_ZIP",
            "targetvalue": "NY-${ZIP}"
        },

        ...
        "content": {
            ...
            "CITY": "New York",
            "ZIP": "10006"
            "PREFIXED_ZIP": "NY-10006",
            ...
        }

Rules can also be chained, allowing them to evaluate multiple comparision before
applying the new target value. If we want a prefixed ZIP code only, if a ZIP
code is provided and if the city is New York, the rule can be written like this::

    {
        "path": "ZIP",
        "op": "exists",
        "subrules": [
            {
                "path": "CITY"
                "op": "equals",
                "value": "New York"
            }
        ]
        "path": "CITY",
        "targetpath": "PREFIXED_ZIP",
        "targetvalue": "NY-${ZIP}"
    },

Of course, any subrule can have subrule on its own. Note, that subrules do not
contain a target path or value, though.
