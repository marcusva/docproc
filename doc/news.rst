Release News
============
This describes the latest changes between the docproc releases.

0.0.2
-----
Released on 2018-XX-XX.

* New :ref:`CommandProc` processor for executing external commands on message
  content.
* Changed in/output identifiers for processors to use more meaningful names.

  * FileWriter, HTTPSender: ``'identifer'`` -> ``'read.from'``
  * TemplateTransformer, HtmlRenderer: ``'identifier'`` -> ``'store.in'``

* Fixed ContentValidator creation bug.
* Fixed a panic on reading empty CSV files in CSVTransformer.

0.0.1
-----
Released on 2018-03-25.

* Initial Release
