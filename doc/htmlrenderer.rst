HTMLRenderer
============

The HTMLRenderer provides templating support via Go's `html/template`_
package. This allows docproc to create complex, message-dependent content to be
stored in the message itself.
It is similar to the :doc:`templatetransformer`, except that ``html/template``
contains some builtin safety nets for HTML content.

Configuration
-------------
The HTMLRenderer requires the following configuration entries:

.. code-block:: ini

    [htmlrenderer-config]
    type = HTMLRenderer
    identifier = path_to_store
    templates = /path/to/all/templates/*.tpl
    templateroot = main

type
   To configure a HTMLRenderer, use ``HTMLRenderer`` as ``type``.

identifier
    The path to use on the message's content to store the transformed result in.

templates
    Location of the template files on disk. This should be a glob pattern match.

templateroot
    The template entry point to use (``{{define "your-entrypoint" }}``).

.. _html/template: https://golang.org/pkg/html/template/