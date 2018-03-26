.. _rulesengine:

Operating the Rules Engine
==========================

docproc ships with a small, easy-to-use rules engine that allows many of its
builtin processors to behave in certain ways, according to your message
contents. The rules are, if not stated otherwise executed against the message's
``content`` section.

Since docproc uses JSON heavily, rules are also expressed in a JSON notation
and are organised as a simple JSON array::

    [
        {
            ... # rule 1
        },
        {
            ... # rule 2
        },
        {
            ... # rule 3
        }
    ]

Rule Configuration
------------------

A rule typically consists of the following fields.

.. code-block:: json

    {
        "name":     "<optional name of the rule>",
        "path":     "<message content path to use for comparing or checking>",
        "op":       "<operator to use>",
        "value":    "<value to use for comparison>",
        "subrules": [ "<more nested rules>" ]
    }

name
   An optional name describing the rule. This is for maintenance purposes and
   does not have any effect on the rule, if provided or absent.

path
   The message's content element to check. Paths can be nested using a dotted
   notation.

op
   The comparison operator to use. If not stated otherwise, the comparison
   will consider path being the left-hand and value the right-hand argument::

        value-of-path <op> rule-value

value
    The value to compare the path's value against. **value** can be omitted, if
    the comparison operator is ``exists`` or ``not exists``. If it is provided
    for those operators, it will be ignored.

subrules
    A list of additional rules that have have to be tested. The rule as well
    as all its sub-rules have to match successfully to consider the rule as a
    whole as successful.

    Please note that all **subrules** are evaluated before the rule itself is
    evaluated. Thus, the most inner subrule is the first being tested.

Setting Paths
-------------

Paths are always relative to the message's content element and can use a
dotted notation to access nested elements of a message. It is also possible
to access array values using brackets and the required index number.

Let's have a look at a few examples of configuring proper paths for rules.
Given the following message

.. code-block:: json

    {
        "content": {
            "name": "John Doe",
            "age": 30,
            "address": {
                "street": "Example Lane 123",
                "zip": "10026",
                "city": "New York"
            },
            "netValues": [
                1000.00,
                453.00,
                -102.00,
                2
            ]
        }
    }

you can access and check the age of John Doe being greater than 20 via

.. code-block:: json

    {
        "path": "age",
        "op": ">",
        "value": 20
    }

Accessing nested elements is done by connecting the element and its sub-element
with a dot. To check, if an address exists and if its city is New York, you can
use ``address.city``.

.. code-block:: json

    {
        "path": "address.city",
        "op": "eq",
        "value": "New York",
        "subrules": [
            {
                "path": "address",
                "op": "exists"
            }
        ]

    }

.. note::
    Subrules are evaluated before the rule itself is evaluated. Thus, if you
    think of multiple conditions that have to apply, you have to build them
    in a reverse order::

        1st condition:      if an address exists
        2nd condition:      and if its city name is "New York"

    thus becomes::

        2nd (outer) rule:   and if its city name is "New York"
        1st (inner) rule:   if an address exists

    Make use of ``name`` to explain more complex rules to keep your maintenance
    efforts at a minimum.

You can access array values using brackets ``[]`` and the value's index.
Indexing starts at zero, so that the first element can be accessed by ``[0]``,
the second by ``[1]`` and so on.

.. code-block:: json

    {
        "path": "netValues[2]",
        "op": ">=",
        "value": 500,
    }


Operators
---------

Existence
    To check, if a given path of a message exists (it may contain nil values
    or empty strings) or not, use the ``exists`` and ``not exists`` operators:

    .. code-block:: json

        {
            "path": "address",
            "op": "exists"
        }

        {
            "path": "alternativeName",
            "op": "not exists"
        }

    Any value provided on the rule, will be ignored.

Equality
    The following operators check, if the provided values are equal:

    ``=``, ``==``, ``eq``, ``equals``

    .. code-block:: json

        {
            "path": "name",
            "op": "=",
            "value": "John Doe",
        }

    Their counterparts, to check for inequality, are:

    ``<>``, ``!=``, ``neq``, ``not equals``

    .. code-block:: json

        {
            "path": "name",
            "op": "neq",
            "value": "Jane Janeson",
        }

Size Comparators
    Values can also be compared by size. This is straightforward for numeric
    values. If you use size comparators on strings, please note that the strings
    are compared lexicographically.

    To check, if the left-hand value is greater than the right-hand value:

    ``>``, ``gt``, ``greater than``

    .. code-block:: json

        {
            "path": "age",
            "op": ">",
            "value": 21,
        }

    To check, if the left-hand value is greater than *or equal to* the
    right-hand value:

    ``>=``, ``gte``, ``greater than or equals``

    .. code-block:: json

        {
            "path": "netValues[0]",
            "op": "gte",
            "value": 500,
        }

    Their counterparts for checking the other way around:

    ``<``, ``lt``, ``less than``

    .. code-block:: json

        {
            "path": "age",
            "op": "<",
            "value": 50,
        }

    and

    ``<=``, ``lte``, ``less than or equals``

    .. code-block:: json

        {
            "path": "netValues[3]",
            "op": "less than or equals",
            "value": -1.0,
        }

String Matching
    To check, if a string contains another string or not, use the following
    operators:

    ``contains``, ``not contains``

    .. code-block:: json

        {
            "path": "name",
            "op": "contains",
            "value": "Doe",
        }
        {
            "path": "name",
            "op": "not contains",
            "value": "Jane",
        }

    As for the size comparators, this checks, if the left-hand value contains
    the right-hand value. To check the other way around, use

    ``in``, ``not in``

    instead.

    .. code-block:: json

        {
            "path": "name",
            "op": "in",
            "value": "John Doe, Jane Doe, or their kids",
        }

        {
            "path": "address.city",
            "op": "not in",
            "value": "London, Vancouver, Washington, Halifax",
        }
