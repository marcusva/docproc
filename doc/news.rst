Release News
============
This describes the latest changes between the docproc releases.

0.0.3
-----
Released on 2018-XX-XX.

* New ``store.base64`` flag for :ref:`CommandProc` to store binary command
  output in messages.
* New ``read.base64`` flag for :ref:`FileWriter` to write binary content, that
  is stored as base64-encoded string in the message.
* Added Apache FOP examples to produce PDF files.
* Fixed :ref:`CommandProc` usage in configuration files.

0.0.2
-----
Released on 2018-03-30.

* New :ref:`CommandProc` processor for executing external commands on message
  content.
* Changed in/output identifiers for processors to use more meaningful names.

  * :ref:`FileWriter`, :ref:`HTTPSender`: ``'identifer'`` -> ``'read.from'``
  * :ref:`TemplateTransformer`, :ref:`HTMLRenderer`: ``'identifier'`` -> ``'store.in'``

* Fixed :ref:`ContentValidator` creation bug.
* Fixed a panic on reading empty CSV files in :ref:`CSVTransformer`.

0.0.1
-----
Released on 2018-03-25.

* Initial Release
