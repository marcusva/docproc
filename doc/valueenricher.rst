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

.. code-block:: json

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

A rule to be used for the ValueEnricher consists of the following fields:

* ``name``, ``path``, ``op``, ``value``
* ``targetpath`` - *specific to the ValueEnricher*
* ``targetvalue`` - *specific to the ValueEnricher*

See :ref:`rulesengine` for more details about how to configure rules.
Rules being used by the ValueEnricher contain two additional fields:

targetpath
    Defines the path to use for writing the provided targetvalue. If the given
    path does not exist, it will be created. Similarily to the ``path``, the
    targetpath can be nested using a dotted notation.

    *Accessing arrays is currently not possible.*

targetvalue
    The value to write into targetpath. The value can contain portions of the
    existing message's content using a ``${<sourcepath>}`` notation.

Defining Target Paths
---------------------

Target paths to write content to can be defined in the same way as the source
paths for comparision. A target path can refer to an existing path, causing it
to be overwritten with the new value on evaluating the rule successfully. The
target path can also be a completely new path, that will be created, if the
rule is successful.

Let's add a city name based on the provided shortcut for the following message.

Message:
    .. code-block:: json

        {
            "content": {
                "city_sc": "NY"
            }
        }

Rule:
    .. code-block:: json

        {
            "path": "city_sc",
            "op": "equals",
            "value": "NY",
            "targetpath": "city",
            "targetvalue": "New York"
        }

Resulting Message:
    .. code-block:: json

        {
            "content": {
                "city_sc": "NY",
                "city": "New York"
            }
        }

Overwrite the city's shortcut with the city name

Message:
    .. code-block:: json

        {
            "content": {
                "city": "NY"
            }
        }

Rule:
    .. code-block:: json

        {
            "path": "city",
            "op": "equals",
            "value": "NY",
            "targetpath": "city",
            "targetvalue": "New York"
        }

Resulting Message:
    .. code-block:: json

        {
            "content": {
                "city": "New York"
            }
        }

Add an address block containing the city name.

Message:
    .. code-block:: json

        {
            "content": {
                "city_sc": "NY"
            }
        }

Rule:
    .. code-block:: json

        {
            "path": "city_sc",
            "op": "equals",
            "value": "NY",
            "targetpath": "address.city",
            "targetvalue": "New York"
        }

Resulting Message:
    .. code-block:: json

        {
            "content": {
                "city_sc": "NY",
                "address": {
                    "city": "New York"
                }
            }
        }

Defining Target Values
----------------------

Target value can be any kind of atomic value types, such as integers, decimal
numbers, boolean values or strings. More complex values, such as JSON objects,
maps or arrays are not supported.

Furthermore, target values can copy the values from existing paths, as long as
those contain atomic value types. To refer to an existing path, use``${}``.

Prefix the ZIP code with state information for New York:

Message:
    .. code-block:: json

        {
            "content": {
                "CITY": "New York",
                "ZIP": "10006",
            }
        }

Rule:
    .. code-block:: json

        {
            "path": "CITY",
            "op": "equals",
            "value": "New York",
            "targetpath": "ZIP",
            "targetvalue": "NY-${ZIP}"
        }

Resulting Message:
    .. code-block:: json

        {
            "content": {
                "CITY": "New York",
                "ZIP": "NY-10006",
            }
        }
