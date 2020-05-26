# docproc - complex content processing made easy

[![Build Status](https://travis-ci.org/marcusva/docproc.svg?branch=master)](https://travis-ci.org/marcusva/docproc)

docproc is a simple content processing pipeline, which allows you to take
arbitrary input data and to transform it to create output data of any kind.

docproc consists of a set of applications, which allow you to perform different
transformation steps one after each other to achieve the desired result. Its
design is based on the functional steps to be taken to get useful output out of
raw data and can be described as follows:

1. consume input content
2. process content based on technical and functional requirements for the
   desired output
3. output the processed content as necessary

![Simple docproc processing layout](https://github.com/marcusva/docproc/blob/master/doc/images/docproc_simple.png "Simple docproc processing layout")

To enable scalability, each of those functional steps can be handled by an
separate application of docproc. The applications are connected via message
queues, they read from and write to. This allows you to scale individual parts
or complete processing pipelines as required by your input and output scenarios.

![Scaled file input, processing and output scenarios](https://github.com/marcusva/docproc/blob/master/doc/images/docproc_scenarios.gif "Scaled file input, processing and output scenarios")

## Features

docproc provides a rich set of features to process content in CSV, SAP RDI and
JSON formats, provided via file exchange or HTTP:

* validation and content enrichment using a simple to maintain rules engine
* text-driven transformation through golang's mighty templating packages, such
  as HTML, XML, JSON, plain text and others
* transforming content easily through external commands
* HTTP transfer, message queue and file-based output

Since docproc uses a simple JSON-based message format internally, applying your
own transformation routines via message queue consumers, HTTP receivers or file
listeners is easily accomplished.

## Documentation

You can find the documentation at doc/html or online at
https://docproc.readthedocs.org.
