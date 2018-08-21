Release News
============
This describes the latest changes between the docproc releases.

0.0.4
-----
Released on 2018-XX-XX.

* New ref:`webinput` application to handle raw JSON input via HTTP and file
  uploads to forward to the ``fileinput`` applicaton
* Updated beanstalkd dependency location

0.0.3
-----
Released on 2018-04-18.

* New ``store.base64`` flag for :ref:`CommandProc` to store binary command
  output in messages.
* New ``read.base64`` flag for :ref:`FileWriter` to write binary content, that
  is stored as base64-encoded string in the message.
* New :ref:`PerformanceChecker` processor to measure the processing times of
  messages.
* Added Apache FOP examples to produce PDF files.
* Fixed :ref:`CommandProc` usage in configuration files.

* :ref:`fileinput` consumes less memory on processing large CSV files now.
* *Test* configurations using  ``memory`` as queue implementation consume
  messages concurrently now. This can be tweaked via ``GOMAXPROCS``.
* Configurations using NSQ as queue implementation consume messages
  concurrently now. This can be tweaked via ``GOMAXPROCS``.

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
