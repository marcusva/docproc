# docproc

docproc is a simple content processing pipeline, which allows you to take
arbitrary input data and to transform it to create output data of any kind.

docproc consists of a set of applications, which allow you to perform different
transformation steps one after each other to achieve the desired result. Its
design is based on the functional steps to be taken to get useful output out of
raw data and can be described as follows:

1. consume input
2. process output based on technical and functional requirements for the desired
   output
3. output the processed content as necessary

To enable scalability, each of those functional steps can be handled by an
separate application of docproc. The applications are connected by message
queues, they read from and write to. This allows you to scale individual parts
or complete processing pipelines as required by your input and output scenarios.
