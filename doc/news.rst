Release News
============
This describes the latest changes between the docproc releases.

0.0.2
-----
Released on 2018-XX-XX.

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
