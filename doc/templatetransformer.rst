.. _templatetransformer:

TemplateTransformer
===================

The TemplateTransformer provides templating support via Go's `text/template`_
package. This allows docproc to create complex, message-dependent content to be
stored in the message itself.

Configuration
-------------
The TemplateTransformer requires the following configuration entries:

.. code-block:: ini

    [templatetransformer-config]
    type = TemplateTransformer
    identifier = path_to_store
    templates = /path/to/all/templates/*.tpl
    templateroot = main

type
    To configure a TemplateTransformer, use ``TemplateTransformer`` as ``type``.

identifier
    The path to use on the message's content to store the transformed result in.

templates
    Location of the template files on disk. This should be a glob pattern match.

templateroot
    The template entry point to use (``{{define "your-templateroot" }}``).

.. _text/template: https://golang.org/pkg/text/template/